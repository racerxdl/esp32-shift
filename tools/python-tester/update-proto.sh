#!/bin/bash

protoc --proto_path=../../pkg/proto --python_out=. shift.proto
