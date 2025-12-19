# Implementasi Soft Delete - Panduan Lengkap

## 📋 Ringkasan Perubahan

Implementasi soft delete telah dibuat untuk menggantikan hard delete di project ini. Perubahan utama:

### ✅ Sudah Dikerjakan:
1. ✅ Migration SQL (`migrations/002_add_soft_delete_columns.sql`)
2. ✅ Repository Pelanggan - Soft Delete Complete
3. ✅ Repository Produk - Soft Delete Complete

### 🔧 Yang Harus Dikerjakan Manual:
4. ⏳ Repository Promo - Perlu update
5. ⏳ Repository Kategori - Perlu update
6. ⏳ Update Schema database/schema_postgres.sql

---

## 🚀 Langkah-Langkah Implementasi

### STEP 1: Jalankan Migration

```bash
# Untuk PostgreSQL
psql -U ritel -d ritel_db -f migrations/002_add_soft_delete_columns.sql

# Untuk SQLite
sqlite3 ritel.db < migrations/002_add_soft_delete_columns.sql
```

**Yang dilakukan migration**:
- Menambah kolom `deleted_at TIMESTAMP NULL` ke tabel: pelanggan, produk, promo, kategori
- Membuat index untuk performa query

---

### STEP 2: Update Promo Repository

**File**: `internal/repository/promo_repository.go`

#### A. Update Fungsi Delete (line ~589)

**SEBELUM:**
```go
func (r *PromoRepository) Delete(id int) error {
	query := `DELETE FROM promo WHERE id = ?`

	result, err := database.DB.Exec(query, id)
	// ...
}
```

**SESUDAH:**
```go
// Delete soft-deletes a promo (sets deleted_at timestamp)
func (r *PromoRepository) Delete(id int) error {
	query := `
		UPDATE promo
		SET deleted_at = datetime('now'), updated_at = datetime('now')
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete promo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("promo not found or already deleted")
	}

	return nil
}

// Restore restores a soft-deleted promo
func (r *PromoRepository) Restore(id int) error {
	query := `
		UPDATE promo
		SET deleted_at = NULL, updated_at = datetime('now')
		WHERE id = ? AND deleted_at IS NOT NULL
	`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to restore promo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("promo not found or not deleted")
	}

	return nil
}

