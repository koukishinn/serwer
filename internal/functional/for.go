package functional

func ForEach[T any](l []T, f func(e T)) {
	for _, e := range l {
		f(e)
	}
}
