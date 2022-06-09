//    Copyright 2021. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package module

import (
	"encoding/json"
	"errors"
	"github.com/go-ceres/ceres/new/rpc/config"
	"github.com/go-ceres/ceres/utils/exec"
	"os"
	"path/filepath"
)

// Module 模块的配置信息
type Module struct {
	Path      string
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}

// LoadMod 获取当前项目的go mod 信息
func LoadMod(workDir string) (*Module, error) {
	if len(workDir) == 0 {
		return nil, errors.New("there is no working path set")
	}
	// 工作路径状态
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}
	data, err := exec.Command("go list -json -m", workDir)
	if err != nil {
		return nil, err
	}
	var m Module
	if err = json.Unmarshal([]byte(data), &m); err != nil {
		return nil, err
	}
	if len(m.GoMod) > 0 {
		return &m, nil
	}
	return nil, nil
}

// InitMod 初始化模块
func InitMod(c *config.Config) (*config.Project, error) {
	// 1.读取go mod信息
	mod, err := LoadMod(c.WorkDir)
	if err != nil {
		return nil, err
	}
	var res config.Project
	// 原来就有go mod
	if mod != nil {
		res.WorkDir = c.WorkDir
		res.Name = filepath.Base(mod.Dir)
		res.Dir = mod.Dir
		res.Path = mod.Path
	} else { // 没有使用go mod
		// 1.写出go.mod文件
		name := filepath.Join(c.Namespace, c.Name)
		_, err := exec.Command("go mod init "+name, c.WorkDir)
		if err != nil {
			return nil, err
		}
		// 2. 再次验证
		return InitMod(c)
	}
	return &res, nil
}
