package gorouter

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
)

type Method string

type ContentType string

const (
	MethodGet  = Method(http.MethodGet)
	MethodPost = Method(http.MethodPost)
	MethodPut  = Method(http.MethodPut)

	ContentTypeJson      = ContentType("application/json")
	ContentTypeFormParam = ContentType("multipart/form-data")
)

type ApiResponse struct {
	Base *http.Response
}

func (api *ApiResponse) Close() {
	if api.Base != nil {
		api.Base.Body.Close()
	}
}

func NewResponse(base *http.Response) *ApiResponse {
	res := &ApiResponse{
		Base: base,
	}
	return res
}

type IRequest interface {
	Request(domain string, path string, indexes map[string]interface{}) (*ApiResponse, error)
}

//ApiForm form the request
type ApiForm struct {
	Method      Method
	ContentType ContentType
	Headers     map[string]string //set in header of request
	Params      map[string]string //append in url of request
	Data        interface{}       //body of request
}

func (frm *ApiForm) Request(domain string, path string, indexes map[string]interface{}) (*ApiResponse, error) {

	params := url.Values{}
	//paramSegments := []string{}
	for key, value := range frm.Params {
		params.Add(key, value)
		//paramSegments = append(paramSegments, fmt.Sprintf("%s=%s", ))
	}

	urlStr := strings.TrimSuffix(domain, "/") + "/" + strings.TrimPrefix(path, "/")
	if len(frm.Params) > 0 {
		urlStr += "?" + params.Encode()
	}
	var req *http.Request = nil
	var err error = nil
	if frm.Method == MethodPost || frm.Method == MethodPut {
		bodyData := ""
		switch frm.ContentType {
		case ContentTypeFormParam:
			//TODO: build from param
			break
		case ContentTypeJson:
			//TODO: build json
			break
		}
		req, err = http.NewRequest(string(frm.Method), urlStr, bytes.NewBuffer([]byte(bodyData)))
		if frm.ContentType != "" {
			req.Header.Set("Content-Type", string(frm.ContentType))
		}
		defer req.Body.Close()
	} else {
		req, err = http.NewRequest(string(frm.Method), urlStr, nil)
	}
	if err != nil {
		return nil, err
	}
	for key, value := range frm.Headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return nil, err
	}

	//defer resp.Body.Close()
	return NewResponse(resp), nil
}

func NewApiForm() *ApiForm {
	return &ApiForm{
		Method:  MethodGet,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
}

//
func NewGetForm() *ApiForm {
	return NewApiForm()
}

func NewPostForm() *ApiForm {
	frm := NewApiForm()
	frm.Method = MethodPost
	return frm
}
