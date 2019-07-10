package main

import (
	"encoding/json"
	"net/http"
)

type loginStruct struct {
	UserName string
	Password string
}

type apiResponse struct {
	w      http.ResponseWriter
	Status int
	Error  string
	Result map[string]string
	Data   map[string]interface{}
}

func newAPIResponse(w http.ResponseWriter) (res apiResponse) {
	res = apiResponse{
		w:      w,
		Status: 1,
		Result: make(map[string]string),
		Data:   make(map[string]interface{}),
	}

	return
}

func (res *apiResponse) Write() {
	// encode to JSON for output
	outbytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		res.w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.Status == 0 {
		// generic failure, write 400 header
		//res.w.WriteHeader(http.StatusBadRequest) // TODO enable once endpoints can all cope with 400 responses
		res.w.Write([]byte(outbytes))
		return
	}

	if res.Status == 1 {
		// generic success, write response
		res.w.Write([]byte(outbytes))
		return
	}

	// other statuses respond with that error code
	res.w.WriteHeader(res.Status)
	res.w.Write([]byte(outbytes))
}
