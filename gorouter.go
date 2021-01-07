package gorouter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//RouteHandle handle a route
type RouteHandle func(*RouteContext, http.ResponseWriter, *http.Request)

//RouteContext context when a route had been parsed
type RouteContext struct {
	Router       *Router
	Parent       *RouteContext
	Action       string            `json:"action"`
	Indexes      map[string]string `json:"indexes"`
	RestPatterns []string          `json:"rest_patterns"`
}

//Router router
type Router struct {
	//Handler
	root RouteDefine
}

//Init init router
func (router *Router) Init(define string, handles map[string]RouteHandle) {

	router.root = RouteDefine{}

	defineSub := map[string]*RouteDefine{}

	err := json.Unmarshal([]byte(define), &defineSub)

	if err != nil {
		fmt.Println("cannot partse define")
		panic(err)
	}

	router.root.Subs = defineSub

	for path, handle := range handles {

		define := router.FindRoute(path)

		if define != nil {

			define.Handle = handle

		}
	}
	router.PrintDebug()
}

//FindRoute find define of a route
func (router *Router) FindRoute(path string) *RouteDefine {

	patterns := []string{}

	if len(path) > 0 {

		patterns = strings.Split(path, "/")
	}

	i := 0
	maxI := len(patterns)

	var routeDefine *RouteDefine = &router.root

	for {
		if i >= maxI {
			break
		}
		pattern := patterns[i]
		subRouteDefine := routeDefine.SubRoute(pattern)
		if subRouteDefine == nil {
			return nil
		}
		routeDefine = subRouteDefine
		i++
	}
	return routeDefine
}

//ServeHTTP handle
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router.Route(r.URL.Path, w, r)
}

//FormatIndex get indexes
func (router *Router) FormatIndex(formats []string, pattern string) map[string]string {

	rs := map[string]string{}
	maxF := len(formats)

	if maxF == 0 {

		return rs
	}

	idicates := strings.Split(pattern, ".")

	maxI := len(idicates)

	if maxI > maxF {

		return rs
	}
	i := 0
	for _, format := range formats {

		if i >= maxI {
			break
		}
		rs[format] = idicates[i]
		i++
	}

	return rs
}

//Route route
func (router *Router) Route(path string, w http.ResponseWriter, r *http.Request) {

	patterns := []string{}

	if len(path) > 0 {

		patterns = strings.Split(path, "/")
	}

	i := 0
	maxI := len(patterns)

	root := RouteContext{
		Action:       "index",
		Indexes:      map[string]string{},
		Parent:       nil,
		RestPatterns: []string{},
	}
	var routeDefine *RouteDefine = &router.root
	var context *RouteContext = &root

	for {
		if i >= maxI {
			break
		}
		pattern := patterns[i]

		subRouteDefine := routeDefine.SubRoute(pattern)

		if subRouteDefine != nil {

			//in case we found sub route
			routeDefine = subRouteDefine
			subContext := &RouteContext{
				Action:       "index",
				Indexes:      map[string]string{},
				Parent:       context,
				RestPatterns: []string{},
			}
			context = subContext

		} else {

			context.Indexes = router.FormatIndex(routeDefine.Indexes, pattern)
			if len(context.Indexes) == 0 {
				context.Action = pattern
			}
			break
		}
		i++
	}

	for {

		if i >= maxI {

			break
		}
		if context.Action == "index" {

			context.Action = patterns[i]
		} else {
			context.RestPatterns = append(context.RestPatterns, patterns[i])
		}
		i++
	}
	if routeDefine.Handle != nil {

		routeDefine.Handle(context, w, r)
	}
}

func printDebug(name string, define *RouteDefine, level int) {

	for i := 0; i < level; i++ {
		fmt.Print(" |")
	}
	fmt.Print(name)
	if len(define.Indexes) > 0 {
		fmt.Print("(", strings.Join(define.Indexes, ","), ")")
	}
	if define.Handle != nil {
		fmt.Print((" handled"))
	} else {
		fmt.Print((" unhandle"))
	}
	fmt.Println("")
	if len(define.Subs) > 0 {
		for subName, subDefine := range define.Subs {

			printDebug(subName, subDefine, level+1)
		}
	}
}

//PrintDebug print structure
func (router *Router) PrintDebug() {

	printDebug("root", &router.root, 0)
}
