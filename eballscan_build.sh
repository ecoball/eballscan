#!/bin/bash
##########################################################################
# Copyright 2018 The eballscan Authors
# This file is part of the eballscan.
#
# The eballscan is free software: you can redistribute it and/or modify
# it under the terms of the GNU Lesser General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# The eballscan is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Lesser General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public License
# along with the eballscan. If not, see <http://www.gnu.org/licenses/>.
############################################################################
#install cockroachdb
wget -qO- https://binaries.cockroachdb.com/cockroach-v2.0.4.linux-amd64.tgz | tar  xvz
if [ 0 -ne $? ]; then
    echo  -e "\033[;31m Unable to download cockroach-v2.0.4.linux-amd64.tgz at this time!!! \033[0m"
    exit 1
fi

sudo cp -i cockroach-v2.0.4.linux-amd64/cockroach /usr/local/bin
if [ 0 -ne $? ]; then
    echo  -e "\033[;31m install cockroach-v2.0.4.linux-amd64 failed!!! \033[0m"
    exit 1
fi

if ! rm -fr "./cockroach-v2.0.4.linux-amd64"
then
    echo  -e "\033[;31m delete cockroach-v2.0.4.linux-amd64 failed!!! \033[0m"
    exit 1
fi

#start cockroachdb
if ! mkdir -p ./build/cockroachdb/store ./build/cockroachdb/log
then
    echo -e "\033[;31m create directory failed!!! \033[0m"
    exit 1
fi

cockroach start --insecure --http-port=8081 --background --store=../build/cockroachdb/store --log-dir=./build/cockroachdb/log
if [ 0 -ne $? ]; then
    echo  -e "\033[;31m start cockroach failed!!! \033[0m"
    exit 1
fi

#create user eballscan
cockroach user set eballscan --insecure
if [ 0 -ne $? ]; then
    echo  -e "\033[;31m create user eballscan failed!!! \033[0m"
    exit 1
fi

#create database blockchain
cockroach sql --insecure -e 'create database blockchain'
if [ 0 -ne $? ]; then
    echo  -e "\033[;31m create database blockchain failed!!! \033[0m"
    exit 1
fi

#grant eballscan
cockroach sql --insecure -e 'GRANT ALL ON DATABASE blockchain TO eballscan'
if [ 0 -ne $? ]; then
    echo  -e "\033[;31m grant eballscan failed!!! \033[0m"
    exit 1
fi

#build project
if ! make
then
    echo  -e "\033[;31m compile eballscan failed!!! \033[0m"
    exit 1
fi

echo -e "\033[;32m build eballscan succeed\033[0m"
