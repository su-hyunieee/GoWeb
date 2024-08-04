package main

import (
	"net/http"

	"web1/myapp"
)

func main() {
		http.ListenAndServe(":3030", myapp.NewHttpHandler())
}

