package main

import

// flags "github.com/jessevdk/go-flags"

(
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/agl/ed25519"
	"github.com/boltdb/bolt"
	"github.com/jtremback/scrooge/neighbor-api"
	"github.com/jtremback/scrooge/types"
)

func main() {

	listen := flag.Bool("l", false, "Listen for hellos")

	controlAddress := flag.String("controlAddress", "", "Control address to listen for communication from other nodes.")
	publicKey := flag.String("publicKey", "", "PublicKey to sign messages to other nodes.")
	privateKey := flag.String("privateKey", "", "PrivateKey to sign messages to other nodes.")

	ifi := flag.String("interface", "", "Physical network interface to operate on.")
	tunnelPublicKey := flag.String("tunnelPublicKey", "", "PublicKey of authenticated tunnel")
	tunnelPrivateKey := flag.String("tunnelPrivateKey", "", "PrivateKey of authenticated tunnel")

	flag.Parse()

	iface, err := net.InterfaceByName(*ifi)
	if err != nil {
		log.Fatalln(err)
	}

	pubKey, err := base64.StdEncoding.DecodeString(*publicKey)
	if err != nil {
		log.Fatalln(err)
	}

	privKey, err := base64.StdEncoding.DecodeString(*privateKey)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := bolt.Open("main.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	neighborAPI := neighborAPI.NeighborAPI{
		Neighbors: map[[ed25519.PublicKeySize]byte]*types.Neighbor{},
		Account: &types.Account{
			PublicKey:        types.BytesToPublicKey(pubKey),
			PrivateKey:       types.BytesToPrivateKey(privKey),
			ControlAddress:   *controlAddress,
			TunnelPublicKey:  *tunnelPublicKey,
			TunnelPrivateKey: *tunnelPrivateKey,
			Seqnum:           0,
		},
	}

	if *listen {
		log.Println("listen")
		err := neighborAPI.McastListen(
			8481,
			iface,
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
			iface,
			func(neighbor *types.Neighbor, err error) {
				if err != nil {
					log.Fatalln(err)
				}
			},
		)
	}
}
