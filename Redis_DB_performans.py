import time
import requests

BASE_URL = "http://localhost:5001"

def measure_response_time(endpoint, count=10):
    durations = []

    for i in range(count):
        start = time.time()
        response = requests.get(f"{BASE_URL}/{endpoint}")
        duration = time.time() - start

        if response.status_code == 200:
            durations.append(duration)
        else:
            print(f"❌ {endpoint} isteği başarısız (status: {response.status_code})")

    avg_duration = sum(durations) / len(durations)
    return avg_duration, durations

if __name__ == "__main__":
    print("🚀 Performans Testi Başlıyor...\n")

    db_avg, db_durations = measure_response_time("locationsDB")
    print(f"🗃️  PostgreSQL (/locationsDB) Ortalama Süre: {db_avg:.6f} saniye")

    redis_avg, redis_durations = measure_response_time("locations")
    print(f"⚡ Redis (/locations) Ortalama Süre: {redis_avg:.6f} saniye")

    print("\n🔍 Karşılaştırma:")
    if redis_avg < db_avg:
        print("✅ Redis daha hızlı.")
    else:
        print("⚠️ PostgreSQL daha hızlı çıktı (beklenmeyen durum).")

