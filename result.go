package fluent

// Result represents the output of a function which may have been computed
// successfully (Ok) or failed with an error (Err).
//
// This implementation is based on Rust's std::result.
type Result[T any] interface {
	// Returns true if the Result is Ok
	IsOk() bool
	// Returns true if the Result is Err
	IsErr() bool
	// Converts the Result into an Option.
	//
	// If the Result is Ok, an Option will be returned, wrapping the result
	// value.
	//
	// If the Result is Err, an empty Option will be returned.
	Ok() Option[T]
	// Converts the Result into an Option.
	//
	// If the Result is Err, an Option will be returned, wrapping the result
	// error.
	//
	// If the Result is Ok, an empty Option will be returned.
	Err() Option[error]
	// Executes the mapper function over the Result value.
	//
	// If the Result is Err, the function is not executed.
	Map(mapper func(T) T) Result[T]

	// Executes the mapper function over the Result error.
	//
	// If the Result is Ok, the function is not executed.
	MapErr(func(error) T) Result[T]

	// Gets the Result value.
	//
	// This function will panic if the Result is Err.
	Get() T

	// Gets the Result error.
	//
	// This function will panic if the Result is Ok.
	GetErr() error

	// Gets the Result value or	the provided value if if the Result is Err.
	OrElse(value T) T

	// Calls the provided function if the Result is Err. Returns Result value
	// if Result is Ok.
	OrElseGet(func() T) T

	// Calls the provided function if the Result is Err. Returns the current
	// Result if Ok.
	Or(func() Result[T]) Result[T]
	String() string
}

// Ok returns a result representing a successful operation resulting in
// the value T.
func Ok[T any](value T) Result[T] {
	return ok[T]{value}
}

// Err returns a result representing a failed operation which resulted in
// on an error.
func Err[T any](e error) Result[T] {
	return err[T]{e}
}

// MapResult executes the mapper function over the Result value if it Ok.
//
// If the Result is Err, the mapper is not executed and an Err Result is
// returned.
func MapResult[T any, R any](r Result[T], mapper func(T) R) Result[R] {
	if r.IsOk() {
		return Ok(mapper(r.Get()))
	} else {
		return Err[R](r.GetErr())
	}
}

// CallResult executes the provided function and returns an Ok result if
// the result error is nil.
//
// If the provided function returns an error, an Err result is returned, even
// if a result is returned.
func CallResult[T any](fn func() (T, error)) Result[T] {
	result, err := fn()
	if err == nil {
		return Ok(result)
	} else {
		return Err[T](err)
	}
}
