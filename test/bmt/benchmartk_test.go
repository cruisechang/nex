package bmt

import (
	"testing"
	"time"
)

func BenchmarkLoops(b *testing.B) {

	customTimerTag := false

	if customTimerTag {
		b.StopTimer()
	}
	b.SetBytes(123456)
	time.Sleep(time.Second)
	if customTimerTag {
		b.StartTimer()
	}
	for i := 0; i < b.N; i++ {

	}

}

func BenchmarkStrConcat3(b *testing.B) {
	str1 := "123456789012345678901234567890"
	str2 := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var str0 = ""
	for idx := 0; idx < b.N; idx++ {
		Strconcat13(str1, str2, &str0)
	}
}

func Strconcat13(str1, str2 string, str0 *string) {
	*str0 = str1 + str2
}
