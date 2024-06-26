package common

const (
	SUCCESS = "success"
	FAIL    = "fail"
)

type ResponseData[T any] struct {
	Status       string `json:"status"`
	Data         T      `json:"data,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func Success[T any](data T) ResponseData[T] {
	return ResponseData[T]{Status: SUCCESS, Data: data}
}

func Fail(errorMessage string) ResponseData[string] {
	return ResponseData[string]{Status: FAIL, ErrorMessage: errorMessage}
}
