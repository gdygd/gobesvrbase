#!/bin/bash

cd ./apisvr
make build
cd ..

cd ./cli
make build
cd ..


cd ./mngprc
make build
