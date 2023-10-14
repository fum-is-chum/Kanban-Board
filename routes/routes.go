package routes

import (
	controller "kanban-board/controllers"
	m "kanban-board/middlewares"
	boardRepo "kanban-board/repository/board"
	boardMemberRepo "kanban-board/repository/board_member"
	userRepo "kanban-board/repository/user"
	authUsecase "kanban-board/usecase/auth"
	boardUsecase "kanban-board/usecase/board"
	boardMemberUsecase "kanban-board/usecase/board_member"
	userUsecase "kanban-board/usecase/user"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRouter(e *echo.Echo, db *gorm.DB) {

	// logger middleware
	m.LoggerMiddleware(e)

	// ------------------------------------------------------------------------------
	// Repos
	userRepo := userRepo.NewUserRepository(db)
	boardRepo := boardRepo.NewBoardRepository(db)
	boardMemberRepo := boardMemberRepo.NewBoardMemberRepostory(db)

	// Services
	authService := authUsecase.NewAuthUseCase(userRepo)
	userService := userUsecase.NewUserUseCase(userRepo)
	boardService := boardUsecase.NewBoardUseCase(userRepo, boardRepo)
	boardMemberService := boardMemberUsecase.NewBoardMemberUseCase(&boardMemberUsecase.BoardMemberMultiRepos{
		UserRepo: userRepo,
		BoardRepo: boardRepo,
		MemberRepo: boardMemberRepo,
	})

	// Controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	boardController := controller.NewBoardController(boardService)
	boardMemberController := controller.NewBoardMemberController(boardMemberService)
	// -------------------------------------------------------------------------------

	// Login
	e.POST("/login", authController.Login)

	// Users
	userGroup := e.Group("/users")
	userGroup.GET("", userController.GetUsers, m.JWTMiddleware())
	userGroup.GET("/:id", userController.GetUser, m.JWTMiddleware())
	userGroup.POST("", userController.CreateUser)
	userGroup.PATCH("/:id", userController.UpdateUser, m.JWTMiddleware())
	// userGroup.DELETE("/:id", userController.DeleteUser, m.JWTMiddleware())

	// Boards
	boardGroup := e.Group("/boards")
	boardGroup.GET("", boardController.GetBoards, m.JWTMiddleware())
	boardGroup.GET("/:id", boardController.GetBoardById, m.JWTMiddleware())
	boardGroup.POST("", boardController.CreateBoard, m.JWTMiddleware())
	boardGroup.PATCH("/:id", boardController.UpdateBoard, m.JWTMiddleware())
	boardGroup.DELETE("/:id", boardController.DeleteBoard, m.JWTMiddleware())

	// Board Member
	boardMemberGroup := e.Group("/board-members")
	boardMemberGroup.POST("/add", boardMemberController.AddNewMember, m.JWTMiddleware())
	boardMemberGroup.POST("/remove", boardMemberController.DeleteMember, m.JWTMiddleware())
}
