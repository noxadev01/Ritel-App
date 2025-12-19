# ⚠️ PENYEBAB DATA HILANG & CARA PENCEGAHAN

## Analisis Lengkap: Mengapa Data Bisa Hilang dari Database

---

## 🔴 **KATEGORI 1: OPERASI DATABASE (Code-Based)**

### 1.1 Hard DELETE - Penghapusan Permanen

**Penyebab**: Eksekusi query `DELETE FROM table`

**Contoh di Code Anda**:

```go
// File: internal/repository/promo_repository.go:590
func (r *PromoRepository) Delete(id int) error {
    query := `DELETE FROM promo WHERE id = ?`  // ← HARD DELETE!
    result, err := database.DB.Exec(query, id)
    // Data HILANG PERMANEN!
}

// File: internal/repository/kategori_repository.go:195
func (r *KategoriRepository) Delete(id int) error {
    query := `DELETE FROM kategori WHERE id = ?`  // ← HARD DELETE!
    // ...
}
```

**Kapan Terjadi**:
- User klik tombol "Hapus" di UI
- Admin manually delete via code
- Scheduled cleanup job
- Testing yang lupa rollback

**Dampak**: ❌ Data hilang PERMANEN, TIDAK BISA di-recover!

**Solusi**: ✅ Gunakan SOFT DELETE (sudah kita implement!)

---

### 1.2 DELETE Tanpa WHERE Clause (Paling Berbahaya!)

**Penyebab**: Query DELETE tanpa filter

```go
// SANGAT BERBAHAYA! ❌❌❌
func (r *KeranjangRepository) Clear() error {
    query := `DELETE FROM keranjang`  // ← TIDAK ADA WHERE!
    _, err := database.DB.Exec(query)
    // Ini akan HAPUS SEMUA keranjang SEMUA user!
}
```

**Yang Terjadi di Code Anda**:

**File**: `internal/repository/keranjang_repository.go:175`
```go
// Clear empties the entire cart
func (r *KeranjangRepository) Clear() error {
    query := `DELETE FROM keranjang`  // ← NO WHERE CLAUSE!
    _, err := database.DB.Exec(query)
    return err
}
```

**Skenario Bahaya**:
```
Kasir A: Checkout transaksi (keranjang A cleared)
System: DELETE FROM keranjang  ← Tanpa WHERE!
Result: SEMUA keranjang (Kasir A, B, C, D) HILANG! 🔥
```

**Dampak**:
- ❌ Multi-user app: Keranjang SEMUA user hilang
- ❌ Desktop app: Aman (1 user), tapi tetap risky

**Solusi**:
```go
// SHOULD BE:
func (r *KeranjangRepository) Clear(userID int) error {
    query := `DELETE FROM keranjang WHERE user_id = ?`
    _, err := database.DB.Exec(query, userID)
    return err
}
```

---

### 1.3 CASCADE DELETE - Efek Domino

**Penyebab**: Foreign key constraint dengan `ON DELETE CASCADE`

**Contoh di Schema Anda**:

**File**: `database/schema_postgres.sql:195-200`
```sql
CREATE TABLE promo_produk (
    promo_id INTEGER NOT NULL,
    produk_id INTEGER NOT NULL,
    CONSTRAINT fk_promo_produk_promo FOREIGN KEY (promo_id)
        REFERENCES promo(id) ON DELETE CASCADE,  -- ← CASCADE!
    CONSTRAINT fk_promo_produk_produk FOREIGN KEY (produk_id)
        REFERENCES produk(id) ON DELETE CASCADE   -- ← CASCADE!
);
```

**Skenario Bahaya**:
```sql
-- Admin hapus 1 promo
DELETE FROM promo WHERE id = 10;

-- CASCADE effect (otomatis dari database):
-- 1. SEMUA promo_produk yang promo_id = 10 ikut TERHAPUS!
-- 2. Bisa ratusan record hilang dalam 1 detik!
```

**Lebih Parah Lagi**:
```sql
-- File: database/schema_postgres.sql:272
REFERENCES produk(id) ON DELETE CASCADE

-- Hapus 1 produk:
DELETE FROM produk WHERE id = 50;

-- Yang ikut TERHAPUS otomatis:
-- ❌ Semua batch produk ini
-- ❌ Semua stok_history
-- ❌ Semua promo_produk
-- ❌ Semua keranjang yang ada produk ini
-- ❌ Ratusan record hilang!
```

