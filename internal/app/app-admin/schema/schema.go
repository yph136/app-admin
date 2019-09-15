package schema

const (
	UserID = "UserID"
)

type ApiData struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func GenerateApiData(msg string, code int, data interface{}) ApiData {
	return ApiData{
		Message: msg,
		Code:    code,
		Data:    data,
	}
}
