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

package gen

import (
	"errors"
	"github.com/go-ceres/ceres/new/model/gorm/parse"
	"github.com/go-ceres/ceres/utils"
	"path/filepath"
	"strings"
)

// Config 配置信息
type Config struct {
	Src         string // sql输入文件
	Dist        string // model文件输出路径
	Pkg         string // 包名
	Prefix      string // 表前缀
	Cache       bool   // 是否使用缓存组件
	AutoPrimary bool   // 当配置为true时，如果创建的sql没有主键，则会主动加上id这个主键
}

// DefaultConfig 默认配置信息
func DefaultConfig() *Config {
	return &Config{
		Src:    "./*.sql",
		Prefix: "",
		Cache:  false,
	}
}

// BuildGenerator 根据配置文件生成代码生成器
func (c *Config) BuildGenerator() (*Generator, error) {
	src := strings.TrimSpace(c.Src)
	if len(src) == 0 {
		return nil, errors.New("expected path or path globbing patterns, but nothing found")
	}
	if c.Dist == "" {
		c.Dist = "."
	}
	distAbs, err := filepath.Abs(c.Dist)
	if err != nil {
		return nil, err
	}
	pkg := filepath.Base(distAbs)
	err = utils.MkdirIfNotFound(distAbs)
	if err != nil {
		return nil, err
	}
	c.Pkg = pkg

	files, err := utils.MatchFiles(src)
	if err != nil {

		return nil, err
	}
	statement, err := parse.Parse(files)
	if err != nil {
		return nil, err
	}

	return &Generator{
		Config:    c,
		Statement: statement,
	}, nil
}
