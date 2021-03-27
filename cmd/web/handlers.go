package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
    // Check if the current request URL path exactly matches "/". If it doesn't, use
    // the http.NotFound() function to send a 404 response to the client.
    // Importantly, we then return form the handler. If we don't return the handler
    // would keep executing and also write the "Hello from Snippetbox" message.
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    w.Write([]byte("Hello from Snippetbox"))
}

// Add a showSnippet handler function
func showSnippet(w http.ResponseWriter, r *http.Request) {
    // Extract the value of the id parameter from the query string and try to
    // convert it to an integer using the strconv.Atoi() function. If it can't
    // be converted to an integer, or the value is less than 1, we return a 404 page
    // not found response.
    id, err := strconv.Atoi(r.URL.Query().Get("id"))

    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    // Use the fmt.Fprintf() function to interpolate the id value with our response
    // and write it to the http.ResponseWriter.
    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a createSnippet handler function
func createSnippet(w http.ResponseWriter, r *http.Request) {
    // Use r.Method to check whether te request is using POST or not. Note that
    // http.MehodPost is a constant equal  to the string "POST".
    if r.Method != http.MethodPost {
        // Use the w.Header().Set() method to add an "Allow: POST" header to the
        // response header map. The first parameter is the header name, and the second 
        // paramter is the header value.
        w.Header().Set("Allow", http.MethodPost)

        // use the http.Error() function to send a 405 status code and
        // "Method Not Allowed" string as the responaw body.
        http.Error(w, "Method Not Allowed", 405)
        return
    }

    w.Write([]byte("Create a new snippet..."))
}
