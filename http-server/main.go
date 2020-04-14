package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func hello(writer http.ResponseWriter, request *http.Request) {
	context := request.Context()
	fmt.Println("server: hello handler started")
	defer fmt.Println("server: hello handler ended")

	select {
	case <-context.Done():
		error := context.Err()
		fmt.Println("server:", error)
		internalError := http.StatusInternalServerError
		http.Error(writer, error.Error(), internalError)
	default:
		fmt.Fprintf(writer, "hello\n")
	}
}

func google(writer http.ResponseWriter, request *http.Request) {
	response, error := http.Get("http://google.com")
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	result, _ := ioutil.ReadAll(response.Body)

	fmt.Fprintf(writer, string(result))

	// scanner := bufio.NewScanner(response.Body)

	// for scanner.Scan() {

	// }
}

func headers(writer http.ResponseWriter, request *http.Request) {
	for name, headers := range request.Header {
		for _, header := range headers {
			fmt.Fprintf(writer, "%s: %s\n", name, header)
		}
	}
}

func main() {
	fmt.Println("starting...")

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/google", google)

	http.ListenAndServe(":8090", nil)
}
