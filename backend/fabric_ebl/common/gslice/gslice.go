package gslice

func ForEach[T any](s []T, f func(T)) {
	for _, v := range s {
		f(v)
	}
}

func Map[T, U any](s []T, f func(T) U) (r []U) {
	r = make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return
}

func Filter[T any](s []T, f func(T) bool) (r []T) {
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return
}

func ToMap[T any, K comparable, V any](s []T, f func(T) (K, V)) (r map[K]V) {
	r = make(map[K]V)
	for _, v := range s {
		k, v := f(v)
		r[k] = v
	}
	return
}

func GroupBy[T any, K comparable](s []T, f func(T) K) (r map[K][]T) {
	r = make(map[K][]T)
	for _, v := range s {
		k := f(v)
		r[k] = append(r[k], v)
	}
	return
}