**Dampak**: ❌ 1 DELETE bisa hapus ratusan/ribuan record terkait!

**Solusi**:
```sql
-- Ganti CASCADE dengan RESTRICT
REFERENCES promo(id) ON DELETE RESTRICT  -- Prevent delete jika ada data terkait
-- Atau pakai soft delete!
```

---

### 1.4 TRUNCATE TABLE - Reset Seluruh Tabel

**Penyebab**: Command TRUNCATE (lebih cepat dari DELETE, tapi lebih berbahaya)

```sql
-- Hapus SEMUA data dari tabel!
TRUNCATE TABLE pelanggan;
-- Semua pelanggan HILANG dalam sekejap! ⚡💥
```

**Kapan Terjadi**:
- Developer salah run script
- Testing database cleanup
- Migration yang salah
- Typo saat development

**Dampak**:
- ❌ SEMUA data di tabel hilang
- ❌ Auto-increment counter di-reset ke 0
- ❌ Lebih cepat dari DELETE tapi TIDAK BISA di-rollback (di beberapa DB)

**Solusi**:
- Jangan pernah pakai TRUNCATE di production
- Selalu backup sebelum run script
- Gunakan soft delete

---

### 1.5 DROP TABLE - Hapus Tabel Beserta Strukturnya

**Penyebab**: Command DROP TABLE

```sql
-- Hapus tabel pelanggan beserta SEMUA datanya!
DROP TABLE pelanggan;
-- Tabel HILANG TOTAL! Structure + Data! 💀
```

**Di Code Anda** (untuk migration):

**File**: `internal/database/database.go:709`
```go
// Migration code - BERBAHAYA jika salah!
if _, err = tx.Exec("DROP TABLE transaksi_item;"); err != nil {
    return fmt.Errorf("failed to drop old table: %w", err)
}
```

**Kapan Terjadi**:
- Migration script yang salah
- Developer typo saat development
- Script cleanup yang over-aggressive
- Database maintenance yang salah

**Dampak**:
- ❌ Tabel + SEMUA data hilang
- ❌ Application crash (table not found)
- ❌ Perlu re-create table & restore dari backup

**Solusi**:
- ALWAYS backup sebelum migration
- Test di development dulu
- Gunakan IF EXISTS dengan hati-hati
- Version control untuk migration scripts

---

### 1.6 UPDATE Tanpa WHERE (Overwrite Semua Data)

**Penyebab**: Query UPDATE tanpa kondisi

```sql
-- BERBAHAYA! ❌
UPDATE pelanggan SET poin = 0;
-- Semua pelanggan poinnya jadi 0!

UPDATE produk SET stok = 0;
-- Semua produk stoknya jadi 0! Seperti data "hilang"!
```

**Kapan Terjadi**:
- Typo/lupa WHERE clause
- Bug di dynamic query builder
- Testing yang tidak rollback

**Dampak**:
- ❌ Semua record ter-update dengan nilai yang sama
- ❌ Data original hilang (overwritten)

**Solusi**:
```sql
-- ALWAYS use WHERE
UPDATE pelanggan SET poin = 0 WHERE id = 5;  -- ✅ Safe

-- Or use transaction for safety
BEGIN;
UPDATE pelanggan SET poin = 0 WHERE id = 5;
-- Check result first
SELECT * FROM pelanggan WHERE id = 5;
-- If OK, commit. If not, rollback.
COMMIT;
```

---

## 🟡 **KATEGORI 2: TRANSACTION & ROLLBACK**

### 2.1 Transaction Rollback (Unintended)

**Penyebab**: Transaction di-rollback karena error

```go
func DoSomething() error {
    tx, _ := db.Begin()
    defer tx.Rollback()  // ← Ini akan rollback jika function return error

    // Insert data
    tx.Exec("INSERT INTO pelanggan ...")

    // Error terjadi di tengah jalan
    if someCondition {
        return fmt.Errorf("error!")  // ← ROLLBACK! Insert di-undo!
    }

    tx.Commit()  // Tidak pernah sampai sini
}
```

**Kapan Terjadi**:
- Validation error di tengah proses
- Network timeout
- Constraint violation
- Panic/crash

**Dampak**:
- ❌ Data yang sudah di-insert/update hilang (rollback)
- User bingung: "Kok data saya hilang?"

