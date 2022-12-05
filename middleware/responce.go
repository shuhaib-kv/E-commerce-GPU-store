package middleware

type Responses struct {
	StatusCode bool     `json:"statusCode"`
	Headers    []string `json:"headers"`
	Body       string   `json:"body"`
}

func Truersponce(s string, p []string) Responses {
	t := Responses{
		StatusCode: true,
		Headers:    p,
		Body:       s,
	}
	return t
}
func Falseresponce(s string, p []string) Responses {
	t := Responses{
		StatusCode: false,
		Headers:    p,
		Body:       s,
	}
	return t
}
