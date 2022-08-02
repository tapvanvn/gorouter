package gorouter

import "strings"

//EndpointDefine ...

type EndpointDefine struct {
	Measurement bool          //should mersuaring the serving of this endpoint
	Handles     []RouteHandle //the handlers will handle the demain
	ApiFormer   IRequest
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
