package util

func PointerOf[A any](a A) *A {
	return &a
}
