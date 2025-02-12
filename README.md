# Go Fiber API Project

Project ini adalah REST API yang dibangun menggunakan Go Fiber framework.

## 🚀 Fitur

- RESTful API endpoints
- JWT Authentication
- Database integration
- Environment configuration
- Middleware support
- Password hashing
- Hot reload development

## 📋 Prasyarat

Sebelum memulai, pastikan sudah terinstall:

- Go (versi 1.16 atau lebih baru)
- MySQL/PostgreSQL
- Air (untuk hot reload) - opsional

## 🛠️ Instalasi

1. Clone repository

```
git clone <repository-url>
cd <project-name>
```

2. Install dependencies
   `go mod tidy`

3. Setup environment variables
   `cp .env.example .env`

Sesuaikan konfigurasi di file `.env`

## 🚀 Menjalankan Aplikasi

### Development (dengan hot reload)

`air`

### Production

`go run main.go`
