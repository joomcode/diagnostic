#!/bin/bash -ex
find -iname "*.go" -exec goimports -w {} \;
