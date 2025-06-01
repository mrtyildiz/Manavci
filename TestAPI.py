import requests
from datetime import datetime, timezone

BASE_URL = "http://localhost:8080"

def post_origin():
    data = {"origin_name": "Türkiye", "description": "Anadolu"}
    res = requests.post(f"{BASE_URL}/origins", json=data)
    print("✅ POST /origins", res.status_code, res.json())
    return res.json().get("origin_id")

def post_location():
    data = {"address": "İstiklal Cad. No:5", "city": "İstanbul", "country": "Türkiye"}
    res = requests.post(f"{BASE_URL}/locations", json=data)
    print("✅ POST /locations", res.status_code, res.json())
    return res.json().get("location_id")

def post_sales_point():
    data = {"name": "Merkez Şube", "address": "Kızılay Ankara"}
    res = requests.post(f"{BASE_URL}/sales-points", json=data)
    print("✅ POST /sales-points", res.status_code, res.json())
    return res.json().get("sales_point_id")

def post_product(origin_id, location_id, sales_point_id):
    data = {
        "product_name": "Test Ürün",
        "price": 99.99,
        "stock": 10,
        "production_date": datetime.now(timezone.utc).isoformat(),
        "expiration_date": datetime.now(timezone.utc).replace(year=datetime.now().year + 1).isoformat(),
        "origin_id": origin_id,
        "current_location_id": location_id,
        "sales_point_id": sales_point_id
    }
    res = requests.post(f"{BASE_URL}/products", json=data)
    print("✅ POST /products", res.status_code, res.json())
    if res.status_code == 201:
        return res.json().get("product_id")
    else:
        return None

def get_all(path):
    res = requests.get(f"{BASE_URL}/{path}")
    print(f"📥 GET /{path}", res.status_code)
    print(res.json())

def get_single(path, id):
    res = requests.get(f"{BASE_URL}/{path}/{id}")
    print(f"📥 GET /{path}/{id}", res.status_code)
    print(res.json())

def put_item(path, id, data):
    res = requests.put(f"{BASE_URL}/{path}/{id}", json=data)
    print(f"✏️ PUT /{path}/{id}", res.status_code)
    print(res.json())

def delete_item(path, id):
    res = requests.delete(f"{BASE_URL}/{path}/{id}")
    print(f"❌ DELETE /{path}/{id}", res.status_code)
    print(res.json())

if __name__ == "__main__":
    print("🚀 API Test Başlıyor...\n")

    # 1. CREATE
    origin_id = post_origin()
    location_id = post_location()
    sales_point_id = post_sales_point()
    product_id = post_product(origin_id, location_id, sales_point_id)

    # 2. READ ALL
    get_all("origins")
    get_all("locations")
    get_all("sales-points")
    get_all("products")

    # 3. READ SINGLE
    get_single("origins", origin_id)
    get_single("locations", location_id)
    get_single("sales-points", sales_point_id)
    if product_id:
        get_single("products", product_id)

        # 4. UPDATE
        updated_product = {
            "product_id": product_id,
            "product_name": "Güncellenmiş Ürün",
            "price": 123.45,
            "stock": 20,
            "production_date": datetime.now(timezone.utc).isoformat(),
            "expiration_date": datetime.now(timezone.utc).replace(year=datetime.now().year + 1).isoformat(),
            "origin_id": origin_id,
            "current_location_id": location_id,
            "sales_point_id": sales_point_id
        }
        put_item("products", product_id, updated_product)

        # 5. DELETE
        delete_item("products", product_id)

        # 6. TEKRAR GET DENEMESİ (404 beklenir)
        get_single("products", product_id)
    else:
        print("❌ Ürün oluşturulamadı, update/delete işlemleri atlandı.")
