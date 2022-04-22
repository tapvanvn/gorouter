package gorouter

func NewEmptyRouteDefine() *RouteDefine {
	return &RouteDefine{
		Indexes: []string{},
		Subs:    map[string]*RouteDefine{},
	}
}

//RouteDefine define a route
type RouteDefine struct {
	Indexes  []string                `json:"indexes,omitempty"`
	Subs     map[string]*RouteDefine `json:"subs,omitempty"`
	Endpoint EndpointDefine          `json:"-"`
}

//SubRoute is subroute existed ?
func (define *RouteDefine) SubRoute(name string) *RouteDefine {

	if sub, ok := define.Subs[name]; ok {

		return sub
	}
	return nil
}
