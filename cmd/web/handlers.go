package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"danschmid.de/snippetbox/pkg/models"
)

// Define a home handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    // Check if the current request URL path exactly matches "/". If it doesn't, use
    // the http.NotFound() function to send a 404 response to the client.
    // Importantly, we then return form the handler.
    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }

    s, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }

    // Use the new render helper.
    app.render(w, r, "home.page.tmpl", &templateData{
        Snippets: s,
    })
}

// Add a showSnippet handler function
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    // Extract the value of the id parameter from the query string and try to
    // convert it to an integer using the strconv.Atoi() function. If it can't
    // be converted to an integer, or the value is less than 1, we return a 404 page
    // not found response.
    id, err := strconv.Atoi(r.URL.Query().Get("id"))

    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    // Use the SnippetModel object's Get method to retrieve the data for a
    // specific record based on its ID. If no matching record is found,
    // return a 404 Not Found response.
    s, err := app.snippets.Get(id)

    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }

        return
    }

    // Use the new render helper.
    app.render(w, r, "show.page.tmpl", &templateData{
        Snippet: s,
    })
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

        app.clientError(w, http.StatusMethodNotAllowed)
        return
    }

    // Create some variables holding dummy data. We'll remove these later on
    // during the build.
    title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := "7"

    // Pass the data to the SnippetModel.Insert() method, receiving the
    // ID of the new record back.
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, err)
    }

    // Redirect the user to the relevant page for the snippet.
    http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
