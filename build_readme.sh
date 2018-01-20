#!/bin/bash

cd maya-cli; go build cd ..
rm -rf cache
./maya-cli/maya-cli -mode=empty -file=README.tpl.md > README.md
