#!/bin/bash

go build
./maya -mode=pelican-md -file=README.tpl.md > README.md
