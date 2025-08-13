package main

import (
	"bytes"
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

func randomExcellentBandwidthRunprofile_pic8() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runprofile_pic8(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	bodyStr := "method=post&pretty=false&format=json&server_timestamps=true&locale=en_US&fb_api_req_friendly_name=FetchProfilePicGraphQL&fb_api_caller_class=graphservice&client_doc_id=332297695115422676241530047494&variables=%7B%22profile_image_big_size%22%3A246%2C%22square_profile_pic_size_big%22%3A246%2C%22square_profile_pic_size_small%22%3A105%2C%22square_profile_pic_size_huge%22%3A1080%7D&fb_api_analytics_tags=%5B%22GraphServices%22%5D"
	body := bytes.NewBufferString(bodyStr)

	//	proxy := fetchProxy()
	host := "graph.facebook.com"

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

	req, err := http.NewRequest("POST", "https://"+host+"/graphql", body)
	if err != nil {
		fmt.Println("❌ Request build fail: " + err.Error())
		return
	}

	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "graph.facebook.com")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "FetchProfilePicGraphQL")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "696323304542553")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)                                                 // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                                 // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthRunprofile_pic8()) //เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

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
