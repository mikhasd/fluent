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
	IsPresent() bool
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
	// Executes the provided function if a value is present
	IfPresent(func(T))
	// If a value is present and the value is accepted by the provided filter,
	// returns an Option representing the value, otherwise will return an empty
	// Option
	Filter(func(T) bool) Option[T]
	String() string
}

// Present returns a new Option wrapping the T value.
func Present[T any](value T) Option[T] {
	return present[T]{value}
}

// Empty returns an empty Option.
func Empty[T any]() Option[T] {
	return empty[T]{}
}

// OfNillable returns an Option describing the provided reference.
//
// If the provided reference of T is non-nil, a Option representing the value
// is returned.
//
// If the value is nil, an empty option is returned.
func OfNillable[T any](ref *T) Option[*T] {
	if ref == nil {
		return Empty[*T]()
	} else {
		return Present(ref)
	}
}

// MapOption executes the mapper function over the Option value if it is not
// empty.
//
// If the Option is empty, the mapper is not executed and an empty Option is
// returned.
func MapOption[T any, R any](o Option[T], mapper func(T) R) Option[R] {
	if o.IsPresent() {
		return Present(mapper(o.Get()))
	} else {
		return Empty[R]()
	}
}
