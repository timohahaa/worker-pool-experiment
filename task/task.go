package task

import (
	"github.com/google/uuid"
	"net/http"
)

type Task struct {
	Id             uuid.UUID
	Request        *http.Request
	Proxy          string
	AdditionalInfo map[string]interface{}
}

func NewTask(req *http.Request, proxy string, addInf map[string]interface{}) *Task {
	id, _ := uuid.NewRandom()
	t := &Task{
		Id:             id,
		Request:        req,
		Proxy:          proxy,
		AdditionalInfo: addInf,
	}
	return t
}

type CompleteTask struct {
	Id       uuid.UUID
	Responce *http.Response
	Err      error
}

func NewCompleteTask(id uuid.UUID, resp *http.Response, err error) *CompleteTask {
	ct := &CompleteTask{
		Id:       id,
		Responce: resp,
		Err:      err,
	}
	return ct
}
