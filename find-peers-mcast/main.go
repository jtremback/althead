package findPeersMCast

import (
	"bufio"
	"fmt"
	"net"
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
			// Make
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
	// l, err := net.Listen(
	// 	"tcp6",
	// 	":"+strconv.Itoa(listenPort),
	// )
	// if err != nil {
	// 	return err
	// }

	// defer l.Close()
	// for {
	// 	// Wait for a connection.
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// Handle the connection in a new goroutine.
	// 	// The loop then returns to accepting, so that
	// 	// multiple connections may be served concurrently.
	// 	go func(c net.Conn) {
	// 		// Echo all incoming data.
	// 		io.Copy(os.Stdout, c)
	// 		// Shut down the connection.
	// 		c.Close()
	// 	}(conn)
	// }

	client, err := net.DialUDP(
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

	client.Write([]byte("althea_service_request " +
		/*l.Addr().String()*/ "shibby" +
		" " +
		strconv.Itoa(listenPort) +
		"\n"))

	return nil
}
