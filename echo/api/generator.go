//+build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func run(workdir, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = workdir
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	fmt.Fprintf(os.Stderr, "↪ %s\n", strings.Join(cmd.Args, " "))
	if err := cmd.Run(); err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			os.Exit(err.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func download(from, to string, mode os.FileMode) (err error) {
	if info, err := os.Stat(to); err == nil && info.Mode() == mode {
		return nil
	}
	fmt.Fprintf(os.Stderr, "⬇ download %s\n", to)
	res, err := http.Get(from)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := res.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()
	_, err = io.Copy(f, res.Body)
	return err
}

func main() {
	// Dependencies
	run(".", "go", "install", "github.com/envoyproxy/protoc-gen-validate")
	run(".", "go", "install", "google.golang.org/protobuf/cmd/protoc-gen-go")
	run(".", "go", "install", "google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	run(".", "go", "install", "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway")
	run(".", "go", "install", "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger")
	run(".", "go", "install", "github.com/mdempsky/unconvert")
	run(".", "go", "install", "mvdan.cc/gofumpt")

	download(
		"https://github.com/gogo/protobuf/raw/master/gogoproto/gogo.proto",
		"third_party/github.com/gogo/protobuf/gogoproto/gogo.proto",
		0644,
	)
	download(
		"https://github.com/googleapis/googleapis/raw/master/google/api/annotations.proto",
		"third_party/google/api/annotations.proto",
		0644,
	)
	download(
		"https://github.com/googleapis/googleapis/raw/master/google/api/http.proto",
		"third_party/google/api/http.proto",
		0644,
	)
	download(
		"https://github.com/envoyproxy/protoc-gen-validate/raw/master/validate/validate.proto",
		"third_party/github.com/envoyproxy/protoc-gen-validate/validate/validate.proto",
		0644,
	)
	download(
		"https://github.com/grpc-ecosystem/grpc-gateway/raw/master/protoc-gen-swagger/options/annotations.proto",
		"third_party/protoc-gen-swagger/options/annotations.proto",
		0644,
	)
	download(
		"https://github.com/grpc-ecosystem/grpc-gateway/raw/master/protoc-gen-swagger/options/openapiv2.proto",
		"third_party/protoc-gen-swagger/options/openapiv2.proto",
		0644,
	)

	args := []string{
		"-I.",
		"-I../third_party",
		fmt.Sprintf("--go_out=."),
		fmt.Sprintf("--go-grpc_out=."),
		fmt.Sprintf("--grpc-gateway_out=."),
		fmt.Sprintf("--swagger_out=."),
		fmt.Sprintf("--validate_out=lang=go:."),
		"--descriptor_set_out=echo.pb",
	}

	for _, version := range []string{"v1alpha1"} {
		path, err := filepath.Abs(version)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		goPackageRoot := filepath.Join("htdvisser.dev", "exp", "echo", "api")

		run(path, "bash", "-c", "rm -f ./*{.pb{,.gw,.validate}.go,.swagger.json}")

		run(path, "mkdir", "-p", goPackageRoot)
		run(path, "ln", "-sf", path, goPackageRoot)

		run(path, "bash", "-c", fmt.Sprintf("protoc %s *.proto", strings.Join(args, " ")))

		run(path, "rm", "-rf", "htdvisser.dev")

		run(path, "gofumpt", "-w", ".")
	}
}
