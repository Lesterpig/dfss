#!/bin/sh

go get -u gopkg.in/mgo.v2
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u google.golang.org/grpc
go get -u github.com/pborman/uuid
go get -u github.com/stretchr/testify/assert
go get -u golang.org/x/crypto/ssh/terminal
go get -u github.com/spf13/viper
go get -u github.com/spf13/cobra

go get -u github.com/inconshreveable/mousetrap # required by cobra for win builds
