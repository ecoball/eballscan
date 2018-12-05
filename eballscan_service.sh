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

SOURCE_DIR=$(cd `dirname $0` && pwd)

# check cockroachdb
if [ ! -e "/usr/local/bin/cockroach" ]
then
    echo -e "\033[;31mPlease install the cockroachdb first!!! \033[0m"
    exit 1
fi

# check eballscan
if [ ! -e "${SOURCE_DIR}/build/eballscan" ] 
then
    echo -e "\033[;31mThe eballscan does not exist!!! \033[0m"
    exit 1
fi

# operation of database
if [ ! -e "${SOURCE_DIR}/build/cockroach-data" -o 2 -ne $(ps -ef | grep cockroach | wc -l) ]
then
    # remove old data
    if [ -e "${SOURCE_DIR}/build/cockroach-data" ]
    then
        rm -fr "${SOURCE_DIR}/build/cockroach-data"
    fi

    # stop cockroachdb
    if [ 2 -eq $(ps -ef | grep cockroach | wc -l) ]
    then
        killall cockroach
	    sleep 5s 
    fi

    # start cockroachdb
    CURRENT_DIR=$(pwd)
    cd ${SOURCE_DIR}/build/ && cockroach start --insecure --http-port=8081 --background && cd $CURRENT_DIR
    if [ 0 -ne $? ]; then
        echo  -e "\033[;31m start cockroach failed!!! \033[0m"
        exit 1
    fi

    # create user eballscan
    cockroach user set eballscan --insecure
    if [ 0 -ne $? ]; then
        echo  -e "\033[;31m create user eballscan failed!!! \033[0m"
        exit 1
    fi

    # create database blockchain
    cockroach sql --insecure -e 'create database blockchain'
    if [ 0 -ne $? ]; then
        echo  -e "\033[;31m create database blockchain failed!!! \033[0m"
        exit 1
    fi

    # grant eballscan
    cockroach sql --insecure -e 'GRANT ALL ON DATABASE blockchain TO eballscan'
    if [ 0 -ne $? ]; then
        echo  -e "\033[;31m grant eballscan failed!!! \033[0m"
        exit 1
    fi
fi

# start eballscan
case $# in
    0)
    ${SOURCE_DIR}/build/eballscan start
    ;;

    1)
    ${SOURCE_DIR}/build/eballscan start -i $1
    ;;

    2)
    ${SOURCE_DIR}/build/eballscan start -i $1 -p $2
    ;;

    *)
    echo "please input eballscan_service | eballscan_service param1(ecoball-ip) | eballscan_service param1(ecoball-ip) param2(ecoball-bystander-port)"
    ;;
esac
