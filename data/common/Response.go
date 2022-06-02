package common

import (
	"net/http"
)

type BasicResponse struct {
	Status   int         `json:"status"`
	Data     interface{} `json:"data,omitempty"`
	Message  string      `json:"message,omitempty"`
	PageInfo *PageInfo   `json:"page_info,omitempty"`
}

type PageInfo struct {
	PageNumber       string `json:"page_number"`
	PageSize         string `json:"page_size"`
	TotalRecordCount string `json:"total_record_count"`
}

func CreateSuccessResponse(data interface{}) BasicResponse {
	return BasicResponse{
		Status: http.StatusOK,
		Data:   data,
	}
}

func CreateErrorResponse(status int, errMsg string) BasicResponse {
	return BasicResponse{
		Status:  status,
		Message: errMsg,
	}
}
