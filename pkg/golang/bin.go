package golang

import (
	"github.com/go-ceres/ceres/utils/pathc"
	"go/build"
	"os"
	"path/filepath"
)

// GoBin returns a path of GOBIN.
func GoBin() string {
	def := build.Default
	goroot := os.Getenv("GOPATH")
	bin := filepath.Join(goroot, "bin")
	if !pathc.FileExists(bin) {
		gopath := os.Getenv("GOROOT")
		bin = filepath.Join(gopath, "bin")
	}
	if !pathc.FileExists(bin) {
		bin = os.Getenv("GOBIN")
	}
	if !pathc.FileExists(bin) {
		bin = filepath.Join(def.GOPATH, "bin")
	}
	return bin
}
