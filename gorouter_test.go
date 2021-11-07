package gorouter_test

import (
	"fmt"
	"testing"

	"github.com/tapvanvn/gorouter/v2"
)

var routeStructure string = `{ 
		"test1":{
			"indexes":["param"],
			"subs":{
				"test_sub_1":{}
			}
		}
	}`

func rootHandler(context *gorouter.RouteContext) {

	fmt.Println("root handle")
}

func test1Handler(context *gorouter.RouteContext) {

	if index, ok := context.Indexes["param"]; ok {

		fmt.Println("test1(", index, ")", context.Action)

	} else {

		fmt.Println("test1 " + context.Action)
	}
}

func TestRoute(t *testing.T) {

	route := gorouter.Router{}
	route.Init("/api", routeStructure, map[string]gorouter.EndpointDefine{
		"":      {true, []gorouter.RouteHandle{rootHandler}},
		"test1": {true, []gorouter.RouteHandle{test1Handler}},
	})

	if route.FindRoute("test1") == nil {

		t.Fail()
	}

	route.Route("/api/test1/action", nil, nil)
	route.Route("/api/test1/index_1/action", nil, nil)
	route.Route("", nil, nil)
}
