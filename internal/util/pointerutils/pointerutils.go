package pointerutils

func AsPointer[T any](value T) *T {
	return &value
}
