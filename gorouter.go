package gorouter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//RouteHandle handle a route
type RouteHandle func(*RouteContext, http.ResponseWriter, *http.Request) bool

//RouteContext context when a route had been parsed
type RouteContext struct {
	Path         string `json:"path"`
	Router       *Router
	Parent       *RouteContext
	Action       string            `json:"action"`
	Indexes      map[string]string `json:"indexes"`
	RestPatterns []string          `json:"rest_patterns"`
}

//Router router
type Router struct {
	//Handler
	root     RouteDefine
	unhandle RouteHandle
}

func defaultUnHandler(ctx *RouteContext, w http.ResponseWriter, r *http.Request) bool {

	w.WriteHeader(http.StatusNotFound)
	return true
}

//Init init router
func (router *Router) Init(define string, handles map[string][]RouteHandle) {

	router.root = RouteDefine{}
	router.unhandle = defaultUnHandler

	defineSub := map[string]*RouteDefine{}

	err := json.Unmarshal([]byte(define), &defineSub)

	if err != nil {
		fmt.Println("cannot partse define")
		panic(err)
	}

	router.root.Subs = defineSub

	for path, handleStack := range handles {

		routeDefine := router.FindRoute(path)

		if routeDefine != nil {

			routeDefine.Handles = handleStack

		}
	}
	if unhandle, ok := handles["unhandle"]; ok && len(unhandle) > 0 {

		router.unhandle = unhandle[0]
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
func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router.Route(r.URL.Path[1:], w, r)

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
		Path:         "/",
		Action:       "index",
		Indexes:      map[string]string{},
		Parent:       nil,
		RestPatterns: []string{},
	}

	var routeDefine *RouteDefine = &router.root
	var context *RouteContext = &root
	var testIndex bool = false

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
				Path:         context.Path + "/" + pattern,
				Action:       "index",
				Indexes:      map[string]string{},
				Parent:       context,
				RestPatterns: []string{},
			}
			context = subContext
			testIndex = false

		} else if !testIndex {

			context.Indexes = router.FormatIndex(routeDefine.Indexes, pattern)

			testIndex = true

			if len(context.Indexes) == 0 {

				context.Action = pattern
			}
		} else if context.Action == "index" {

			context.Action = pattern

		} else {

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
	var handled = false

	for _, handler := range routeDefine.Handles {

		if handler(context, w, r) {
			handled = true
			break
		}
	}

	if !handled {

		router.unhandle(context, w, r)
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
	if len(define.Handles) > 0 {

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
