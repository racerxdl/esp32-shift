#!/bin/bash

protoc --proto_path=../../pkg/proto --plugin=protoc-gen-nanopb=../nanopb-gen/protoc-gen-nanopb --nanopb_out=. shift.proto
