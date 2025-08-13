import requests
import uuid
import random
import time
import gzip
import json
from io import BytesIO

def random_bandwidth():
    return str(random.randint(20000000, 35000000))

def generate_session_id(cid):
    return f"nid=TSRsHSL+wunc;tid=206;nc=0;fc=0;bc=0;cid={cid}"

def generate_hex32():
    return ''.join(random.choice('abcdef0123456789') for _ in range(32))

def generate_trace_id():
    return str(uuid.uuid4())

def check_facebook_account(access_token, user_agent, proxy, device_group, net_hni, sim_hni):
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

    try:
        response = requests.post(
            "https://graph.facebook.com/graphql",
            headers=headers,
            data=data,
            proxies=proxies,
            timeout=10
        )

        # ตรวจสอบว่าตอบกลับมาแบบ gzip หรือไม่
        if response.headers.get("Content-Encoding") == "gzip":
            buf = BytesIO(response.content)
            f = gzip.GzipFile(fileobj=buf)
            body = f.read()
        else:
            body = response.content

        body_json = json.loads(body)
        print("✅ Response:", json.dumps(body_json, indent=2))

        if "error" in body_json:
            print("❌ Token น่าจะหมดอายุ หรือถูกแบน:", body_json["error"]["message"])
            return "❌ Token น่าจะหมดอายุ หรือถูกแบน:", body_json["error"]["message"]
        elif "data" in body_json:
            print("✅ Token ยังใช้งานได้")
            return "✅ Token ยังใช้งานได้"
        else:
            print("⚠️ ไม่สามารถตีความผลลัพธ์ได้")
            return "⚠️ ไม่สามารถตีความผลลัพธ์ได้"

    except requests.exceptions.RequestException as e:
        print("❌ Request ล้มเหลว:", e)
        return False

# ===== Example Usage =====
# Call the function with necessary arguments
data = check_facebook_account(
    access_token="EAAAAUaZA8jlABPALa8RgOFdOvpOYfaB2mI505LJbrrsMfdQqyAs5XOoJcgbd3l7rZA3Q48k2yZBfCqqvY3lWbZBYF0nrYHxqq6W1I7Kxwn1xTd0arCpLoDc3OO0Gfk7rrBN9FZCUZAdWdL7CU84l5L0PGM4TnT6TZBLHZBY8MXpiw9vDplp9OM8Ny7r5ta6cCUZCUssYVQQZDZD",  # token จริง
    user_agent="[FBAN/FB4A;FBAV/424.0.0.26.84;FBBV/522888777;FBDM/{density=3.0,width=1080,height=2400};FBLC/th_TH;FBRV/522444333;FBCR/TRUE-H;FBMF/vivo;FBBD/vivo;FBPN/com.facebook.katana;FBDV/Y20;FBSV/11;FBOP/1;FBCA/arm64-v8a:]",
    proxy="http://gamezaba:JQk2ywNFd7@51.79.228.210:11868",  # หรือ None
    device_group="2573",
    net_hni="52001",
    sim_hni="52001"
)

# Use the response
print("Account check result:", data)
