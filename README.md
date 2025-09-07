# Rout CLI

Rout adalah antarmuka baris perintah (CLI) serbaguna yang kini dilengkapi dengan integrasi AI Gemini. Selain berfungsi sebagai *wrapper* untuk Zsh, Rout juga menyediakan perintah interaktif untuk berinteraksi langsung dengan model bahasa besar.

## Fitur Utama

-   **Integrasi AI Gemini**: Perintah `rcli` memungkinkan Anda berinteraksi dengan model Gemini AI secara interaktif.
-   **Animasi Menunggu Respons**: Menampilkan animasi visual yang menarik saat menunggu respons dari AI.
-   **Output AI yang Rapi**: Jawaban dari AI diformat dengan *word wrap* otomatis agar mudah dibaca.
-   **Konfigurasi Bahasa AI**: AI diinstruksikan untuk selalu menggunakan Bahasa Indonesia (campuran formal dan non-formal).
-   **Penanganan API Key Aman**: Menggunakan file `.env` untuk menyimpan API Key Gemini dengan aman.
-   **Wrapper Zsh (Opsional)**: Masih dapat berfungsi sebagai *wrapper* untuk Zsh dengan konfigurasi terisolasi (fitur lama, bisa disebutkan singkat).

## Instalasi

Untuk menginstal `rout`, Anda perlu meng-compile kode sumbernya. Pastikan Anda memiliki Go (versi 1.24.5 atau lebih baru) terinstal di sistem Anda.

1.  **Kloning Repositori:**
    ```bash
    git clone git@github.com:IfanID/rout-cli.git
    cd rout-cli
    ```

2.  **Instal Dependensi Go:**
    ```bash
    go mod tidy
    ```

3.  **Siapkan API Key Gemini:**
    Buat file `.env` di root direktori proyek (`rout-cli/`) dan tambahkan API Key Gemini Anda:
    ```
    GEMINI_API_KEY=YOUR_GEMINI_API_KEY_HERE
    ```
    Ganti `YOUR_GEMINI_API_KEY_HERE` dengan API Key Gemini Anda yang sebenarnya.

4.  **Compile Aplikasi:**
    ```bash
    go build -o rout
    ```

5.  **Pindahkan Binary (Opsional):**
    Untuk dapat menjalankan `rout` dari mana saja, pindahkan *binary* yang sudah di-compile ke direktori `$PATH` Anda (misalnya `/usr/local/bin` atau `/data/data/com.termux/files/usr/bin` untuk Termux).
    ```bash
    mv rout /data/data/com.termux/files/usr/bin/
    ```

## Penggunaan

### Perintah `rcli` (Interaksi AI)

Untuk memulai sesi chat interaktif dengan AI Gemini:

```bash
rout rcli
```

Setelah menjalankan perintah di atas, Anda akan masuk ke mode chat. Ketik pertanyaan Anda dan tekan Enter. Ketik `exit` atau `keluar` untuk mengakhiri sesi.

### Mode Zsh (Fitur Lama)

Untuk meluncurkan sesi Zsh baru dengan konfigurasi terisolasi (fitur asli `rout`):

```bash
rout
```

Ini akan memuat konfigurasi Zsh dari `~/.rout/.zshrc`.

## Konfigurasi Zsh Kustom (Fitur Lama)

`rout` dirancang untuk menggunakan file `.zshrc` yang terpisah agar konfigurasi Anda tidak bercampur.

-   **Lokasi File Konfigurasi**: File `.zshrc` yang digunakan oleh `rout` berada di `~/.rout/.zshrc`.
-   **Mengedit Konfigurasi**: Anda bisa mengedit file ini menggunakan editor teks favorit Anda. Perubahan akan diterapkan saat Anda menjalankan `rout` berikutnya.

## Catatan Penting

-   Pastikan Anda memiliki koneksi internet untuk berinteraksi dengan AI Gemini.
-   Model AI yang digunakan saat ini adalah `gemini-2.5-flash`.