package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
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

func randomExcellentBandwidthfriend_accept1() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func extractFriendIDsfriend_accept1(bodyResp []byte) {

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "friend.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
		return
	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/friend.db")

	//	dir, _ := os.Getwd()
	//	dbPath := filepath.Join(dir, "friend.db")
	//	db, err := sql.Open("sqlite3", dbPath)
	//	if err != nil {
	//		panic("❌ เปิดฐานข้อมูลไม่สำเร็จ: " + err.Error())
	//	}
	//	defer db.Close()

	var jsonResp map[string]interface{}
	if err := json.Unmarshal(bodyResp, &jsonResp); err != nil {
		fmt.Println("❌ แปลง JSON ไม่สำเร็จ: " + err.Error())
		return
	}

	dataRaw, ok := jsonResp["data"]
	if !ok {
		fmt.Println("❌ ไม่พบ key 'data'")
		return
	}
	data, ok := dataRaw.(map[string]interface{})
	if !ok {
		fmt.Println("❌ data ไม่ใช่ map[string]interface{}")
		return
	}

	viewerRaw, ok := data["viewer"]
	if !ok {
		fmt.Println("❌ ไม่พบ key 'viewer'")
		return
	}
	viewer, ok := viewerRaw.(map[string]interface{})
	if !ok {
		fmt.Println("❌ viewer ไม่ใช่ map[string]interface{}")
		return
	}

	tabRaw, ok := viewer["dynamic_friending_tab"]
	if !ok {
		fmt.Println("❌ ไม่พบ key 'dynamic_friending_tab'")
		return
	}
	tab, ok := tabRaw.(map[string]interface{})
	if !ok {
		fmt.Println("❌ dynamic_friending_tab ไม่ใช่ map[string]interface{}")
		return
	}

	edgesRaw, ok := tab["edges"]
	if !ok {
		fmt.Println("❌ ไม่พบ key 'edges'")
		return
	}
	edges, ok := edgesRaw.([]interface{})
	if !ok {
		fmt.Println("❌ edges ไม่ใช่ []interface{}")
		return
	}

	if len(edges) == 0 {
		fmt.Println("❌ ไม่มีคำขอเป็นเพื่อน")
		return
	}

	for _, item := range edges {
		node := item.(map[string]interface{})["node"].(map[string]interface{})
		typename, ok := node["__typename"].(string)
		if !ok || typename != "FriendRequestsFriendingTabRow" {
			continue
		}

		user, ok := node["user"].(map[string]interface{})
		if !ok {
			continue
		}

		friendID, ok := user["id"].(string)
		if !ok {
			continue
		}

		_, err := db.Exec(`INSERT INTO friend_info (friend_requester_id) VALUES (?)`, friendID)
		if err != nil {
			fmt.Println("❌ INSERT FAIL:", err)
			return
		} else {
			fmt.Println("✅ INSERT OK:", friendID)
		}
	}
}

func Runfriend_accept1(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	body := `method=post&pretty=false&format=json&server_timestamps=true&locale=en_US&purpose=fetch&fb_api_req_friendly_name=FriendingJewelContentQuery&fb_api_caller_class=graphservice&client_doc_id=274349594917117610276491443888&variables=%7B%22profile_picture_normal_size%22%3A242%2C%22profile_picture_small_size%22%3A158%2C%22pivot_link_options%22%3A%22default%22%2C%22nt_render_id%22%3A%220%22%2C%22nt_context%22%3A%7B%22using_white_navbar%22%3Atrue%2C%22pixel_ratio%22%3A3%2C%22is_push_on%22%3Atrue%2C%22styles_id%22%3A%22196702b4d5dfb9dbf1ded6d58ee42767%22%2C%22bloks_version%22%3A%22c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722%22%7D%2C%22scale%22%3A%223%22%2C%22supported_features%22%3A%7B%22client_ccu_status%22%3A%22DISABLED%22%7D%2C%22receiver_friction_enabled%22%3Atrue%2C%22pivot_links_enabled%22%3Atrue%2C%22dynamic_friending_tab_paginating_first%22%3A20%7D&fb_api_analytics_tags=%5B%22At_Connection%22%2C%22GraphServices%22%5D`
	host := "graph.facebook.com"
	//	address := host + ":443"

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

	// proxy := os.Getenv("USE_PROXY")
	// auth := os.Getenv("USE_PROXY_AUTH")

	// conn, err := net.DialTimeout("tcp", proxy, 10*time.Second)
	// if err != nil {
	// 	panic("❌ Proxy fail: " + err.Error())
	// }

	// reqLine := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n", address, host)
	// if auth != "" {
	// 	reqLine += "Proxy-Authorization: Basic " + auth + "\r\n"
	// }
	// reqLine += "\r\n"
	// fmt.Fprintf(conn, reqLine)

	// br := bufio.NewReader(conn)
	// respLine, _ := br.ReadString('\n')
	// if !strings.Contains(respLine, "200") {
	// 	panic("❌ CONNECT fail: " + respLine)
	// }
	// for {
	// 	line, _ := br.ReadString('\n')
	// 	if line == "\r\n" || line == "" {
	// 		break
	// 	}
	// }

	// utlsConn := utls.UClient(conn, &utls.Config{ServerName: host}, utls.HelloAndroid_11_OkHttp)
	// if err := utlsConn.Handshake(); err != nil {
	// 	panic("❌ TLS handshake fail: " + err.Error())
	// }

	req, _ := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Host", host)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "FriendingJewelContentQuery")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"fetch","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-graphql-request-purpose", "fetch")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthfriend_accept1())
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

	extractFriendIDsfriend_accept1(bodyResp)
}
