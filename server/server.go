package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/go-ignite/ignite/config"
	"github.com/go-ignite/ignite/service"
)

var Set = wire.NewSet(wire.Struct(new(Options), "*"), New)

type Server struct {
	engine *gin.Engine
	opts   *Options
}

type Options struct {
	Config  config.Server
	Service *service.Service
}

func New(opts *Options) *Server {
	s := &Server{
		engine: gin.Default(),
		opts:   opts,
	}

	userRouter := s.engine.Group("/api/user")
	{
		userRouter.POST("/login", s.opts.Service.UserLogin)
		userRouter.POST("/register", s.opts.Service.UserRegister)

		userRouter.Use(s.opts.Service.Auth(false))
		{
			userRouter.POST("/sync", s.opts.Service.Sync)
			userRouter.POST("/services", s.opts.Service.CreateService)
			//authRouter.GET("/info", userHandler.UserInfoHandler)

			// nodes
			//authRouter.GET("/nodes", userHandler.ListNodes)

			// services
			//authRouter.GET("/services/config", userHandler.GetServiceConfig)

			//authRouter.GET("/services", userHandler.ListServices)
			//authRouter.DELETE("/services/:id", userHandler.RemoveService)
			//authRouter.POST("/nodes/:nodeID/services", userHandler.CreateService)
		}
	}

	adminRouter := s.engine.Group("/api/admin")
	{
		adminRouter.POST("/login", s.opts.Service.AdminLogin)
		adminRouter.Use(s.opts.Service.Auth(true))
		{
			//user account related operations
			adminRouter.GET("/accounts", s.opts.Service.GetAccountList)
			//adminRouter.PUT("/accounts/:id/reset", s.opts.Service.ResetAccountHand)
			adminRouter.DELETE("/accounts/:id", s.opts.Service.DestroyAccount)
			// authRouter.PUT("/:id/stop", r.StopServiceHandler)
			// authRouter.PUT("/:id/start", r.StartServiceHandler)
			// authRouter.PUT("/:id/renew", r.RenewServiceHandler)

			// codes
			adminRouter.GET("/codes", s.opts.Service.GetInviteCodeList)
			adminRouter.DELETE("/codes/:id", s.opts.Service.RemoveInviteCode)
			adminRouter.POST("/codes_batch", s.opts.Service.GenerateInviteCodes)

			// nodes
			adminRouter.GET("/nodes", s.opts.Service.GetAllNodes)
			adminRouter.POST("/nodes", s.opts.Service.AddNode)
			adminRouter.PUT("/nodes/:id", s.opts.Service.UpdateNode)
			adminRouter.DELETE("/nodes/:id", s.opts.Service.DeleteNode)

			// services
			//authRouter.GET("/services", adminHandler.ListServices)
			//authRouter.DELETE("/services/:id", adminHandler.RemoveService)
		}
	}

	return s
}

func (s *Server) Start() error {
	return s.engine.Run(s.opts.Config.Address)
}
