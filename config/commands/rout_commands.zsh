# File ini menangani semua perintah Go biasa (non-cd).

# Dapatkan root direktori proyek secara dinamis berdasarkan lokasi file ini.
PROJECT_ROOT=$(dirname $(dirname $(dirname "${(%):-%x}")))

# Fungsi pembantu untuk menjalankan perintah biasa
_rout_run_generic_command() {
  go run "$1" "${@:2}"
}

COMMANDS_DIR="$PROJECT_ROOT/core/system/command"

if [ -d "$COMMANDS_DIR" ]; then
  for cmd_file in "$COMMANDS_DIR"/*.go; do
    if [ -f "$cmd_file" ]; then
      local base_name=$(basename "$cmd_file" .go)

      # Abaikan file yang diawali dengan 'cd_' karena sudah ditangani file lain
      if [[ "$base_name" != cd_* ]]; then
        local cmd_name=$base_name
        
        # Definisikan fungsi yang memanggil helper biasa
        eval "$cmd_name() { _rout_run_generic_command \"$cmd_file\" \"\$@\" }"
      fi
    fi
  done
fi