# 🔄 Panduan Lengkap Restore Data

## Cara Restore Data yang Sudah Di-Soft Delete

---

## 1️⃣ SERVICE LAYER

### Tambahkan Fungsi Restore ke Service

**File**: `internal/service/pelanggan_service.go`

Tambahkan fungsi ini ke PelangganService:

```go
// Restore restores a soft-deleted pelanggan
func (s *PelangganService) Restore(id int) error {
    // Validate ID
    if id <= 0 {
        return fmt.Errorf("invalid pelanggan ID")
    }

    // Call repository restore
    err := s.pelangganRepo.Restore(id)
    if err != nil {
        return fmt.Errorf("failed to restore pelanggan: %w", err)
    }

    log.Printf("Successfully restored pelanggan ID: %d", id)
    return nil
}

// GetDeleted gets all soft-deleted pelanggan (for admin)
func (s *PelangganService) GetDeleted() ([]*models.Pelanggan, error) {
    return s.pelangganRepo.GetDeleted()
}
```

**File**: `internal/service/produk_service.go`

```go
// Restore restores a soft-deleted product
func (s *ProdukService) Restore(id int) error {
    if id <= 0 {
        return fmt.Errorf("invalid product ID")
    }

    err := s.produkRepo.Restore(id)
    if err != nil {
        return fmt.Errorf("failed to restore product: %w", err)
    }

    log.Printf("Successfully restored product ID: %d", id)
    return nil
}

// GetDeleted gets all soft-deleted products (for admin)
func (s *ProdukService) GetDeleted() ([]*models.Produk, error) {
    return s.produkRepo.GetDeleted()
}
```

---

## 2️⃣ WAILS APP LAYER (Desktop Mode)

### Tambahkan Method ke app.go

**File**: `app.go`

Tambahkan method Wails untuk expose ke frontend:

```go
// ==================== PELANGGAN RESTORE ====================

// RestorePelanggan restores a soft-deleted customer
func (a *App) RestorePelanggan(id int) error {
    return a.services.PelangganService.Restore(id)
}

// GetDeletedPelanggan gets all soft-deleted customers
func (a *App) GetDeletedPelanggan() ([]*models.Pelanggan, error) {
    return a.services.PelangganService.GetDeleted()
}

// ==================== PRODUK RESTORE ====================

// RestoreProduk restores a soft-deleted product
func (a *App) RestoreProduk(id int) error {
    return a.services.ProdukService.Restore(id)
}

// GetDeletedProduk gets all soft-deleted products
func (a *App) GetDeletedProduk() ([]*models.Produk, error) {
    return a.services.ProdukService.GetDeleted()
}

// ==================== PROMO RESTORE ====================

// RestorePromo restores a soft-deleted promo
func (a *App) RestorePromo(id int) error {
    return a.services.PromoService.Restore(id)
}

// GetDeletedPromo gets all soft-deleted promos
func (a *App) GetDeletedPromo() ([]*models.Promo, error) {
    return a.services.PromoService.GetDeleted()
}

// ==================== KATEGORI RESTORE ====================

// RestoreKategori restores a soft-deleted category
func (a *App) RestoreKategori(id int) error {
    return a.services.KategoriService.Restore(id)
}

// GetDeletedKategori gets all soft-deleted categories
func (a *App) GetDeletedKategori() ([]*models.Kategori, error) {
    return a.services.KategoriService.GetDeleted()
}
```

Setelah tambah ini, jalankan:
```bash
wails generate module
```

---

## 3️⃣ HTTP HANDLER LAYER (Web Mode)

### Tambahkan Handler untuk Restore

**File**: `internal/http/handlers/pelanggan_handler.go`

Tambahkan method restore:

```go
// Restore restores a soft-deleted pelanggan
func (h *PelangganHandler) Restore(c *gin.Context) {
    // Get ID from URL parameter
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        response.BadRequest(c, "Invalid pelanggan ID", err)
        return
    }

    // Call service restore
    err = h.services.PelangganService.Restore(id)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, nil, "Pelanggan berhasil di-restore")
}

// GetDeleted gets all soft-deleted pelanggan
func (h *PelangganHandler) GetDeleted(c *gin.Context) {
    pelanggans, err := h.services.PelangganService.GetDeleted()
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, pelanggans, "Success")
}
```

