package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"danschmid.de/snippetbox/pkg/models"
)

// Define a home handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
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
    id, err := strconv.Atoi(r.URL.Query().Get(":id"))

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

// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
    app.render(w, r, "create.page.tmpl", nil)
}

// Add a createSnippet handler function
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    // First we call r.ParseForm() which adds any data in POST request bodies
    // to the r.PostForm map. This also works in the same way for PUT and PATCH
    // requests. If there are any errors, we use our app.ClientError helper to send
    // a 400 Bad Request response to the user.
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    // Use the r.PostForm.Get() method to retrieve the relevant data fields
    // from the r.PostForm map.
    title := r.PostForm.Get("title")
    content := r.PostForm.Get("content")
    expires := r.PostForm.Get("expires")

    // Initialize a map to hold any validation errors.
    errors := make(map[string]string)

    // Check that the title field is not blank and is not more than 100 characters
    // long. If it fails either of those checks, add a message to the errors
    // map using the field name as the key.
    if strings.TrimSpace(title) == "" {
        errors["title"] = "This field cannot be blank"
    } else if utf8.RuneCountInString(title) > 100 {
        errors["title"] = "This field is too long (maximum is 100 characters)"
    }

    // Check that the Content field isn't blank.
    if strings.TrimSpace(content) == "" {
        errors["content"] = "This field cannot be blank"
    }

    // Check the expires field isn't blank and matches one of the permitted
    // values ("1", "7" or "365").
    if strings.TrimSpace(expires) == "" {
        errors["expires"] = "This field cannot be blank"
    } else if expires != "365" && expires != "7" && expires != "1" {
        errors["expires"] = "This field is invalid"
    }

    // If there are any validation errors, re-display the create.page.tmpl
    // template passing in the validation errors and previously submitted
    // r.PostForm data.
    if len(errors) > 0 {
        app.render(w, r, "create.page.tmpl", &templateData{
            FormErrors: errors,
            FormData: r.PostForm,
        })
        return
    }

    // Create a new snippet record in the database using the form data.
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, err)
    }

    // Redirect the user to the relevant page for the snippet.
    http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
} 
