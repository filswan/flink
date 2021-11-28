#!/bin/bash

CONF_FILE_DIR=${HOME}/.swan/flink/data
mkdir -p ${CONF_FILE_DIR}

CONF_FILE_PATH=${CONF_FILE_DIR}/config.toml
echo $CONF_FILE_PATH

if [ -f "${CONF_FILE_PATH}" ]; then
    echo "${CONF_FILE_PATH} exists"
else
    cp ./config/config.toml.example $CONF_FILE_PATH
    echo "${CONF_FILE_PATH} created"
fi

BINARY_NAME=flink-data
make
chmod +x ./build/${BINARY_NAME}
./build/${BINARY_NAME}                         # Run swan provider


