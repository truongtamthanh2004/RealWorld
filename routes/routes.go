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
	favoriteRepo := repositories.NewFavoriteRepository()
	articleRepo := repositories.NewArticleRepository()
	articleService := services.NewArticleService(articleRepo, followRepo, favoriteRepo)
	articleHandler := handlers.NewArticleHandler(articleService)
	commentRepo := repositories.NewCommentRepository()
	commentService := services.NewCommentService(commentRepo, articleRepo, followRepo)
	commentHandler := handlers.NewCommentHandler(commentService)
	tagRepo := repositories.NewTagRepository()
	tagService := services.NewTagService(tagRepo)
	tagHandler := handlers.NewTagHandler(tagService)

	api := r.Group("/api")
	api.POST("/users", userHandler.Register)
	api.POST("/users/login", userHandler.Login)
	api.GET("/users", middlewares.AuthMiddleware(), userHandler.GetCurrentUser)
	api.PUT("/users", middlewares.AuthMiddleware(), userHandler.UpdateUser)

	profiles := api.Group("/profiles")
	profiles.GET("/:username", profileHandler.GetProfile)
	profiles.POST("/:username/follow", middlewares.AuthMiddleware(), profileHandler.FollowUser)
	profiles.DELETE("/:username/follow", middlewares.AuthMiddleware(), profileHandler.UnfollowUser)

	articles := api.Group("/articles")
	articles.GET("/", articleHandler.ListArticles)
	articles.GET("/feed", middlewares.AuthMiddleware(), articleHandler.FeedArticles)
	articles.GET("/:slug", articleHandler.GetArticle)
	articles.POST("/", middlewares.AuthMiddleware(), articleHandler.CreateArticle)
	articles.PUT("/:slug", middlewares.AuthMiddleware(), articleHandler.UpdateArticle)
	articles.DELETE("/:slug", middlewares.AuthMiddleware(), articleHandler.DeleteArticle)
	articles.POST("/:slug/comments", middlewares.AuthMiddleware(), commentHandler.AddComment)
	articles.GET("/:slug/comments", middlewares.AuthMiddleware(), commentHandler.GetComments)
	articles.DELETE("/:slug/comments/:id", middlewares.AuthMiddleware(), commentHandler.DeleteComment)
	articles.POST("/:slug/favorite", middlewares.AuthMiddleware(), articleHandler.FavoriteArticle)
	articles.DELETE("/:slug/favorite", middlewares.AuthMiddleware(), articleHandler.UnfavoriteArticle)

	api.GET("/tags", tagHandler.GetTags)

	return r
}
