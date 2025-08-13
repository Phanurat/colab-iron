package main

import (
	"bytes"
	"compress/gzip"
	"crypto/rand" // ← เพิ่มอันนี้
	"database/sql"
	"fmt"
	"io"
	"math/big"
	mrand "math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ------------ CONFIG ----------------
var (
	hostfetch_feed2 = "web.facebook.com"
	urlfetch_feed2  = "https://web.facebook.com/messaging/lightspeed/request"
	//userIDfetch_feed2 = "61562198647863"
)

func randomExcellentBandwidthfetch_feed2() string {
	mrand.Seed(time.Now().UnixNano()) // ✅ ต้องใช้ mrand ไม่ใช่ rand (crypto)
	min := 20000000
	max := 35000000
	return strconv.Itoa(mrand.Intn(max-min+1) + min)
}

// ------------ MAIN ------------------
func Runfetch_feed2_lightspeed(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	//	address := hostfetch_feed2 + ":443"
	boundary := genHexfetch_feed2(16)
	deviceID := genUUIDfetch_feed2()
	familyID := genUUIDfetch_feed2()
	requestToken := genUUIDfetch_feed2()
	epochID := genUint64fetch_feed2()
	cursor := genBase64fetch_feed2(32)
	traceID := genUint64fetch_feed2()
	traceType := genIntfetch_feed2(0, 1)

	body := buildMultipartfetch_feed2(boundary, epochID, cursor, traceID, traceType)
	contentLen := len(body)

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

	req, _ := http.NewRequest("POST", urlfetch_feed2, bytes.NewReader(body))
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	req.Header.Set("Host", hostfetch_feed2)
	req.Header.Set("device_id", deviceID)
	req.Header.Set("family_device_id", familyID)
	req.Header.Set("request_token", requestToken)
	req.Header.Set("user_id", userID)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "msysDataTask0")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","retry_attempt":"0"},"application_tags":"unknown"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-zero-eh", "2,,AS_f4gKEkvTH3SfgCYJVzirQUTh0TLneLfitj0JQoYbj90OG3tBV9erihGNv-2EK4YE")
	req.Header.Set("Zero-Rated", "0")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthfetch_feed2())
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

	bw := tlsConns.RWWeb.Writer
	br := tlsConns.RWWeb.Reader

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

// ------------ HELPERS ------------------

func buildMultipartfetch_feed2(boundary string, epochID uint64, cursor string, traceID uint64, traceType int) []byte {
	delimiter := "--" + boundary
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s\r\n", delimiter)
	fmt.Fprintf(&b, `Content-Disposition: form-data; name="request_payload"`+"\r\n\r\n")
	fmt.Fprintf(&b, `{"database":"44","epoch_id":%d,"format":"flatbuffer","last_applied_cursor":"%s","network_epoch_failure_count":"1","propagated_trace_id":%d,"propagated_trace_type":%d,"version":"7019082018129972"}`+"\r\n", epochID, cursor, traceID, traceType)
	fmt.Fprintf(&b, "%s\r\n", delimiter)
	fmt.Fprintf(&b, `Content-Disposition: form-data; name="request_type"`+"\r\n\r\n")
	fmt.Fprintf(&b, "2\r\n")
	fmt.Fprintf(&b, "%s--\r\n", delimiter)
	return b.Bytes()
}

func genBase64fetch_feed2(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	b := make([]byte, n)
	for i := range b {
		r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[r.Int64()]
	}
	return string(b)
}

func genUint64fetch_feed2() uint64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(0).Lsh(big.NewInt(1), 63))
	return n.Uint64()
}

func genIntfetch_feed2(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}

func genUUIDfetch_feed2() string {
	b := make([]byte, 16)
	mrand.Read(b) // ใช้ math/rand แบบไม่ปลอดภัยก็ได้
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func genHexfetch_feed2(n int) string {
	b := make([]byte, n)
	mrand.Read(b)
	return fmt.Sprintf("%x", b)
}
