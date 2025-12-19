# Dokumentasi Sample Data Toko Sayuran dan Buah

## 📋 Overview

File `sample_data_toko_sayuran_buah.sql` berisi data sample yang realistis untuk bisnis toko sayuran dan buah. Data ini dirancang khusus untuk aplikasi POS dengan karakteristik bisnis sayuran dan buah segar.

## 🗂️ Struktur Data

### 1. **USERS** (4 records)
Data kasir dan admin toko:
- **Admin**: 1 user (admin)
- **Staff/Kasir**: 3 users (Siti, Budi, Rina)
- Password default: `password123` (hash bcrypt)
- Status: semua aktif

**Kolom yang diisi:**
- `username`, `password`, `nama_lengkap`, `role`, `status`, `created_at`

---

### 2. **KATEGORI** (7 records)
Kategori produk sayuran dan buah:
- Sayuran Daun (🌿)
- Sayuran Umbi (🥕)
- Sayuran Buah (🍅)
- Buah Lokal (🍌)
- Buah Import (🍎)
- Bumbu Dapur (🧄)
- Sayuran Organik (🌱)

**Kolom yang diisi:**
- `nama`, `deskripsi`, `icon`, `created_at`

---

### 3. **PRODUK** (68 records)
Total 68 produk dengan distribusi:
- **Sayuran Daun**: 8 produk (Bayam, Kangkung, Sawi, Selada, dll)
- **Sayuran Umbi**: 7 produk (Wortel, Kentang, Bawang Merah/Putih, dll)
- **Sayuran Buah**: 14 produk (Tomat, Terong, Timun, Cabai, dll)
- **Buah Lokal**: 18 produk (Pisang, Jeruk, Mangga, Semangka, dll)
- **Buah Import**: 7 produk (Apel, Anggur, Stroberi, Kiwi, dll)
- **Bumbu Dapur**: 8 produk (Jahe, Kunyit, Serai, Daun Bawang, dll)
- **Sayuran Organik**: 4 produk (Bayam, Kangkung, Tomat, Wortel organik)

**Karakteristik Data:**
- **SKU**: Format `SV-XXX`, `BU-XXX`, `BM-XXX`, `OR-XXX`
- **Barcode**: Format `899123456XXXX`
- **Satuan**: `kg`, `ikat`, `buah`, `pack`, `sisir`
- **Masa Simpan**: 
  - Sayuran daun: 2-3 hari
  - Sayuran buah: 3-5 hari
  - Buah lokal: 3-7 hari
  - Buah import: 5-14 hari
  - Bumbu: 14 hari
- **Harga Beli vs Jual**: Margin 30-50%
- **Harga Realistis**: 
  - Sayuran daun: Rp 8.000 - Rp 18.000/ikat
  - Cabai: Rp 25.000 - Rp 50.000/kg (fluktuatif)
  - Buah lokal: Rp 8.000 - Rp 35.000/kg
  - Buah import: Rp 35.000 - Rp 75.000/kg

**Kolom yang diisi:**
- `sku`, `barcode`, `nama`, `kategori`, `berat`, `harga_beli`, `harga_jual`, `stok`, `satuan`, `masa_simpan_hari`, `hari_pemberitahuan_kadaluarsa`, `deskripsi`, `created_at`

---

### 4. **PELANGGAN** (13 records)
Tiga tipe pelanggan:

**a. Ibu Rumah Tangga (5 records):**
- Level 1-2
- Total belanja: Rp 750.000 - Rp 3.250.000
- Diskon: 0-5%

**b. Pedagang/Reseller (4 records):**
- Level 1-3
- Total belanja: Rp 2.250.000 - Rp 12.500.000
- Diskon: 0-15% (untuk pembelian besar)

**c. Restoran (4 records):**
- Level 2-3
- Total belanja: Rp 6.000.000 - Rp 17.500.000
- Diskon: 10-15% (pembelian rutin)

**Kolom yang diisi:**
- `nama`, `telepon`, `email`, `tipe`, `level`, `poin`, `diskon_persen`, `total_transaksi`, `total_belanja`, `alamat`, `created_at`

---

### 5. **BATCH** (16 records)
Tracking kadaluarsa untuk produk perishable:

