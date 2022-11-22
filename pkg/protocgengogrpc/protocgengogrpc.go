package protocgengogrpc

import (
	"github.com/go-ceres/ceres/pkg/golang"
	"github.com/go-ceres/ceres/utils/env"
	"github.com/go-ceres/ceres/utils/execc"
	"github.com/go-ceres/ceres/utils/installer"
	"strings"
)

const (
	Name = "protoc-gen-go-grpc"
	url  = "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
)

func Install(cacheDir string) (string, error) {
	return installer.Install(cacheDir, Name, func(dest string) (string, error) {
		err := golang.Install(url)
		return dest, err
	})
}

func Exists() bool {
	_, err := env.LookUpProtocGenGoGrpc()
	return err == nil
}

// Version is used to get the version of the protoc-gen-go-grpc plugin.
func Version() (string, error) {
	path, err := env.LookUpProtocGenGoGrpc()
	if err != nil {
		return "", err
	}
	version, err := execc.Command(path+" --version", "")
	if err != nil {
		return "", err
	}
	fields := strings.Fields(version)
	if len(fields) > 1 {
		return fields[1], nil
	}
	return "", nil
}
