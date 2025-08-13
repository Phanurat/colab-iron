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

	_ "github.com/mattn/go-sqlite3"
)

func randomExcellentBandwidthlock_profile2() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// ---------- Config ---------
var (
	imageURLlock_profile2   = "https://static.xx.fbcdn.net/rsrc.php/v4/yu/r/FHVwy5Z6lXT.png"
	hostlock_profile2       = "static.xx.fbcdn.net"
	refererlock_profile2    = "fbapp://350685531728/unknown"
	fbFriendlylock_profile2 = "image"
	zeroEhlock_profile2     = "2,,ASAySH2oMJyaEZ8_gJqsBkU9XiRTbXAuIxwORpw9LUATP-gBzQrZQZq8gPwqoERuScM"
)

// ---------------------------

func Runlock_profile2(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	connToken := generateConnectionTokenlock_profile2()
	sessionID := generateSessionIDlock_profile2(connToken)

	analyticsTags := `{"network_tags":{"product":"350685531728","purpose":"prefetch","request_category":"image","retry_attempt":"0"},"application_tags":"unknown"}`

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

	req, _ := http.NewRequest("GET", imageURLlock_profile2, nil)
	req.Header.Set("Accept-Encoding", "zstd, gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostlock_profile2)
	req.Header.Set("Priority", "u=3")
	req.Header.Set("Referer", refererlock_profile2)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-token", connToken)
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", fbFriendlylock_profile2)
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", analyticsTags)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-session-id", sessionID)
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-zero-eh", zeroEhlock_profile2)
	req.Header.Set("Zero-Rated", "0")
	req.Header.Set("x-fb-net-hni", netHni)                                               // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                               // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthlock_profile2()) //เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

	// ---------- SEND ----------
	bw := tlsConns.RWstatic.Writer
	br := tlsConns.RWstatic.Reader

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

func generateConnectionTokenlock_profile2() string {
	const chars = "abcdef0123456789"
	rnd := make([]byte, 32)
	for i := range rnd {
		rnd[i] = chars[rand.Intn(len(chars))]
	}
	return string(rnd)
}

func generateSessionIDlock_profile2(cid string) string {
	return fmt.Sprintf("nid=TSRsHSL+wunc;tid=205;nc=0;fc=0;bc=0;cid=%s", cid)
}
