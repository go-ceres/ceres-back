package protoc

import (
	"archive/zip"
	"fmt"
	"github.com/go-ceres/ceres/utils/env"
	"github.com/go-ceres/ceres/utils/execc"
	"github.com/go-ceres/ceres/utils/installer"
	"github.com/go-ceres/ceres/utils/vars"
	"github.com/go-ceres/ceres/utils/zipc"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var url = map[string]string{
	"linux_32":   "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-linux-x86_32.zip",
	"linux_64":   "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-linux-x86_64.zip",
	"darwin":     "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-osx-x86_64.zip",
	"windows_32": "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-win32.zip",
	"windows_64": "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-win64.zip",
}

const (
	Name        = "protoc"
	ZipFileName = Name + ".zip"
)

func Install(cacheDir string) (string, error) {
	return installer.Install(cacheDir, Name, func(dest string) (string, error) {
		goos := runtime.GOOS
		tempFile := filepath.Join(os.TempDir(), ZipFileName)
		bit := 32 << (^uint(0) >> 63)
		var downloadUrl string
		switch goos {
		case vars.OsMac:
			downloadUrl = url[vars.OsMac]
		case vars.OsWindows:
			downloadUrl = url[fmt.Sprintf("%s_%d", vars.OsWindows, bit)]
		case vars.OsLinux:
			downloadUrl = url[fmt.Sprintf("%s_%d", vars.OsLinux, bit)]
		default:
			return "", fmt.Errorf("unsupport OS: %q", goos)
		}

		err := installer.Download(downloadUrl, tempFile)
		if err != nil {
			return "", err
		}

		return dest, zipc.Unpacking(tempFile, filepath.Dir(dest), func(f *zip.File) bool {
			return filepath.Base(f.Name) == filepath.Base(dest)
		})
	})
}

func Exists() bool {
	_, err := env.LookUpProtoc()
	return err == nil
}

func Version() (string, error) {
	path, err := env.LookUpProtoc()
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
