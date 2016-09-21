package findPeersBabel

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
)

func Find(babelPort int) ([]net.IP, error) {
	conn, err := net.Dial("tcp", "[::1]:8481")
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(conn)

	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Println(status)

	matched, err := regexp.MatchString("BABEL", status)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("Could not connect to Babel config server.")
	}

	conn.Write([]byte("dump\n"))

	for scanner.Scan() {
		text := scanner.Text()
		matched, err := regexp.MatchString("add neighbour", text)
		if err != nil {
			return nil, err
		}
		if matched {
			re := regexp.MustCompile("address (\\S*) if (\\S*)")
			matches := re.FindStringSubmatch(text)
			fmt.Printf("%v%%%v", matches[1], matches[2])
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return nil, nil
}
