package routes

import (
	controller "kanban-board/controllers"
	userRepo "kanban-board/repository/user"
	userUsecase "kanban-board/usecase/user"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRouter(e *echo.Echo, db *gorm.DB) {

	// Users
	userRepo := userRepo.NewUserRepository(db)
	userService := userUsecase.NewUserUseCase(userRepo)
	userController := controller.NewUserController(userService)

	userGroup := e.Group("/users")
	userGroup.GET("", userController.GetUsers)
	userGroup.GET("/:id", userController.GetUser)
	userGroup.POST("", userController.CreateUser)
	userGroup.PATCH("/:id", userController.UpdateUser)
	userGroup.DELETE("/:id", userController.DeleteUser)
}