import os
import sqlite3
import requests
import uuid
import random
import time
import gzip
import json
from io import BytesIO
from datetime import datetime
from zoneinfo import ZoneInfo
from requests.exceptions import RequestException

def random_bandwidth():
    return str(random.randint(20000000, 35000000))

def generate_session_id(cid):
    return f"nid=TSRsHSL+wunc;tid=206;nc=0;fc=0;bc=0;cid={cid}"

def generate_hex32():
    return ''.join(random.choice('abcdef0123456789') for _ in range(32))

def generate_trace_id():
    return str(uuid.uuid4())

def check_facebook_account(access_token, user_agent, proxy, device_group, net_hni, sim_hni, max_retries=3):
    conn_token = generate_hex32()
    session_id = generate_session_id(conn_token)
    trace_id = generate_trace_id()

    headers = {
        "Authorization": f"OAuth {access_token}",
        "User-Agent": user_agent,
        "Content-Type": "application/x-www-form-urlencoded",
        "Accept-Encoding": "gzip, deflate",
        "X-FB-Friendly-Name": "FetchMessengerJewelCount",
        "x-fb-connection-token": conn_token,
        "x-fb-session-id": session_id,
        "x-fb-net-hni": net_hni,
        "x-fb-sim-hni": sim_hni,
        "x-fb-device-group": device_group,
        "x-fb-connection-bandwidth": random_bandwidth(),
        "x-fb-connection-quality": "EXCELLENT",
        "x-fb-background-state": "1",
        "X-FB-HTTP-Engine": "Liger",
        "X-FB-Request-Analytics-Tags": '{"network_tags":{"product":"350685531728"}}',
    }

    data = {
        "method": "post",
        "pretty": "false",
        "format": "json",
        "fb_api_req_friendly_name": "FetchMessengerJewelCount",
        "client_doc_id": "232448440414169222349211474621",
        "variables": "{}"
    }

    proxies = {
        "http": proxy,
        "https": proxy
    }

    attempt = 0
    while attempt < max_retries:
        try:
            response = requests.post(
                "https://graph.facebook.com/graphql",
                headers=headers,
                data=data,
                proxies=proxies,
                timeout=10
            )

            # ตรวจสอบ gzip content
            try:
                if response.headers.get("Content-Encoding") == "gzip":
                    try:
                        buf = BytesIO(response.content)
                        f = gzip.GzipFile(fileobj=buf)
                        body = f.read()
                    except OSError:
                        print("⚠️ Gzip decode failed. Using raw content.")
                        body = response.content
                else:
                    body = response.content

                body_json = json.loads(body)
            except Exception as e:
                return f"❌ Failed to parse response JSON: {e}"

            print("📨 Raw Response:", json.dumps(body_json, indent=2))

            if "error" in body_json:
                return f"❌ Token น่าจะหมดอายุ หรือถูกแบน: {body_json['error']['message']}"
            elif "data" in body_json:
                return "✅ Token ยังใช้งานได้"
            else:
                return "⚠️ ไม่สามารถตีความผลลัพธ์ได้"

        except RequestException as e:
            attempt += 1
            print(f"⚠️ Attempt {attempt}: Request failed with {e}")
            if attempt < max_retries:
                print("🔁 Retrying in 2 seconds...")
                time.sleep(2)
            else:
                print("❌ Failed after 3 attempts.")
                return f"❌ Request ล้มเหลวหลังจากพยายาม {max_retries} ครั้ง: {e}"

# === DATABASE PART ===
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
BASE_DIR = os.path.abspath(os.path.join(BASE_DIR, ".."))
db_files = {}

def scan_dbs():
    for folder in os.listdir(BASE_DIR):
        path = os.path.join(BASE_DIR, folder, "fb_comment_system.db")
        if os.path.isfile(path):
            db_files[folder] = path
            print(f"📦 Found database: {path}")

def check_db():
    valid_count = 0
    invalid_count = 0

    for folder, db_path in db_files.items():
        try:
            conn = sqlite3.connect(db_path)
            cursor = conn.cursor()
            cursor.execute("SELECT name FROM sqlite_master WHERE type='table'")
            tables = [t[0] for t in cursor.fetchall()]

            if "app_profiles" in tables:
                cursor.execute("SELECT * FROM app_profiles")
                profiles = cursor.fetchall()
                if not profiles:
                    print(f"⚠️ No profiles in {folder}")
                    continue

                for row in profiles:
                    try:
                        actor_token = row[1]
                        access_token = row[2]
                        user_agent = row[3]

                        cursor.execute("SELECT * FROM proxy_table LIMIT 1")
                        proxy_row = cursor.fetchone()
                        if not proxy_row:
                            print(f"⚠️ No proxy found for {actor_token}")
                            continue

                        proxy = f"http://{proxy_row[1]}"
                        if not proxy.startswith("http://"):
                            print(f"⚠️ Invalid proxy: {proxy}")
                            continue

                        print(f"🔍 Checking token for actor: {actor_token}")
                        result = check_facebook_account(
                            access_token,
                            user_agent,
                            proxy,
                            device_group="2573",
                            net_hni="52001",
                            sim_hni="52001"
                        )

                        print("📣 ตรวจสอบแล้ว:", result)

                        bangkok_time = datetime.now(ZoneInfo("Asia/Bangkok"))

                        log_res = requests.post(
                            "http://127.0.0.1:5050/api/insert/check-acc",
                            params={
                                "actor_id": actor_token,
                                "access_token": access_token,
                                "proxy": proxy,
                                "response": result,
                                "log": "ok",
                                "timestamp": bangkok_time
                            },
                            timeout=5
                        )

                        if log_res.status_code == 200:
                            print("✅ Logged to API")
                        else:
                            print(f"⚠️ Logging failed: {log_res.status_code}")

                    except Exception as e:
                        print(f"❌ Error checking profile in {folder}: {e}")
                        continue

                valid_count += 1
            else:
                print(f"❌ Invalid database: no 'app_profiles' table in {folder}")
                invalid_count += 1

        except sqlite3.Error as e:
            print(f"❌ SQLite error in {folder}: {e}")
            invalid_count += 1
        finally:
            conn.close()

    print(f"\n🔍 Total scanned: {len(db_files)} databases")
    print(f"✅ Valid: {valid_count}")
    print(f"❌ Invalid: {invalid_count}")

# === MAIN EXECUTION ===
if __name__ == "__main__":
    try:
        res = requests.delete("http://127.0.0.1:5050/api/delete/check-acc")
        if res.status_code == 200:
            print("🧹 Cleared old check logs.")
        else:
            print(f"⚠️ Could not clear logs: {res.status_code}")
    except Exception as e:
        print(f"⚠️ Failed to connect to clear logs: {e}")

    scan_dbs()
    check_db()
