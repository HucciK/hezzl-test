package itemHandler

import "encoding/json"

type BadRequest struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (i ItemHandler) responseWithError(statusCode int, message string) []byte {
	response := BadRequest{StatusCode: statusCode, Message: message}
	reponseJSON, _ := json.Marshal(&response)
	return reponseJSON
}
