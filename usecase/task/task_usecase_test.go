package usecase

import (
	"errors"
	"kanban-board/dto"
	"kanban-board/model"
	repo "kanban-board/repository/task"
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

func TestGetTasks(t *testing.T) {
	boardId := uint(1)
	t.Run("Success Get Tasks", func(t *testing.T) {
		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("Get", boardId).Return(mockTasksData, nil).Once()

		service := NewTaskUseCase(mockRepo)
		tasks, err := service.GetTasks(boardId)

		assert.NoError(t, err)
		assert.NotNil(t, tasks)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Get Tasks (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("Get", boardId).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockRepo)
		tasks, err := service.GetTasks(boardId)

		assert.Error(t, err)
		assert.Nil(t, tasks)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetTaskById(t *testing.T) {
	issuerId := uint(1)
	t.Run("Success get task by id", func(t *testing.T) {
		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("GetById", mockTasksData[0].ID, issuerId).Return(&mockTasksData[0], nil).Once()

		service := NewTaskUseCase(mockRepo)
		task, err := service.GetTaskById(mockTasksData[0].ID, issuerId)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed get task by id (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")
		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("GetById", mockTasksData[0].ID, issuerId).Return(nil, expectedErr).Once()

		service := NewTaskUseCase(mockRepo)
		task, err := service.GetTaskById(mockTasksData[0].ID, issuerId)

		assert.Error(t, err)
		assert.Nil(t, task)
		mockRepo.AssertExpectations(t)
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

		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("Create", issuerId, taskModel).Return(nil).Once()

		service := NewTaskUseCase(mockRepo)
		err := service.CreateTask(issuerId, createTaskRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Create Task (Struct Invalid)", func(t *testing.T) {
		createTaskRequest := &dto.TaskCreateRequest{
			Title:         "Task 1",
			Desc:          "Task 1 Desc",
			BoardColumnID: 1,
		}

		mockRepo := repo.NewMockTaskRepo()

		service := NewTaskUseCase(mockRepo)
		err := service.CreateTask(issuerId, createTaskRequest)

		assert.Error(t, err)
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

		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("Create", issuerId, taskModel).Return(expectedErr).Once()

		service := NewTaskUseCase(mockRepo)
		err := service.CreateTask(issuerId, createTaskRequest)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateTask(t *testing.T) {
	issuerId := uint(1)
	updateTaskRequest := &dto.TaskUpdateRequest{
		Title:         "Update Title",
		Desc:          "Update Desc",
		BoardColumnID: uint(1),
	}
	t.Run("Success Update Task", func(t *testing.T) {
		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("Update", mockTasksData[0].ID, issuerId, updateTaskRequest).Return(nil).Once()

		service := NewTaskUseCase(mockRepo)
		err := service.UpdateTask(mockTasksData[0].ID, issuerId, updateTaskRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Update Task (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")

		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("Update", mockTasksData[0].ID, issuerId, updateTaskRequest).Return(expectedErr).Once()

		service := NewTaskUseCase(mockRepo)
		err := service.UpdateTask(mockTasksData[0].ID, issuerId, updateTaskRequest)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteTask(t *testing.T) {
	issuerId := uint(1)
	t.Run("Success Delete Task", func(t *testing.T) {
		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("Delete", mockTasksData[0].ID, issuerId).Return(nil).Once()

		service := NewTaskUseCase(mockRepo)
		err := service.DeleteTask(mockTasksData[0].ID, issuerId)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Delete Task (Internal Server Error)", func(t *testing.T) {
		expectedErr := errors.New("Internal Server Error")

		mockRepo := repo.NewMockTaskRepo()
		mockRepo.On("Delete", mockTasksData[0].ID, issuerId).Return(expectedErr).Once()

		service := NewTaskUseCase(mockRepo)
		err := service.DeleteTask(mockTasksData[0].ID, issuerId)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