**File**: `internal/http/handlers/produk_handler.go`

```go
// Restore restores a soft-deleted product
func (h *ProdukHandler) Restore(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        response.BadRequest(c, "Invalid product ID", err)
        return
    }

    err = h.services.ProdukService.Restore(id)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, nil, "Produk berhasil di-restore")
}

// GetDeleted gets all soft-deleted products
func (h *ProdukHandler) GetDeleted(c *gin.Context) {
    produks, err := h.services.ProdukService.GetDeleted()
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, produks, "Success")
}
```

---

## 4️⃣ ROUTER (HTTP Routes)

### Tambahkan Routes untuk Restore

**File**: `internal/http/router.go`

Tambahkan routes di bagian yang sesuai:

```go
// Di dalam protected group

// ==================== PELANGGAN ====================
pelanggan := protected.Group("/pelanggan")
{
    pelanggan.GET("", pelangganHandler.GetAll)
    pelanggan.GET("/:id", pelangganHandler.GetByID)
    pelanggan.POST("", pelangganHandler.Create)
    pelanggan.PUT("", pelangganHandler.Update)
    pelanggan.DELETE("/:id", pelangganHandler.Delete)

    // ✨ TAMBAHKAN INI (New routes untuk restore)
    pelanggan.POST("/:id/restore", pelangganHandler.Restore)
    pelanggan.GET("/deleted", pelangganHandler.GetDeleted)
}

// ==================== PRODUK ====================
produk := protected.Group("/produk")
{
    produk.GET("", produkHandler.GetAll)
    produk.GET("/:id", produkHandler.GetByID)
    produk.POST("", produkHandler.Create)
    produk.PUT("", produkHandler.Update)
    produk.DELETE("/:id", produkHandler.Delete)

    // ✨ TAMBAHKAN INI (New routes untuk restore)
    produk.POST("/:id/restore", produkHandler.Restore)
    produk.GET("/deleted", produkHandler.GetDeleted)
}

// ==================== PROMO ====================
promo := protected.Group("/promo")
{
    // ... existing routes ...

    // ✨ TAMBAHKAN INI
    promo.POST("/:id/restore", promoHandler.Restore)
    promo.GET("/deleted", promoHandler.GetDeleted)
}

// ==================== KATEGORI ====================
kategori := protected.Group("/kategori")
{
    // ... existing routes ...

    // ✨ TAMBAHKAN INI
    kategori.POST("/:id/restore", kategoriHandler.Restore)
    kategori.GET("/deleted", kategoriHandler.GetDeleted)
}
```

---

## 5️⃣ FRONTEND API CLIENT

### Update API Client untuk Restore

**File**: `frontend/src/api/pelanggan.js`

Tambahkan fungsi restore:

```javascript
import client from './client';
import { isWebMode } from '../utils/environment';

const pelangganAPI = {
  // ... existing functions ...

  /**
   * Restore soft-deleted customer
   * @param {number} id
   * @returns {Promise<void>}
   */
  restore: async (id) => {
    if (isWebMode()) {
      await client.post(`/api/pelanggan/${id}/restore`);
    } else {
      const { RestorePelanggan } = await import('../../wailsjs/go/main/App');
      return await RestorePelanggan(id);
    }
  },

  /**
   * Get all deleted customers (for admin)
   * @returns {Promise<Array>}
   */
  getDeleted: async () => {
    if (isWebMode()) {
      const response = await client.get('/api/pelanggan/deleted');
      return response.data.data;
    } else {
      const { GetDeletedPelanggan } = await import('../../wailsjs/go/main/App');
      return await GetDeletedPelanggan();
    }
  },
};

export default pelangganAPI;
```

**File**: `frontend/src/api/produk.js`

```javascript
const produkAPI = {
  // ... existing functions ...

  /**
   * Restore soft-deleted product
   * @param {number} id
   * @returns {Promise<void>}
   */
  restore: async (id) => {
    if (isWebMode()) {
      await client.post(`/api/produk/${id}/restore`);
    } else {
      const { RestoreProduk } = await import('../../wailsjs/go/main/App');
      return await RestoreProduk(id);
    }
  },

  /**
   * Get all deleted products (for admin)
   * @returns {Promise<Array>}
   */
  getDeleted: async () => {
    if (isWebMode()) {
      const response = await client.get('/api/produk/deleted');
      return response.data.data;
    } else {
      const { GetDeletedProduk } = await import('../../wailsjs/go/main/App');
      return await GetDeletedProduk();
    }
  },
};

export default produkAPI;
```

