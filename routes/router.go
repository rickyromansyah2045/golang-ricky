package routes

import (
	"content_portal/controllers"
	"content_portal/database"
	"content_portal/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.2"})

	db := database.ConnectDB()
	userController := controllers.NewUserController(db)
	contentController := controllers.NewContentController(db)
	pointController := controllers.NewPointController(db)
	companyController := controllers.NewCompanyController(db)

	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", userController.Register)
		userGroup.POST("/login", userController.Login)
		userGroup.PUT("/", middlewares.Auth(), userController.Update)
		userGroup.DELETE("/", middlewares.Auth(), userController.Delete)
		userGroup.GET("/point", userController.GetPoint)
	}

	contentGroup := router.Group("/contents")
	{
		contentGroup.POST("/", middlewares.Auth(), contentController.Create)
		contentGroup.GET("/", middlewares.Auth(), contentController.Get)
		contentGroup.GET("/:contentId", middlewares.Auth(), contentController.GetById)
		contentGroup.PUT("/:contentId", middlewares.Auth(), contentController.Update)
		contentGroup.DELETE("/:contentId", middlewares.Auth(), contentController.Delete)
	}

	pointGroup := router.Group("/setpoint")
	{
		pointGroup.POST("/", middlewares.Auth(), pointController.Create)
		pointGroup.GET("/", middlewares.Auth(), pointController.Get)
		pointGroup.PUT("/:pointId", middlewares.Auth(), pointController.Update)
		pointGroup.DELETE("/:pointId", middlewares.Auth(), pointController.Delete)
	}

	companyGroup := router.Group("/company")
	{
		companyGroup.POST("/register", companyController.Register)
		companyGroup.POST("/login", companyController.Login)
		companyGroup.PUT("/", middlewares.Auth(), companyController.Update)
		companyGroup.DELETE("/", middlewares.Auth(), companyController.Delete)
		companyGroup.GET("/", middlewares.Auth(), companyController.GetUserRegister)

		// companyGroup.POST("/", middlewares.Auth(), companyController.Create)
		// companyGroup.PUT("/:socialMediaId", middlewares.Auth(), companyController.Update)
		// companyGroup.DELETE("/:socialMediaId", middlewares.Auth(), companyController.Delete)
	}

	return router

}
