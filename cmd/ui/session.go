package ui

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

func StartShellSession() {
	// Siapkan lingkungan Zsh kustom
	customZdotdir, err := setupZshEnvironment()
	if err != nil {
		panic(err)
	}

	// Dapatkan path executable saat ini untuk portabilitas
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	c := exec.Command("zsh")
	c.Env = os.Environ()
	c.Env = append(c.Env, "ROUT_SESSION=1")
	c.Env = append(c.Env, "ZDOTDIR="+customZdotdir)
	c.Env = append(c.Env, "ROUT_EXECUTABLE_PATH="+exePath)

	ptmx, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	// Atur agar ukuran window terminal kita disinkronkan dengan pty.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			_ = pty.InheritSize(os.Stdin, ptmx)
		}
	}()
	ch <- syscall.SIGWINCH

	// Atur terminal ke mode "raw" menggunakan pustaka 'term'.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	// Salin output dari zsh (pty) ke output standar (layar pengguna).
	go func() { _, _ = io.Copy(os.Stdout, ptmx) }()

	// Salin input dari pengguna (stdin) ke zsh (pty).
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()

	// Tunggu proses zsh selesai.
	_ = c.Wait()

	// Penting: Tutup ptmx setelah perintah selesai untuk menghentikan
	// goroutine io.Copy yang membaca dari stdin.
	_ = ptmx.Close()
}
