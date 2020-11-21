#!/bin/bash

version=$1
packages=(
  "godbledger-linux-arm-5"
  # "godbledger-linux-arm-7"
  # "godbledger-linux-arm-64"
  "godbledger-darwin-10.6-amd64"
  # "godbledger-darwin-10.6-amd64"
)

if [ -z "$version" ]
then
  echo "Usage: $0 <version>"
  exit
fi

make VERSION=$version release

WORKING_DIR=bin/release/

echo "Working in $WORKING_DIR..."

mkdir -p $WORKING_DIR
cd $WORKING_DIR

for package in packages; do
  tar -czvf "$package-v$version.tar.gz" "$package"
done

echo '#### sha256sum'
sha256sum godbledger-*-v$version.tar.gz
