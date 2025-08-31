package core

import (
	"fmt"
	"runtime"
	"rout/core/system/util"
)

// MOTD mencetak Message Of The Day (MOTD) aplikasi ROut.
func MOTD() {
	fmt.Println(`
┌──────────────────────────────────────────┐
│                                          │
│   ██████╗  ██████╗ ██╗   ██╗ ████████╗   │
│   ██╔══██╗██╔═══██╗██║   ██║╚══██╔══╝   │
│   ██████╔╝██║   ██║██║   ██║   ██║      │
│   ██╔══██╗██║   ██║██║   ██║   ██║      │
│   ██║  ██║╚██████╔╝╚██████╔╝   ██║      │
│   ╚═╝  ╚═╝ ╚═════╝  ╚═════╝    ╚═╝      │
│                                          │
└──────────────────────────────────────────┘
`)
	util.TypeOut(fmt.Sprintf("  Version: %s | OS: %s | Arch: %s", "v0.0.3", runtime.GOOS, runtime.GOARCH))
	util.TypeOut("  Ketik 'help' untuk bantuan, 'exit' untuk keluar.")
	fmt.Println()
}

// Logout mencetak pesan keluar dari aplikasi ROut.
func Logout() {
	util.TypeOut("Terima kasih telah menggunakan ROut!")
}