**Solusi**:
```go
func DoSomething() error {
    tx, _ := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // ... operations ...

    if err != nil {
        tx.Rollback()
        return err
    }

    // Only commit if NO errors
    return tx.Commit()  // ✅ Explicit commit
}
```

---

### 2.2 Dual Database Rollback

**Penyebab**: Rollback di dual database mode (PostgreSQL + SQLite)

**Di Code Anda**:

```go
// internal/database/database.go - Dual transaction
func (dt *DualTx) Rollback() error {
    // Rollback PostgreSQL
    dt.PostgresTx.Rollback()

    // Rollback SQLite
    dt.SQLiteTx.Rollback()

    // Data di KEDUA database hilang!
}
```

**Skenario**:
```
1. Insert data ke PostgreSQL ✅
2. Insert data ke SQLite ❌ (error)
3. Rollback KEDUA database
4. Data hilang dari PostgreSQL juga!
```

**Dampak**: ❌ Data bisa hilang dari kedua database karena partial failure

**Solusi**: Already handled with transaction, tapi tetap risky di network issues

---

## 🟠 **KATEGORI 3: CONSTRAINT VIOLATIONS**

### 3.1 Foreign Key Constraint RESTRICT

**Penyebab**: Delete parent record yang punya child

```sql
-- Schema:
FOREIGN KEY (produk_id) REFERENCES produk(id) ON DELETE RESTRICT

-- Coba delete produk:
DELETE FROM produk WHERE id = 50;
-- ERROR: cannot delete, ada transaksi_item yang pakai produk ini

-- Tapi jika TIDAK ada constraint:
-- Data bisa "orphan" (produk_id ada tapi produknya sudah hilang)
```

**Dampak**:
- ✅ RESTRICT: Mencegah data loss (bagus!)
- ❌ NO CONSTRAINT: Data jadi orphan (seperti "hilang")

---

### 3.2 Unique Constraint Violations

**Penyebab**: Insert duplicate, old data di-overwrite

```sql
-- Jika ada UNIQUE constraint:
INSERT INTO pelanggan (telepon, nama) VALUES ('08123', 'Budi');
-- ERROR: duplicate

-- Jika pakai INSERT ... ON CONFLICT UPDATE (PostgreSQL):
INSERT INTO pelanggan (telepon, nama) VALUES ('08123', 'Andi')
ON CONFLICT (telepon) DO UPDATE SET nama = 'Andi';
-- Data lama (Budi) di-OVERWRITE jadi Andi!
```

**Dampak**: ❌ Data lama ter-overwrite oleh data baru

---

## 🔵 **KATEGORI 4: BUG DI CODE**

### 4.1 Logic Error - Salah Kondisi

```go
// BUG: Salah logic
func DeleteExpiredProducts() {
    // SHOULD BE: WHERE kadaluarsa < NOW()
    // ACTUAL: WHERE kadaluarsa > NOW()  ← SALAH!
    query := `DELETE FROM produk WHERE kadaluarsa > NOW()`
    db.Exec(query)
    // Yang ke-delete: produk yang BELUM expired!
    // Yang expired malah TIDAK ke-delete!
}
```

**Kapan Terjadi**:
- Typo operator (>, <, >=, <=, =, !=)
- Logic terbalik
- Timezone issues
- Date comparison salah

**Dampak**: ❌ Data yang salah terhapus

---

### 4.2 Loop Delete (Accidental Mass Delete)

```go
// BUG: Delete semua produk!
func DeleteProducts(ids []int) {
    for _, id := range ids {
        // BUG: Seharusnya WHERE id = ?
        query := `DELETE FROM produk`  // ← LUPA WHERE!
        db.Exec(query)  // Setiap loop, hapus SEMUA produk!
    }
}
```

**Dampak**: ❌ Mass deletion tidak disengaja

---

### 4.3 Race Condition (Concurrency Issues)

```go
// Thread 1: Check stok
stok := GetStok(produkID)  // Stok = 10

// Thread 2: Check stok (same time)
stok := GetStok(produkID)  // Stok = 10

// Thread 1: Update stok
UpdateStok(produkID, stok - 5)  // Stok jadi 5

// Thread 2: Update stok
UpdateStok(produkID, stok - 5)  // Stok jadi 5 (SHOULD BE 0!)

// Data "hilang": Seharusnya 10 - 5 - 5 = 0, tapi jadi 5
```

**Dampak**: ❌ Data inconsistency, seperti "hilang"

**Solusi**: Use database locks or atomic operations

---

