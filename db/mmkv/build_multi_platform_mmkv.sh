#!/bin/bash

# requirements:
# gcc
# clang
# gcc-multilib
# gcc-aarch64-linux-gnu
# gcc-arm-linux-gnueabihf
# gcc-aarch64-linux-gnu
# arm-linux-gnueabihf-gcc
# gcc-mingw-w64-x86-64
# gcc-mingw-w64-i686

# output path
OUTPUT_DIR="./build"
mkdir -p ${OUTPUT_DIR}

platforms=(
  "linux/amd64"
  "linux/arm64"
  "linux/arm"
  "windows/amd64"
  "windows/386"
  "darwin/amd64"
  "darwin/arm64"
)

for platform in "${platforms[@]}"; do
  GOOS=${platform%/*}
  GOARCH=${platform#*/}
  OUTPUT_FILE="${OUTPUT_DIR}/mmkv_${GOOS}_${GOARCH}"

  if [ "${GOOS}" == "windows" ]; then
    OUTPUT_FILE+=".dll"
  elif [ "${GOOS}" == "darwin" ]; then
    OUTPUT_FILE+=".dylib"
  else
    OUTPUT_FILE+=".so"
  fi

  echo "Building for ${GOOS}/${GOARCH}..."
  CC=""
  if [ "${GOOS}" == "linux" ] && [ "${GOARCH}" == "arm" ]; then
    CC="arm-linux-gnueabihf-gcc"
  elif [ "${GOOS}" == "linux" ] && [ "${GOARCH}" == "arm64" ]; then
    CC="aarch64-linux-gnu-gcc"
  elif [ "${GOOS}" == "windows" ] && [ "${GOARCH}" == "386" ]; then
    CC="i686-w64-mingw32-gcc"
  elif [ "${GOOS}" == "windows" ] && [ "${GOARCH}" == "amd64" ]; then
    CC="x86_64-w64-mingw32-gcc"
  elif [ "${GOOS}" == "darwin" ] && [ "${GOARCH}" == "amd64" ]; then
    CC="clang"
  elif [ "${GOOS}" == "darwin" ] && [ "${GOARCH}" == "arm64" ]; then
    CC="clang"
  fi
  GOOS=${GOOS} GOARCH=${GOARCH} CC=${CC} go build -buildmode=c-shared -o ${OUTPUT_FILE} || {
    echo "Failed to build for ${platform}!"
    exit 1
  }
done

echo "Build completed! Files are in ${OUTPUT_DIR}"
