package rand

import (
	"math/rand"
)

// letterBytes — строка, содержащая все допустимые символы для генерации случайных строк.
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytes генерирует случайную строку длиной n символов, состоящую из латинских букв (строчных и заглавных).
// Внутри строки используются только символы из константы letterBytes.
//
// Параметры:
//   - n: длина генерируемой строки.
//
// Возвращаемое значение:
//   - строка длиной n, содержащая случайно выбранные символы из letterBytes.
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