---

## 6️⃣ FRONTEND UI COMPONENT

### Contoh Component untuk Restore

Buat component untuk lihat & restore data yang dihapus:

**File**: `frontend/src/components/pages/pelanggan/DeletedPelanggan.jsx`

```jsx
import React, { useState, useEffect } from 'react';
import { pelangganAPI } from '../../../api';
import { useToast } from '../../common/ToastContainer';

const DeletedPelanggan = () => {
  const [deletedPelanggan, setDeletedPelanggan] = useState([]);
  const [loading, setLoading] = useState(false);
  const { showToast } = useToast();

  // Load deleted customers
  const loadDeleted = async () => {
    setLoading(true);
    try {
      const data = await pelangganAPI.getDeleted();
      setDeletedPelanggan(data);
    } catch (error) {
      showToast('Gagal memuat data pelanggan terhapus', 'error');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadDeleted();
  }, []);

  // Restore customer
  const handleRestore = async (pelanggan) => {
    const confirm = window.confirm(
      `Restore pelanggan "${pelanggan.nama}"?\n\n` +
      `Data akan aktif kembali dan muncul di daftar pelanggan.`
    );

    if (!confirm) return;

    try {
      await pelangganAPI.restore(pelanggan.id);
      showToast(`Pelanggan "${pelanggan.nama}" berhasil di-restore!`, 'success');

      // Reload list
      await loadDeleted();
    } catch (error) {
      showToast('Gagal me-restore pelanggan: ' + error.message, 'error');
      console.error(error);
    }
  };

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800">
          Pelanggan yang Dihapus
        </h1>
        <p className="text-gray-600">
          Daftar pelanggan yang sudah dihapus. Anda bisa me-restore jika diperlukan.
        </p>
      </div>

      {loading ? (
        <div className="text-center py-8">
          <p>Loading...</p>
        </div>
      ) : deletedPelanggan.length === 0 ? (
        <div className="bg-gray-50 rounded-lg p-8 text-center">
          <p className="text-gray-500">
            Tidak ada pelanggan yang dihapus
          </p>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  Nama
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  Telepon
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  Poin
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  Dihapus Pada
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {deletedPelanggan.map((pelanggan) => (
                <tr key={pelanggan.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {pelanggan.id}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {pelanggan.nama}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {pelanggan.telepon}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {pelanggan.poin} poin
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {new Date(pelanggan.deleted_at).toLocaleString('id-ID')}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <button
                      onClick={() => handleRestore(pelanggan)}
                      className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded transition-colors"
                    >
                      🔄 Restore
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default DeletedPelanggan;
```

---

## 7️⃣ CARA PAKAI (User Perspective)

### Skenario: User salah hapus pelanggan

**STEP 1: Hapus Pelanggan (Soft Delete)**
```
User: Klik tombol "Hapus" di DaftarPelanggan
→ Pelanggan hilang dari list
→ Data masih ada di database (deleted_at = NOW())
```

**STEP 2: Lihat Data yang Dihapus**
```
Admin: Buka menu "Pelanggan" → "Data Terhapus"
→ Muncul list pelanggan yang sudah dihapus
→ Ada tombol "Restore" di setiap row
```

**STEP 3: Restore Pelanggan**
```
Admin: Klik tombol "Restore" pada pelanggan yang salah dihapus
→ Konfirmasi: "Restore pelanggan 'Budi'?"
→ Klik OK
→ Pelanggan kembali aktif di DaftarPelanggan
→ Data lengkap (poin, transaksi history) tetap utuh!
```

---

## 8️⃣ MANUAL RESTORE VIA DATABASE

### Jika perlu restore langsung via database:

