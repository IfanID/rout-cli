package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

// CommandHelp mendefinisikan struktur untuk menampilkan bantuan perintah.
type CommandHelp struct {
	Name        string
	Description string
	Flags       []string
}

var HelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Menampilkan daftar perintah dan bantuan ini.",
	Long:  `Menampilkan daftar semua perintah khusus yang tersedia beserta dokumentasinya.`,
	Run:   runHelpCommand,
}

func runHelpCommand(cmd *cobra.Command, args []string) {
	// Daftar semua perintah khusus yang tersedia beserta dokumentasinya.
	commands := []CommandHelp{
		{
			Name:        "sub",
			Description: "Pindah cepat ke direktori proyek yang sering digunakan.",
			Flags: []string{
				"-lokasi : Menampilkan lokasi direktori yang tersimpan.",
				"-ganti  : Memperbarui atau mengganti lokasi direktori.",
			},
		},
		{
			Name:        "help",
			Description: "Menampilkan daftar perintah dan bantuan ini.",
			Flags:       []string{},
		},
		// Tambahkan perintah baru di sini di masa depan.
	}

	fmt.Println("\nPerintah Khusus ROut")
	fmt.Println("---------------------")
	fmt.Println("Perintah-perintah ini adalah skrip mandiri yang dijalankan melalui ROut.")

	for _, cmd := range commands {
		fmt.Printf("\n  ■ %s\n", cmd.Name)
		fmt.Printf("    %s\n", cmd.Description)
		if len(cmd.Flags) > 0 {
			fmt.Println("    Opsi:")
			for _, flag := range cmd.Flags {
				fmt.Printf("      %s\n", flag)
			}
		}
	}
	fmt.Println()
}