**Karakteristik:**
- **ID Format**: `BATCH-{produk_id}-{YYYYMMDD}`
- **Status**: `fresh`, `expired`
- **Supplier**: Nama supplier (dalam keterangan)
- **Pattern Restok**: 
  - Sayuran daun: Setiap 2-3 hari
  - Sayuran buah: Setiap hari
  - Buah import: Setiap 2 minggu

**Sample Batch:**
- Bayam: 4 batch (restok setiap 2-3 hari)
- Cabai Rawit: 4 batch (harga fluktuatif - naik saat musim hujan)
- Tomat: 4 batch (restok harian)
- Apel Fuji: 4 batch (restok bulanan)

**Kolom yang diisi:**
- `id`, `produk_id`, `qty`, `qty_tersisa`, `tanggal_restok`, `masa_simpan_hari`, `tanggal_kadaluarsa`, `status`, `supplier`, `keterangan`, `created_at`

---

### 6. **TRANSAKSI** (26 records)
Data transaksi 6 bulan terakhir (Januari - Juni 2024):

**Pattern Realistis:**
- **Pagi**: 06:00 - 09:00 (ramai)
- **Sore**: 16:00 - 19:00 (ramai)
- **Weekend**: Lebih ramai dari weekday
- **Akhir Bulan**: Peningkatan transaksi

**Distribusi:**
- Januari: 8 transaksi
- Februari: 3 transaksi
- Maret: 2 transaksi
- April: 2 transaksi
- Mei: 2 transaksi
- Juni: 9 transaksi

**Nilai Transaksi:**
- Minimal: Rp 35.000 (pelanggan reguler)
- Maksimal: Rp 520.000 (restoran)
- Rata-rata: Rp 150.000 - Rp 300.000

**Kolom yang diisi:**
- `nomor_transaksi`, `tanggal`, `pelanggan_nama`, `pelanggan_telp`, `pelanggan_id`, `staff_id`, `staff_nama`, `subtotal`, `diskon`, `diskon_promo`, `diskon_pelanggan`, `poin_ditukar`, `diskon_poin`, `total`, `total_bayar`, `kembalian`, `status`, `kasir`, `created_at`

---

### 7. **TRANSAKSI_ITEM** (Sample untuk beberapa transaksi)
Detail item dalam transaksi:

**Karakteristik:**
- Setiap transaksi memiliki 2-5 item
- Produk yang sering dibeli: Bayam, Kangkung, Tomat, Cabai, Wortel
- Restoran membeli dalam jumlah besar (10-20 ikat/kg)

**Kolom yang diisi:**
- `transaksi_id`, `produk_id`, `produk_sku`, `produk_nama`, `produk_kategori`, `harga_satuan`, `jumlah`, `subtotal`, `created_at`

---

### 8. **PEMBAYARAN** (11 records)
Metode pembayaran:
- **Tunai**: Untuk transaksi kecil (< Rp 100.000)
- **Transfer**: Untuk transaksi besar (restoran, pedagang)

**Kolom yang diisi:**
- `transaksi_id`, `metode`, `jumlah`, `referensi`, `created_at`

---

### 9. **STOK_HISTORY** (15 records)
History perubahan stok untuk tracking:

**Jenis Perubahan:**
- `restok`: Penambahan stok dari supplier
- `penjualan`: Pengurangan karena terjual
- `kadaluarsa`: Pengurangan karena kadaluarsa

**Sample untuk:**
- Bayam: 6 history (restok, penjualan, kadaluarsa)
- Cabai Rawit: 9 history (dengan catatan harga fluktuatif)
- Tomat: 9 history (restok harian)

**Kolom yang diisi:**
- `produk_id`, `stok_sebelum`, `stok_sesudah`, `perubahan`, `jenis_perubahan`, `keterangan`, `created_at`

---

### 10. **PROMO** (4 records)
Promo dan diskon:

1. **Diskon Sayuran Daun 10%**: Untuk produk kategori sayuran daun
2. **Diskon Buah Import 15%**: Untuk buah import
3. **Buy 2 Get 1 Cabai**: Promo khusus cabai
4. **Diskon Produk Mendekati Kadaluarsa 20%**: Untuk produk yang akan kadaluarsa dalam 1 hari

**Kolom yang diisi:**
- `nama`, `kode`, `tipe`, `tipe_promo`, `nilai`, `min_quantity`, `max_diskon`, `tanggal_mulai`, `tanggal_selesai`, `status`, `deskripsi`, `created_at`

**PROMO_PRODUK**: Link promo ke produk tertentu (12 records)

---

