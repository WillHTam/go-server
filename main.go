package main

import "net/http"

func main() {
	// install a handler function at the root path of our webserver.
	// http.HandleFunc operates on the default HTTP router, officially called a ServeMux.
	http.HandleFunc("/", hello)
	http.ListenAndServe(":3000", nil)
}

// Every time a new request comes into the HTTP server matching the root path,
// the server will spawn a new goroutine executing the hello function.
// And the hello function simply uses the http.ResponseWriter to write a response to the client.
// Since http.ResponseWriter.Write takes the more general []byte, or byte-slice, as a parameter,
// we do a simple type conversion of our string.
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hellou!"))
}
