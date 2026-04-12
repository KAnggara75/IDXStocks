# Fitur: Sinkronisasi Data Saham Delisting (Berdasarkan Periode)

## Deskripsi
Fitur ini bertujuan untuk melakukan penarikan data jadwal *delisting* saham secara periodik (berdasarkan bulan dan tahun). Data tersebut diambil dari sistem statistik digital Bursa Efek Indonesia (IDX) dan dimasukkan ke dalam basis data agar status saham tetap mutakhir.

## Spesifikasi API Eksternal (IDX)
- **Method**: `GET`
- **Tujuan**: `https://idx.co.id/primary/DigitalStatistic/GetApiDataPaginated?urlName=LINK_DELISTING&periodYear={year}&periodMonth={month}`
- **Struktur Response**:
  ```json
  {
      "data": [
          {
              "code": "JKSW",
              "DeListingDate": "18 July 2025"
          }
      ]
  }
  ```

## Spesifikasi API Internal Target
- **Method**: `PUT`
- **Endpoint**: `/api/v1/stocks/delisting/sync` *(Disarankan memakai bentuk jamak `stocks` agar konsisten)*
- **Request Body**:
  ```json
  {
      "year": 2025,
      "month": 7
  }
  ```
- **Response**: Menampilkan sekumpulan objek hasil (code, name, dll.) dari saham yang sukses diperbarui datanya.

---

## Tahapan Implementasi (Standard Operating Procedure)

SOP ini diperuntukkan bagi *Junior Programmer* atau *AI Assistant* dalam menyelesaikan fitur secara bertahap dan terstruktur tanpa merusak kode / *flow* yang sudah berjalan baik.

### 1. Persiapan Struktur Data (Model Layer)
- Deklarasikan struct untuk menampung HTTP JSON Request Body pada bagian *handler* atau *model*.
  ```go
  type SyncDelistingRequest struct {
      Year  int `json:"year" validate:"required"`
      Month int `json:"month" validate:"required"`
  }
  ```
- Deklarasikan struct untuk *parsing* data respons dari API IDX.

### 2. Implementasi Client API IDX (Service Layer)
- Buat atau tambahkan method baru di service (bisa membuat `idx_service.go` baru atau menggabungkannya ke service serba guna).
  - Contoh fungsi: `FetchDelistedStocks(year, month int) ([]models.IdxDelistedStock, error)`.
- **Waspada Blokir WAF**: Seringkali Endpoint dari domain `idx.co.id` memblokir *request* yang dilakukan oleh *script* tanpa identitas. Pastikan menambahkan *Header* minimal seperti `User-Agent` (menyerupai *browser*) ke dalam *http request* Anda.

### 3. Parsing dan Formatting Tanggal (Helper)
- API IDX mengembalikan tipe string seperti `"18 July 2025"`.
- Kolom tabel `idxstock.stocks` eksisting untuk `delisting_date` menggunakan tipe data `DATE`.
- Buat *util function* yang bertugas mengonversi `"18 July 2025"` (format Indonesia Timur/Bahasa Inggris lokal) agar secara solid terkonversi menjadi format standard SQL `YYYY-MM-DD`.

### 4. Konstruksi Logika Bisnis (Usecase Layer)
- Buka `internal/usecases/stock_usecase.go` dan tambahkan fungsionalitas yang menjembatani Controller dan Repository.
  - `SyncDelistingStocks(ctx context.Context, year, month int) ([]models.StockResponse, error)`
- Alur Kerjanya:
  1. Validasi parameter masukan (jika dirasa validasi level Usecase dibutuhkan).
  2. Panggil *Service Layer* `FetchDelistedStocks(year, month)`.
  3. Lakukan konversi (parsing tanggal) dari array *response* API. Apabila format *invalid*, cukup log *warning* dan lompati *record* tersebut.
  4. Panggil *Repository Layer* untuk melakukan *batch update*.

### 5. Penulisan Query Update Performa Tinggi (Repository Layer)
- Masuk ke `internal/repositories/stock_repository.go` dan buat fungsi untuk melakukan penyimpanan (*Update*).
- Hindari *looping query* `UPDATE table SET ... WHERE code = ...` secara satu-per-satu karena sangat membebani I/O *database*.
- Kerjakan *Bulk Update*. Anda dapat memanfaatkan kemampuan PostgreSQL dengan sintaks seperti:
  ```sql
  UPDATE idxstock.stocks AS s
  SET delisting_date = tmp.delisting_date::DATE, 
      last_modified = now()
  FROM (VALUES 
     ('JKSW', '2025-07-18'),
     ('CODE', 'YYYY-MM-DD')
  ) AS tmp(code, delisting_date)
  WHERE s.code = tmp.code;
  ```
  Atau gunakan pgx Batch Queue `UPDATE idxstock.stocks SET delisting_date = $1 WHERE code = $2;` secara dibungkus *Transaction* batching. 
- Harus mengembalikan baris yang terpengaruh (RETURNING clauses).

### 6. Endpoint Routing & Response Handler
- Buat `SyncDelistedHandler` di dalam `internal/handlers/stock_handler.go`.
- Tangkap payload body menggunakan `c.BodyParser(&request)`.
- Serahkan beban eksekusi kepada *Usecase*.
- Registrasi alamat url: `v1.Put("/stocks/delisting/sync", stockHandler.SyncDelistedHandler)` pada blok rute di `internal/routes/routes.go`.
- Jalankan pengecekan kompilasi `go build ./...` dan uji skenario berhasil/gagal via REST client.
