package wireguard

import (
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"errors"

	"github.com/incentivized-mesh-infrastructure/scrooge/types"
)

func CreateTunnel(
	tunnel *types.Tunnel,
	tunnelAddress net.UDPAddr,
	tunnelPrivateKey string,
) error {
	exec.Command("ip", "link", "add", "dev", tunnel.VirtualInterface.Name, "type", "wireguard").Run()
	exec.Command("ip", "address", "add", "dev", tunnel.VirtualInterface.Name, tunnelAddress.IP.String()).Run()

	privateKeyFile, err := ioutil.TempFile("", "example")
	if err != nil {
		return err
	}

	defer os.Remove(privateKeyFile.Name()) // clean up

	privateKeyFile.Chmod(0700)
	privateKeyFile.Chown(0, 0)

	_, err = privateKeyFile.Write([]byte(tunnelPrivateKey))
	if err != nil {
		return err
	}

	err = privateKeyFile.Close()
	if err != nil {
		return err
	}

	exec.Command(
		"wg", "set", tunnel.VirtualInterface.Name,
		"listen-port", string(tunnelAddress.Port),
		"private-key", privateKeyFile.Name(),
		"peer", tunnel.PublicKey,
		"allowed-ips", "0.0.0.0",
		"endpoint", tunnel.Endpoint,
	).Run()

	exec.Command("ip", "link", "set", "up", tunnel.VirtualInterface.Name).Run()

	out, err := exec.Command("wg", "showconf", tunnel.VirtualInterface.Name).Output()
	if err != nil {
		return err
	}

	config := parseConfig(string(out))

	if config.PrivateKey != tunnelPrivateKey ||
		config.ListenPort != tunnelAddress.Port {
		return errors.New("")
	}

	return nil
}

type WireguardConfig struct {
	PrivateKey string
	ListenPort int
	Peer       struct {
		PublicKey  string
		AllowedIPs []net.IP
		Endpoint   string
	}
}

func parseConfig(s string) WireguardConfig {
	var config WireguardConfig

	config.PrivateKey = findFirstSubmatch(s, "PrivateKey")
	config.ListenPort = findFirstSubmatch(s, "ListenPort")
	config.Peer.PublicKey = findFirstSubmatch(s, "PublicKey")
	config.Peer.AllowedIPs = parseIPs(findFirstSubmatch(s, "AllowedIPs"))
	config.Peer.Endpoint = findFirstSubmatch(s, "Endpoint")

	return config
}

func parseIPs(s string) []net.IP {
	fields := strings.Split(s, ",")

	var ips []net.IP
	for _, field := range fields {
		ips = append(ips, net.ParseIP(field))
	}

	return ips
}

func findFirstSubmatch(s string, name string) string {
	re := regexp.MustCompile(name + " = (.*)")
	res := re.FindAllStringSubmatch(s, 1)
	return res[0][1]
}

// func parseConfig(config string) {
// 	re := regexp.MustCompile(`\[Interface\]`)
// 	idxs := re.FindAllIndex([]byte(config), -1)
// 	// strs := re.FindAllString(config, -1)

// 	var prevEnd int

// 	for _, idx := range idxs {
// 		start := idx[0]
// 		end := idx[1]

// 		var body string
// 		if prevEnd != 0 {
// 			body = config[prevEnd:end]
// 		}
// 		prevEnd = end
// 	}
// }

// [Interface]
// PrivateKey = yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=
// ListenPort = 51820

// [Peer]
// PublicKey = xTIBA5rboUvnH4htodjb6e697QjLERt1NAB4mZqp8Dg=
// Endpoint = 192.95.5.67:1234
// AllowedIPs = 10.192.122.3/32, 10.192.124.1/24

// [Interface]
// PrivateKey = yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=
// ListenPort = 51820

// [Peer]
// PublicKey = xTIBA5rboUvnH4htodjb6e697QjLERt1NAB4mZqp8Dg=
// Endpoint = 192.95.5.67:1234
// AllowedIPs = 10.192.122.3/32, 10.192.124.1/24
