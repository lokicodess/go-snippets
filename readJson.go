func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Decode the request body into the target destination.
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {

		// Prepare error variables for comparison.
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		// JSON is badly formed → syntax error with offset.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// Unexpected EOF (common JSON syntax issue).
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Wrong type for a field → inform the client.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// Body is empty → no JSON sent.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Developer mistake: passing something that isn't a valid pointer.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// Default: return the basic error.
		default:
			return err
		}
	}

	// No error
	return nil
}