## 🟢 **KATEGORI 5: INFRASTRUCTURE & HARDWARE**

### 5.1 Database Corruption

**Penyebab**:
- Disk failure
- Power outage saat write operation
- File system errors
- Database crash

**Dampak**:
- ❌ Partial data loss
- ❌ Table corruption
- ❌ Index corruption

**Solusi**:
- Regular backups
- RAID configuration
- UPS (Uninterruptible Power Supply)
- Database replication

---

### 5.2 Disk Full / Out of Space

**Penyebab**: Disk penuh, database tidak bisa write

```bash
df -h
# /dev/sda1  100G  100G  0  100%  /var/lib/postgresql
```

**Dampak**:
- ❌ Insert/Update GAGAL
- ❌ Transaction rollback
- ❌ Database crash
- Data seperti "hilang" karena tidak ter-save

**Solusi**: Monitor disk space, cleanup logs

---

### 5.3 Database Server Crash

**Penyebab**:
- Out of memory (OOM)
- Process killed
- Server restart
- Network failure

**Dampak**:
- ❌ Uncommitted transactions hilang
- ❌ In-memory data loss
- ❌ Partial writes

**Solusi**:
- Database high availability (HA)
- Replication
- Automatic failover

---

## 🟣 **KATEGORI 6: HUMAN ERROR**

### 6.1 Salah Execute Script

```sql
-- Developer mau run di development:
DELETE FROM pelanggan WHERE id = 5;

-- Tapi ternyata connect ke PRODUCTION! 💀
-- Data production HILANG!
```

**Kapan Terjadi**:
- Wrong database connection
- Copy-paste script salah
- Testing di production (big mistake!)

**Dampak**: ❌ Data production hilang

**Solusi**:
- Always check connection before run script
- Use transaction + verify before commit
- Separate credentials for dev/prod
- Read-only access for non-admin

---

### 6.2 Backup Restore Salah

```bash
# Mau restore backup kemarin
pg_restore backup_2025_12_18.sql

# Tapi ternyata restore backup lama (overwrite data terbaru)
pg_restore backup_2025_12_01.sql  # ← SALAH!

# Data dari 18 Des hilang, balik ke 1 Des!
```

**Dampak**: ❌ Data terbaru ter-overwrite oleh data lama

---

### 6.3 Salah Klik (UI/Frontend)

```javascript
// User klik "Delete All" tidak sengaja
const handleDeleteAll = async () => {
    // SHOULD: Show confirmation
    // ACTUAL: Langsung delete tanpa confirm
    await produkAPI.deleteAll();  // SEMUA produk HILANG!
};
```

**Dampak**: ❌ Mass deletion karena UI tidak ada confirmation

**Solusi**:
- Always show confirmation dialog
- Implement soft delete
- Undo functionality

---

## 🔴 **KATEGORI 7: SECURITY BREACH**

### 7.1 SQL Injection

```go
// VULNERABLE CODE! ❌
func DeleteUser(username string) {
    query := fmt.Sprintf("DELETE FROM users WHERE username = '%s'", username)
    db.Exec(query)
}

// Attacker input:
username := "admin' OR '1'='1"

// Query jadi:
// DELETE FROM users WHERE username = 'admin' OR '1'='1'
// Semua user TERHAPUS! 💀
```

**Dampak**: ❌ Mass deletion oleh attacker

**Solusi**: Use prepared statements!
```go
query := `DELETE FROM users WHERE username = ?`
db.Exec(query, username)  // ✅ Safe
```

---

### 7.2 Unauthorized Access

**Penyebab**:
- Leaked credentials
- Weak password
- No authentication
- Missing authorization checks

```bash
# Attacker dapat access database:
psql -U ritel -d ritel_db
> DROP TABLE pelanggan;  # GONE! 💀
```

**Dampak**: ❌ Intentional data deletion oleh attacker

**Solusi**:
- Strong passwords
- Network firewall
- Database user permissions (read-only for app)
- Audit logging

---

## 📊 **RINGKASAN PENYEBAB DATA LOSS**

