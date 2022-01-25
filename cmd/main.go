package main

import (
	"fmt"
	"net/http"
	"shorturl"
)

func main() {
	fmt.Println("Start Server...")
	http.HandleFunc("/", shorturl.Redirect)
	http.HandleFunc("/add", shorturl.Add)
	http.ListenAndServe(":8080", nil)
}
