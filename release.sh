#!/bin/sh
#
# 2017-05-28 adbr

prog=pogoda
version=1.0.1

build() {
	os=$1
	arch=$2
	ext=""
	if [ $os = "windows" ]; then
	    ext=".exe"
	fi
	dir=$prog-$version-$os-$arch
	bdir=build/$dir

	echo "building: $bdir"
	mkdir -p $bdir
	env CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o $bdir/pogoda$ext
	cp README.org $bdir
	cp LICENCE $bdir

	mkdir -p release
	tarfile=$dir.tar.gz
	echo "creating: $tarfile"
	(cd build; tar -czf $tarfile $dir)
	mv build/$tarfile release
}

build openbsd amd64
build windows amd64
build linux amd64
build linux arm
