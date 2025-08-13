package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func randomExcellentBandwidthchange_profile1_details() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// ================= CONFIG ====================
var (
	//	token     = "EAAAAUaZA8jlABOZBf5ZAtsPTMPIO3fxCr0jBK4qZCZBXffYZBEQzDtLJZA3tfXA0IKs4w6FoI715YNZCr3ZCURErwCZAwuryz2TZBRHS0Orb11KC5no5vFyOIdfrWJRrpTrQVSGuNXmGSQa3Fv8z2io0ftToBZBZBWs96fnavAycMVzYJrqDkB6xkquRTru0SQhGCUz3kM0BPxAZDZD" // <- ใส่ token เอง
	//	proxy     = ""                                                                                                                                                                                                                     // เช่น "127.0.0.1:8888"
	//	userAgent = "[FBAN/FB4A;FBAV/443.0.0.23.229;FBBV/543547945;FBDM={density=2.625,width=1080,height=1920};FBLC=en_US;FBRV/546817856;FBCR=;FBMF=samsung;FBBD=samsung;FBPN=com.facebook.katana;FBDV/SM-J730G;FBSV/9;FBOP/1;FBCA/arm64-v8a:;]"
	hostchange_profile1_details = "graph.facebook.com"
)

// =============================================

func Runchange_profile1_details(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	rand.Seed(time.Now().UnixNano())

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "fb_comment_system.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
		return
	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/fb_comment_system.db")

	var accessToken, userID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	// ====== FORM BODY (gzip'ed form-urlencoded) ======
	form := url.Values{}
	form.Set("fields", "resource{download_url,delta_download_url,uncompressed_file_sha256_checksum,uncompressed_file_size,compression_format,dod_version_number,js_segment_hash}")
	form.Set("native_build", "543547945")
	form.Set("ota_build", "546817856")
	form.Set("resource_name", "main.jsbundle")
	form.Set("resource_flavor", "hbc-seg-1047")
	form.Set("prefer_compressed", "true")
	form.Set("locale", "en_US")
	form.Set("client_country_code", "TH")
	form.Set("method", "GET")
	form.Set("fb_api_req_friendly_name", "get_on_demand_resource_metadata")
	form.Set("fb_api_caller_class", "Fb4aGraphApiDownloader")

	encoded := form.Encode()
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte(encoded))
	gzipWriter.Close()

	// ====== HTTP REQUEST ======
	req, _ := http.NewRequest("POST", "https://"+hostchange_profile1_details+"/ota_resource", &buf)
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostchange_profile1_details)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthchange_profile1_details())
	req.Header.Set("X-FB-Connection-Quality", "EXCELLENT")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "get_on_demand_resource_metadata")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","retry_attempt":"0"},"application_tags":"unknown"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("Zero-Rated", "0")
	req.Header.Set("x-fb-net-hni", netHni) // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni) // เพิ่มเข้าไป
	//เพิ่มเข้าไป

	// ---------- SEND ----------
	bw := tlsConns.RWGraph.Writer
	br := tlsConns.RWGraph.Reader

	err = req.Write(bw)
	if err != nil {
		fmt.Println("❌ Write fail: " + err.Error())
		return

	}
	bw.Flush() // ✅ ต้อง flush เพื่อให้ข้อมูลถูกส่งออกจริง ๆ

	// ✅ ใช้ reader ตัวเดียวกับที่รับมาจาก utls
	resp, err := http.ReadResponse(br, req)
	if err != nil {
		fmt.Println("❌ Read fail: " + err.Error())
		return

	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("❌ GZIP decompress fail: " + err.Error())
			return

		}
		defer reader.Close()
	} else {
		reader = resp.Body
	}

	bodyResp, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("❌ Body read fail: " + err.Error())
		return

	}

	fmt.Println("✅ Status:", resp.Status)
	fmt.Println("📦 Response:", string(bodyResp))
}
