#!/bin/bash

cd templates/
qtc >/dev/null 2>&1 
cd ..
go run cmd/torch/main.go $@