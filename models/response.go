package models

type Response struct {
	Code    int
	Message string
	Headers map[string][]string
	Body    string
}
