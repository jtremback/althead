package findPeersMCast

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type Service struct {
	Denom        string
	Rate         int
	TunnelIP     net.IP
	TunnelPort   int
	TunnelPubkey string
}

func Advertise(
	iface *net.Interface,
	listenPort int,
) error {

	conn, err := net.ListenMulticastUDP(
		"udp6",
		iface,
		&net.UDPAddr{
			IP:   net.ParseIP("ff02::1"),
			Port: listenPort,
			Zone: iface.Name,
		},
	)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		msg := strings.Split(scanner.Text(), " ")

		if msg[0] == "althea_service_request" {
			fmt.Println("right kind of message")
			conn, err := net.Dial("tcp6", msg[1])
			if err != nil {
				return err
			}
			conn.Write([]byte("derp"))
		}
		fmt.Println(msg)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func QueryPeers(
	iface *net.Interface,
	listenPort int,
	mCastPort int,
) error {
	addrs, err := iface.Addrs()

	var ip net.IP
	for _, a := range addrs {
		addr, _, err := net.ParseCIDR(a.String())
		if err != nil {
			return err
		}

		fmt.Println("addr", addr.String())

		if addr.IsLinkLocalUnicast() {
			ip = net.ParseIP(addr.String())
		}
	}
	if len(ip) == 0 {
		return errors.New("could not find link local ipv6 address for interface")
	}

	// ip, err := func() (*net.IP, error) {
	// 	for _, a := range addrs {
	// 		p := net.ParseIP(a.String())
	// 		if p.IsLinkLocalUnicast() {
	// 			return &p, nil
	// 		}
	// 	}
	// 	return nil, errors.New("could not find link local ipv6 address for interface")
	// }()
	// if err != nil {
	// 	return err
	// }
	fmt.Println("listen", ip.String())
	// l, err := net.Listen(
	// 	"tcp6",
	// 	"["+ip+"]:"+strconv.Itoa(listenPort),
	// )
	l, err := net.ListenTCP("tcp6", &net.TCPAddr{
		IP:   ip,
		Port: listenPort,
		Zone: iface.Name,
	})
	if err != nil {
		return err
	}

	defer l.Close()

	ch := make(chan error)

	go func() {
		for {
			fmt.Println("foo")
			// Wait for a connection.
			conn, err := l.Accept()
			fmt.Println("doo")
			if err != nil {
				ch <- err
			}
			// Handle the connection in a new goroutine.
			// The loop then returns to accepting, so that
			// multiple connections may be served concurrently.
			go func(c net.Conn) {
				// Echo all incoming data.
				io.Copy(os.Stdout, c)
				// Shut down the connection.
				c.Close()
			}(conn)
		}
	}()

	conn, err := net.DialUDP(
		"udp6",
		nil,
		&net.UDPAddr{
			IP:   net.ParseIP("ff02::1"),
			Port: mCastPort,
			Zone: iface.Name,
		})
	if err != nil {
		return err
	}

	fmt.Println(l.Addr())

	conn.Write([]byte("althea_service_request " +
		l.Addr().String() +
		" " +
		strconv.Itoa(listenPort) +
		"\n"))

	err = <-ch
	if err != nil {
		return err
	}

	return nil
}
