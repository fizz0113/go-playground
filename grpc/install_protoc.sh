#!/bin/bash

set -e

OS_NAME="$(uname -s)"
OS_ARCH="$(uname -m)"
PROTOC_VERSION="3.7.1"
PROTOC_ZIP_NAME="protoc-${PROTOC_VERSION}-${OS_NAME}-${OS_ARCH}.zip"

echo "Download ${PROTOC_ZIP_NAME}"
if [ ! -f ${PROTOC_ZIP_NAME} ]; then
  curl -sOL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP_NAME}
fi

echo "Unzip protoc binary to /usr/local/bin"
unzip -o ${PROTOC_ZIP_NAME} -d /usr/local bin/protoc
echo "Unzip protoc type to /usr/local/include"
unzip -o ${PROTOC_ZIP_NAME} -d /usr/local 'include/*'

echo "Remove download file"
rm -f ${PROTOC_ZIP_NAME}

echo "Done"
