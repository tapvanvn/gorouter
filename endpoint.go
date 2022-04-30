package gorouter

import "strings"

//EndpointDefine ...
type EndpointDefine struct {
	Measurement bool
	Handles     []RouteHandle
}

func GetBlueprint(endpointDefine string) string {
	parts := strings.Split(endpointDefine, "/")
	remain := []string{}
	for _, part := range parts {
		if part[0] != ':' {
			remain = append(remain, part)
		}
	}
	return strings.Join(remain, "/")
}
