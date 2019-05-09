package service

import (
	"encoding/json"
	"time"

	"github.com/go-ignite/ignite/config"
	_ "github.com/go-ignite/ignite/docs"
	"github.com/go-ignite/ignite/handler"
	"github.com/go-ignite/ignite/middleware"
	"github.com/go-ignite/ignite/state"
	"github.com/go-ignite/ignite/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type Service struct {
	*gin.Engine
	userHandler  *handler.UserHandler
	adminHandler *handler.AdminHandler
}

func New(userHandler *handler.UserHandler, adminHandler *handler.AdminHandler) *Service {
	return &Service{
		Engine:       gin.New(),
		userHandler:  userHandler,
		adminHandler: adminHandler,
	}
}

func (s *Service) Run() {
	s.Engine.Run(config.C.App.Address)
}

func (s *Service) Init() *Service {
	//if gin.Mode() == gin.DebugMode {
	//	s.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//
	//	s.Use(cors.New(cors.Config{
	//		AllowAllOrigins: true,
	//	}))
	//}

	m := melody.New()
	s.GET("/api/ws/nodes", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		if !utils.VerifyToken(string(msg), nil) {
			return
		}
		for {
			nam := state.GetLoader().NodeAvailableMap()
			msg, _ := json.Marshal(nam)
			if err := s.Write(msg); err != nil {
				break
			}
			// TODO should be configurable
			time.Sleep(3 * time.Second)
		}
	})

	userRouter := s.Group("/api/user")
	{
		userHandler := s.userHandler
		userRouter.POST("/login", userHandler.Login)
		userRouter.POST("/register", userHandler.Register)

		authRouter := userRouter.Group("/auth")
		authRouter.Use(middleware.NewUserAuthHandler().Auth())
		{
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

	adminRouter := s.Group("/admin/api")
	{
		adminHandler := s.adminHandler
		adminRouter.POST("/login", adminHandler.Login)
		adminRouter.Use(middleware.NewAdminAuthHandler().Auth())
		{
			//user account related operations
			//authRouter.GET("/status_list", adminHandler.PanelStatusListHandler)
			// authRouter.PUT("/:id/reset", r.ResetAccountHandler)
			// authRouter.PUT("/:id/destroy", r.DestroyAccountHandler)
			// authRouter.PUT("/:id/stop", r.StopServiceHandler)
			// authRouter.PUT("/:id/start", r.StartServiceHandler)
			// authRouter.PUT("/:id/renew", r.RenewServiceHandler)

			// codes
			adminRouter.GET("/codes", adminHandler.GetInviteCodeList)
			adminRouter.DELETE("/codes/:id", adminHandler.RemoveInviteCode)
			adminRouter.POST("/codes_batch", adminHandler.GenerateInviteCodes)

			// nodes
			adminRouter.GET("/nodes", adminHandler.GetAllNodes)
			adminRouter.POST("/nodes", adminHandler.AddNode)
			adminRouter.PUT("/nodes/:id", adminHandler.UpdateNode)
			adminRouter.DELETE("/nodes/:id", adminHandler.DeleteNode)

			// services
			//authRouter.GET("/services", adminHandler.ListServices)
			//authRouter.DELETE("/services/:id", adminHandler.RemoveService)
		}
	}
	return s
}
