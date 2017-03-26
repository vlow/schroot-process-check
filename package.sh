#!/bin/bash
VERSION="1.1"

SOURCEDIR=$(pwd)
TEMPDIR=$(mktemp -d)
GO_TEMPPATH=$(mktemp -d)

# build binary
(export GOPATH=$GO_TEMPPATH; go get github.com/go-ini/ini)
(export GOPATH=$GO_TEMPPATH; go build $SOURCEDIR/main.go)

# create target layout
BIN_DIR=$TEMPDIR/usr/bin

mkdir -p $BIN_DIR
chmod -R 755 $TEMPDIR/usr

# copy files to target layout
cp -r $SOURCEDIR/main $BIN_DIR/schroot-process-check
chmod 6711 $BIN_DIR/schroot-process-check

# build packages
fpm -s dir -t rpm -n schroot-process-check -v $VERSION --after-install set-permission.sh -C $TEMPDIR .
fpm -s dir -t deb -n schroot-process-check -v $VERSION --after-install set-permission.sh -C $TEMPDIR .
fpm -s dir -t pacman -n schroot-process-check -v $VERSION --after-install set-permission.sh -C $TEMPDIR .

rm -rf $TEMPDIR
rm -rf $GOPATH
