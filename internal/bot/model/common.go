package model

type TelegramID = uint64

type Range[T any] struct {
	Start T
	End   T
}
