package main

import

// flags "github.com/jessevdk/go-flags"

(
	"flag"
	"log"
	"net"

	"github.com/jtremback/althea/find-peers-mcast"
)

func main() {

	server := flag.Bool("s", false, "Run a server")

	flag.Parse()

	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		log.Fatalln(err)
	}

	if *server {
		log.Println("listen")
		err := findPeersMCast.Listen(
			iface,
			8481,
		)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println("hello")
		findPeersMCast.Hello(
			iface,
			8481,
			func(ip net.IP, err error) {
				if err != nil {
					log.Fatalln(err)
				}
			},
		)
	}
}
