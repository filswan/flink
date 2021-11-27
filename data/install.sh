#!/bin/bash
BINARY_NAME=filink-data-0.2.0-unix
TAG_NAME=v0.2.0

wget https://github.com/filswan/go-swan-provider/releases/download/${TAG_NAME}/${BINARY_NAME}
wget https://github.com/filswan/go-swan-provider/releases/download/${TAG_NAME}/aria2.conf
wget https://github.com/filswan/go-swan-provider/releases/download/${TAG_NAME}/aria2c.service

CONF_FILE_DIR=${HOME}/.swan/filink/data
mkdir -p ${CONF_FILE_DIR}

CONF_FILE_PATH=${CONF_FILE_DIR}/config.toml
echo $CONF_FILE_PATH

if [ -f "${CONF_FILE_PATH}" ]; then
    echo "${CONF_FILE_PATH} exists"
else
    cp ./config/config.toml.example $CONF_FILE_PATH
    echo "${CONF_FILE_PATH} created"
fi

chmod +x ./build/${BINARY_NAME}
./build/${BINARY_NAME}                         # Run swan provider


