package fn

// 类似于三元函数
func If[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}
