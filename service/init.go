package main

import (
	"fmt"
	"github.com/csby/wsf/logger"
	"github.com/csby/wsf/server/host"
	"github.com/csby/wsf/types"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	moduleType    = "server"
	moduleName    = "wsf-svc-example"
	moduleRemark  = "WEB服务示例"
	moduleVersion = "1.0.1.0"
)

var (
	cfg              = NewConfig()
	log              = &logger.Writer{Level: logger.LevelAll}
	svr types.Server = nil
)

func init() {
	moduleArgs := &types.Args{}
	serverArgs := &types.SvcArgs{}
	moduleArgs.Parse(os.Args, moduleType, moduleName, moduleVersion, moduleRemark, serverArgs)
	now := time.Now()
	cfg.Module.Type = moduleType
	cfg.Module.Name = moduleName
	cfg.Module.Version = moduleVersion
	cfg.Module.Remark = moduleRemark
	cfg.Module.Path = moduleArgs.ModulePath()
	cfg.Service.BootTime = now

	rootFolder := filepath.Dir(moduleArgs.ModuleFolder())
	cfgFolder := filepath.Join(rootFolder, "cfg")
	cfgName := fmt.Sprintf("%s.json", moduleName)
	if serverArgs.Help {
		serverArgs.ShowHelp(cfgFolder, cfgName)
		os.Exit(11)
	}

	// init config
	svcArgument := ""
	cfgPath := serverArgs.Cfg
	if cfgPath != "" {
		svcArgument = fmt.Sprintf("-cfg=%s", cfgPath)
	} else {
		cfgPath = filepath.Join(cfgFolder, cfgName)
	}
	_, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		err = cfg.SaveToFile(cfgPath)
		if err != nil {
			fmt.Println("generate configure file fail: ", err)
		}
	} else {
		err = cfg.LoadFromFile(cfgPath)
		if err != nil {
			fmt.Println("load configure file fail: ", err)
		}
	}

	// init certificate
	if cfg.Https.Enabled {
		certFilePath := cfg.Https.Cert.Server.File
		if certFilePath == "" {
			certFilePath = filepath.Join(rootFolder, "crt", "server.pfx")
			cfg.Https.Cert.Server.File = certFilePath
		}
	}

	// init path of site
	if cfg.Root == "" {
		cfg.Root = filepath.Join(rootFolder, "site", "root")
	}
	if cfg.Document.Root == "" {
		cfg.Document.Root = filepath.Join(rootFolder, "site", "doc")
	}
	if cfg.Operation.Root == "" {
		cfg.Operation.Root = filepath.Join(rootFolder, "site", "opt")
	}
	if cfg.Webapp.Root == "" {
		cfg.Webapp.Root = filepath.Join(rootFolder, "site", "webapp")
	}

	// init service
	if strings.TrimSpace(cfg.Service.Name) == "" {
		cfg.Service.Name = moduleName
	}
	serviceName := cfg.Service.Name
	log.Init(cfg.Log.Level, serviceName, cfg.Log.Folder)
	ext := &Ext{}
	svrRouter := NewRouter(log, cfg, ext)
	svrHost := host.NewHost(log, &cfg.Configure, svrRouter, nil)
	svr, err = host.NewServer(log, svrHost, serviceName, svcArgument)
	if err != nil {
		fmt.Println("init service fail: ", err)
		os.Exit(12)
	}
	if !svr.Interactive() {
		ext.restart = svr.Restart
	}
	serverArgs.Execute(svr)

	// information
	log.Std = true
	zoneName, zoneOffset := now.Zone()
	LogInfo("start at: ", moduleArgs.ModulePath())
	LogInfo("version: ", moduleVersion)
	LogInfo("zone: ", zoneName, "-", zoneOffset/int(time.Hour.Seconds()))
	LogInfo("log path: ", cfg.Log.Folder)
	LogInfo("configure path: ", cfgPath)
	LogInfo("configure info: ", cfg)
}
