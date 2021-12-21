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
	Present() bool
	Get() T
	Map(func(T) T) Option[T]
	OrElse(T) T
	OrElseGet(func() T) T
	Or(func() Option[T]) Option[T]
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

func MapOption[T any, R any](o Option[T], mapper func(T) R) Option[R] {
	if o.Present() {
		return OptionPresent(mapper(o.Get()))
	} else {
		return OptionEmpty[R]()
	}
}
