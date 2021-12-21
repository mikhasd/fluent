package fluent

type Result[T any] interface {
	IsOk() bool
	IsErr() bool
	Ok() Option[T]
	Err() Option[error]
	Map(func(T) T) Result[T]
	MapErr(func(error) T) Result[T]
	Get() T
	GetErr() error
	OrElse(T) T
	OrElseGet(func() T) T
	String() string
}

func ResultOk[T any](value T) Result[T] {
	return ok[T]{value}
}

func ResultErr[T any](e error) Result[T] {
	return err[T]{e}
}

func ResultMap[T any, R any](r Result[T], mapper func(T) R) Result[R] {
	if r.IsOk() {
		return ResultOk(mapper(r.Get()))
	} else {
		return ResultErr[R](r.GetErr())
	}
}

func ResultFromCall[T any](fn func() (T, error)) Result[T] {
	result, err := fn()
	if err == nil {
		return ResultOk(result)
	} else {
		return ResultErr[T](err)
	}
}
