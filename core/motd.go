package core

import (
	"fmt"
	"runtime"
)

// PrintMOTD mencetak Message Of The Day (MOTD) aplikasi ROut.
func MOTD() {
	fmt.Println(`
  _____   ____        _   
 |  __ \ / __ \      | |  
 | |__) | |  | |_   _| |_ 
 |  _  /| |  | | | | | __|
 | | \ \| |__| | |_| | |_ 
 |_|  \_\____/ \__,_|\__| 
`)
	fmt.Printf("Version: %s\n", "v0.0.1")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Arch: %s\n", runtime.GOARCH)
	fmt.Println("Ketik 'exit' atau 'quit' untuk keluar.")
}

// Logout mencetak pesan keluar dari aplikasi ROut.
func Logout() {
	fmt.Println("Keluar dari ROut.")
}