## 📊 Sample Queries untuk Laporan

File SQL juga menyertakan sample queries yang bisa digunakan untuk:

1. **Penjualan Harian**: Total transaksi dan penjualan per hari
2. **Produk Terlaris**: Ranking produk berdasarkan jumlah terjual
3. **Analisis Kadaluarsa**: Produk yang akan/m sudah kadaluarsa
4. **Profitabilitas Per Produk**: Profit margin per produk
5. **Trend Musiman**: Penjualan per bulan
6. **Penjualan Per Jam**: Pattern waktu ramai
7. **Pelanggan Terbaik**: Ranking pelanggan berdasarkan total belanja
8. **Produk Mendekati Kadaluarsa**: Untuk aplikasi diskon otomatis

---

## 🎯 Karakteristik Data Realistis

### 1. **Harga Fluktuatif**
- **Cabai**: Harga naik saat musim hujan (lihat batch history)
- **Buah Import**: Harga lebih mahal dan stabil
- **Sayuran Lokal**: Harga relatif stabil

### 2. **Masa Simpan Pendek**
- Sayuran daun: 2-3 hari
- Sayuran buah: 3-5 hari
- Buah: 3-14 hari (tergantung jenis)

### 3. **Pattern Penjualan**
- **Pagi (06-09)**: Ramai (ibu rumah tangga belanja)
- **Sore (16-19)**: Ramai (setelah kerja)
- **Weekend**: Lebih ramai dari weekday
- **Akhir Bulan**: Peningkatan transaksi

### 4. **Tipe Pelanggan**
- **Ibu Rumah Tangga**: Belanja kecil, sering
- **Pedagang**: Belanja besar, diskon 10-15%
- **Restoran**: Belanja rutin, diskon 10-15%

### 5. **Restok Pattern**
- **Sayuran Daun**: Setiap 2-3 hari
- **Sayuran Buah**: Setiap hari
- **Buah Import**: Setiap 2 minggu/bulan

---

## 🚀 Cara Menggunakan

### 1. **Import ke Database SQLite**

```bash
# Menggunakan sqlite3 command line
sqlite3 ritel.db < sample_data_toko_sayuran_buah.sql
```

### 2. **Atau melalui aplikasi Go**

File SQL bisa dijalankan langsung melalui aplikasi dengan menambahkan fungsi import, atau dijalankan manual menggunakan tool SQLite.

### 3. **Verifikasi Data**

Setelah import, verifikasi dengan query:

```sql
-- Cek jumlah data
SELECT 'Users' as tabel, COUNT(*) as jumlah FROM users
UNION ALL
SELECT 'Kategori', COUNT(*) FROM kategori
UNION ALL
SELECT 'Produk', COUNT(*) FROM produk
UNION ALL
SELECT 'Pelanggan', COUNT(*) FROM pelanggan
UNION ALL
SELECT 'Transaksi', COUNT(*) FROM transaksi
UNION ALL
SELECT 'Batch', COUNT(*) FROM batch;
```

---

## ⚠️ Catatan Penting

1. **Password Default**: Semua user memiliki password `password123` (hash bcrypt)
2. **Harga**: Semua harga dalam Rupiah (integer, tanpa koma)
3. **Stok**: Dalam satuan sesuai produk (kg, ikat, buah, pack)
4. **Tanggal**: Format `YYYY-MM-DD HH:MM:SS`
5. **Foreign Keys**: Pastikan data diimport sesuai urutan (users → kategori → produk → pelanggan → transaksi → dll)
6. **ID Auto Increment**: ID akan otomatis di-generate, pastikan tidak ada konflik

---

## 📈 Statistik Data

- **Total Users**: 4
- **Total Kategori**: 7
- **Total Produk**: 68
- **Total Pelanggan**: 13
- **Total Transaksi**: 26
- **Total Batch**: 16
- **Total Stok History**: 15
- **Total Promo**: 4

---

## 🔄 Update & Maintenance

Data ini bisa diperluas dengan:
- Menambah lebih banyak transaksi (6 bulan penuh)
- Menambah lebih banyak batch untuk tracking yang lebih detail
- Menambah data stok_history untuk analisis yang lebih lengkap
- Menambah promo dan diskon lainnya

---

**Dibuat untuk**: Aplikasi POS Toko Sayuran dan Buah  
**Format**: SQLite SQL  
**Versi**: 1.0  
**Tanggal**: 2024

