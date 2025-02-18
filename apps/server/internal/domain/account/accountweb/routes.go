package accountweb

import (
	"github.com/antoniosarro/gosvelte/apps/server/internal/web"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/logger"
)

type Options struct {
	Log            *logger.Log
	AccountUseCase iUsecase
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.AccountUseCase, opts.Log)

	g := web.Echo.Group("/api/v1/account")
	g.POST("/register", con.register)
	g.GET("/me", con.me, web.Mid.Authenticated)
}
