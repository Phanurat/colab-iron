package main

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func randomExcellentBandwidthcover_pic2() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// func getLatestPhotoIDcover_pic2() string {
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

func generateUploadIDscover_pic2() (string, string, string) {
	sessionID := uuid.New().String()
	random := time.Now().UnixNano() % 1000000000
	idempotence := fmt.Sprintf("%s_%d_0", sessionID, random)
	waterfallID := fmt.Sprintf("%s_996C60C8D20C_Mixed_0", sessionID)
	return sessionID, idempotence, waterfallID
}

func Runcover_pic2(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	host := "rupload.facebook.com"
	uuidPart := uuid.New().String()[:32]        // เอาแค่ 32 ตัวแรก
	randPart := strconv.Itoa(rand.Intn(100000)) // เช่น 67846
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	uploadID := fmt.Sprintf("%s-0-%s-%s-%s", uuidPart, randPart, timestamp, timestamp)
	url := "https://" + host + "/fb_photo/" + uploadID
	//photoID := getLatestPhotoIDcover_pic2()

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

	composerSessionID, idempotenceToken, waterfallID := generateUploadIDscover_pic2()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("❌ Request build fail: " + err.Error())
		return
	}

	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("add_to_photo_id", photoID)
	req.Header.Set("audience_exp", "true")
	req.Header.Set("composer_session_id", composerSessionID)
	req.Header.Set("idempotence_token", idempotenceToken)
	req.Header.Set("published", "false")
	req.Header.Set("qn", composerSessionID)
	req.Header.Set("Segment-Start-Offset", "0")
	req.Header.Set("Segment-Type", "3")
	req.Header.Set("target_id", userID)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthcover_pic2())
	req.Header.Set("X-FB-Connection-Quality", "EXCELLENT")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "Resumable-Upload-Get")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","retry_attempt":"0"},"application_tags":"unknown"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-zero-eh", "2,,ASAySH2oMJyaEZ8_gJqsBkU9XiRTbXAuIxwORpw9LUATP-gBzQrZQZq8gPwqoERuScM")
	req.Header.Set("X_FB_PHOTO_WATERFALL_ID", waterfallID)
	req.Header.Set("Zero-Rated", "0")
	req.Header.Set("x-fb-net-hni", netHni) // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)

	// ---------- SEND ----------
	bw := tlsConns.RWrupload.Writer
	br := tlsConns.RWrupload.Reader

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
