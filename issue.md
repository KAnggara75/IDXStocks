# Fitur: Sinkronisasi Data Industri dan Sub-Industri Saham (Pasardana)

## Deskripsi
Tugas ini bertujuan untuk membangun fitur sinkronisasi data **Industri** dan **Sub-Industri** dari API pihak ketiga (Pasardana). Data tersebut akan diambil, diproses untuk menghilangkan duplikasi (karena API mengembalikan data per saham), dan kemudian disimpan/diperbarui (*upsert*) ke dalam database PostgreSQL secara efisien.

## Spesifikasi API

### API Eksternal (Pasardana)
- **URL**: `GET https://www.pasardana.id/api/StockSearchResult/GetAll`
- **Contoh Response JSON**:
  ```json
  [
    {
      "NewSubIndustryId": 47,
      "NewSubIndustryName": "Perkebunan & Tanaman Pangan",
      "NewIndustryId": 22,
      "NewIndustryName": "Produk Makanan Pertanian",
      "...": "..." // Kolom lain abaikan
    }
  ]
  ```
  *(Catatan: Karena API mengembalikan hasil pencarian saham, akan ada banyak entri dengan Industry dan Sub-Industry yang berulang. Data harus dideduplikasi di level kode sebelum di-insert.)*

### API Internal Baru
- **Method**: `PUT`
- **Endpoint**: `/api/v1/industry/sync`
- **Response**: Mengembalikan daftar objek `id` dan `name` dari tabel/entitas yang berhasil di-update. (Anda dapat menyatukan responnya atau mengembalikan struktur gabungan yang merepresentasikan industri dan sub-industri yang tersinkron).
  ```json
  {
    "industries": [
      { "id": 22, "name": "Produk Makanan Pertanian" }
    ],
    "sub_industries": [
      { "id": 47, "name": "Perkebunan & Tanaman Pangan" }
    ]
  }
  ```

---

## Tahapan Implementasi (Wajib Diikuti Secara Berurutan)

### 1. Buat Tabel Database (Migrasi)
- Buat file SQL baru di folder `migrations/`, contoh: `003_create_industry_tables.sql`.
- Buat 2 tabel baru di skema `idxstock`:
  1. **`idxstock.industry`**
     - `id INT PRIMARY KEY`
     - `name VARCHAR(200) NOT NULL`
     - `last_modified TIMESTAMPTZ DEFAULT now()`
  2. **`idxstock.sub_industry`**
     - `id INT PRIMARY KEY`
     - `name VARCHAR(200) NOT NULL`
     - `industry_id INT NOT NULL` (Foreign Key merujuk ke `idxstock.industry(id)`)
     - `last_modified TIMESTAMPTZ DEFAULT now()`
- Pastikan memberikan akses / kepemilikan (*owner*) `pakaiwa_app` pada kedua tabel tersebut.

### 2. Definisikan Model Data (Structs)
- Di dalam folder `internal/models/` (misalnya `industry.go`), buat struct untuk memetakan JSON API Eksternal:
  ```go
  type PasardanaSearchResult struct {
      NewSubIndustryId   int    `json:"NewSubIndustryId"`
      NewSubIndustryName string `json:"NewSubIndustryName"`
      NewIndustryId      int    `json:"NewIndustryId"`
      NewIndustryName    string `json:"NewIndustryName"`
  }
  ```
- Buat juga struct respons minimalis untuk hasil output endpoint internal.

### 3. Ekstraksi dan Pemrosesan Data (Handling Duplikasi)
- Karena respons dari pihak ketiga berbentuk deretan panjang (berpotensi memiliki ribuan objek dengan ratusan id industri yang terduplikasi), aplikasikan teknik *Data Deduplication* *(menggunakan Go `map[int]Struct`)* sebelum memanggil Repository. Pastikan Anda mendaftar ID unit (Unik) untuk batch insert yang lebih ringan.

### 4. Implementasi Repository Layer
- Buat `internal/repositories/industry_repository.go` beserta interfacenya.
- Buat fungsi *batch upsert* memanfaatkan `pgx.Batch` dan klausa `ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name ... RETURNING id, name`.
- Anda memerlukan dua fungsi dasar (atau satu fungsi menggunakan transaksi beruntun):
  - `UpsertIndustries(ctx context.Context, industries []models.Industry) ([]models.BasicResponse, error)`
  - `UpsertSubIndustries(ctx context.Context, subIndustries []models.SubIndustry) ([]models.BasicResponse, error)`
- Pastikan penyimpanan Industri dijalankan lebih dulu sebelum Sub-Industri untuk menghindari constraint *Foreign Key* yang gagal.

### 5. Tambahkan Method HTTP Service
- Di dalam `internal/services/stock_service.go` (atau service baru), buat fungsi `FetchPasardanaStockSearchResult() ([]models.PasardanaSearchResult, error)`.
- Menggunakan `http.Get`, tarik data JSON dan *decode* ke *struct slice* yang disediakan. Jangan lupa menutup `resp.Body` setelahnya.

### 6. Koordinasi pada Usecase Layer
- Pada `internal/usecases/` (buat baru atau satukan ke stock usecase bila relevan/disetujui).
- Logikanya:
  1. *Fetch* data mentah dari Service.
  2. *Extract & Deduplicate* List Unique Industry dan Unique Sub-Industry.
  3. Panggil Repo `UpsertIndustries`.
  4. Panggil Repo `UpsertSubIndustries`.
  5. Bungkus dan return data hasil update-nya.

### 7. Integrasi Handler & Router (API Exposure)
- Di endpoint *handler*, misal `IndustrySyncHandler(c *fiber.Ctx) error`.
- Akses dan map endpoint tersebut di `internal/routes/routes.go` di bawah blok `v1` group dengan verb `PUT`, URL: `/api/v1/industry/sync`.
- Daftarkan inject *dependency* repository yang baru pada setting router *(Setup)*.

### 8. Testing Kelayakan
- Pastikan aplikasi berjalan tanpa error migrasi.
- Tes *Trigger API* menggunakkan *cURL* atau *Postman* untuk route `/api/v1/industry/sync`.
- Verifikasi tabel database bertambah barisnya. Ulangi menembak API tersebut dan pastikan tak ada *error duplikasi* ID *(Idempotent behaviour)*.
