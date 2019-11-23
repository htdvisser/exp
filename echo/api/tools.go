//+build tools

package api

import (
	_ "github.com/envoyproxy/protoc-gen-validate"                      // Tool dependenncy.
	_ "github.com/gogo/protobuf/protoc-gen-gogofast"                   // Tool dependenncy.
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway" // Tool dependenncy.
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"      // Tool dependenncy.
	_ "github.com/mdempsky/unconvert"                                  // Tool dependenncy.
	_ "mvdan.cc/gofumpt"                                               // Tool dependenncy.
)
