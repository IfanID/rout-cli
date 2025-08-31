package util

import (
	"fmt"
	"time"
)

// TypeOut mencetak teks dengan efek mengetik dan diakhiri baris baru.
func TypeOut(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(10 * time.Millisecond) // Jeda 10 milidetik antar karakter
	}
	fmt.Println()
}

// TypeOutInline mencetak teks dengan efek mengetik tanpa baris baru di akhir.
func TypeOutInline(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(10 * time.Millisecond)
	}
}
