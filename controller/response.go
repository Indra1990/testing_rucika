package controller

import (
	"encoding/json"
	"net/http"
)

type responseCode struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Desc string `json:"description"`
}
type responseResult struct {
	ResponseCode responseCode `json:"responseCode"`
	Result       interface{}  `json:"result"`
}
type responseWithResult struct {
	ResponseWithResult interface{} `json:"responseWithResult"`
}

func responseWithJson(w http.ResponseWriter, code int, name string, desc string, result interface{}) {
	response := responseWithResult{
		responseResult{
			ResponseCode: responseCode{
				Code: code,
				Name: name,
				Desc: desc,
			},
			Result: result,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
