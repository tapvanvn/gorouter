package gorouter

import "net/http"

//RouteContext context when a route had been parsed
type RouteContext struct {
	Path         string `json:"path"`
	Router       *Router
	Parent       *RouteContext
	Action       string            `json:"action"`
	Indexes      map[string]string `json:"indexes"`
	RestPatterns []string          `json:"rest_patterns"`
	Handled      bool
	Dictionary   map[string]interface{}
	W            http.ResponseWriter
	R            *http.Request
}

//Trail pass context to handle in trail of handler
func (context *RouteContext) Trail(trails []RouteHandle) {
	for _, handler := range trails {
		if context.Handled {
			return
		}
		handler(context)
	}
}

//AnyIndex find if any index with name exists in current stack of context
func (context *RouteContext) AnyIndex(name string) (string, bool) {

	if index, ok := context.Indexes[name]; ok {

		return index, true

	} else if context.Parent != nil {

		return context.Parent.AnyIndex(name)
	}
	return "", false
}
