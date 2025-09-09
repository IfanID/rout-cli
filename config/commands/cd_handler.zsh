# File ini khusus menangani perintah 'cd' yang terkait dengan 'rout'.

# Dapatkan root direktori proyek secara dinamis berdasarkan lokasi file ini.
PROJECT_ROOT=$(dirname $(dirname $(dirname "${(%):-%x}")))

# Fungsi pembantu untuk menjalankan perintah 'cd'
_rout_run_cd_sub() {
  local target_dir
  # Jalankan 'rout sub' dan teruskan semua argumen (misal: -ganti, -lokasi)
  target_dir=$($PROJECT_ROOT/rout sub "$@")

  # Jika output tidak kosong (artinya bukan hanya pesan status), coba cd
  if [ -n "$target_dir" ]; then
    # Pastikan kita hanya mengambil baris terakhir jika ada output lain
    target_dir=$(echo "$target_dir" | tail -n 1)
    if [ -d "$target_dir" ]; then
        cd "$target_dir" || echo "zsh: no such file or directory: $target_dir"
    # else
        # Jika output bukan direktori yang valid, cetak saja (misal: output dari -lokasi)
        # echo "$target_dir"
    fi
  fi
}

# Definisikan fungsi 'sub' secara eksplisit
sub() {
  _rout_run_cd_sub "$@"
}

# Fungsi pembantu untuk rman
_rout_run_cd_rman() {
  local target_dir
  target_dir=$($PROJECT_ROOT/rout rman "$@")
  if [ -n "$target_dir" ]; then
    target_dir=$(echo "$target_dir" | tail -n 1)
    if [ -d "$target_dir" ]; then
        cd "$target_dir" || echo "zsh: no such file or directory: $target_dir"
    fi
  fi
}

# Definisikan fungsi 'rman' secara eksplisit
rman() {
  _rout_run_cd_rman "$@"
}
