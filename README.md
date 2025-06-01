# 🧪 API Test Scripti

Bu Python scripti, `http://localhost:8080` adresinde çalışan bir RESTful API'yi test etmek amacıyla hazırlanmıştır. Amaç, ürün, konum, menşei ve satış noktaları gibi veriler üzerinde CRUD işlemlerinin otomatik olarak test edilmesini sağlamaktır.

## 🚀 İşlem Akışı

### 1. 🟢 CREATE – Veri Oluşturma

- `POST /origins`: Yeni bir menşei kaydı oluşturur (örneğin: Türkiye - Anadolu).
- `POST /locations`: Yeni bir konum (adres, şehir, ülke) oluşturur.
- `POST /sales-points`: Yeni bir satış noktası oluşturur.
- `POST /products`: Yukarıdaki verileri kullanarak yeni bir ürün oluşturur.

### 2. 📥 READ – Veri Listeleme ve Getirme

- `GET /origins`, `GET /locations`, `GET /sales-points`, `GET /products`: Tüm kayıtları listeler.
- `GET /<endpoint>/<id>`: Belirli bir ID'ye sahip kaydı getirir.

### 3. ✏️ UPDATE – Veri Güncelleme

- `PUT /products/<id>`: Ürün bilgilerini (ad, fiyat, stok vb.) günceller.

### 4. ❌ DELETE – Veri Silme

- `DELETE /products/<id>`: Ürünü sistemden siler.

### 5. ✅ DOĞRULAMA – Silinen Ürünü Tekrar Getirme

- `GET /products/<id>`: Silinen ürün sorgulanır, beklenen sonuç `404 Not Found`.

---

## 🔧 Fonksiyonlar

| Fonksiyon              | Açıklama                                         |
|------------------------|--------------------------------------------------|
| `post_origin()`        | Yeni bir origin (menşei) oluşturur.              |
| `post_location()`      | Yeni bir konum (şehir/adres) oluşturur.          |
| `post_sales_point()`   | Yeni bir satış noktası oluşturur.                |
| `post_product()`       | Ürün oluşturur.                                  |
| `get_all(path)`        | Verilen endpoint için tüm kayıtları getirir.     |
| `get_single(path, id)` | Verilen endpoint ve ID için tek bir kayıt getirir.|
| `put_item(path, id, data)` | Belirtilen kaydı günceller.               |
| `delete_item(path, id)`| Belirtilen kaydı siler.                          |

---

## ▶️ Kullanım

```bash
python test_api.py
