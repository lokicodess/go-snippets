// logError logs an error message along with the request method and URI.
func (app *application) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// errorResponse sends a JSON-formatted error message to the client
// with the provided status code.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	// Attempt to write the JSON response.
	// If this fails, log the error and send a generic 500 status.
	if err := app.writeJSON(w, status, env, nil); err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse logs the detailed error and sends a generic 500 error to the client.
// This is used when the application encounters an unexpected runtime problem.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// Log the detailed error.
	app.logError(r, err)

	// Generic message for the client (don’t expose internal errors).
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// badRequestResponse sends a 400 Bad Request error with the given message.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// notFoundResponse sends a 404 Not Found error.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse sends a 405 Method Not Allowed error.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
