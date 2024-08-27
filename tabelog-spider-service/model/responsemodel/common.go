package responsemodel

import "tabelog-spider/model/enum"

type CommonErrorResponse struct {
	ErrorCode    enum.ErrorCode `json:"error_code"`
	ErrorMessage string         `json:"error_message"`
	Result       *interface{}   `json:"result"`
}

type CommonResponse struct {
	Result interface{} `json:"result"`
}
