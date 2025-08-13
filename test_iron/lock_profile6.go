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

func randomExcellentBandwidthlock_profile6() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

var (
	hostlock_profile6           = "static.xx.fbcdn.net"
	imagePathlock_profile6      = "/rsrc.php/v4/yq/r/wt7Cb0lR2LN.png"
	connTokenlock_profile6      = generateHex32lock_profile6()
	xFbSessionIDlock_profile6   = fmt.Sprintf("nid=TSRsHSL+wunc;tid=209;nc=0;fc=0;bc=0;cid=%s", connTokenlock_profile6)
	xFbRequestTagslock_profile6 = `{"network_tags":{"product":"350685531728","purpose":"prefetch","request_category":"image","retry_attempt":"0"},"application_tags":"unknown"}`
	xZeroEhlock_profile6        = "2,,ASAySH2oMJyaEZ8_gJqsBkU9XiRTbXAuIxwORpw9LUATP-gBzQrZQZq8gPwqoERuScM"
)

func Runlock_profile6(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

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

	urlImg := "https://" + hostlock_profile6 + imagePathlock_profile6
	req, _ := http.NewRequest("GET", urlImg, nil)
	req.Header.Set("Accept-Encoding", "zstd, gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostlock_profile6)
	req.Header.Set("Priority", "u=3")
	req.Header.Set("Referer", "fbapp://350685531728/unknown")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-token", connTokenlock_profile6)
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "image")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", xFbRequestTagslock_profile6)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-session-id", xFbSessionIDlock_profile6)
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-zero-eh", xZeroEhlock_profile6)
	req.Header.Set("Zero-Rated", "0")
	req.Header.Set("x-fb-net-hni", netHni)                                               // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                               // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthlock_profile6()) //เพิ่มเข้าไป
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

func generateHex32lock_profile6() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
