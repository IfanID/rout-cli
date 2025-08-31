/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	
	"os"
	"strings"

	"github.com/spf13/cobra"

	"rout/core"
	"rout/core/system/command" // Added for ls, cd, pwd functions
)

// rootCmd merepresentasikan perintah dasar ketika dipanggil tanpa subperintah
var rootCmd = &cobra.Command{
	Use:   "rout",
	Short: "Deskripsi singkat aplikasi Anda",
	Long: `Deskripsi yang lebih panjang yang mencakup beberapa baris dan kemungkinan berisi
contoh dan penggunaan aplikasi Anda. Contoh:

Cobra adalah pustaka CLI untuk Go yang memberdayakan aplikasi.
Aplikasi ini adalah alat untuk menghasilkan file yang dibutuhkan
untuk membuat aplikasi Cobra dengan cepat.`,
	// Hapus komentar baris berikut jika aplikasi dasar Anda
	// memiliki tindakan yang terkait dengannya:
	Run: func(cmd *cobra.Command, args []string) {
			core.MOTD() // Cetak MOTD saat aplikasi dimulai
			reader := bufio.NewReader(os.Stdin)
			for {
				core.Prompt() // Panggil fungsi prompt dari paket core

				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				if input == "exit" || input == "quit" {
					core.Logout()
					break
				}

				command.RegisterCommands(input) // Call the command registration function

			}
		},
}

// Execute menambahkan semua subperintah ke perintah root dan mengatur flag dengan tepat.
// Ini dipanggil oleh main.main(). Ini hanya perlu terjadi sekali pada rootCmd.
func Execute(initialCwd string) {
	os.Chdir(initialCwd) // Set working directory to where rout was launched
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Di sini Anda akan mendefinisikan flag dan pengaturan konfigurasi Anda.
	// Cobra mendukung flag persisten, yang, jika didefinisikan di sini,
	// akan bersifat global untuk aplikasi Anda.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "file konfigurasi (default adalah $HOME/.rout.yaml)")

	// Cobra juga mendukung flag lokal, yang hanya akan berjalan
	// ketika tindakan ini dipanggil secara langsung.
	rootCmd.Flags().BoolP("toggle", "t", false, "Pesan bantuan untuk toggle")
}
