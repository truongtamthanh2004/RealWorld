package routes

import (
	"Netlfy/handlers"
	"Netlfy/middlewares"
	"Netlfy/repositories"
	"Netlfy/services"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	followRepo := repositories.NewFollowRepository()
	profileService := services.NewProfileService(userRepo, followRepo)
	profileHandler := handlers.NewProfileHandler(profileService)
	articleRepo := repositories.NewArticleRepository()
	articleService := services.NewArticleService(articleRepo, followRepo)
	articleHandler := handlers.NewArticleHandler(articleService)

	api := r.Group("/api")
	api.POST("/users", userHandler.Register)
	api.POST("/users/login", userHandler.Login)
	api.GET("/users", middlewares.AuthMiddleware(), userHandler.GetCurrentUser)
	api.PUT("/users", middlewares.AuthMiddleware(), userHandler.UpdateUser)

	profiles := api.Group("/profiles")
	profiles.GET("/:username", profileHandler.GetProfile)
	profiles.POST("/:username/follow", middlewares.AuthMiddleware(), profileHandler.FollowUser)
	profiles.DELETE("/:username/follow", middlewares.AuthMiddleware(), profileHandler.UnfollowUser)

	api.GET("/articles", articleHandler.ListArticles)

	return r
}