```sql
-- 1. Lihat data yang sudah dihapus
SELECT id, nama, telepon, deleted_at
FROM pelanggan
WHERE deleted_at IS NOT NULL;

-- Output:
-- id | nama | telepon    | deleted_at
-- 5  | Budi | 08123456   | 2025-12-19 14:30:00

-- 2. Restore specific pelanggan
UPDATE pelanggan
SET deleted_at = NULL, updated_at = NOW()
WHERE id = 5;

-- 3. Verify sudah aktif kembali
SELECT id, nama, telepon, deleted_at
FROM pelanggan
WHERE id = 5;

-- Output:
-- id | nama | telepon    | deleted_at
-- 5  | Budi | 08123456   | NULL  ← Sudah NULL, data aktif!
```

### Restore semua data yang dihapus hari ini

```sql
-- HATI-HATI! Ini akan restore SEMUA yang dihapus hari ini
UPDATE pelanggan
SET deleted_at = NULL, updated_at = NOW()
WHERE DATE(deleted_at) = CURRENT_DATE;
```

---

## 9️⃣ TESTING RESTORE

### Test Script Lengkap

```go
package main

import "testing"

func TestRestoreWorkflow(t *testing.T) {
    repo := NewPelangganRepository()

    // 1. Create pelanggan
    pelanggan := &models.Pelanggan{
        Nama:    "Test User",
        Telepon: "08123456789",
        Poin:    100,
    }
    err := repo.Create(pelanggan)
    assert.Nil(t, err)
    assert.NotEqual(t, 0, pelanggan.ID)

    // 2. Verify dapat di-get
    result, err := repo.GetByID(pelanggan.ID)
    assert.Nil(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Test User", result.Nama)

    // 3. Delete (soft delete)
    err = repo.Delete(pelanggan.ID)
    assert.Nil(t, err)

    // 4. Verify TIDAK bisa di-get lagi (karena deleted)
    result, err = repo.GetByID(pelanggan.ID)
    assert.Nil(t, err)
    assert.Nil(t, result) // Not found

    // 5. Verify ada di GetDeleted
    deleted, err := repo.GetDeleted()
    assert.Nil(t, err)
    assert.Len(t, deleted, 1)
    assert.Equal(t, pelanggan.ID, deleted[0].ID)

    // 6. 🔄 RESTORE
    err = repo.Restore(pelanggan.ID)
    assert.Nil(t, err)

    // 7. Verify bisa di-get lagi (data aktif kembali)
    result, err = repo.GetByID(pelanggan.ID)
    assert.Nil(t, err)
    assert.NotNil(t, result) // Found!
    assert.Equal(t, "Test User", result.Nama)
    assert.Equal(t, 100, result.Poin) // Data lengkap tetap ada!

    // 8. Verify TIDAK ada lagi di GetDeleted
    deleted, err = repo.GetDeleted()
    assert.Nil(t, err)
    assert.Len(t, deleted, 0) // Empty, karena sudah di-restore
}
```

---

## 🎯 QUICK CHECKLIST

Untuk implement restore lengkap, Anda perlu:

- [ ] ✅ Repository `Restore()` - **SUDAH ADA!**
- [ ] ✅ Repository `GetDeleted()` - **SUDAH ADA!**
- [ ] ⏳ Service `Restore()` - Tambahkan fungsi di service
- [ ] ⏳ Service `GetDeleted()` - Tambahkan fungsi di service
- [ ] ⏳ Wails `RestorePelanggan()` - Tambahkan di app.go (desktop mode)
- [ ] ⏳ Handler `Restore()` - Tambahkan di handler (web mode)
- [ ] ⏳ Router - Tambahkan route POST `/:id/restore`
- [ ] ⏳ Frontend API `restore()` - Tambahkan di API client
- [ ] ⏳ Frontend UI - Buat component untuk lihat & restore data deleted

---

## 🎉 HASIL AKHIR

Setelah implement semua ini, user bisa:

1. **Hapus data** (soft delete) - Data hilang dari list tapi masih di database
2. **Lihat data yang dihapus** - Admin bisa lihat history data yang dihapus
3. **Restore data** - Klik tombol restore, data aktif kembali
4. **Data lengkap tetap utuh** - Poin, transaksi, history semua tetap ada!

**No more "Waduh data hilang!"** ✨

---

Ada pertanyaan lain tentang implementasi restore?
