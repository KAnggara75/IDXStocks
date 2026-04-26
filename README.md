# IDXStock

A robust Go application to fetch and manage Indonesian stock market data.

## 📊 Data Sources

This application leverages data from the following reliable sources:
- **IDX (Indonesia Stock Exchange)**: [idx.co.id](https://www.idx.co.id)
- **Pasardana**: [pasardana.id](https://pasardana.id)

## 🚀 Features

- **Stock History**: Retrieve historical price data via RESTful API.
- **Dynamic Filtering**: Support for `fields` selection, `start_date`, and `end_date`.
- **Delisting Data**: Track delisted stocks.
- **Containerized**: Ready for deployment with Docker/Podman.
- **Code Coverage**: 7.2% (Target: 75%)

## 🛠 Tech Stack

- **Go**: 1.26.2
- **Framework**: Fiber v3
- **Database**: PostgreSQL (pgx)
- **Logging**: Logrus

## 📝 Planning & Issues

*See [issue.md](issue.md) for current development planning and roadmap.*

### Prompt

buatkan issue.md yang berisi perencanaan untuk nanti di implementasikan oleh junior programmer atau ai model yang lebih murah

isi dari planning nya sebagai berikut

buat endpoint untuk get all histori berdasarkan stock code


jelaskan tahpan-tahapan yangg harus dilakukan untuk implementasikan fitur ini, anggap nanti yang mengimplementasikan adalah junior programmer atau ai model yang lebih murah

---

## 📚 API Documentation

### System
- `GET /health` : Check API and Database connection status.

### Stocks Management
- `POST /api/v1/stocks/upload` : Preview parsed JSON from a stock file upload.
- `PATCH /api/v1/stocks/upload` : Process and save uploaded stock JSON file.
- `PUT /api/v1/stocks/sync` : Sync detailed company profiles and statistics for all active stocks.
- `PUT /api/v1/stocks/delisting/sync` : Fetch and sync delisted stocks from Pasardana.
- `PUT /api/v1/stocks/history/sync` : Download and sync EOD historical price data for stocks.
- `GET /api/v1/stocks/:code/history` : Retrieve EOD historical data based on stock code. Supported query params: `start_date`, `end_date`, `fields`.

### Master Data (Sectors & Industries)
- `PUT /api/v1/sectors/sync` : Sync sectors and sub-sectors from Pasardana search results.
- `PUT /api/v1/industries/sync` : Sync and map industry classifications to the local database.

### Broker Activities
- `GET /api/v1/broker/sync` : Sync daily broker summary and activity transactions using an external token. Supported params: `broker_code`, `from`, `to`, `transaction_type`, `market_board`, `investor_type`.
- `PUT /api/v1/partition/broker-activity` : Maintenance endpoint to automatically calculate, prepare, and allocate PostgreSQL partitions for the next 2 ISO weeks of broker activities.

### Assets & Media
- `GET /api/v1/assets/logos/companies/:code` : Proxy endpoint for retrieving company logo images securely, applying caching and streaming without local storage (e.g., `/api/v1/assets/logos/companies/BBCA.png`).
