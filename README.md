# Rout CLI

Rout adalah sebuah antarmuka baris perintah (CLI) yang kini berfungsi sebagai *wrapper* untuk Zsh. Ini memungkinkan Anda untuk menjalankan sesi Zsh yang terisolasi dengan konfigurasi kustom, terpisah dari konfigurasi Zsh utama sistem Anda.

## Fitur Utama

-   **Wrapper Zsh**: Menjalankan sesi Zsh yang sepenuhnya fungsional.
-   **Konfigurasi Terisolasi**: Menggunakan file konfigurasi Zsh (`.zshrc`) yang terpisah di `~/.rout/`, sehingga tidak mengganggu setup Zsh utama Anda.
-   **Oh My Zsh Default**: Secara default, `rout` akan memuat tema dan plugin standar dari Oh My Zsh, memberikan pengalaman shell yang kaya fitur.

## Instalasi

Untuk menginstal `rout`, Anda perlu meng-compile kode sumbernya. Pastikan Anda memiliki Go (versi 1.24.5 atau lebih baru) terinstal di sistem Anda.

1.  **Kloning Repositori:**
    ```bash
    git clone git@github.com:IfanID/rout-cli.git
    cd rout-cli
    ```

2.  **Compile Aplikasi:**
    ```bash
    go build -o rout
    ```

3.  **Pindahkan Binary (Opsional):**
    Untuk dapat menjalankan `rout` dari mana saja, pindahkan *binary* yang sudah di-compile ke direktori `$PATH` Anda (misalnya `/usr/local/bin` atau `/data/data/com.termux/files/usr/bin` untuk Termux).
    ```bash
    mv rout /data/data/com.termux/files/usr/bin/
    ```

## Penggunaan

Setelah instalasi, Anda bisa menjalankan `rout` dari terminal Anda:

```bash
rout
```

Ini akan meluncurkan sesi Zsh baru dengan konfigurasi yang dimuat dari `~/.rout/.zshrc`.

## Konfigurasi Zsh Kustom

`rout` dirancang untuk menggunakan file `.zshrc` yang terpisah agar konfigurasi Anda tidak bercampur.

-   **Lokasi File Konfigurasi**: File `.zshrc` yang digunakan oleh `rout` berada di `~/.rout/.zshrc`.
-   **Mengedit Konfigurasi**: Anda bisa mengedit file ini menggunakan editor teks favorit Anda. Perubahan akan diterapkan saat Anda menjalankan `rout` berikutnya.
    ```bash
    nano ~/.rout/.zshrc
    # atau
    vim ~/.rout/.zshrc
    ```
-   **Tema dan Plugin**: Anda bisa mengubah `ZSH_THEME` atau menambahkan `plugins` di file `~/.rout/.zshrc` ini, sama seperti Anda mengonfigurasi Oh My Zsh biasa.

## Catatan Penting

-   `rout` tidak lagi menyediakan perintah CLI internal (seperti `ls`, `cd`, `mkdir`, dll.) yang ditulis dalam Go. Semua perintah ditangani oleh Zsh itu sendiri.
-   Pastikan Oh My Zsh terinstal di sistem Anda, karena `rout` bergantung padanya untuk memuat tema dan plugin.
