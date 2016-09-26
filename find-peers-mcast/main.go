package findPeersMCast

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func setup(
	iface *net.Interface,
	mcastPort int,
	secret string,
	pubkey string,
) error {
	conf := fmt.Sprintf(`
bind %v:%v;
mode tap;
method "xsalsa20-poly1305";
mtu 1426;
secret "%v";
interface "%k";

include peers from "/tmp/althea/fastd_peers";
	`, iface, mcastPort, secret)

	confFile, err := ioutil.TempFile("", "althead_fastd_conf")
	if err != nil {
		return err
	}

	defer os.Remove(confFile.Name())

	_, err = confFile.Write([]byte(conf))
	if err != nil {
		return err
	}

	err = confFile.Close()
	if err != nil {
		return err
	}

	peersDir, err := ioutil.TempDir("", "althead_fastd_peers")
	if err != nil {
		return err
	}

}

func Advertise(
	iface *net.Interface,
	mcastPort int,
	secret string,
	pubkey string,
) error {

	conn, err := net.ListenMulticastUDP(
		"udp6",
		iface,
		&net.UDPAddr{
			IP:   net.ParseIP("ff02::1"),
			Port: mcastPort,
			Zone: iface.Name,
		},
	)
	if err != nil {
		return err
	}

	for {
		b := make([]byte, 64)
		_, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			return err
		}

		msg := strings.Split(string(b), " ")

		if msg[0] == "althea_hello" {
			conn, err := net.DialUDP(
				"udp6",
				nil,
				addr,
			)
			if err != nil {
				return err
			}

			// write peer file
			peerFilename := peersDir + msg[1]
			peerConf := []byte(fmt.Sprintf(`key %v;`, msg[1]))
			err = ioutil.WriteFile(peerFilename, peerConf, 0600)
			if err != nil {
				return err
			}

			defer os.Remove(peerFilename)

			conn.Write([]byte(
				"althea_ihu " +
					"",
			))

			conn.Close()

		}
		fmt.Println("got msg", string(b), addr)
	}

	return nil
}

func firstLinkLocalUnicast(iface *net.Interface) (*net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		}

		if ip.IsLinkLocalUnicast() {
			return &ip, nil
		}
	}
	return nil, errors.New("Could not find link local unicast ipv6 address for interface " + iface.Name)
}

func QueryPeers(
	iface *net.Interface,
	mCastPort int,
	cb func(net.IP, error),
) {
	ip, err := firstLinkLocalUnicast(iface)
	if err != nil {
		cb(nil, err)
	}

	conn, err := net.ListenUDP("udp6", &net.UDPAddr{
		IP:   *ip,
		Port: 0,
		Zone: iface.Name,
	})
	if err != nil {
		cb(nil, err)
	}

	defer conn.Close()

	conn.WriteToUDP([]byte("althea_hello "+
		"b023b49ef699da32cef9046434ab5f73d0ea7b8ee4e4bb31b8f4f2e27d138e54\n",
	), &net.UDPAddr{
		IP:   net.ParseIP("ff02::1"),
		Port: mCastPort,
		Zone: iface.Name,
	})

	for {
		b := make([]byte, 10)
		_, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			cb(nil, err)
		}
		fmt.Println(addr.String())
		cb(net.ParseIP(addr.String()), nil)
	}
}
