package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/KonstantinZy/gin-gorm-api/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine
var elId uint

func init() {
	if err := os.Chdir(".."); err != nil { // run from root folder of app
		panic(err)
	}

	models.StartDB("app_test.db")
	models.MigrateOrm()

	r = gin.Default()
	RegisterTaskURI(r)
}

func TestCreateRoute(t *testing.T) {

	taskToCreate := CreateTaskInput{
		Name:        "test task name",
		Description: "test task description",
		Start:       1662631155,
		Finish:      1662717555,
		Subtask: []CreateTaskSubtask{
			{Description: "test subtask 1"},
			{Description: "test subtask 2"},
			{Description: "test subtask 3"},
		},
	}

	w := httptest.NewRecorder()

	jsn, err := json.Marshal(taskToCreate)
	if err != nil {
		t.Error("Can't marshal for sending to route: " + err.Error())
	}
	buf := bytes.NewReader(jsn)
	req, _ := http.NewRequest("POST", "/task", buf)
	r.ServeHTTP(w, req)

	if !assert.Equal(t, 200, w.Code) {
		t.Error("Wrong response code:", w.Code)
	}

	var data map[string]map[string]uint
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Error("Can't unmarshal json: ", err.Error(), "body", w.Body)
	}

	if int(data["data"]["id"]) == 0 {
		t.Error("Wrong response, can not get id of element : " + err.Error())
	}

	elId = data["data"]["id"] // save for next tests
}

func TestUpdateRoute(t *testing.T) {

	newName := "name updated"

	taskToUpdate := UpdateTaskInput{
		Name: newName,
	}

	w := httptest.NewRecorder()

	jsn, err := json.Marshal(taskToUpdate)
	if err != nil {
		t.Fatal("Can't marshal for sending to route: " + err.Error())
	}
	buf := bytes.NewReader(jsn)
	req, _ := http.NewRequest("PATCH", "/task/"+strconv.Itoa(int(elId)), buf)
	r.ServeHTTP(w, req)

	if !assert.Equal(t, 200, w.Code) {
		t.Fatal("Wrong response code:", w.Code)
	}

	var data map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fatal("Can't unmarshal json: ", err.Error(), "body", w.Body)
	}

	res, ok := data["data"]
	if !ok {
		t.Fatal("Got error response:", w.Body)
	}

	if !assert.Equal(t, res, "Ok") {
		t.Fatal("Got wrong response:", w.Body)
	}

	var taskToCheck models.Task
	if err := models.DB.Where("id = ?", elId).Find(&taskToCheck).Error; err != nil {
		t.Fatal("Can't find task for cheking in orm:", err.Error())
	}

	if !assert.Equal(t, taskToCheck.Name, newName) {
		t.Fatal("Element did not updated in ORM:", taskToCheck.Name, "!=", newName)
	}
}

func TestGetRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/task", nil)
	r.ServeHTTP(w, req)

	if !assert.Equal(t, 200, w.Code) {
		t.Error("Wrong response code:", w.Code)
	}

	var data map[string][]models.Task
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Error("Can't unmarshal json: " + err.Error())
	}

	var tasksToCheck []models.Task
	if err := models.DB.Preload("SubTasks").Find(&tasksToCheck).Error; err != nil {
		t.Error("Can't find task for cheking in orm:", err.Error())
	}

	diff := deep.Equal(data["data"], tasksToCheck)
	if len(diff) > 0 {
		t.Error("Not equals tasks got from API and orm - differense", diff)
	}
}

func TestGetOneRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/task/"+strconv.Itoa(int(elId)), nil)
	r.ServeHTTP(w, req)

	if !assert.Equal(t, 200, w.Code) {
		t.Error("Wrong response code:", w.Code)
	}

	var data map[string]models.Task
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Error("Can't unmarshal json:", err.Error(), "body", w.Body)
	}

	var taskToCheck models.Task
	if err := models.DB.Where("id = ?", elId).Preload("SubTasks").Find(&taskToCheck).Error; err != nil {
		t.Error("Can't find task for cheking in orm:", err.Error())
	}

	diff := deep.Equal(data["data"], taskToCheck)
	if len(diff) > 0 {
		t.Error("Not equals tasks got from API and orm - differense", diff)
	}
}

func TestDeleteRoute(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/task/"+strconv.Itoa(int(elId)), nil)
	r.ServeHTTP(w, req)

	if !assert.Equal(t, 200, w.Code) {
		t.Fatal("Wrong response code:", w.Code)
	}

	var data map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
		t.Fatal("Can't unmarshal json: ", err.Error(), "body", w.Body)
	}

	res, ok := data["data"]
	if !ok {
		t.Fatal("Got error response:", w.Body)
	}

	if !assert.Equal(t, res, "Ok") {
		t.Fatal("Got wrong response result:", w.Body)
	}

	var taskToCheck models.Task
	if err := models.DB.Where("id = ?", elId).First(&taskToCheck).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Fatal("Element wasnt delete from DB")
		}
	} else {
		t.Fatal("Element wasnt delete from DB")
	}
}
