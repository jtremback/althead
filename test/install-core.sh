#!bash
# This should install the CORE network simulator on ubuntu and allow you to run tests.

apt-get update
apt-get install -y wget automake gcc pkg-config make libev-dev python-dev bridge-utils ebtables
wget -O - https://github.com/coreemu/core/archive/release-4.8.tar.gz | tar zxf -
cd core-release-4.8 || exit
./bootstrap.sh
./configure
make 
make install