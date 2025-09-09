# Aturan Komunikasi
**Bahasa Wajib**: Semua interaksi dengan Qwen Code harus menggunakan Bahasa Indonesia. Aturan ini bersifat absolut dan tanpa pengecualian untuk memastikan konsistensi, kenyamanan, dan relevansi dalam konts pengguna Indonesia.

---

# Alur Kerja Git & Versioning

Bagian ini merangkum aturan wajib untuk setiap pembaruan kode, memastikan semua perubahan tercatat dengan jelas, konsisten, dan sistematis.

---

### Aturan Komit
1.  **Bahasa Wajib**: Semua pesan komit harus ditulis dalam **Bahasa Indonesia**.
2.  **Format Pesan**: Setiap pesan harus diakhiri dengan **nomor versi proyek yang baru**.
3.  **Komit Awal**: Untuk komit pertama di repositori, gunakan format pesan: `upload pertama v0.0.0`.

### Aturan Versioning
1.  **Kenaikan Otomatis**: Sebelum `commit` dan `push` saat Anda perintahkan untuk "unggah", versi *patch* proyek akan dinaikkan (misal: `v0.1.4` menjadi `v0.1.5`).
2.  **Pencantuman Versi**: Versi baru inilah yang akan dicantumkan di akhir pesan komit.
3.  **Sinkronisasi Versi Kode**: Pastikan nomor versi yang ditampilkan dalam kode (misalnya, di `motd.py`) diperbarui agar sesuai dengan versi proyek yang baru sebelum melakukan `push` ke GitHub.

### Contoh Skenario
-   **Versi Saat Ini**: `v0.1.4`
-   **Perintah Anda**: "Perbaiki bug login dan unggah."
-   **Proses Saya**:
    1.  Menaikkan versi proyek ke `v0.1.5`.
    2.  Menjalankan `git commit -m "memperbaiki bug login v0.1.5"`.
    3.  Menjalankan `git push`.

---

# Panduan Git & Repositori

Bagian ini mencakup semua konfigurasi dan perintah Git yang relevan untuk proyek ini.

---

### Otentikasi via SSH
**PENTING**: Untuk berinteraksi dengan repositori remote di GitHub, kita akan selalu menggunakan otentikasi berbasis kunci SSH, bukan kata sandi. Ini adalah metode yang lebih aman dan diutamakan oleh GitHub.

Kunci SSH yang digunakan untuk proyek ini disimpan di:
-   `/data/data/com.termux/files/home/.ssh`

Pastikan kunci ini sudah terkonfigurasi dengan benar di akun GitHub Anda.

### Perintah Umum Git

#### 1. Inisialisasi Repositori Baru
Gunakan langkah-langkah berikut untuk proyek yang dimulai dari awal.

```bash
# 1. Buat file README sebagai penanda awal
echo "# rout" >> README.md

# 2. Inisialisasi repositori Git lokal
git init

# 3. Tambahkan semua file ke staging area
git add .

# 4. Lakukan komit pertama sesuai format yang ditentukan
git commit -m "upload pertama v0.0.0"

# 5. Atur nama branch utama menjadi 'main'
git branch -M main

# 6. Hubungkan repositori lokal ke remote GitHub via SSH
git remote add origin git@github.com:IfanID/rout-cli.git

# 7. Unggah (push) branch 'main' ke remote
git push -u origin main
```

#### 2. Mendorong dari Repositori yang Sudah Ada
Gunakan perintah ini jika Anda bekerja pada salinan lokal yang belum terhubung ke remote.

```bash
# 1. Hubungkan repositori lokal ke remote GitHub via SSH
git remote add origin git@github.com:IfanID/rout-cli.git

# 2. (Opsional) Pastikan nama branch utama adalah 'main'
git branch -M main

# 3. Unggah (push) branch 'main' ke remote
git push -u origin main
```

#### 3. Mengembalikan Repositori ke Kondisi Bersih
Gunakan perintah ini untuk mengembalikan repositori lokal ke kondisi bersih, menghapus semua perubahan lokal yang belum di-commit dan file yang belum terlacak.

```bash
# 1. Hapus semua perubahan pada file yang sudah terlacak (modified tracked files)
git reset --hard HEAD

# 2. Hapus semua file dan direktori baru yang belum terlacak (untracked files/dirs)
git clean -fd

# 3. Unduh versi terbaru dari remote untuk memastikan sinkronisasi total
git pull origin main
```

#### 4. Mengubah Komit Terakhir (Amending Last Commit)
Gunakan perintah ini untuk mengubah komit terakhir, misalnya untuk menambahkan perubahan kecil yang terlupakan atau mengoreksi pesan komit.

```bash
# 1. Lakukan perubahan yang diinginkan pada file.
# 2. Tambahkan perubahan ke area staging.
git add .

# 3. Ubah komit terakhir.
#    - Untuk menambahkan perubahan ke komit terakhir tanpa mengubah pesan:
git commit --amend --no-edit

#    - Untuk mengubah pesan komit terakhir (akan membuka editor teks):
git commit --amend

# 4. Jika komit yang diubah sudah didorong ke remote, Anda perlu melakukan force push.
#    PERINGATAN: Ini akan menulis ulang riwayat dan berisiko bagi kolaborator lain.
git push --force origin main
```

---

# Aturan Teknis Proyek

### Lokasi Eksekusi Rout
Untuk dapat memanggil `rout` dari mana saja di sistem, pastikan executable `rout` berada di path berikut:
- `/data/data/com.termux/files/usr/bin/rout`
Perlu diperhatikan bahwa repositori ini berisi *kode sumber* untuk `rout`. Executable `rout` harus dibangun dari kode sumber ini dan ditempatkan di lokasi yang ditentukan di atas.