package api

import "net/http"

func CreateWorkOrderHandler(w http.ResponseWriter, r *http.Request) {
	msg := "hello"
	_, err := w.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
}