| Kategori | Penyebab | Risk Level | Ada di Code? |
|----------|----------|------------|--------------|
| **Code-Based** | Hard DELETE | 🔴 HIGH | ✅ Yes (promo, kategori) |
| | DELETE tanpa WHERE | 🔴 CRITICAL | ✅ Yes (keranjang) |
| | CASCADE DELETE | 🟡 MEDIUM | ✅ Yes (schema) |
| | TRUNCATE TABLE | 🟠 HIGH | ❌ No |
| | DROP TABLE | 🔴 CRITICAL | ⚠️ Yes (migration) |
| | UPDATE tanpa WHERE | 🟡 MEDIUM | ❌ No |
| **Transaction** | Rollback unintended | 🟡 MEDIUM | Possible |
| | Dual DB rollback | 🟡 MEDIUM | ✅ Yes (dual mode) |
| **Constraint** | FK violation | 🟢 LOW | Protected |
| **Bug** | Logic error | 🟡 MEDIUM | Possible |
| | Race condition | 🟡 MEDIUM | Possible |
| **Infrastructure** | DB corruption | 🟠 MEDIUM | Possible |
| | Disk full | 🟡 MEDIUM | Possible |
| | Server crash | 🟡 MEDIUM | Possible |
| **Human Error** | Wrong script | 🔴 HIGH | Possible |
| | Wrong restore | 🟠 MEDIUM | Possible |
| | Salah klik | 🟡 MEDIUM | Possible |
| **Security** | SQL injection | 🔴 CRITICAL | ❌ Protected |
| | Unauthorized access | 🔴 HIGH | Depends |

---

## 🛡️ **SOLUSI PENCEGAHAN**

### ✅ Sudah Implemented
- [x] Prepared statements (no SQL injection)
- [x] Foreign key constraints
- [x] Transaction safety
- [x] Soft delete (pelanggan, produk)

### ⏳ Perlu Ditambahkan
- [ ] Soft delete untuk promo & kategori
- [ ] Fix: Clear keranjang with user_id filter
- [ ] Confirmation dialogs di frontend
- [ ] Database backups (scheduled)
- [ ] Audit logging (who deleted what & when)
- [ ] Read-only database user for app
- [ ] Point-in-time recovery setup

### 🎯 Best Practices
```sql
-- 1. Always use WHERE clause
DELETE FROM table WHERE id = ?  -- ✅ Good
DELETE FROM table              -- ❌ Bad

-- 2. Use transactions
BEGIN;
DELETE FROM table WHERE id = ?;
-- Verify result
SELECT * FROM table WHERE id = ?;
COMMIT;  -- or ROLLBACK if wrong

-- 3. Backup before dangerous operations
pg_dump ritel_db > backup_before_delete.sql
-- Then do the operation

-- 4. Use soft delete
UPDATE table SET deleted_at = NOW() WHERE id = ?  -- ✅ Recoverable
DELETE FROM table WHERE id = ?                     -- ❌ Permanent
```

---

## 🚨 **EMERGENCY: Data Sudah Hilang, Apa Yang Harus Dilakukan?**

### Step 1: STOP Immediately!
```bash
# Stop application
systemctl stop ritel-app

# Prevent further writes
chmod 000 /var/lib/postgresql/data  # Make DB read-only
```

### Step 2: Check Backups
```bash
# List available backups
ls -lah /backup/database/

# Check latest backup
ls -lah backup_*.sql | tail -1
```

### Step 3: Restore dari Backup
```bash
# Restore PostgreSQL
psql -U ritel -d ritel_db < backup_latest.sql

# Restore SQLite
cp backup_ritel.db ritel.db
```

### Step 4: Check Transaction Logs (if available)
```sql
-- PostgreSQL WAL (Write-Ahead Log)
-- Bisa recover data dari log jika WAL archiving enabled
```

### Step 5: Soft Delete Recovery
```sql
-- Jika pakai soft delete, EASY!
UPDATE pelanggan SET deleted_at = NULL WHERE id = 5;
-- Data kembali! ✅
```

---

## 📝 **CHECKLIST KEAMANAN DATA**

- [ ] Implement soft delete untuk semua tabel penting
- [ ] Backup database otomatis (daily)
- [ ] Test restore procedure
- [ ] Add confirmation dialogs untuk delete operations
- [ ] Audit logging (track who deleted what)
- [ ] Separate credentials dev/staging/prod
- [ ] Read-only database user untuk aplikasi
- [ ] Foreign key constraints check
- [ ] WHERE clause validation
- [ ] Transaction safety review
- [ ] Disk space monitoring
- [ ] Database replication (HA)

---

**Kesimpulan**: Data bisa hilang dari banyak penyebab, tapi dengan **soft delete + regular backups + proper validation**, Anda bisa mencegah 90% data loss! 🛡️
