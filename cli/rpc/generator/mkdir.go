//    Copyright 2022. Go-Ceres
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

package generator

import (
	"github.com/go-ceres/ceres/cli/rpc/config"
	"github.com/go-ceres/ceres/cli/rpc/parser/model"
	"github.com/go-ceres/ceres/ctx"
	"github.com/go-ceres/ceres/utils/pathx"
	"github.com/go-ceres/ceres/utils/stringx"
	"path/filepath"
	"strings"
)

const (
	wdKey             = "wd"
	configKey         = "config"
	bootKey           = "boot"
	manifestKey       = "manifest"
	globalKey         = "global"
	dtoKey            = "dto"
	serviceKey        = "service"
	packageKey        = "pkg"
	logicKey          = "logic"
	serverKey         = "server"
	applicationKey    = "application"
	infrastructureKey = "infrastructure"
	domainKey         = "domain"
	interfacesKey     = "interfaces"
	grpcKey           = "grpc"
	pbKey             = "pb"
	protoGoKey        = "proto-go"
)

type (
	// DirContext 文件夹上下文接口
	DirContext interface {
		GetMain() Dir
		GetConfig() Dir
		GetBoot() Dir
		GetGlobal() Dir
		GetApplication() Dir
		GetDomain() Dir
		GetInfrastructure() Dir
		GetInterfaces() Dir
		GetGrpc() Dir
		GetServer() Dir
		GetDto() Dir
		GetService() Dir
		GetLogic() Dir
		GetPb() Dir
		GetProtoGo() Dir
		GetServiceName() stringx.String
		SetPbDir(pbDir, grpcDir string)
	}
	// Dir 文件路径
	Dir struct {
		Package         string                                 // 文件完整包名
		Base            string                                 // 文件最后一级包名
		Filename        string                                 // 文件完整路径
		GetChildPackage func(childPath string) (string, error) // 文件夹完整路径
	}
	// defaultDirContext 文件夹管理上下文
	defaultDirContext struct {
		dirMap      map[string]Dir // 文件夹集合
		serviceName stringx.String // 服务名，该类型方便对字符串进行处理
		project     *ctx.Project   // 项目上下文
	}
)

func (d defaultDirContext) GetBoot() Dir {
	return d.dirMap[bootKey]
}

func (d defaultDirContext) GetGlobal() Dir {
	return d.dirMap[globalKey]
}

func (d defaultDirContext) GetDto() Dir {
	return d.dirMap[dtoKey]
}

func (d defaultDirContext) GetGrpc() Dir {
	return d.dirMap[grpcKey]
}

func (d defaultDirContext) GetMain() Dir {
	return d.dirMap[wdKey]
}

func (d defaultDirContext) GetConfig() Dir {
	return d.dirMap[configKey]
}

func (d defaultDirContext) GetApplication() Dir {
	return d.dirMap[applicationKey]
}

func (d defaultDirContext) GetDomain() Dir {
	return d.dirMap[domainKey]
}

func (d defaultDirContext) GetInfrastructure() Dir {
	return d.dirMap[infrastructureKey]
}

func (d defaultDirContext) GetInterfaces() Dir {
	return d.dirMap[interfacesKey]
}

func (d defaultDirContext) GetServer() Dir {
	return d.dirMap[serverKey]
}

func (d defaultDirContext) GetService() Dir {
	return d.dirMap[serviceKey]
}

func (d defaultDirContext) GetLogic() Dir {
	return d.dirMap[logicKey]
}

func (d defaultDirContext) GetPb() Dir {
	return d.dirMap[pbKey]
}

func (d defaultDirContext) GetProtoGo() Dir {
	return d.dirMap[protoGoKey]
}

func (d defaultDirContext) GetServiceName() stringx.String {
	return d.serviceName
}

func (d defaultDirContext) SetPbDir(pbDir, grpcDir string) {
	d.dirMap[pbKey] = Dir{
		Filename: pbDir,
		Package:  filepath.ToSlash(filepath.Join(d.project.Path, strings.TrimPrefix(pbDir, d.project.Dir))),
		Base:     filepath.Base(pbDir),
	}

	d.dirMap[protoGoKey] = Dir{
		Filename: grpcDir,
		Package: filepath.ToSlash(filepath.Join(d.project.Path,
			strings.TrimPrefix(grpcDir, d.project.Dir))),
		Base: filepath.Base(grpcDir),
	}
}

func (d *Dir) Valid() bool {
	return len(d.Filename) > 0 && len(d.Package) > 0
}

