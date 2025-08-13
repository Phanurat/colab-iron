import requests
from datetime import datetime
import time

# URL Backend ของคุณ
API_ENDPOINT = "http://localhost:5000/api/insert/news"

# URL API ดึงข้อมูลลิงก์ (Google Apps Script)
SOURCE_URL = "https://script.googleusercontent.com/macros/echo?user_content_key=AehSKLh3Jsu5Smn9d1t-tN-ipe25kTJwWB53auv1tXfatmAe_s7OAvL9dZvWGjQYLk9foxYnr4Ig_hh59nHJIiyTxbBhJ8rGGW1ufYlsW-INXekM6VWfL0TOMizTAhLtiVSW73ml0vJoHOIAwP_h85t8TvXb-sm4roqWVMCM_Pl_gOj_GK0O30EB9CZSysIQSryq0W90a985Em3tqyhEE13Zb7J8LUnhe-azZQWJQlR4r7r7qke6q6j_NQr8xeenfPtPlqSzt-Mc3leFg_6Nw2J9Y-2dT0NJi7IoXs5a4mR2&lib=Mx3V237tqGGRF5UjJv2nSnh-6qQotOF5P"

def fetch_links():
    try:
        res = requests.get(SOURCE_URL)
        res.raise_for_status()
        data = res.json()
        for item in data:
            item["dt"] = datetime.fromisoformat(item["Timestamp"].replace("Z", "+00:00"))
        return data
    except Exception as e:
        print(f"❌ ดึงลิงก์ล้มเหลว: {e}")
        return []

def send_to_backend(item):
    try:
        topic = str(item.get("Topic", "")).strip()
        link = item.get("Link", "").strip()
        reaction = "like"
        like_value = int(item.get("Like") or 1)
        comment_value = int(item.get("Comment") or 0)
        timestamp = item.get("Timestamp")
        status = "like_only"
        log = "unused"
        status_code = "wait..."

        payload = {
            "topic": topic,
            "link": link,
            "reaction": reaction,
            "likeValue": like_value,
            "commentValue": comment_value,
            "timestamp": timestamp,
            "status": status,
            "log": log,
            "status_code": status_code
        }

        res = requests.post(API_ENDPOINT, json=payload)
        if res.status_code == 200:
            print(f"✅ POST: {link}")
        else:
            print(f"❌ ERROR {res.status_code}: {res.text}")

    except Exception as e:
        print(f"❌ POST Exception: {e}")

# ✅ ดึงลิงก์ทั้งหมดแล้วส่งเข้า backend
if __name__ == "__main__":
    print("🚀 เริ่มส่งข้อมูลทั้งหมดแบบ POST JSON...")
    data = fetch_links()

    for item in data:
        send_to_backend(item)
        time.sleep(0.1)  # เพิ่ม delay เล็กน้อยกันโหลดหนัก
