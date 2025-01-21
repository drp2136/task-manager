package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/mocks"
	"task-manager/internal/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errMock         = errors.New("internal error")
	errMockNotFound = errors.New("record not found")
	uuid1, _        = uuid.NewV7()
)

func Test_GetTasks(t *testing.T) {
	// Create a new http request
	req, err := http.NewRequest(http.MethodGet, "/tasks/", nil)
	require.Nil(t, err)
	req = req.WithContext(context.Background())

	taskService := new(mocks.ITaskService)
	taskHandler := NewTaskHandler(taskService)

	// Test case 1
	t.Run("GetTasks: error", func(t *testing.T) {
		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		taskService.On("GetAllTasks", mock.Anything).
			Return(nil, errMock).Once()

		// Call the GetTasks function
		taskHandler.GetTasks(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusInternalServerError, w.Code)
		// Define the expected response
		expectedResp := `{"message":"Internal Server Error"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 2
	t.Run("GetTasks: success", func(t *testing.T) {
		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		var tasks = []model.Task{{ID: uuid1, Title: "Task 1", Description: "Description 1", Status: "pending"}}
		taskService.On("GetAllTasks", mock.Anything).
			Return(tasks, nil).Once()

		// Call the GetTasks function
		taskHandler.GetTasks(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusOK, w.Code)
		// Define the expected response

		var respObj []model.Task
		err := json.Unmarshal([]byte(resp), &respObj)
		require.Nil(t, err)
		require.Equal(t, tasks, respObj)
	})
}

func Test_CreateTask(t *testing.T) {
	taskService := new(mocks.ITaskService)
	taskHandler := NewTaskHandler(taskService)

	// Test case 1
	t.Run("CreateTask: input validation error", func(t *testing.T) {
		var tasks = model.Task{Description: "Description 1", Status: "pending"}
		body, err := json.Marshal(tasks)
		require.Nil(t, err)

		// Create a new http request
		req, err := http.NewRequest(http.MethodPost, "/tasks/", bytes.NewReader(body))
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the CreateTask function
		taskHandler.CreateTask(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusBadRequest, w.Code)
		fmt.Println(resp)
	})

	// Test case 2
	t.Run("CreateTask: error", func(t *testing.T) {
		var tasks = model.Task{ID: uuid1, Title: "Task 1", Description: "Description 1", Status: "pending"}
		body, err := json.Marshal(tasks)
		require.Nil(t, err)

		// Create a new http request
		req, err := http.NewRequest(http.MethodPost, "/tasks/", bytes.NewReader(body))
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		taskService.On("CreateTask", mock.Anything, mock.AnythingOfType("*model.Task")).
			Return(errMock).Once()

		// Call the CreateTask function
		taskHandler.CreateTask(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusInternalServerError, w.Code)
		// Define the expected response
		expectedResp := `{"message":"Internal Server Error"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 3
	t.Run("CreateTask: success", func(t *testing.T) {
		var task = model.Task{ID: uuid1, Title: "Task 1", Description: "Description 1", Status: "pending"}
		body, err := json.Marshal(task)
		require.Nil(t, err)

		// Create a new http request
		req, err := http.NewRequest(http.MethodPost, "/tasks/", bytes.NewReader(body))
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		taskService.On("CreateTask", mock.Anything, mock.AnythingOfType("*model.Task")).
			Return(nil).Once()

		// Call the CreateTask function
		taskHandler.CreateTask(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusCreated, w.Code)
		var respObj model.Task
		err = json.Unmarshal([]byte(resp), &respObj)
		require.Nil(t, err)
		require.Equal(t, task.Title, respObj.Title)
		require.Equal(t, task.Description, respObj.Description)
	})
}

func Test_GetTaskByID(t *testing.T) {
	taskService := new(mocks.ITaskService)
	taskHandler := NewTaskHandler(taskService)

	// Test case 1
	t.Run("GetTaskByID: invalid task id", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodGet, "/tasks/abcd1234", nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: "abcd1234"})

		// Call the GetTaskByID function
		taskHandler.GetTaskByID(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusBadRequest, w.Code)
		// Define the expected response
		expectedResp := `{"message":"Bad Request"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 2
	t.Run("GetTaskByID: record not found", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodGet, "/tasks/"+uuid1.String(), nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		taskService.On("GetTaskByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(nil, errMockNotFound).Once()

		// Call the GetTaskByID function
		taskHandler.GetTaskByID(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusNotFound, w.Code)
		// Define the expected response
		expectedResp := `{"message":"task not found"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 3
	t.Run("GetTaskByID: errors", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodGet, "/tasks/"+uuid1.String(), nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		taskService.On("GetTaskByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(nil, errMock).Once()

		// Call the GetTaskByID function
		taskHandler.GetTaskByID(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusInternalServerError, w.Code)
		// Define the expected response
		expectedResp := `{"message":"Internal Server Error"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 4
	t.Run("GetTaskByID: success", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodGet, "/tasks/"+uuid1.String(), nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		var task = model.Task{ID: uuid1, Title: "Task 1", Description: "Description 1", Status: "pending"}
		taskService.On("GetTaskByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(&task, nil).Once()

		// Call the GetTaskByID function
		taskHandler.GetTaskByID(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusOK, w.Code)
		// Define the expected response

		var respObj model.Task
		err = json.Unmarshal([]byte(resp), &respObj)
		require.Nil(t, err)
		require.Equal(t, task.Title, respObj.Title)
		require.Equal(t, task.Description, respObj.Description)
	})
}

func Test_UpdateTaskByID(t *testing.T) {
	taskService := new(mocks.ITaskService)
	taskHandler := NewTaskHandler(taskService)

	// Test case 1
	t.Run("UpdateTaskByID: invalid task id", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodPatch, "/tasks/abcd1234", nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: "abcd1234"})

		// Call the UpdateTaskByID function
		taskHandler.UpdateTaskByID(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusBadRequest, w.Code)
		// Define the expected response
		expectedResp := `{"message":"Bad Request"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 2
	t.Run("UpdateTaskByID: record not found", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodPatch, "/tasks/"+uuid1.String(), nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		taskService.On("GetTaskByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(nil, errMockNotFound).Once()

		// Call the UpdateTaskByID function
		taskHandler.UpdateTaskByID(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusNotFound, w.Code)
		// Define the expected response
		expectedResp := `{"message":"task not found"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 3
	t.Run("UpdateTaskByID: error find by id", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodPatch, "/tasks/"+uuid1.String(), nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		taskService.On("GetTaskByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(nil, errMock).Once()

		// Call the UpdateTaskByID function
		taskHandler.UpdateTaskByID(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusInternalServerError, w.Code)
		// Define the expected response
		expectedResp := `{"message":"Internal Server Error"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 4
	t.Run("UpdateTaskByID: input validation error", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodPatch, "/tasks/"+uuid1.String(), nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		var task = model.Task{ID: uuid1, Description: "Description 1", Status: "pending"}
		taskService.On("GetTaskByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(&task, nil).Once()

		// Call the UpdateTaskByID function
		taskHandler.UpdateTaskByID(c)

		// Check the status code
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case 5
	t.Run("UpdateTaskByID: error", func(t *testing.T) {
		var task = model.Task{ID: uuid1, Title: "Task 1", Description: "Description 1", Status: "pending"}
		body, err := json.Marshal(task)
		require.Nil(t, err)

		// Create a new http request
		req, err := http.NewRequest(http.MethodPatch, "/tasks/"+uuid1.String(), bytes.NewReader(body))
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		taskService.On("GetTaskByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(&task, nil).Once()
		taskService.On("UpdateTask", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*model.Task")).
			Return(errMock).Once()

		// Call the UpdateTaskByID function
		taskHandler.UpdateTaskByID(c)

		// Check the status code
		require.Equal(t, http.StatusInternalServerError, w.Code)
		// Define the expected response
		resp := w.Body.String()
		expectedResp := `{"message":"Internal Server Error"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 6
	t.Run("UpdateTaskByID: success", func(t *testing.T) {
		var task = model.Task{ID: uuid1, Title: "Task 1", Description: "Description 1", Status: "pending"}
		body, err := json.Marshal(task)
		require.Nil(t, err)

		// Create a new http request
		req, err := http.NewRequest(http.MethodPatch, "/tasks/"+uuid1.String(), bytes.NewReader(body))
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		taskService.On("GetTaskByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(&task, nil).Once()
		taskService.On("UpdateTask", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*model.Task")).
			Return(nil).Once()

		// Call the UpdateTaskByID function
		taskHandler.UpdateTaskByID(c)

		// Check the status code
		require.Equal(t, http.StatusOK, w.Code)
		// Define the expected response
		resp := w.Body.String()
		expectedResp := `{"message":"Task updated successfully"}`
		require.Equal(t, expectedResp, resp)
	})
}

func Test_DeleteTaskByID(t *testing.T) {
	taskService := new(mocks.ITaskService)
	taskHandler := NewTaskHandler(taskService)

	// Test case 1
	t.Run("DeleteTaskByID: invalid task id", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodGet, "/tasks/abcd1234", nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: "abcd1234"})

		// Call the DeleteTaskByID function
		taskHandler.DeleteTaskByID(c)

		resp := w.Body.String()
		// Check the status code
		require.Equal(t, http.StatusBadRequest, w.Code)
		// Define the expected response
		expectedResp := `{"message":"Bad Request"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 2
	t.Run("DeleteTaskByID: error", func(t *testing.T) {
		// Create a new http request
		req, err := http.NewRequest(http.MethodGet, "/tasks/"+uuid1.String(), nil)
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		taskService.On("DeleteTask", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(errMock).Once()

		// Call the DeleteTaskByID function
		taskHandler.DeleteTaskByID(c)

		// Check the status code
		require.Equal(t, http.StatusInternalServerError, w.Code)
		// Define the expected response
		resp := w.Body.String()
		expectedResp := `{"message":"Internal Server Error"}`
		require.Equal(t, expectedResp, resp)
	})

	// Test case 3
	t.Run("DeleteTaskByID: success", func(t *testing.T) {
		var task = model.Task{ID: uuid1, Title: "Task 1", Description: "Description 1", Status: "pending"}
		body, err := json.Marshal(task)
		require.Nil(t, err)

		// Create a new http request
		req, err := http.NewRequest(http.MethodGet, "/tasks/"+uuid1.String(), bytes.NewReader(body))
		require.Nil(t, err)
		req = req.WithContext(context.Background())

		// Create a new gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = append(c.Params, gin.Param{Key: "taskId", Value: uuid1.String()})

		taskService.On("DeleteTask", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(nil).Once()

		// Call the DeleteTaskByID function
		taskHandler.DeleteTaskByID(c)

		// Check the status code
		require.Equal(t, http.StatusOK, w.Code)
		// Define the expected response
		resp := w.Body.String()
		expectedResp := `{"message":"Task deleted successfully"}`
		require.Equal(t, expectedResp, resp)
	})
}
