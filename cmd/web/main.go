package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
    // Define a new command-line flag with the name 'addr', a default value of ":4000"
    // and some short help text explaining what the flag controls. The value of the
    // flag will be stored in the addr variable at runtime.
    addr := flag.String("addr", ":4000", "HTTP network address")

    // Importantly, we use the flag.Parse() function to parse the command-line flag.
    // This reads in the command-line flag value and assigns it to the addr
    // variable. You need to call this *before* you use the addr variable
    // otherwise it will always contain the default value of ":4000". If any errors are
    // encountered during parsing the application will be terminated.
    flag.Parse()

    // Use log.New() to create a logger for writing information messages. This takes
    // three parameters: the destination to write the logs to (os.Stdout), a string
    // prefix for message (INFO followed by a tab), and flags to indicate what
    // additional information to include (local date and time). Note that the flags
    // are joined using the bitwise OR operator |.
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

    // Create a logger for writing error messages in the same way, but use stderr as
    // the destination and use the log.Lshortfile flag to include the relevant
    // file name and line number.
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

    // Use the http.NewServeMux() function to initialize a new servemux, then
    // register the home function as the handler for the "/" URL pattern.
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)

    // Create a file server which serves files out of the "./ui/static" directory.
    // Note that the path given to the http.Dir function is relative to the project
    // directory root.
    fileServer := http.FileServer(http.Dir("./ui/static"))

    // Use the mux.Handle() function to register the file server as the handler for
    // all URL paths that start with "/static/". For matching paths, we strip the
    // "/static" prefix before the request reaches the file server.
    mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))

    // Initialize a new http.Server struct. We set the Addr and Handler fields so
    // that the server uses the same network address and routes as before, and set
    // the ErrorLog field so that the server now uses the custom errorLog logger in
    // the event of any problems.
    srv := &http.Server{
        Addr: *addr,
        ErrorLog: errorLog,
        Handler: mux,
    }

    // Use the http.ListenAndServe() function to start a new web server. We pass in
    // two parameters: the TCP network address to listen on (in this case ":4000")
    // and the servemux we just created. If http.ListenAndServe() returns an error
    // we use the log.Fatal() function to log the error message and exit. Note that
    // any error returned by http.ListenAndServe() is always non-nil.
    // The value returned from the flag.String() function is a pointer to the flag
    // value, not the value itself. So we need to dereference the pointer (i.e.
    // prefix it with the * symbol) before using it. Note that we're using the
    // log.Printf() function to interpolate the address with the log message.
    infoLog.Printf("Starting server on %s", *addr)
    err := srv.ListenAndServe()
    errorLog.Fatal(err)
}

// A middleware to check if the request URL ends with a / character, and if it does,
// return a 404 Not Found response instead of passing on the request to the http.FileServer.
func neuter(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasSuffix(r.URL.Path, "/") {
            http.NotFound(w, r)
            return
        }

        next.ServeHTTP(w, r)
    })
}

