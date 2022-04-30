package gorouter

import "reflect"

type ReflectRoute struct {
	objectMap map[string]reflect.Type
}

func (router *ReflectRoute) route(endpointBlueprint string, context *RouteContext) {

	
}
