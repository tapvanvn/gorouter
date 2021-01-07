package gorouter_test

import (
	"fmt"
	"gorouter"
	"net/http"
	"testing"
)

var routeStructure string = `{ 
		"test1":{
			"subs":{
				"test_sub_1":{}
			}
		}
	}`

func rootHandler(context *gorouter.RouteContext, w http.ResponseWriter, r *http.Request) {

	fmt.Println("root handle")
}

func test1Handler(context *gorouter.RouteContext, w http.ResponseWriter, r *http.Request) {

	fmt.Println("test1 " + context.Action)
}

func TestRoute(t *testing.T) {

	route := gorouter.Router{}
	route.Init(routeStructure, map[string]gorouter.RouteHandle{
		"":      rootHandler,
		"test1": test1Handler,
	})

	if route.FindRoute("test1") == nil {

		t.Fail()
	}

	route.Route("test1/action", nil, nil)
	route.Route("", nil, nil)
}
