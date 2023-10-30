package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	boardRepo "kanban-board/repository/board"
	taskRepo "kanban-board/repository/task"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var mockTasksData = []model.Task{
	{
		Model:         gorm.Model{ID: 1},
		Title:         "Task 1",
		Desc:          "Task 1 Desc",
		BoardColumnID: 1,
		BoardID:       1,
	},
	{
		Model:         gorm.Model{ID: 2},
		Title:         "Task 2",
		Desc:          "Task 2 Desc",
		BoardColumnID: 2,
		BoardID:       1,
	},
}

var mockMembersData = []model.BoardMember{
	{
		BoardID: uint(1),
		UserID:  uint(1),
	},
	{
		BoardID: uint(1),
		UserID:  uint(2),
	},
}

func TestIsOwner(t *testing.T) {
	var issuerId = uint(1)
	var boardId = uint(1)
	t.Run("Success isOwner", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardOwner", boardId).Return(&issuerId, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.isOwner(issuerId, boardId)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isOwner (Not owner)", func(t *testing.T) {
		expectedErr := errors.New("User is not owner of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardOwner", boardId).Return(&issuerId, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.isOwner(uint(10), boardId)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isOwner (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardOwner", boardId).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.isOwner(issuerId, boardId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestIsMember(t *testing.T) {
	var issuerId = uint(1)
	var boardId = uint(1)
	t.Run("Success isMember", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(mockMembersData, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockMemberRepo)
		err := service.isMember(issuerId, boardId)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isMember (Not a member)", func(t *testing.T) {
		expectedErr := errors.New("User is not member of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(mockMembersData, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockMemberRepo)
		err := service.isMember(uint(10), boardId)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed isMember (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockMemberRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockMemberRepo)
		err := service.isMember(issuerId, boardId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestGetTasks(t *testing.T) {
	boardId := uint(1)
	issuerId := uint(1)
	t.Run("Success Get Tasks", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("Get", boardId).Return(mockTasksData, nil).Once()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(mockMembersData, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		tasks, err := service.GetTasks(boardId, issuerId)

		assert.NoError(t, err)
		assert.NotNil(t, tasks)
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Tasks (isMember false)", func(t *testing.T) {
		expectedErr := errors.New("User is not member of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		tasks, err := service.GetTasks(boardId, issuerId)

		assert.Error(t, err)
		assert.Nil(t, tasks)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Tasks (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("Get", boardId).Return(nil, expectedErr).Once()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(mockMembersData, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		tasks, err := service.GetTasks(boardId, issuerId)

		assert.Error(t, err)
		assert.Nil(t, tasks)
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestGetTaskById(t *testing.T) {
	issuerId := uint(1)
	boardId := uint(1)
	t.Run("Success get task by id", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByTaskId", mockTasksData[0].ID).Return(&boardId, nil).Once()
		mockBoardRepo.On("GetBoardMembers", uint(1)).Return(mockMembersData, nil).Once()
		mockTaskRepo.On("GetById", mockTasksData[0].ID, issuerId).Return(&mockTasksData[0], nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		task, err := service.GetTaskById(mockTasksData[0].ID, issuerId)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed get task by id (error on getting boardId)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByTaskId", mockTasksData[0].ID).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		task, err := service.GetTaskById(mockTasksData[0].ID, issuerId)

		assert.Error(t, err)
		assert.Nil(t, task)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("Failed get task by id (isMember is false)", func(t *testing.T) {
		expectedErr := errors.New("User is not member of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByTaskId", mockTasksData[0].ID).Return(&boardId, nil).Once()
		mockBoardRepo.On("GetBoardMembers", uint(1)).Return(mockMembersData, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		task, err := service.GetTaskById(mockTasksData[0].ID, uint(10))

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Equal(t, expectedErr.Error(), err.Error())
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed get task by id (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByTaskId", mockTasksData[0].ID).Return(&boardId, nil).Once()
		mockBoardRepo.On("GetBoardMembers", uint(1)).Return(mockMembersData, nil).Once()
		mockTaskRepo.On("GetById", mockTasksData[0].ID, issuerId).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		task, err := service.GetTaskById(mockTasksData[0].ID, issuerId)

		assert.Error(t, err)
		assert.Nil(t, task)
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestCreateTask(t *testing.T) {
	issuerId := uint(1)

	t.Run("Success Create Task", func(t *testing.T) {
		createTaskRequest := &dto.TaskCreateRequest{
			Title:         "Task 1",
			Desc:          "Task 1 Desc",
			BoardColumnID: 1,
			BoardID:       1,
		}

		taskModel := &model.Task{
			Title:         createTaskRequest.Title,
			Desc:          createTaskRequest.Desc,
			BoardColumnID: createTaskRequest.BoardColumnID,
			BoardID:       createTaskRequest.BoardID,
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("Create", issuerId, taskModel).Return(nil).Once()
		mockBoardRepo.On("GetBoardMembers", createTaskRequest.BoardID).Return(mockMembersData, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.CreateTask(issuerId, createTaskRequest)

		assert.NoError(t, err)
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Create Task (Struct Invalid)", func(t *testing.T) {
		createTaskRequest := &dto.TaskCreateRequest{
			Title:         "Task 1",
			Desc:          "Task 1 Desc",
			BoardColumnID: 1,
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardMembers", createTaskRequest.BoardID).Return(mockMembersData, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.CreateTask(issuerId, createTaskRequest)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Create Task (isMember false)", func(t *testing.T) {
		createTaskRequest := &dto.TaskCreateRequest{
			Title:         "Task 1",
			Desc:          "Task 1 Desc",
			BoardColumnID: 1,
		}

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardMembers", createTaskRequest.BoardID).Return(mockMembersData, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.CreateTask(uint(10), createTaskRequest)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Create Task (Internal Server Error)", func(t *testing.T) {
		createTaskRequest := &dto.TaskCreateRequest{
			Title:         "Task 1",
			Desc:          "Task 1 Desc",
			BoardColumnID: 1,
			BoardID:       1,
		}
		taskModel := &model.Task{
			Title:         createTaskRequest.Title,
			Desc:          createTaskRequest.Desc,
			BoardColumnID: createTaskRequest.BoardColumnID,
			BoardID:       createTaskRequest.BoardID,
		}
		expectedErr := errors.New("Internal Server Error")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("Create", issuerId, taskModel).Return(expectedErr).Once()
		mockBoardRepo.On("GetBoardMembers", createTaskRequest.BoardID).Return(mockMembersData, nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.CreateTask(issuerId, createTaskRequest)

		assert.Error(t, err)
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestUpdateTask(t *testing.T) {
	issuerId := uint(1)
	boardId := uint(1)
	updateTaskRequest := &dto.TaskUpdateRequest{
		Title:         "Update Title",
		Desc:          "Update Desc",
		BoardColumnID: uint(1),
	}
	t.Run("Success Update Task", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByColumnId", updateTaskRequest.BoardColumnID).Return(&boardId, nil).Once()
		mockBoardRepo.On("GetBoardMembers", uint(1)).Return(mockMembersData, nil).Once()
		mockTaskRepo.On("Update", mockTasksData[0].ID, issuerId, updateTaskRequest).Return(nil).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.UpdateTask(mockTasksData[0].ID, issuerId, updateTaskRequest)

		assert.NoError(t, err)
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Task (isMember false)", func(t *testing.T) {
		expectedErr := errors.New("User is not member of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByColumnId", updateTaskRequest.BoardColumnID).Return(&boardId, nil).Once()
		mockBoardRepo.On("GetBoardMembers", uint(1)).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.UpdateTask(mockTasksData[0].ID, issuerId, updateTaskRequest)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Task (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByColumnId", updateTaskRequest.BoardColumnID).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.UpdateTask(mockTasksData[0].ID, issuerId, updateTaskRequest)

		assert.Error(t, err)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Task (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByColumnId", updateTaskRequest.BoardColumnID).Return(&boardId, nil).Once()
		mockBoardRepo.On("GetBoardMembers", uint(1)).Return(mockMembersData, nil).Once()
		mockTaskRepo.On("Update", mockTasksData[0].ID, issuerId, updateTaskRequest).Return(expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.UpdateTask(mockTasksData[0].ID, issuerId, updateTaskRequest)

		assert.Error(t, err)
		mockTaskRepo.AssertExpectations(t)
		mockBoardRepo.AssertExpectations(t)
	})
}

func TestDeleteTask(t *testing.T) {
	issuerId := uint(1)
	boardId := uint(1)

	t.Run("Success Delete Task", func(t *testing.T) {
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(mockMembersData, nil).Once()
		mockTaskRepo.On("GetBoardIdByTaskId", mockTasksData[0].ID).Return(&boardId, nil).Once()
		mockTaskRepo.On("Delete", mockTasksData[0].ID, issuerId).Return(nil).Once()
		
		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.DeleteTask(mockTasksData[0].ID, issuerId)

		assert.NoError(t, err)
		mockBoardRepo.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Task (isMember false)", func(t *testing.T) {
		expectedErr := errors.New("User is not member of this board!")
		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByTaskId", mockTasksData[0].ID).Return(&boardId, nil).Once()
		mockBoardRepo.On("GetBoardMembers", uint(1)).Return(nil, expectedErr).Once()
		
		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.DeleteTask(mockTasksData[0].ID, issuerId)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedErr.Error())
		mockBoardRepo.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Task (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByTaskId", mockTasksData[0].ID).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.DeleteTask(mockTasksData[0].ID, issuerId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
	})
	t.Run("Failed Delete Task (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByTaskId", mockTasksData[0].ID).Return(&boardId, nil).Once()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.DeleteTask(mockTasksData[0].ID, issuerId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
	})
	t.Run("Failed Delete Task (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")

		mockBoardRepo := boardRepo.NewMockBoardRepo()
		mockTaskRepo := taskRepo.NewMockTaskRepo()
		mockTaskRepo.On("GetBoardIdByTaskId", mockTasksData[0].ID).Return(&boardId, nil).Once()
		mockBoardRepo.On("GetBoardMembers", boardId).Return(mockMembersData, nil).Once()
		mockTaskRepo.On("Delete", mockTasksData[0].ID, issuerId).Return(expectedErr).Once()

		service := NewTaskUseCase(mockBoardRepo, mockTaskRepo)
		err := service.DeleteTask(mockTasksData[0].ID, issuerId)

		assert.Error(t, err)
		mockBoardRepo.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
	})
}
