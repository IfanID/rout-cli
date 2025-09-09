/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"rout/cmd/core"
	"rout/cmd/ui"
	"rout/cmd/core/system/command"
	"rout/cmd/rcli"
)

var version = "v0.0.11"

// rootCmd merepresentasikan perintah dasar ketika dipanggil tanpa subperintah
var rootCmd = &cobra.Command{
	Use:   "rout",
	Short: "Deskripsi singkat aplikasi Anda",
	Long: `Deskripsi yang lebih panjang yang mencakup beberapa baris dan kemungkinan berisi
contoh dan penggunaan aplikasi Anda. Contoh:

Cobra adalah pustaka CLI untuk Go yang memberdayakan aplikasi.
Aplikasi ini adalah alat untuk menghasilkan file yang dibutuhkan
untuk membuat aplikasi Cobra dengan cepat.`,
	Run: func(cmd *cobra.Command, args []string) {
		core.MOTD(version) // Cetak MOTD saat aplikasi dimulai
		ui.StartShellSession()
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
	rootCmd.AddCommand(rcli.RcliCmd)
	rootCmd.AddCommand(command.ConvCmd)
	rootCmd.AddCommand(command.SubCmd)
	rootCmd.AddCommand(command.RmanCmd)
	rootCmd.AddCommand(command.HelpCmd)

	command.SubCmd.Flags().BoolVarP(&command.ForceChange, "ganti", "g", false, "Paksa untuk mengubah lokasi 'sub'")
	command.SubCmd.Flags().BoolVarP(&command.ShowLocation, "lokasi", "l", false, "Tampilkan lokasi 'sub' saat ini")

	command.RmanCmd.Flags().BoolVarP(&command.ForceChange, "ganti", "g", false, "Paksa untuk mengubah lokasi 'rman'")
	command.RmanCmd.Flags().BoolVarP(&command.ShowLocation, "lokasi", "l", false, "Tampilkan lokasi 'rman' saat ini")

	rootCmd.Flags().BoolP("toggle", "t", false, "Pesan bantuan untuk toggle")
}