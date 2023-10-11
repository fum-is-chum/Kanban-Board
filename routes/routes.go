package routes

import (
	controller "kanban-board/controllers"
	m "kanban-board/middlewares"
	userRepo "kanban-board/repository/user"
	authUsecase "kanban-board/usecase/auth"
	userUsecase "kanban-board/usecase/user"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRouter(e *echo.Echo, db *gorm.DB) {

	// logger middleware
	m.LoggerMiddleware(e)

	// Login
	userRepo := userRepo.NewUserRepository(db)
	authService := authUsecase.NewAuthUseCase(userRepo)
	authController := controller.NewAuthController(authService)

	e.POST("/login", authController.Login)

	// Users
	userService := userUsecase.NewUserUseCase(userRepo)
	userController := controller.NewUserController(userService)

	userGroup := e.Group("/users")
	userGroup.GET("", userController.GetUsers, m.JWTMiddleware())
	userGroup.GET("/:id", userController.GetUser, m.JWTMiddleware())
	userGroup.POST("", userController.CreateUser)
	userGroup.PATCH("/:id", userController.UpdateUser, m.JWTMiddleware())
	// userGroup.DELETE("/:id", userController.DeleteUser, m.JWTMiddleware())
}
