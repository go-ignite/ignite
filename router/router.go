package router

import (
	_ "github.com/go-ignite/ignite/docs"
	"github.com/go-ignite/ignite/handler"
	"github.com/go-ignite/ignite/middleware"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Router struct {
	*gin.Engine
	*handler.UserHandler
	*handler.AdminHandler
}

func (r *Router) Init() {
	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
		}))
	}

	userRouter := r.Group("/api/user")
	{
		userRouter.POST("/login", r.LoginHandler)
		userRouter.POST("/signup", r.SignupHandler)

		authRouter := userRouter.Group("/auth")
		authRouter.Use(middleware.Auth(false))
		{
			authRouter.GET("/info", r.UserInfoHandler)
			authRouter.GET("/service/config", r.ServiceConfigHandler)
			authRouter.POST("/service/create", r.CreateServiceHandler)

			// nodes
			authRouter.GET("/nodes", r.UserHandler.ListNodes)
		}
	}

	adminRouter := r.Group("/api/admin")
	{
		adminRouter.POST("/login", r.PanelLoginHandler)
		authRouter := adminRouter.Group("/auth")
		authRouter.Use(middleware.Auth(true))
		{
			//user account related operations
			authRouter.GET("/status_list", r.PanelStatusListHandler)
			authRouter.PUT("/:id/reset", r.ResetAccountHandler)
			authRouter.PUT("/:id/destroy", r.DestroyAccountHandler)
			authRouter.PUT("/:id/stop", r.StopServiceHandler)
			authRouter.PUT("/:id/start", r.StartServiceHandler)
			authRouter.PUT("/:id/renew", r.RenewServiceHandler)

			//invite code related operations
			authRouter.GET("/code_list", r.InviteCodeListHandler)
			authRouter.PUT("/:id/remove", r.RemoveInviteCodeHandler)
			authRouter.POST("/code_generate", r.GenerateInviteCodeHandler)

			// nodes
			authRouter.GET("/nodes", r.AdminHandler.ListNodes)
			authRouter.POST("/nodes", r.AddNode)
			authRouter.PUT("/nodes/:id", r.UpdateNode)
			authRouter.DELETE("/nodes/:id", r.DeleteNode)
		}
	}
}