// GetDeleted retrieves all soft-deleted promos (for admin/audit purposes)
func (r *PromoRepository) GetDeleted() ([]*models.Promo, error) {
	query := `
		SELECT
			p.id, p.nama, p.kode, p.tipe, p.tipe_promo, p.tipe_produk_berlaku, p.nilai,
			p.min_quantity, p.max_diskon, p.buy_quantity, p.get_quantity,
			p.harga_bundling, p.diskon_bundling, p.produk_x, p.produk_y,
			p.tanggal_mulai, p.tanggal_selesai, p.status,
			p.created_at, p.updated_at
		FROM promo p
		WHERE p.deleted_at IS NOT NULL
		ORDER BY p.deleted_at DESC
	`

	// Implementation sama seperti GetAll, tapi untuk yang sudah dihapus
	// [Copy implementation dari GetAll]

	return promos, nil
}
```

#### B. Update Semua Fungsi Get (tambahkan `AND deleted_at IS NULL`)

**GetAll** (line ~80):
```go
func (r *PromoRepository) GetAll() ([]*models.Promo, error) {
	query := `
		SELECT ...
		FROM promo p
		LEFT JOIN promo_produk pp ON p.id = pp.promo_id
		WHERE p.deleted_at IS NULL  -- ← TAMBAHKAN INI
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`
```

**GetByID** (line ~195):
```go
func (r *PromoRepository) GetByID(id int) (*models.Promo, error) {
	query := `
		SELECT ...
		FROM promo p
		LEFT JOIN promo_produk pp ON p.id = pp.promo_id
		WHERE p.id = ? AND p.deleted_at IS NULL  -- ← TAMBAHKAN INI
		GROUP BY p.id
	`
```

**GetByKode** (line ~311):
```go
func (r *PromoRepository) GetByKode(kode string) (*models.Promo, error) {
	query := `
		SELECT ...
		FROM promo p
		WHERE p.kode = ? AND p.deleted_at IS NULL  -- ← TAMBAHKAN INI
	`
```

---

### STEP 3: Update Kategori Repository

**File**: `internal/repository/kategori_repository.go`

#### A. Update Fungsi Delete (line ~194)

**SEBELUM:**
```go
func (r *KategoriRepository) Delete(id int) error {
	query := `DELETE FROM kategori WHERE id = ?`
	// ...
}
```

**SESUDAH:**
```go
// Delete soft-deletes a kategori (sets deleted_at timestamp)
func (r *KategoriRepository) Delete(id int) error {
	query := `
		UPDATE kategori
		SET deleted_at = datetime('now'), updated_at = datetime('now')
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete kategori: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("kategori not found or already deleted")
	}

	return nil
}

// Restore restores a soft-deleted kategori
func (r *KategoriRepository) Restore(id int) error {
	query := `
		UPDATE kategori
		SET deleted_at = NULL, updated_at = datetime('now')
		WHERE id = ? AND deleted_at IS NOT NULL
	`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to restore kategori: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("kategori not found or not deleted")
	}

	return nil
}
```

#### B. Update Semua Fungsi Get

**GetAll** (line ~44):
```go
func (r *KategoriRepository) GetAll() ([]*models.Kategori, error) {
	query := `
		SELECT k.id, ...
		FROM kategori k
		WHERE k.deleted_at IS NULL  -- ← TAMBAHKAN INI
		ORDER BY k.nama ASC
	`
```

**GetByID** (line ~88):
```go
func (r *KategoriRepository) GetByID(id int) (*models.Kategori, error) {
	query := `
		SELECT k.id, ...
		FROM kategori k
		WHERE k.id = ? AND k.deleted_at IS NULL  -- ← TAMBAHKAN INI
	`
```

**GetByNama** (line ~126):
```go
func (r *KategoriRepository) GetByNama(nama string) (*models.Kategori, error) {
	query := `
		SELECT k.id, ...
		FROM kategori k
		WHERE k.nama = ? AND k.deleted_at IS NULL  -- ← TAMBAHKAN INI
	`
```

---

### STEP 4: Update Schema Database

**File**: `database/schema_postgres.sql`

Tambahkan kolom `deleted_at` di setiap tabel yang belum ada:

```sql
-- Pelanggan table (line ~70-90)
CREATE TABLE IF NOT EXISTS pelanggan (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    telepon VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255),
    alamat TEXT,
    level INTEGER DEFAULT 1,
    tipe VARCHAR(50) DEFAULT 'reguler',
    poin INTEGER DEFAULT 0,
    total_transaksi INTEGER DEFAULT 0,
    total_belanja INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL  -- ← TAMBAHKAN INI
);

-- Produk table (line ~25-45)
CREATE TABLE IF NOT EXISTS produk (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(255) UNIQUE NOT NULL,
    barcode VARCHAR(255) UNIQUE,
    nama VARCHAR(255) NOT NULL,
    kategori VARCHAR(255),
    berat REAL DEFAULT 0,
    harga_beli INTEGER DEFAULT 0,
    harga_jual INTEGER NOT NULL,
    stok REAL DEFAULT 0,
    satuan VARCHAR(50) DEFAULT 'pcs',
    jenis_produk VARCHAR(50) DEFAULT 'satuan',
    kadaluarsa VARCHAR(100),
    tanggal_masuk TIMESTAMP,
    deskripsi TEXT,
    gambar VARCHAR(500),
    hari_pemberitahuan_kadaluarsa INTEGER DEFAULT 0,
    masa_simpan_hari INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL  -- ← TAMBAHKAN INI
);

-- Promo table (line ~145-175)
CREATE TABLE IF NOT EXISTS promo (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    kode VARCHAR(100) UNIQUE NOT NULL,
    tipe VARCHAR(50) DEFAULT 'persen',
    tipe_promo VARCHAR(50) DEFAULT 'diskon_produk',
    tipe_produk_berlaku VARCHAR(50) DEFAULT 'semua',
    nilai INTEGER DEFAULT 0,
    min_quantity INTEGER DEFAULT 1,
    max_diskon INTEGER DEFAULT 0,
    buy_quantity INTEGER DEFAULT 0,
    get_quantity INTEGER DEFAULT 0,
    harga_bundling INTEGER DEFAULT 0,
    diskon_bundling INTEGER DEFAULT 0,
    produk_x INTEGER,
    produk_y INTEGER,
    tanggal_mulai TIMESTAMP,
    tanggal_selesai TIMESTAMP,
    status VARCHAR(50) DEFAULT 'aktif',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL  -- ← TAMBAHKAN INI
);

-- Kategori table (perlu ditambahkan)
CREATE TABLE IF NOT EXISTS kategori (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) UNIQUE NOT NULL,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL  -- ← TAMBAHKAN INI
);
```

---

## 🎯 Cara Menggunakan Soft Delete

### 1. Delete (Soft Delete)

```go
// Di service layer atau handler
err := pelangganRepo.Delete(pelangganID)
// Data TIDAK dihapus permanen, hanya ditandai dengan deleted_at
```

### 2. Restore (Kembalikan Data yang Dihapus)

```go
// Restore pelanggan yang sudah dihapus
err := pelangganRepo.Restore(pelangganID)
// deleted_at di-set NULL, data aktif kembali
```

### 3. Lihat Data yang Sudah Dihapus

```go
// Untuk admin/audit
deletedPelanggans, err := pelangganRepo.GetDeleted()
// Menampilkan semua pelanggan yang deleted_at != NULL
```

### 4. Query Normal (Otomatis Exclude Deleted)

```go
// Semua fungsi Get otomatis exclude data yang sudah dihapus
allPelanggans := pelangganRepo.GetAll()
// Hanya return yang deleted_at IS NULL
```

---

## 📊 Perbandingan Before/After

### BEFORE (Hard Delete):
```sql
DELETE FROM pelanggan WHERE id = 5;
-- Data HILANG PERMANEN! ❌
```

### AFTER (Soft Delete):
```sql
UPDATE pelanggan
SET deleted_at = NOW()
WHERE id = 5;
-- Data MASIH ADA, hanya ditandai! ✅
```

---

## ⚙️ Testing Soft Delete

```go
// Test scenario
func TestSoftDelete(t *testing.T) {
    // 1. Create pelanggan
    pelanggan := &models.Pelanggan{Nama: "Test", Telepon: "08123"}
    repo.Create(pelanggan)

    // 2. Verify dapat di-get
    result, _ := repo.GetByID(pelanggan.ID)
    assert.NotNil(t, result)

    // 3. Delete (soft delete)
    repo.Delete(pelanggan.ID)

    // 4. Verify TIDAK bisa di-get lagi
    result, _ = repo.GetByID(pelanggan.ID)
    assert.Nil(t, result) // Not found karena deleted

    // 5. Verify ada di GetDeleted
    deleted, _ := repo.GetDeleted()
    assert.Len(t, deleted, 1)

    // 6. Restore
    repo.Restore(pelanggan.ID)

    // 7. Verify bisa di-get lagi
    result, _ = repo.GetByID(pelanggan.ID)
    assert.NotNil(t, result) // Found lagi!
}
```

---

## 🔍 Troubleshooting

### Q: Data tidak muncul setelah delete
**A**: Ini normal! Data soft-deleted tidak muncul di query GetAll/GetByID. Gunakan `GetDeleted()` untuk melihatnya.

### Q: Bagaimana hard delete jika benar-benar perlu?
**A**: Buat fungsi terpisah `HardDelete()` atau `PermanentDelete()`:
```go
func (r *PelangganRepository) HardDelete(id int) error {
    query := `DELETE FROM pelanggan WHERE id = ?`
    // HANYA untuk kasus khusus (GDPR, cleanup, dll)
}
```

### Q: Kolom deleted_at belum ada
**A**: Jalankan migration: `psql -U ritel -d ritel_db -f migrations/002_add_soft_delete_columns.sql`

---

## ✅ Checklist Implementasi

- [x] Migration SQL created
- [x] Pelanggan Repository updated
- [x] Produk Repository updated
- [ ] Promo Repository updated (KERJAKAN MANUAL)
- [ ] Kategori Repository updated (KERJAKAN MANUAL)
- [ ] Schema database updated (KERJAKAN MANUAL)
- [ ] Testing soft delete functionality
- [ ] Update API handlers jika perlu (untuk restore endpoint)

---

## 🚨 PENTING!

**Backup database sebelum apply perubahan!**
```bash
# PostgreSQL
pg_dump -U ritel ritel_db > backup_before_soft_delete.sql

# SQLite
cp ritel.db ritel.db.backup
```

---

Implementasi soft delete ini memberikan:
- ✅ Data safety (tidak hilang permanen)
- ✅ Audit trail (kapan dihapus, bisa di-restore)
- ✅ Compliance (GDPR, SOC2)
- ✅ Recovery capability

**Selamat mengimplementasikan! 🎉**
