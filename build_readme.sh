#!/bin/bash

go build
./maya -mode=empty -file=README.tpl.md > README.md
