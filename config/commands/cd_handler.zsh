# File ini khusus menangani perintah 'cd' yang dibuat dari file Go.

# Dapatkan root direktori proyek secara dinamis berdasarkan lokasi file ini.
# (%):-%x adalah path ke file ini. dirname akan naik satu level direktori.
PROJECT_ROOT=$(dirname $(dirname $(dirname "${(%):-%x}")))

# Fungsi pembantu untuk menjalankan perintah 'cd'
_rout_run_cd_command() {
  local target_dir
  target_dir=$(go run "$1" "${@:2}")

  if [ -n "$target_dir" ]; then
    cd "$target_dir" || echo "zsh: no such file or directory: $target_dir"
  else
    echo "Error: Perintah '$(basename "$1")' tidak menghasilkan output direktori."
  fi
}

COMMANDS_DIR="$PROJECT_ROOT/core/system/command"

if [ -d "$COMMANDS_DIR" ]; then
  # Loop hanya untuk file yang cocok dengan pola 'cd_*.go'
  for cmd_file in "$COMMANDS_DIR"/cd_*.go;
  do
    if [ -f "$cmd_file" ]; then
      local base_name=$(basename "$cmd_file" .go)
      local cmd_name=${base_name#cd_}
      
      # Definisikan fungsi yang memanggil helper 'cd'
      eval "$cmd_name() { _rout_run_cd_command \"$cmd_file\" \"\$@\" }"
    fi
  done
fi
