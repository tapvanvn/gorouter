package gorouter

import (
	"fmt"
	"strings"
)

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
func (define *RouteDefine) BuildRequestSegment(indexes map[string]interface{}) (string, error) {
	var segments = []string{}
	//TODO: format of index
	for _, index := range define.Indexes {
		if value, has := indexes[index]; has {
			segments = append(segments, fmt.Sprintf("%v", value))
		}
	}
	//TODO: should by pass error
	return strings.Join(segments, ";"), nil
}
