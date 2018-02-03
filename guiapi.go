package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mbertschler/blocks/html"
)

func guiAPI() Handler {
	handler := Handler{
		Functions: map[string]Callable{
			"hello": helloHandler,
		},
	}
	// testHello(handler)
	// testHelloHandler()
	return handler
}

func testHello(handler Handler) {
	req := &Request{
		Actions: []Action{
			{
				ID:   123,
				Name: "hello",
				Args: []byte(`{"name":"Martin"}`),
			},
			{
				ID:   124,
				Name: "hellos",
				Args: []byte(`{"name":"Lucas"}`),
			},
		},
	}
	resp := handler.Handle(req)
	log.Printf("%+v", resp)
}

func testHelloHandler() {
	args := map[string]interface{}{
		"name": "Martin",
		"age":  3,
	}
	arg, err := json.Marshal(args)
	if err != nil {
		log.Fatal(err)
	}
	res, err := helloHandler(arg)
	log.Printf("%+v %#v", res, err)
}

func helloHandler(in json.RawMessage) (*Result, error) {
	type argType struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	args := argType{}
	err := json.Unmarshal(in, &args)
	if err != nil {
		return nil, err
	}
	block := html.H1(nil, html.Text("Hello "+args.Name))
	out, err := html.RenderString(block)
	if err != nil {
		return nil, err
	}
	ret := Result{
		HTML: []HTMLUpdate{
			{
				Operation: HTMLReplace,
				Selector:  "#container",
				Content:   out,
			},
		},
	}
	return &ret, nil
}

// ============================================
// Logic
// ============================================

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req Request
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	resp := h.Handle(&req)
	enc := json.NewEncoder(w)
	err = enc.Encode(resp)
	if err != nil {
		log.Println("encoding error:", err)
	}
}

func (h Handler) Handle(req *Request) *Response {
	var resp Response
	for _, action := range req.Actions {
		var res = Result{
			ID:   action.ID,
			Name: action.Name,
		}
		fn, ok := h.Functions[action.Name]
		if !ok {
			res.Error = &Error{
				Code:    "undefinedFunction",
				Message: fmt.Sprint(action.Name, " is not defined"),
			}

		} else {
			r, err := fn(action.Args)
			if err != nil {
				res.Error = &Error{
					Code:    "error",
					Message: err.Error(),
				}
			}
			res.HTML = r.HTML
			res.JS = r.JS
		}
		resp.Results = append(resp.Results, res)
	}
	return &resp
}

// ============================================
// Types
// ============================================

// Request is the sent body of a GUI API call
type Request struct {
	Actions []Action
}

type Action struct {
	ID   int    `json:",omitempty"` // ID can be used from the client to identify responses
	Name string // Name of the action that is called
	// Args as object, gets parsed by the called function
	Args json.RawMessage `json:",omitempty"`
}

type Handler struct {
	Functions map[string]Callable
}

type Callable func(args json.RawMessage) (*Result, error)

// Response is the returned body of a GUI API call
type Response struct {
	Results []Result
}

type Result struct {
	ID    int          `json:",omitempty"` // ID from the calling action is returned
	Name  string       // Name of the action that was called
	Error *Error       `json:",omitempty"`
	HTML  []HTMLUpdate `json:",omitempty"` // DOM updates to apply
	JS    []JSCall     `json:",omitempty"` // JS calls to execute
}

type Error struct {
	Code    string
	Message string
}

type HTMLOp int8

const (
	HTMLReplace HTMLOp = iota + 1
	HTMLDelete
	HTMLAppend
	HTMLPrepend
)

type HTMLUpdate struct {
	Operation HTMLOp // how to apply this update
	Selector  string // jQuery style selector: #id .class
	Content   string `json:",omitempty"` // inner HTML
	// Init calls are executed after the HTML is added
	Init []JSCall `json:",omitempty"`
	// Destroy calls are executed before the HTML is removed
	Destroy []JSCall `json:",omitempty"`
}

type JSCall struct {
	Name string // name of the function to call
	// Arguments as object, gets encoded by the called function
	Arguments json.RawMessage `json:",omitempty"`
}
