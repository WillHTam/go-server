package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

//The type keyword defines a new type, which we call weatherData, and declare as a struct.
// Each field in the struct has a name (e.g. Name, Main), a type (string, another anonymous struct), and what’s known as a tag.
//Tags are like metadata, and allow us to use the encoding/json package to directly unmarshal the API’s response into our struct.
type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json: "main"`
}

func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=<~APPID HERE~>&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	// If no errors, defer call to close response body
	// "https://golang.org/doc/effective_go.html#defer"
	defer resp.Body.Close()

	// allocate weatherData struct
	var d weatherData

	// use a json.Decoder to unmarshal from the response body directly into our struct
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	// If the decode succeeds, we return the weatherData to the caller, with a nil error to indicate success.
	return d, nil
}

func main() {
	// install a handler function at the root path of our webserver.
	// http.HandleFunc operates on the default HTTP router, officially called a ServeMux.
	http.HandleFunc("/", hello)

	// Request Handler for /weather route
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		// takes the string after /weather/ and assigns it to city
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		// Make query with city
		// If error, report with http.Error helper function
		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// if error, return out of function to complete HTTP request
			return
		}

		// Otherwise, tell Client that JSON is being sent
		// use json.NewEncoder to encode weatherData as JSON directly
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

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
