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

package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	ceresDir = ".ceres"
)

// MatchFiles 搜索匹配文件
func MatchFiles(src string) ([]string, error) {
	dir, pattern := filepath.Split(src)
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(abs)
	if err != nil {
		return nil, err
	}
	var res []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		match, err := filepath.Match(pattern, name)
		if err != nil {
			return nil, err
		}

		if !match {
			continue
		}

		res = append(res, filepath.Join(abs, name))
	}
	return res, nil
}

// MkdirIfNotFound 如果文件夹不存在则创建文件夹
func MkdirIfNotFound(dir string) error {
	if len(dir) == 0 {
		return nil
	}
	// 创建文件夹
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// FileExists 判断文件是否存在
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// LoadTpl 加载模板
func LoadTpl(category, filename, def string) (string, error) {
	dir, err := GetTplDir(category)
	if err != nil {
		return "", err
	}

	filename = filepath.Join(dir, filename)
	if !FileExists(filename) {
		return def, nil
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// GetTplDir 获取指定模板路径
func GetTplDir(category string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, category, ceresDir), nil
}
