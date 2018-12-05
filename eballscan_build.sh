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

# clean old install
if [ -d "${SOURCE_DIR}/build" ]; then
    echo -e "\033[;33m old eballscan and cockroach-data needs to be removed.\033[0m"
    echo -e "\033[;33m Do you wish to remove this install?\033[0m"
    select yn in "Yes" "No"; do
        echo $yn
        case $yn in
            [Yy]* )
                if ! rm -fr "${SOURCE_DIR}/build/"
                then
                    echo  -e "\033[;31m remove ./build/ failed!!! \033[0m"
                    exit 1
                fi
                break;;
            [Nn]* )
                echo -e "\033[;33m Aborting uninstall\033[0m"
                exit 0;;
        esac
   done
fi

# install cockroachdb
if [ ! -e "/usr/local/bin/cockroach" ]; then
    wget -qO- https://binaries.cockroachdb.com/cockroach-v2.0.6.linux-amd64.tgz | tar  xvz
    if [ 0 -ne $? ]; then
        echo  -e "\033[;31m Unable to download cockroach-v2.0.6.linux-amd64.tgz at this time!!! \033[0m"
        exit 1
    fi

    sudo cp -i cockroach-v2.0.6.linux-amd64/cockroach /usr/local/bin
    if [ 0 -ne $? ]; then
        echo  -e "\033[;31m install cockroach-v2.0.6.linux-amd64 failed!!! \033[0m"
        exit 1
    fi

    if ! rm -fr "./cockroach-v2.0.6.linux-amd64"
    then
        echo  -e "\033[;31m remove cockroach-v2.0.6.linux-amd64 failed!!! \033[0m"
        exit 1
    fi
fi

# build project
if ! make -C ${SOURCE_DIR}
then
    echo  -e "\033[;31m compile eballscan failed!!! \033[0m"
    exit 1
fi

echo -e "\033[;32mbuild eballscan succeed\033[0m"
