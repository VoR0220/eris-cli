#!/usr/bin/env bash
# ----------------------------------------------------------
# PURPOSE

# This is a helper script to get a cleaned vendor and to work around 
# the problem of the vendor cleaner plugin deleting a required crypto package

# ----------------------------------------------------------
# REQUIREMENTS

# glide installed locally

# ----------------------------------------------------------
# USAGE

# ./cleanVendor.sh in the root directory

go get -u github.com/sgotti/glide-vc
go get -u github.com/ethereum/go-ethereum
glide vc
rm -rf vendor/github.com/ethereum/go-ethereum/crypto/secp256k1
mv $GOPATH/src/github.com/ethereum/go-ethereum/crypto/secp256k1 $GOPATH/src/github.com/eris-ltd/eris-cli/vendor/github.com/ethereum/go-ethereum/crypto/
