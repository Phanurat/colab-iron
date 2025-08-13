import requests
from urllib.parse import urlparse, parse_qs

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
                return f"✅ ลิงก์จริง (ไม่สามารถแปลงให้สวย): {real_url}"
        else:
            return f"✅ ลิงก์จริง: {real_url}"
    except Exception as e:
        return f"❌ Error: {e}"

# ใช้งาน
short_link = "https://www.facebook.com/share/p/19Zbdi2L8C/?mibextid=wwXIfr"
print(resolve_fb_shortlink(short_link))
