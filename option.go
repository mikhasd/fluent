package fluent

// Option represents an optional value.
//
// Option can be used as alternative for:
//   - Uninitialized values
//   - Invalid/empty return values
//   - os.IsNotExist(e error)
//
// This implementation is based on Java's java.util.Optional and Rust's
// std::option.
type Option[T any] interface {
	// Returns true if a value is present.
	Present() bool
	// Gets the Option value or panics if empty
	Get() T
	// Applies the provided mapper function over the Option value, if present.
	Map(mapper func(T) T) Option[T]
	// Returns the Option value if present, or the provided value it empty.
	OrElse(other T) T
	// Returns the Option value if present, or invoke the provided function if
	// empty.
	OrElseGet(func() T) T
	// If empty, calls the provided function and returns it result or returns
	// the current option if present.
	Or(func() Option[T]) Option[T]
	// Returns the Option value wrapped into a result if present. If empty,
	// returns a Result with a wrapped error.
	OrError(e error) Result[T]
	String() string
}

// OptionPresent returns a new Option wrapping the T value.
func OptionPresent[T any](value T) Option[T] {
	return present[T]{value}
}

// OptionEmpty returns a new empty Option.
func OptionEmpty[T any]() Option[T] {
	return empty[T]{}
}

// OptionFromReference return an Option describing the provided reference.
//
// If the provided reference of T is non-nil, a Option representing the value
// is returned.
//
// If the value is nil, an empty option is returned.
func OptionFromReference[T any](ref *T) Option[T] {
	if ref == nil {
		return OptionEmpty[T]()
	} else {
		return OptionPresent(*ref)
	}
}

// MapOption executes the mapper function over the Option value if it is not
// empty.
//
// If the Option is empty, the mapper is not executed and an empty Option is
// returned.
func MapOption[T any, R any](o Option[T], mapper func(T) R) Option[R] {
	if o.Present() {
		return OptionPresent(mapper(o.Get()))
	} else {
		return OptionEmpty[R]()
	}
}
