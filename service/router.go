package main

import (
	"github.com/csby/database/memdb"
	"github.com/csby/wsf/opt"
	"github.com/csby/wsf/types"
	"net/http"
	"os"
)

func NewRouter(log types.Log, cfg *Config, ext types.HttpHandlerExtend, site types.Site) types.HttpHandler {
	instance := &Router{}
	instance.SetLog(log)
	instance.cfg = cfg
	instance.ext = ext
	instance.site = site

	instance.optToken = memdb.NewToken(cfg.Operation.Api.Token.Expiration, "opt")
	instance.optWSChannels = types.NewSocketChannelCollection()
	instance.opt = opt.NewHandler(log, &cfg.Configure, instance.optToken, instance.optWSChannels, site)

	return instance
}

type Router struct {
	types.Base

	cfg           *Config
	ext           types.HttpHandlerExtend
	opt           opt.Handler
	site          types.Site
	optToken      memdb.Token
	optWSChannels types.SocketChannelCollection
}

func (s *Router) Map(router types.Router) {
	// 后台服务管理
	if s.opt != nil {
		err := s.opt.Init(router, func(path types.Path, router types.Router, tokenChecker types.RouterPreHandle) error {
			return nil
		})
		if err != nil {
			s.LogError(err)
			os.Exit(8)
		}
	}
}

func (s *Router) PreRouting(w http.ResponseWriter, r *http.Request, a types.Assistant) bool {
	// enable across access
	if r.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type,token")
		return true
	}

	return false
}

func (s *Router) PostRouting(w http.ResponseWriter, r *http.Request, a types.Assistant) {

}

func (s *Router) NotFound() func(http.ResponseWriter, *http.Request, types.Assistant) {
	if s.opt == nil {
		return nil
	}

	return s.opt.NotFound
}

func (s *Router) Extend() types.HttpHandlerExtend {
	return s.ext
}