// mkdir 创建文件夹
func (g *Generator) mkdir(project *ctx.Project, proto model.Proto, conf *config.Config) (DirContext, error) {
	dirMap := make(map[string]Dir)
	bootDir := filepath.Join(project.WorkDir, "bootstrap")                // boot启动文件夹
	manifestDir := filepath.Join(project.WorkDir, "manifest")             // 配置部属相关
	configDir := filepath.Join(manifestDir, "config")                     // 配置文件文件夹
	globalDir := filepath.Join(project.WorkDir, "global")                 // 全局变量文件夹
	applicationDir := filepath.Join(project.WorkDir, "application")       // 应用层
	domainDir := filepath.Join(project.WorkDir, "domain")                 // 领域文件夹
	infrastructureDir := filepath.Join(project.WorkDir, "infrastructure") // 基础设施层
	interfacesDir := filepath.Join(project.WorkDir, "interfaces")         // 接口层
	grpcDir := filepath.Join(interfacesDir, "grpc")                       // grpc服务
	serverDir := filepath.Join(grpcDir, "server")                         // 服务层，与grpc之间的联系
	dtoDir := filepath.Join(applicationDir, "dto")                        // dto文件夹
	serviceDir := filepath.Join(applicationDir, "service")                // 服务层
	logicDir := filepath.Join(grpcDir, "logic")                           // 逻辑层处理
	pbDir := filepath.Join(project.WorkDir, proto.GoPackage)              // pb文件夹
	protoGoDir := pbDir                                                   // go文件
	if conf != nil {
		pbDir = conf.ProtoGenGrpcDir
		protoGoDir = conf.ProtoGenGoDir
	}
	getChildPackage := func(parent, childPath string) (string, error) {
		child := strings.TrimPrefix(childPath, parent)
		abs := filepath.Join(parent, strings.ToLower(child))
		if conf.Multiple {
			if err := pathx.MkdirIfNotExist(abs); err != nil {
				return "", err
			}
		}
		childPath = strings.TrimPrefix(abs, project.Dir)
		pkg := filepath.Join(project.Path, childPath)
		return filepath.ToSlash(pkg), nil
	}

	dirMap[wdKey] = Dir{
		Filename: project.WorkDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(project.WorkDir, project.Dir))),
		Base:     filepath.Base(project.WorkDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(project.WorkDir, childPath)
		},
	}
	dirMap[configKey] = Dir{
		Filename: configDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(configDir, project.Dir))),
		Base:     filepath.Base(configDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(configDir, childPath)
		},
	}
	dirMap[bootKey] = Dir{
		Filename: bootDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(bootDir, project.Dir))),
		Base:     filepath.Base(bootDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(bootDir, childPath)
		},
	}
	dirMap[globalKey] = Dir{
		Filename: globalDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(globalDir, project.Dir))),
		Base:     filepath.Base(globalDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(globalDir, childPath)
		},
	}
	dirMap[manifestKey] = Dir{
		Filename: manifestDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(manifestDir, project.Dir))),
		Base:     filepath.Base(manifestDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(manifestDir, childPath)
		},
	}
	dirMap[infrastructureKey] = Dir{
		Filename: infrastructureDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(infrastructureDir, project.Dir))),
		Base:     filepath.Base(infrastructureDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(infrastructureDir, childPath)
		},
	}
	dirMap[applicationKey] = Dir{
		Filename: applicationDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(applicationDir, project.Dir))),
		Base:     filepath.Base(applicationDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(applicationDir, childPath)
		},
	}
	dirMap[domainKey] = Dir{
		Filename: domainDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(domainDir, project.Dir))),
		Base:     filepath.Base(domainDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(domainDir, childPath)
		},
	}
	dirMap[interfacesKey] = Dir{
		Filename: interfacesDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(interfacesDir, project.Dir))),
		Base:     filepath.Base(interfacesDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(interfacesDir, childPath)
		},
	}
	dirMap[grpcKey] = Dir{
		Filename: grpcDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(grpcDir, project.Dir))),
		Base:     filepath.Base(grpcDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(grpcDir, childPath)
		},
	}
	dirMap[serverKey] = Dir{
		Filename: serverDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(serverDir, project.Dir))),
		Base:     filepath.Base(serverDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(serverDir, childPath)
		},
	}
	dirMap[dtoKey] = Dir{
		Filename: dtoDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(dtoDir, project.Dir))),
		Base:     filepath.Base(dtoDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(dtoDir, childPath)
		},
	}
	dirMap[serviceKey] = Dir{
		Filename: serviceDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(serviceDir, project.Dir))),
		Base:     filepath.Base(serviceDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(serviceDir, childPath)
		},
	}
	dirMap[logicKey] = Dir{
		Filename: logicDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(logicDir, project.Dir))),
		Base:     filepath.Base(logicDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(logicDir, childPath)
		},
	}
	dirMap[pbKey] = Dir{
		Filename: pbDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(pbDir, project.Dir))),
		Base:     filepath.Base(pbDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(pbDir, childPath)
		},
	}
	dirMap[protoGoKey] = Dir{
		Filename: protoGoDir,
		Package:  filepath.ToSlash(filepath.Join(project.Path, strings.TrimPrefix(protoGoDir, project.Dir))),
		Base:     filepath.Base(protoGoDir),
		GetChildPackage: func(childPath string) (string, error) {
			return getChildPackage(protoGoDir, childPath)
		},
	}
	for _, dir := range dirMap {
		err := pathx.MkdirIfNotExist(dir.Filename)
		if err != nil {
			return nil, err
		}
	}
	serviceName := strings.TrimSuffix(proto.Name, filepath.Ext(proto.Name))
	return &defaultDirContext{
		dirMap:      dirMap,
		project:     project,
		serviceName: stringx.NewString(strings.ReplaceAll(serviceName, "-", "")),
	}, nil
}
