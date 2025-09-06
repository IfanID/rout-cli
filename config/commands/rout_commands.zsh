# File ini menangani semua perintah Go biasa (non-cd).

# Dapatkan root direktori proyek secara dinamis berdasarkan lokasi file ini.
PROJECT_ROOT=$(dirname $(dirname $(dirname "${(%):-%x}")))

# Fungsi pembantu untuk menjalankan perintah biasa
_rout_run_generic_command() {
  $PROJECT_ROOT/rout "$1" "${@:2}"
}

conv() {
    _rout_run_generic_command "conv" "$@"
}

help() {
    _rout_run_generic_command "help" "$@"
}
