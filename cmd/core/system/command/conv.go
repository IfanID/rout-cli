package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"rout/cmd/core/system/command/conv"
	"rout/cmd/core/system/util"
)

var ConvCmd = &cobra.Command{
	Use:   "conv [file.ts|all]",
	Short: "Konversi file .ts ke .vtt dengan CUE cerdas.",
	Long: `Alat ini mengonversi file video .ts menjadi file subtitle .vtt.\nFitur utamanya adalah penambahan atau pembaruan blok CUE (logo) secara cerdas\ndengan menganalisis konten subtitle menggunakan Gemini API untuk penempatan waktu yang optimal.\n\nPenggunaan:\n  rout conv               // Menampilkan informasi penggunaan dan lokasi 'sub'.\n  rout conv <namafile.ts>  // Mengonversi satu file.\n  rout conv all             // Mengonversi semua file .ts di direktori saat ini.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if conv.CheckSubDir() {
				configPath, _ := util.GetConfigPath()
				config, err := util.ReadConfig(configPath)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error membaca konfigurasi:", err)
					return
				}
			lokasiCdSub, exists := config["lokasiCdSub"]
			if exists && lokasiCdSub != "" {
					fmt.Printf("Perintah 'conv' siap digunakan di direktori: %s\n", lokasiCdSub)
					fmt.Println("Gunakan 'conv all' untuk mengonversi semua file .ts, atau 'conv <namafile.ts>' untuk satu file.")
				}
			}
		} else {
			conv.HandleConversion(args[0])
		}
	},
}
