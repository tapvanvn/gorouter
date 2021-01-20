package gorouter_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/tapvanvn/gorouter"
)

var routeStructure string = `{ 
		"test1":{
			"indexes":["param"],
			"subs":{
				"test_sub_1":{}
			}
		}
	}`

func rootHandler(context *gorouter.RouteContext, w http.ResponseWriter, r *http.Request) bool {

	fmt.Println("root handle")
	return true
}

func test1Handler(context *gorouter.RouteContext, w http.ResponseWriter, r *http.Request) bool {

	if index, ok := context.Indexes["param"]; ok {

		fmt.Println("test1(", index, ")", context.Action)

	} else {

		fmt.Println("test1 " + context.Action)
	}
	return true
}

func TestRoute(t *testing.T) {

	route := gorouter.Router{}
	route.Init(routeStructure, map[string][]gorouter.RouteHandle{
		"":      {rootHandler},
		"test1": {test1Handler},
	})

	if route.FindRoute("test1") == nil {

		t.Fail()
	}

	route.Route("test1/action", nil, nil)
	route.Route("test1/index_1/action", nil, nil)
	route.Route("", nil, nil)
}
