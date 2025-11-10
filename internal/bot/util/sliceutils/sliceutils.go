package sliceutils

func First[T any](slice []T) *T {
	if len(slice) == 0 {
		return nil
	}

	return &slice[0]
}

func Last[T any](slice []T) *T {
	if len(slice) == 0 {
		return nil
	}

	return &slice[len(slice)-1]
}
