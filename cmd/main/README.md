# İnsani Yardım Web Sitesi Projesi

Bu proje, Go dili ve Gin framework'ü kullanılarak geliştirilmiş bir insani yardım kampanyası sitesidir.

## Özellikler

- PostgreSQL veritabanı ile kampanya yönetimi
- GORM ile kolay veritabanı etkileşimi
- Herkese açık kampanya listeleme ve detay sayfaları
- Basic Auth ile korunan Admin Paneli
- Tam CRUD (Create, Read, Update, Delete) işlevselliği

## Projeyi Çalıştırma

Bu projeyi kendi bilgisayarında çalıştırmak için aşağıdaki adımları izle:

### Gereksinimler

- [Go](https://go.dev/dl/) (1.20 veya üstü)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

### Adımlar

1.  **Veritabanını Başlat:**
    Proje ana dizinindeyken terminali aç ve aşağıdaki komutu çalıştır. Bu komut, Docker üzerinde PostgreSQL veritabanını arka planda başlatacaktır.
    ```bash
    docker-compose up -d
    ```

2.  **Uygulamayı Başlat:**
    Aynı terminalde, aşağıdaki komutla Go web sunucusunu başlat.
    ```bash
    go run cmd/main/main.go
    ```

3.  **Siteye Eriş:**
    - Herkese açık site için: `http://localhost:8080`
    - Admin paneli için: `http://localhost:8080/admin/dashboard`

### Admin Giriş Bilgileri

- **Kullanıcı Adı:** `admin`
- **Şifre:** `12345`

---