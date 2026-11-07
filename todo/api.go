package todo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall/js"

	"github.com/augustoroman/promise"
	"github.com/google/uuid"
)

const baseURL = ""

func PromiseToken(p *promise.Promise) {
	go func() {
		t, err := TokenAPI()
		if err != nil {
			p.Reject(err)
		} else {
			p.Resolve(t)
		}
	}()
}

func TokenAPI() (string, error) {
	type accessToken struct {
		Token string `json:"token"`
	}
	resp, err := http.Get(baseURL + "/tokenz")
	if err != nil {
		js.Global().Call("alert", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	var t accessToken
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return "", err
	}

	return t.Token, nil
}

func PromiseTodoList(token string, p *promise.Promise) {
	go func() {
		t, err := TodoListAPI(token)
		if err != nil {
			p.Reject(err.Error())
		} else {
			p.Resolve(t)
		}
	}()
}

func TodoListAPI(token string) ([]Todo, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/todos", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	var t []Todo
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}

	return t, nil
}

func PromiseAddTodo(token, text string, p *promise.Promise) {
	go func() {
		t, err := AddTodoAPI(token, text)
		if err != nil {
			p.Reject(err.Error())
			println(err.Error())
		} else {
			p.Resolve(t)
		}
	}()
}

func AddTodoAPI(token string, text string) (Todo, error) {
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(&Todo{Text: text})

	req, err := http.NewRequest(http.MethodPost, baseURL+"/todos", buf)
	if err != nil {
		return Todo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("TransactionID", uuid.New().String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return Todo{}, err
	}
	if resp.StatusCode >= 400 {
		return Todo{}, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	var t Todo
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return Todo{}, err
	}

	return t, nil
}

func PromiseRemoveTodo(token string, id uint, p *promise.Promise) {
	go func() {
		err := RemoveTodoAPI(token, id)
		if err != nil {
			println("error:", err.Error())
			p.Reject(err.Error())
		} else {
			p.Resolve("")
		}
	}()
}

func PromiseUpdateTodo(token string, id uint, completed bool, p *promise.Promise) {
	go func() {
		t, err := UpdateTodoAPI(token, id, completed)
		if err != nil {
			p.Reject(err.Error())
		} else {
			p.Resolve(t)
		}
	}()
}

func UpdateTodoAPI(token string, id uint, completed bool) (Todo, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(struct {
		Completed bool `json:"completed"`
	}{Completed: completed}); err != nil {
		return Todo{}, err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/todos/%d", baseURL, id), buf)
	if err != nil {
		return Todo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("TransactionID", uuid.New().String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return Todo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return Todo{}, errors.New(resp.Status)
	}

	var t Todo
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return Todo{}, err
	}
	return t, nil
}

func RemoveTodoAPI(token string, id uint) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/todos/%d", baseURL, id), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("TransactionID", uuid.New().String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()

	return nil
}
