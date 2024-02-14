package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf("Hello, world! %v", time.Now().Format("2006-01-02 15:04:05")))
}

func main() {
	http.HandleFunc("/", IndexHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
