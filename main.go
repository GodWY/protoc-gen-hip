package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/GodWY/protoc-gen-hip/internal_genhi"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/compiler/protogen"
)

const Version = "v1.2.0"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "%v %v\n", filepath.Base(os.Args[0]), Version)
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "--help" {
		os.Exit(0)
	}

	var (
		flags   flag.FlagSet
		plugins = flags.String("plugins", "", "deprecated option")
	)

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		if *plugins != "" {
			return errors.New("protoc-gen-go: plugins are not supported; use 'protoc --go-grpc_out=...' to generate gRPC\n\n")
		}
		for _, f := range gen.Files {
			if f.Generate {
				internal_genhi.GenerateFile(gen, f)
			}
		}
		gen.SupportedFeatures = internal_genhi.SupportedFeatures
		return nil
	})
}
