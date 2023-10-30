package routes

import (
	controller "kanban-board/controllers"
	m "kanban-board/middlewares"
	boardRepo "kanban-board/repository/board"
	boardColumnRepo "kanban-board/repository/board_column"
	boardMemberRepo "kanban-board/repository/board_member"
	taskRepo "kanban-board/repository/task"
	assigneeRepo "kanban-board/repository/task_assignee"
	userRepo "kanban-board/repository/user"
	authUsecase "kanban-board/usecase/auth"
	boardUsecase "kanban-board/usecase/board"
	boardColumnUsecase "kanban-board/usecase/board_column"
	boardMemberUsecase "kanban-board/usecase/board_member"
	taskUsecase "kanban-board/usecase/task"
	assigneeUsecase "kanban-board/usecase/task_assignee"
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
	boardMemberRepo := boardMemberRepo.NewBoardMemberRepository(db)
	boardColumnRepo := boardColumnRepo.NewBoardColumnRepository(db)
	taskRepo := taskRepo.NewTaskRepository(db)
	assigneeRepo := assigneeRepo.NewTaskAssigneeRepository(db)

	// Services
	authService := authUsecase.NewAuthUseCase(userRepo)
	userService := userUsecase.NewUserUseCase(userRepo)
	boardService := boardUsecase.NewBoardUseCase(boardRepo)
	boardMemberService := boardMemberUsecase.NewBoardMemberUseCase(boardRepo, boardMemberRepo)
	boardColumnService := boardColumnUsecase.NewBoardColumnUseCase(boardColumnRepo)
	taskService := taskUsecase.NewTaskUseCase(boardRepo, taskRepo)
	assigneeService := assigneeUsecase.NewTaskAssigneeUseCase(boardRepo, assigneeRepo)

	// Controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	boardController := controller.NewBoardController(boardService)
	boardMemberController := controller.NewBoardMemberController(boardMemberService)
	boardColumnController := controller.NewBoardColumnController(boardColumnService)
	taskController := controller.NewTaskController(taskService)
	assigneeController := controller.NewTaskAssigneeController(assigneeService)
	// -------------------------------------------------------------------------------

	// Login
	e.POST("/login", authController.Login)

	// Users
	userGroup := e.Group("/users")
	userGroup.GET("", userController.GetUsers, m.JWTMiddleware())
	userGroup.GET("/:id", userController.GetUser, m.JWTMiddleware())
	userGroup.POST("", userController.CreateUser)
	userGroup.PUT("/:id", userController.UpdateUser, m.JWTMiddleware())
	// userGroup.DELETE("/:id", userController.DeleteUser, m.JWTMiddleware())

	// Boards
	boardGroup := e.Group("/boards")
	boardGroup.GET("", boardController.GetBoards, m.JWTMiddleware())
	boardGroup.GET("/:id", boardController.GetBoardById, m.JWTMiddleware())
	boardGroup.POST("", boardController.CreateBoard, m.JWTMiddleware())
	boardGroup.PUT("/:id", boardController.UpdateBoard, m.JWTMiddleware())
	boardGroup.DELETE("/:id", boardController.DeleteBoard, m.JWTMiddleware())

	// Board Members
	boardMemberGroup := e.Group("/board-members")
	boardMemberGroup.POST("/add", boardMemberController.AddMember, m.JWTMiddleware())
	boardMemberGroup.POST("/remove", boardMemberController.RemoveMember, m.JWTMiddleware())
	boardMemberGroup.POST("/exit", boardMemberController.ExitBoard, m.JWTMiddleware())

	// Board Column
	columnGroup := e.Group("/board-columns")
	columnGroup.GET("", boardColumnController.GetColumns, m.JWTMiddleware())
	columnGroup.GET("/:id", boardColumnController.GetColumn, m.JWTMiddleware())
	columnGroup.POST("", boardColumnController.CreateColumn, m.JWTMiddleware())
	columnGroup.PUT("/:id", boardColumnController.UpdateColumn, m.JWTMiddleware())
	columnGroup.DELETE("/:id", boardColumnController.DeleteColumn, m.JWTMiddleware())

	// Tasks
	taskGroup := e.Group("/tasks")
	taskGroup.GET("", taskController.GetTasks, m.JWTMiddleware())
	taskGroup.GET("/:id", taskController.GetTaskById, m.JWTMiddleware())
	taskGroup.POST("", taskController.CreateTask, m.JWTMiddleware())
	taskGroup.PUT("/:id", taskController.UpdateTask, m.JWTMiddleware())
	taskGroup.DELETE("/:id", taskController.DeleteTask, m.JWTMiddleware())

	// Task Assignee
	assigneeGroup := e.Group("/task-assignees")
	assigneeGroup.POST("/add", assigneeController.AddAssignee, m.JWTMiddleware())
	assigneeGroup.POST("/remove", assigneeController.RemoveAssignee, m.JWTMiddleware())
	assigneeGroup.POST("/exit", assigneeController.ExitTask, m.JWTMiddleware())
}
