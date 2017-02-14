package wireguard

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/incentivized-mesh-infrastructure/scrooge/types"
)

func CreateTunnel(
	tunnel *types.Tunnel,
	tunnelAddress net.UDPAddr,
) {
	exec.Command("ip", "link", "add", "dev", tunnel.VirtualInterface.Name, "type", "wireguard").Run()
	exec.Command("ip", "address", "add", "dev", tunnel.VirtualInterface.Name, tunnelAddress.IP.String()).Run()

	content := []byte("temporary file's content")
	privateKeyFile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(privateKeyFile.Name()) // clean up

	privateKeyFile.Chmod(0700)
	privateKeyFile.Chown(0, 0)

	if _, err := privateKeyFile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := privateKeyFile.Close(); err != nil {
		log.Fatal(err)
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
}
