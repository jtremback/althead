#!bash

export GOPATH=/vagrant

# sudo modprobe dummy numdummies=1
# sudo ip addr add 10.0.1.2/24 dev "$(sudo brctl show | awk '/dummy0/ { print $1 }')"
/usr/local/go/bin/go build && /usr/bin/python ./test/smoketest.py