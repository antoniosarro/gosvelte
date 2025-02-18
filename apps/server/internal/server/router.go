package server

import (
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/account/accountrepo"
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/account/accountuc"
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/account/accountweb"
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/auth/authrepo"
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/auth/authuc"
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/auth/authweb"
	"github.com/antoniosarro/gosvelte/apps/server/internal/web"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func router(w *web.Web, conf *Options) {
	w.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	accountDBRepository := accountrepo.NewDB(conf.DB)
	accountCacheRepository := accountrepo.NewCache(conf.Cache)
	accountUseCase := accountuc.New(conf.ServConf, conf.Log, accountDBRepository, accountCacheRepository)
	accountweb.Route(w, &accountweb.Options{
		Log:            conf.Log,
		AccountUseCase: accountUseCase,
	})

	authCacheRepository := authrepo.NewCache(conf.Cache)
	authUC := authuc.New(conf.ServConf, conf.Log, authCacheRepository, accountDBRepository)
	authweb.Route(w, &authweb.Options{
		Log:    conf.Log,
		AuthUC: authUC,
	})
}
