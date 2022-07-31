package gorouter

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var __serviced uint64 = 0
var __serviced_mutex sync.Mutex

//RouteHandle handle a route
type RouteHandle func(*RouteContext)

//Router router
type Router struct {
	//Handler
	root           RouteDefine
	unhandle       RouteHandle
	maintainHandle RouteHandle
	prefix         string
	prefixLen      int
	measurement    bool
	maintainTo     int64
}

func defaultUnHandler(ctx *RouteContext) {

	if ctx.W != nil {
		ctx.W.WriteHeader(http.StatusNotFound)
	}
}

func defaultMaintainHandler(ctx *RouteContext) {

	if ctx.W != nil {

		ctx.W.WriteHeader(http.StatusServiceUnavailable)
		ctx.W.Write([]byte("server maintaining"))
	}
}

//SetMeasurement set measurement flag
func (router *Router) SetMeasurement(active bool) {
	router.measurement = active
}

//SetMaintainTime set maintaining time to time point
func (router *Router) SetMaintainTime(timestamp int64) {
	router.maintainTo = timestamp
}

//Init init router
func (router *Router) Init(prefix string, define string, endpoints map[string]EndpointDefine) {

	router.root = RouteDefine{}
	router.unhandle = defaultUnHandler
	router.maintainHandle = defaultMaintainHandler
	router.prefix = prefix
	router.prefixLen = len(prefix)

	defineSub := map[string]*RouteDefine{}

	err := json.Unmarshal([]byte(define), &defineSub)

	if err != nil {
		fmt.Println("cannot parse define")
		panic(err)
	}

	router.root.Subs = defineSub

	for path, endpoint := range endpoints {

		routeDefine := router.FindRoute(path)

		if routeDefine != nil {

			routeDefine.Endpoint = endpoint

		}
	}
	if unhandleEndpoint, ok := endpoints["unhandle"]; ok && len(unhandleEndpoint.Handles) > 0 {

		router.unhandle = unhandleEndpoint.Handles[0]
	}
	if maintainEndpoint, ok := endpoints["maintain"]; ok && len(maintainEndpoint.Handles) > 0 {

		router.maintainHandle = maintainEndpoint.Handles[0]
	}
	router.PrintDebug()
}

//Init router in with root handler also has indexes
func (router *Router) InitWithIndexes(prefix string, indexes []string, define string, endpoints map[string]EndpointDefine) {

	router.root = RouteDefine{}
	router.unhandle = defaultUnHandler
	router.maintainHandle = defaultMaintainHandler
	router.prefix = prefix
	router.prefixLen = len(prefix)
	router.root.Indexes = indexes

	defineSub := map[string]*RouteDefine{}

	err := json.Unmarshal([]byte(define), &defineSub)

	if err != nil {
		fmt.Println("cannot partse define")
		panic(err)
	}

	router.root.Subs = defineSub

	for path, endpoint := range endpoints {

		routeDefine := router.FindRoute(path)

		if routeDefine != nil {

			routeDefine.Endpoint = endpoint
		}
	}
	if unhandleEndpoint, ok := endpoints["unhandle"]; ok && len(unhandleEndpoint.Handles) > 0 {

		router.unhandle = unhandleEndpoint.Handles[0]
	}
	if maintainEndpoint, ok := endpoints["maintain"]; ok && len(maintainEndpoint.Handles) > 0 {

		router.maintainHandle = maintainEndpoint.Handles[0]
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

	idicates := strings.Split(pattern, ";")

	maxI := len(idicates)

	if maxI > maxF {

		return rs
	}
	i := 0
	for _, format := range formats {

		if i >= maxI {

			break
		}
		decodedValue, err := url.QueryUnescape(idicates[i])
		if err == nil {
			rs[format] = decodedValue
		}
		i++
	}

	return rs
}

//Route route
func (router *Router) Route(path string, w http.ResponseWriter, r *http.Request) {

	now := time.Now()

	var context *RouteContext = &RouteContext{
		Path:         "/",
		Action:       "index",
		Indexes:      map[string]string{},
		Parent:       nil,
		RestPatterns: []string{},
		Handled:      false,
		Dictionary:   map[string]interface{}{},
		W:            w,
		R:            r,
	}

	if router.maintainTo > now.Unix() {

		router.maintainHandle(context)
		return
	}

	if router.prefixLen > 0 && !strings.HasPrefix(path, router.prefix) {

		router.unhandle(context)
		return
	}

	routePath := path[router.prefixLen:]

	patterns := []string{}

	if len(routePath) > 0 {

		patterns = strings.Split(routePath, "/")
	}

	i := 0
	maxI := len(patterns)

	var routeDefine *RouteDefine = &router.root

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
				Dictionary:   map[string]interface{}{},
				W:            w,
				R:            r,
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

	for _, handler := range routeDefine.Endpoint.Handles {

		handler(context)
		if context.Handled {

			break
		}
	}

	if !context.Handled {

		router.unhandle(context)
	}

	if (router.measurement || routeDefine.Endpoint.Measurement) && r != nil {

		processTime := time.Now().Sub(now).Nanoseconds()

		__serviced_mutex.Lock()
		__serviced++

		if __serviced == math.MaxUint64 {
			__serviced = 0
		}
		__serviced_mutex.Unlock()

		fmt.Printf("mersure: %s %0.2fms serviced:%d\n", r.URL.Path, float32(processTime/1_000_000), __serviced)
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
	if len(define.Endpoint.Handles) > 0 {

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
