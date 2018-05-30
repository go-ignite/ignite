package router

import (
	"github.com/go-ignite/ignite/config"
	_ "github.com/go-ignite/ignite/docs"
	"github.com/go-ignite/ignite/handler"
	"github.com/go-ignite/ignite/middleware"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Router struct {
	*gin.Engine
}

func New(engine *gin.Engine) *Router {
	return &Router{
		Engine: engine,
	}
}

func (r *Router) InitGeneral() {
	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
		}))
	}
}

func (r *Router) InitUser(l *logrus.Logger) {
	userHandler := handler.NewUserHandler(l)
	userRouter := r.Group("/api/user")
	{
		userRouter.POST("/login", userHandler.LoginHandler)
		userRouter.POST("/signup", userHandler.SignupHandler)

		authRouter := userRouter.Group("/auth")
		authRouter.Use(middleware.Auth(config.C.Secret.User))
		{
			authRouter.GET("/info", userHandler.UserInfoHandler)
			authRouter.GET("/service/config", userHandler.ServiceConfigHandler)
			authRouter.POST("/service/create", userHandler.CreateServiceHandler)
		}
	}
}

func (r *Router) InitAdmin(l *logrus.Logger) {
	adminHandler := handler.NewAdminHandler(l)
	adminRouter := r.Group("/api/admin")
	{
		adminRouter.POST("/login", adminHandler.PanelLoginHandler)
		authRouter := adminRouter.Group("/auth")
		authRouter.Use(middleware.Auth(config.C.Secret.Admin))
		{
			//user account related operations
			authRouter.GET("/status_list", adminHandler.PanelStatusListHandler)
			authRouter.PUT("/:id/reset", adminHandler.ResetAccountHandler)
			authRouter.PUT("/:id/destroy", adminHandler.DestroyAccountHandler)
			authRouter.PUT("/:id/stop", adminHandler.StopServiceHandler)
			authRouter.PUT("/:id/start", adminHandler.StartServiceHandler)
			authRouter.PUT("/:id/renew", adminHandler.RenewServiceHandler)

			//invite code related operations
			authRouter.GET("/code_list", adminHandler.InviteCodeListHandler)
			authRouter.PUT("/:id/remove", adminHandler.RemoveInviteCodeHandler)
			authRouter.POST("/code_generate", adminHandler.GenerateInviteCodeHandler)
		}
	}
}
