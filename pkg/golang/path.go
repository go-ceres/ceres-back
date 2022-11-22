package golang

import (
	"github.com/go-ceres/ceres/ctx"
	"github.com/go-ceres/ceres/utils/pathc"
	"path/filepath"
	"strings"
)

func GetParentPackage(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	projectConf, err := ctx.PrepareProject(abs)
	if err != nil {
		return "", err
	}

	// fix https://github.com/zeromicro/go-zero/issues/1058
	wd := projectConf.WorkDir
	d := projectConf.Dir
	same, err := pathc.SameFile(wd, d)
	if err != nil {
		return "", err
	}

	trim := strings.TrimPrefix(projectConf.WorkDir, projectConf.Dir)
	if same {
		trim = strings.TrimPrefix(strings.ToLower(projectConf.WorkDir), strings.ToLower(projectConf.Dir))
	}

	return filepath.ToSlash(filepath.Join(projectConf.Path, trim)), nil
}
