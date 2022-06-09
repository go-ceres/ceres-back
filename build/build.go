package build

import (
	"fmt"
	"github.com/go-ceres/cli/v2"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "name",
		Usage:   "build application name",
		Aliases: []string{"n"},
		EnvVars: []string{"BUILD_NAME"},
	},
	&cli.StringFlag{
		Name:    "version",
		Usage:   "build application version",
		Aliases: []string{"v"},
		EnvVars: []string{"BUILD_VERSION"},
	},
	&cli.StringFlag{
		Name:    "path",
		Aliases: []string{"p"},
		Usage:   "Application output path",
		EnvVars: []string{"BUILD_PATH"},
	},
	&cli.StringFlag{
		Name:    "status",
		Aliases: []string{"s"},
		Usage:   "Application output path",
		EnvVars: []string{"BUILD_STATUS"},
	},
}

// Run 当前命令行执行工具
func Run(ctx *cli.Context) error {
	// 应用名称
	appName := ctx.String("name")
	if len(appName) == 0 {
		fmt.Println("please ")
		return nil
	}
	path := ctx.String("path")
	if len(path) == 0 {
		path = "./bin"
	}
	// 构建时间
	buildTime := time.Now().Format("2006-01-02--15:04:05")
	// 构建版本
	version := ctx.String("version")
	if len(version) == 0 {
		versionCmd := exec.Command("git", "describe", "--tags")
		versionOut, err := versionCmd.CombinedOutput()
		if err != nil {
			version = "latest"
		} else {
			version = string(versionOut)
		}
	}
	// 构建时的host
	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	// 构建用户
	userInfo, err := user.Current()
	buildUser := ""
	if err == nil {
		buildUser = userInfo.Username
	}
	// 构建状态
	buildStatus := ctx.String("status")
	buildStatusScript := `
#!/bin/bash
if git diff-index --quiet HEAD --; then
  tree_status="Clean"
else
  tree_status="Modified"
fi
echo $tree_status
`
	if len(buildStatus) == 0 {
		statusCmd := exec.Command("bash", "-c", buildStatusScript)
		statusOut, outErr := statusCmd.CombinedOutput()
		if outErr == nil {
			buildStatus = string(statusOut)
		}
	}
	// git分支信息
	branch := ""
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, branchErr := branchCmd.CombinedOutput()
	if branchErr == nil {
		branch = string(branchOut)
	}

	// git提交信息
	commit := ""
	commitCmd := exec.Command("git", "rev-parse", "HEAD")
	commitOut, commitErr := commitCmd.CombinedOutput()
	if commitErr == nil {
		commit = string(commitOut)
	}

	var ags []string
	var params []string
	ags = append(ags, "go")
	ags = append(ags, "build")
	ags = append(ags, "-x")
	ags = append(ags, "-trimpath")
	params = append(params, "-X github.com/go-ceres/go-ceres/cmd.appName="+appName)
	params = append(params, "-X github.com/go-ceres/go-ceres/cmd.version="+version)
	params = append(params, "-X github.com/go-ceres/go-ceres/cmd.branch="+branch)
	params = append(params, "-X github.com/go-ceres/go-ceres/cmd.commit="+commit)
	params = append(params, "-X github.com/go-ceres/go-ceres/cmd.buildHost="+name)
	params = append(params, "-X github.com/go-ceres/go-ceres/cmd.buildTime="+buildTime)
	params = append(params, "-X github.com/go-ceres/go-ceres/cmd.buildUser="+buildUser)
	params = append(params, "-X github.com/go-ceres/go-ceres/cmd.buildStatus="+buildStatus)
	ags = append(ags, "-ldflags")
	ags = append(ags, strings.Join(params, " "))
	ags = append(ags, "-v")
	ags = append(ags, "-o")
	ags = append(ags, filepath.Join(path, appName))
	// 如果有文件
	file := ctx.Args().First()
	if len(file) != 0 {
		ags = append(ags, file)
	}
	// 组装
	cmd := exec.Command("time", ags...)
	envs := os.Environ()
	fmt.Println(envs)
	// 设置环境变量
	cmd.Env = envs
	// 执行命令
	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return nil
}
