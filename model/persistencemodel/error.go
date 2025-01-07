package persistencemodel

// PersistenceError is returned when error in persistence has occurred.
// It uses http error codes but the `Message` is a human readable message.
type PersistenceError struct {
	// Code is a _HTTP_ error code.
	Code int
	// Custom is a custom code to be more precise.
	Custom int
	// Message is a custom message that is human readable.
	Message string
}

func (e PersistenceError) Error() string {
	return e.Message
}

func Error400(message string, custom ...int) PersistenceError {
	if len(custom) > 0 {
		return PersistenceError{Code: 400, Custom: custom[0], Message: message}
	}
	return PersistenceError{Code: 400, Message: message}
}

func Error404(message string, custom ...int) PersistenceError {
	if len(custom) > 0 {
		return PersistenceError{Code: 404, Custom: custom[0], Message: message}
	}
	return PersistenceError{Code: 404, Message: message}
}

func Error409(message string, custom ...int) PersistenceError {
	if len(custom) > 0 {
		return PersistenceError{Code: 409, Custom: custom[0], Message: message}
	}
	return PersistenceError{Code: 409, Message: message}
}

func Error500(message string, custom ...int) PersistenceError {

	if len(custom) > 0 {
		return PersistenceError{Code: 500, Custom: custom[0], Message: message}
	}
	return PersistenceError{Code: 500, Message: message}
}
