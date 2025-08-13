import requests
import time
from datetime import datetime, timezone
from urllib.parse import urlparse, parse_qs


# 🔗 URL จาก Google Apps Script
url = "https://script.googleusercontent.com/macros/echo?user_content_key=AehSKLh3Jsu5Smn9d1t-tN-ipe25kTJwWB53auv1tXfatmAe_s7OAvL9dZvWGjQYLk9foxYnr4Ig_hh59nHJIiyTxbBhJ8rGGW1ufYlsW-INXekM6VWfL0TOMizTAhLtiVSW73ml0vJoHOIAwP_h85t8TvXb-sm4roqWVMCM_Pl_gOj_GK0O30EB9CZSysIQSryq0W90a985Em3tqyhEE13Zb7J8LUnhe-azZQWJQlR4r7r7qke6q6j_NQr8xeenfPtPlqSzt-Mc3leFg_6Nw2J9Y-2dT0NJi7IoXs5a4mR2&lib=Mx3V237tqGGRF5UjJv2nSnh-6qQotOF5P"

link_buffer = []
start_time = None  # เวลาที่เริ่มดึงลิงก์ (ตอนรันสคริปต์)

def resolve_fb_shortlink(short_url):
    try:
        res = requests.get(short_url, allow_redirects=False)
        if 'Location' not in res.headers:
            return "❌ ไม่พบ redirect"
        real_url = res.headers['Location']

        parsed = urlparse(real_url)
        if parsed.path == "/story.php":
            qs = parse_qs(parsed.query)
            story_fbid = qs.get("story_fbid", [None])[0]
            page_id = qs.get("id", [None])[0]
            if story_fbid and page_id:
                # สร้างลิงก์เต็ม
                return f"https://www.facebook.com/{page_id}/posts/{story_fbid}"
            else:
                return f"{real_url}"
        else:
            return f"{real_url}"
    except Exception as e:
        return f"❌ Error: {e}"



def fetch_links():
    try:
        response = requests.get(url)
        response.raise_for_status()
        data = response.json()
        for item in data:
            item["dt"] = datetime.fromisoformat(item["Timestamp"].replace("Z", "+00:00"))
        return sorted(data, key=lambda x: x["dt"], reverse=True)
    except Exception as e:
        print("❌ ดึงข้อมูลล้มเหลว:", e)
        return []

def update_buffer():
    global link_buffer
    latest_data = fetch_links()
    if not latest_data:
        return

    # กรองเฉพาะลิงก์ที่ Timestamp > เวลารัน และยังไม่อยู่ใน buffer
    new_items = [
        item for item in latest_data
        if item["dt"] > start_time and item["Link"] not in [i["Link"] for i in link_buffer]
    ]

    if new_items:
        link_buffer.extend(new_items)
        print(f"\n🆕 พบลิงก์ใหม่ {len(new_items)} รายการ:")
        for item in new_items:
            topic = str(item.get("Topic", "")).strip()
            like = item.get("Like", "")
            comment = item.get("Comment", "")
            link = item["Link"]
            link = resolve_fb_shortlink(link)  # แปลงลิงก์สั้นเป็นลิงก์จริง
            timestamp = item["dt"].isoformat()

            print(f"➤ {timestamp} | {item['UserId']} ➜ {link}")

            if topic and like and comment:
                try:
                    payload = {
                        "topic": topic,
                        "link": link,
                        "reaction": "like",
                        "likeValue": 1,
                        "commentValue": int(comment),
                        "timestamp": timestamp,
                        "status": "like_only",
                        "log": "unused",
                        "status_code": "wait..."
                    }

                    response = requests.post("http://localhost:5000/api/insert/news", json=payload)
                    if response.status_code == 200:
                        print(f"✅ INSERT: {link}")
                    else:
                        print(f"❌ ERROR {response.status_code}: {response.text}")
                except Exception as e:
                    print(f"❌ Exception: {e}")
            else:
                print(f"⚠️ ข้ามลิงก์ (ข้อมูลไม่ครบ): {link}")
    else:
        print("⏳ ไม่มีลิงก์ใหม่")

if __name__ == "__main__":
    start_time = datetime.now(timezone.utc)
    print(f"🚀 เริ่มดึงลิงก์แบบ Realtime ตั้งแต่ {start_time.isoformat()}")

    while True:
        update_buffer()
        time.sleep(5)
