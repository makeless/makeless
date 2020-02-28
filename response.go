package main

type response struct {
	Data   string `json:"data"`
	Error  string `json:"error"`
	Base64 bool   `json:"base64"`
}
