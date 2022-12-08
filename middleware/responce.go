package middleware

type Responses struct {
	StatusCode bool `json:"statusCode"`

	Body string `json:"body"`
}

func Truersponce(s string) Responses {
	t := Responses{
		StatusCode: true,

		Body: s,
	}
	return t
}
func Falseresponce(s string) Responses {
	t := Responses{
		StatusCode: false,
		Body:       s,
	}
	return t
}
