// Package protoc-gen-hugodata contains a protoc plugin and a Hugo module for building
// documentation for your protos.
//
// Command protoc-gen-hugodata is the protoc plugin that generates yaml files
// that can be used by the shortcodes in this Hugo module.
//
// Typical usage of the protoc plugin:
//
//     protoc -I [your imports ...] --hugodata_out=output_path=path/to/data:path/to/data /path/to/*.proto
//
// To use this as a Hugo module, add the following to your config.toml:
//
//     [module]
//       [[module.imports]]
//         path = "htdvisser.dev/exp/protoc-gen-hugodata"
//
// And then run:
//
//     hugo mod get htdvisser.dev/exp/protoc-gen-hugodata
//
// After this, you can use the shortcodes in your markdown files:
//
//     {{< proto/method package="your.proto.v1" service="ServiceName" method="MethodName" >}}
//     {{< proto/message package="your.proto.v1" message="ServiceMethodRequest" >}}
//     {{< proto/enum package="your.proto.v1" enum="EnumName" >}}
package main

import pgs "github.com/lyft/protoc-gen-star"

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		HugoData(),
	).Render()
}
