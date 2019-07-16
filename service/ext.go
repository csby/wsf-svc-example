package main

import "github.com/csby/wsf/types"

type Ext struct {
	restart func() error
}

func (s *Ext) Restart() func() error {
	return s.restart
}

func (s *Ext) RedirectToHttps() bool {
	return cfg.Http.RedirectToHttps
}

func (s *Ext) DocumentEnabled() bool {
	return cfg.Document.Enabled
}

func (s *Ext) DocumentRoot() string {
	return cfg.Document.Root
}

func (s *Ext) ServerInfo() *types.ServerInformation {
	return &types.ServerInformation{
		Name:    moduleRemark,
		Version: moduleVersion,
	}
}
