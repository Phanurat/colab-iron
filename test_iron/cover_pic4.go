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

func randomExcellentBandwidthcover_pic4() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// func getLatestPhotoIDcover_pic4() string {
// 	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
// 	if folder == "" {
// 		folder = "."
// 	}

// 	dbPath := filepath.Join(folder, "cover_photo_id.db")
// 	fmt.Println("📂 DB PATH:", dbPath)

// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
// 	}
// 	defer db.Close()

// 	fmt.Println("📂 DB PATH:", folder+"/cover_photo_id.db")

// 	row := db.QueryRow("SELECT cover_pic_id FROM cover_photo_id_table ORDER BY id DESC LIMIT 1")
// 	var photoID string
// 	err = row.Scan(&photoID)
// 	if err != nil {
// 		fmt.Println("❌ ดึง pic_id ไม่สำเร็จ: " + err.Error())
// 	}
// 	return photoID
// }

func Runcover_pic4(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	//	photoID := getLatestPhotoIDcover_pic4()

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

	var photoID string
	err = db.QueryRow("SELECT pic_id FROM cover_photo_id_table LIMIT 1").Scan(
		&photoID)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล cover_photo_id_table ไม่สำเร็จ: " + err.Error())
		return
	}

	form := url.Values{}
	form.Set("focus_y", "0.49925926")
	form.Set("focus_x", "0.5")
	form.Set("photo", photoID)
	form.Set("no_feed_story", "false")
	form.Set("locale", "en_US")
	form.Set("client_country_code", "TH")
	form.Set("fb_api_req_friendly_name", "set_cover_photo")
	form.Set("fb_api_caller_class", "SetCoverPhotoHandlerImpl")

	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	_, _ = gz.Write([]byte(form.Encode()))
	_ = gz.Close()
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	req, _ := http.NewRequest("POST", "https://graph.facebook.com/me/cover", &compressed)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "graph.facebook.com")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthcover_pic4())
	req.Header.Set("X-FB-Connection-Quality", "EXCELLENT")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-net-hni", netHni) // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "set_cover_photo")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","retry_attempt":"0"},"application_tags":"unknown"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-zero-eh", "2,,ASAySH2oMJyaEZ8_gJqsBkU9XiRTbXAuIxwORpw9LUATP-gBzQrZQZq8gPwqoERuScM")
	req.Header.Set("X-ZERO-STATE", "unknown")

	// ---------- SEND ----------
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

	_, err = db.Exec("DELETE FROM cover_photo_id_table WHERE pic_id = ?", photoID) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ cover_photo_id_table ออกจากฐานข้อมูลแล้ว:", photoID)
	}

}
