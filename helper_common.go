package main

import (
	"encoding/json"
	//. "github.com/gtechx/base/common"
	//"github.com/gtechx/chatserver/db"
)

func jsonMarshal(data interface{}, out *[]byte, perrcode *uint16) bool {
	databytes, err := json.Marshal(data)
	if err != nil {
		*perrcode = ERR_JSON_SERIALIZE
	} else {
		*out = databytes
		return true
	}

	return false
}

func jsonUnMarshal(data []byte, out interface{}, perrcode *uint16) bool {
	err := json.Unmarshal(data, out)
	if err == nil {
		return true
	}
	*perrcode = ERR_INVALID_JSON
	return false
}
