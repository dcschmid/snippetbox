package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"danschmid.de/snippetbox/pkg/forms"
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
    app.render(w, r, "create.page.tmpl", &templateData{
        Form: forms.New(nil),
    })
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

    // Create a new forms.Form struct containing the POSTed data from the
    // form, then use the validation methods to check the content.
    form := forms.New(r.PostForm)
    form.Required("title", "content", "expires")
    form.MaxLength("title", 100)
    form.PermittedValues("expires", "365", "7", "1")

    // If the form isn't valid, redisplay the template passing in the
    // form.Form object as the data.
    if !form.Valid() {
        app.render(w, r, "create.page.tmpl", &templateData{Form: form})
        return
    }

    // Because the form data (with type url.Values) has been anonymously embedded
    // in the form.Form struct, we can use the Get() method to retrieve
    // the validated value for a particular form field.
    id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
    if err != nil {
        app.serverError(w, err)
    }

    // Redirect the user to the relevant page for the snippet.
    http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
} 
