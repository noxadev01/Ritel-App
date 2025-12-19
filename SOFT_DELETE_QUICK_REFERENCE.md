# 🚀 Soft Delete Quick Reference

## Cheat Sheet untuk Developer

### 🔹 Konsep Dasar

| Operasi | Hard Delete (OLD) | Soft Delete (NEW) |
|---------|-------------------|-------------------|
| **Command** | `DELETE FROM table` | `UPDATE table SET deleted_at = NOW()` |
| **Data** | ❌ Hilang permanen | ✅ Masih ada di DB |
| **Recovery** | ❌ Tidak bisa | ✅ Bisa restore |
| **Query** | `SELECT *` | `SELECT * WHERE deleted_at IS NULL` |

---

## 📝 Code Templates

### Delete (Soft Delete)
```go
func (r *Repository) Delete(id int) error {
    query := `
        UPDATE table_name
        SET deleted_at = datetime('now'), updated_at = datetime('now')
        WHERE id = ? AND deleted_at IS NULL
    `
    result, err := database.DB.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("record not found or already deleted")
    }

    return nil
}
```

### Restore
```go
func (r *Repository) Restore(id int) error {
    query := `
        UPDATE table_name
        SET deleted_at = NULL, updated_at = datetime('now')
        WHERE id = ? AND deleted_at IS NOT NULL
    `
    result, err := database.DB.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to restore: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("record not found or not deleted")
    }

    return nil
}
```

### Get Deleted
```go
func (r *Repository) GetDeleted() ([]*Model, error) {
    query := `
        SELECT column1, column2, ...
        FROM table_name
        WHERE deleted_at IS NOT NULL
        ORDER BY deleted_at DESC
    `
    // ... scan rows ...
}
```

### Update Get Functions
```diff
  func (r *Repository) GetAll() ([]*Model, error) {
      query := `
          SELECT column1, column2, ...
          FROM table_name
+         WHERE deleted_at IS NULL
          ORDER BY created_at DESC
      `
  }

  func (r *Repository) GetByID(id int) (*Model, error) {
      query := `
          SELECT column1, column2, ...
          FROM table_name
-         WHERE id = ?
+         WHERE id = ? AND deleted_at IS NULL
      `
  }
```

---

## 🗄️ Database Schema

### Add Deleted At Column
```sql
ALTER TABLE table_name ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_table_name_deleted_at ON table_name(deleted_at);
```

### Table Definition
```sql
CREATE TABLE IF NOT EXISTS table_name (
    id SERIAL PRIMARY KEY,
    -- ... other columns ...
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL  -- ← Add this
);
```

---

## 💻 Usage Examples

### Basic CRUD

```go
// CREATE
repo.Create(&model)

// READ (auto-exclude deleted)
model, err := repo.GetByID(id)
allModels, err := repo.GetAll()

// UPDATE (normal)
repo.Update(&model)

// DELETE (soft delete)
repo.Delete(id)

// RESTORE (undelete)
repo.Restore(id)

// GET DELETED (for admin)
deletedModels, err := repo.GetDeleted()
```

### Frontend Integration

```javascript
// Delete (soft delete)
await pelangganAPI.delete(id)
// User tidak melihat data lagi, tapi data masih di DB

// Restore (jika ada endpoint)
await pelangganAPI.restore(id)
// Data muncul kembali di list
```

---

## 🧪 Testing Checklist

```go
// Test scenario
✅ Create record
✅ Verify GetByID returns record
✅ Delete record (soft delete)
✅ Verify GetByID returns nil (not found)
✅ Verify GetDeleted returns the record
✅ Restore record
✅ Verify GetByID returns record again
✅ Verify GetDeleted is empty
```

---

## 🔧 Troubleshooting

| Problem | Solution |
|---------|----------|
| Data tidak muncul setelah create | Check apakah deleted_at = NULL |
| Delete tidak bekerja | Check WHERE deleted_at IS NULL di query Delete |
| GetAll return data yang sudah dihapus | Tambah `WHERE deleted_at IS NULL` |
| Restore tidak bekerja | Check WHERE deleted_at IS NOT NULL di query Restore |
| Column deleted_at not found | Run migration SQL |

---

## 📊 SQL Queries

### Manual Check
```sql
-- Lihat semua data (termasuk deleted)
SELECT id, nama, deleted_at FROM pelanggan;

-- Lihat data aktif saja
SELECT id, nama FROM pelanggan WHERE deleted_at IS NULL;

-- Lihat data yang dihapus saja
SELECT id, nama, deleted_at FROM pelanggan WHERE deleted_at IS NOT NULL;

-- Restore manual (jika perlu)
UPDATE pelanggan SET deleted_at = NULL WHERE id = 5;

-- Hard delete (jika BENAR-BENAR perlu)
DELETE FROM pelanggan WHERE id = 5 AND deleted_at IS NOT NULL;
```

### Statistics
```sql
-- Count data aktif
SELECT COUNT(*) FROM pelanggan WHERE deleted_at IS NULL;

-- Count data yang dihapus
SELECT COUNT(*) FROM pelanggan WHERE deleted_at IS NOT NULL;

-- Data dihapus hari ini
SELECT id, nama, deleted_at
FROM pelanggan
WHERE DATE(deleted_at) = CURRENT_DATE;

-- Data dihapus minggu ini
SELECT id, nama, deleted_at
FROM pelanggan
WHERE deleted_at >= NOW() - INTERVAL '7 days';
```

---

## 🎯 Best Practices

### ✅ DO
- Selalu tambah `WHERE deleted_at IS NULL` di semua query SELECT
- Gunakan soft delete untuk data bisnis penting
- Buat fungsi `GetDeleted()` untuk audit
- Buat fungsi `Restore()` untuk recovery
- Tambah index di kolom `deleted_at`
- Backup database sebelum implement

### ❌ DON'T
- Jangan lupa filter `deleted_at IS NULL`
- Jangan hard delete data bisnis
- Jangan gunakan soft delete untuk data temporary (session, cache)
- Jangan expose data deleted ke user biasa
- Jangan lupa test restore functionality

---

## 🚀 Quick Start (Copy-Paste Ready)

### 1. Tambah kolom ke tabel
```sql
ALTER TABLE your_table ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_your_table_deleted_at ON your_table(deleted_at);
```

### 2. Update Repository Delete
```go
// Change this:
query := `DELETE FROM your_table WHERE id = ?`

// To this:
query := `UPDATE your_table SET deleted_at = datetime('now') WHERE id = ? AND deleted_at IS NULL`
```

### 3. Update Repository GetAll/GetByID
```go
// Add this to WHERE clause:
AND deleted_at IS NULL
```

### 4. Add Restore function
```go
func (r *YourRepository) Restore(id int) error {
    query := `UPDATE your_table SET deleted_at = NULL WHERE id = ? AND deleted_at IS NOT NULL`
    _, err := database.DB.Exec(query, id)
    return err
}
```

**Done!** 🎉

---

## 📚 References

- `SOFT_DELETE_SUMMARY.md` - Overview lengkap
- `SOFT_DELETE_IMPLEMENTATION.md` - Step-by-step guide
- `migrations/002_add_soft_delete_columns.sql` - Migration script
- `internal/repository/pelanggan_repository.go` - Example implementation

---

**Happy Coding! No more data loss panic! 🎊**
