from selenium import webdriver
from selenium.webdriver.chrome.options import Options

options = Options()
options.add_argument("--headless")  # ถ้าต้องการรันแบบไม่แสดง browser
driver = webdriver.Chrome(options=options)

# จากนั้นเข้า profile URL
driver.get("https://www.facebook.com/profile.php?id=61577125341551")

html = driver.page_source
print(html)
