#!/bin/bash

go build
rm -rf cache
./maya -mode=empty -file=README.tpl.md > README.md
