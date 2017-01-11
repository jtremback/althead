package main

import

// flags "github.com/jessevdk/go-flags"

(
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/boltdb/bolt"
	"github.com/jtremback/althea/neighbor-api"
	"github.com/jtremback/althea/types"
)

func main() {

	listen := flag.Bool("l", false, "Listen for hellos")

	controlAddress := flag.String("controlAddress", "", "Control address to listen for communication from other nodes.")
	controlPubkey := flag.String("controlPubkey", "", "Pubkey to sign messages to other nodes.")
	controlPrivkey := flag.String("controlPrivkey", "", "Privkey to sign messages to other nodes.")

	ifi := flag.String("interface", "", "Physical network interface to operate on.")
	tunnelAddress := flag.String("tunnelAddress", "", "Address that authenticated tunnel will be served on.")
	tunnelPubkey := flag.String("tunnelPubkey", "", "Pubkey of authenticated tunnel")
	tunnelPrivkey := flag.String("tunnelPrivkey", "", "Privkey of authenticated tunnel")

	flag.Parse()

	iface, err := net.InterfaceByName(*ifi)
	if err != nil {
		log.Fatalln(err)
	}

	cpubkey, err := base64.StdEncoding.DecodeString(*controlPubkey)
	if err != nil {
		log.Fatalln(err)
	}

	cprivkey, err := base64.StdEncoding.DecodeString(*controlPrivkey)
	if err != nil {
		log.Fatalln(err)
	}

	tpubkey, err := base64.StdEncoding.DecodeString(*tunnelPubkey)
	if err != nil {
		log.Fatalln(err)
	}

	tprivkey, err := base64.StdEncoding.DecodeString(*tunnelPrivkey)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := bolt.Open("main.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	tunnel := &types.Tunnel{
		ControlAddress: *controlAddress,
		ControlPubkey:  cpubkey,
		ControlPrivkey: cprivkey,
		TunnelAddress:  *tunnelAddress,
		TunnelPubkey:   tpubkey,
		TunnelPrivkey:  tprivkey,
		Interface:      iface,
	}

	neighborAPI := neighborAPI.NeighborAPI{
		DB:     db,
		Tunnel: tunnel,
	}

	if *listen {
		log.Println("listen")
		err := neighborAPI.McastListen(
			8481,
			func(neighbor *types.Neighbor, err error) {
				if err != nil {
					log.Fatalln(err)
				}
			},
		)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println("hello")
		neighborAPI.McastHello(
			8481,
			func(neighbor *types.Neighbor, err error) {
				if err != nil {
					log.Fatalln(err)
				}
			},
		)
	}
}
