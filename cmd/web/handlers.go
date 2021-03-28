package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Define a home handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    // Check if the current request URL path exactly matches "/". If it doesn't, use
    // the http.NotFound() function to send a 404 response to the client.
    // Importantly, we then return form the handler.
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    // Initialize a slice containing the paths to the two files. Note that the
    // home.page.tmpl file must be the *first* file in the slice.
    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }

    // Use the template.ParseFiles() function to read the template file into a
    // template set. Notice that we can pass the slice of file paths
    // as a variadic parameter? If there's an error, we log the detailed error message and use
    // the http.Error() function to send a generic 500 Internal Server Error
    // response to the user.
    ts, err := template.ParseFiles(files...)

    if err != nil {
        app.errorLog.Println(err.Error())
        http.Error(w, "Internal Server error", 500)
        return
    }

    // We then use the Execute() method on the template set to write the template
    // content as the response body. The last parameter to Execute() represents any
    // dynamic data that we want to pass in, which for now we'll leave as nil.
    err = ts.Execute(w, nil)

    if err != nil {
        app.errorLog.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
    }
}

// Add a showSnippet handler function
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
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
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
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
