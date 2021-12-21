package fluent

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

func OptionPresent[T any](value T) Option[T] {
	return present[T]{value}
}

func OptionEmpty[T any]() Option[T] {
	return empty[T]{}
}

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
