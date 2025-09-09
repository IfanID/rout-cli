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
			Name:        "rman",
			Description: "Mengatur atau pindah ke direktori 'rman'.",
			Flags: []string{
				"-ganti  : Paksa untuk mengatur ulang lokasi direktori.",
				"-lokasi : Tampilkan lokasi direktori yang saat ini disimpan.",
			},
		},
	}

	convCmd := CommandHelp{
		Name:        "conv",
		Description: "Konversi file .ts ke .vtt dengan CUE cerdas.",
		Flags: []string{
			"<namafile.ts> : Mengonversi satu file.",
			"all           : Mengonversi semua file .ts di direktori.",
		},
	}

	helpCmd := CommandHelp{
		Name:        "help",
		Description: "Menampilkan daftar perintah dan bantuan ini.",
		Flags:       []string{},
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

		if cmd.Name == "sub" {
			fmt.Printf("\n    └── ■ %s\n", convCmd.Name)
			fmt.Printf("        %s\n", convCmd.Description)
			if len(convCmd.Flags) > 0 {
				fmt.Println("        Opsi:")
				for _, flag := range convCmd.Flags {
					fmt.Printf("          %s\n", flag)
				}
			}
		}
	}

	fmt.Println("\n---------------------")
	fmt.Printf("\n  ■ %s\n", helpCmd.Name)
	fmt.Printf("    %s\n", helpCmd.Description)

	fmt.Println()
}
