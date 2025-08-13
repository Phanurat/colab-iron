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

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func randomExcellentBandwidthfollow() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runfollow(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var subscribee string
	err = db.QueryRow("SELECT subscribee_id FROM subscribee_id_table LIMIT 1").Scan(
		&subscribee)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	// ----------------- Build variables -----------------
	variables := fmt.Sprintf(`{
		"input":{
			"tracking":[],
			"subscribee_id":"%s",
			"subscribe_location":"PROFILE",
			"story_id":null,
			"client_mutation_id":"%s",
			"is_tracking_encrypted":false,
			"attribution_id_v2":"%s",
			"actor_id":"%s"
		},
		"fetch_profile_context_row":true,
		"nt_context":{}
	}`,
		subscribee,
		generateUUIDfollow(),
		generateAttributionIDfollow(),
		userID,
	)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "ActorSubscribeCoreMutation")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "333384874815355113079294933428")
	form.Set("variables", variables)
	form.Set("fb_api_analytics_tags", fmt.Sprintf(`["visitation_id=4748854339:%s:0:%d","GraphServices"]`, shortIDfollow(), time.Now().Unix()))
	form.Set("client_trace_id", generateUUIDfollow())

	// ----------------- Gzip compress -----------------
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	_, _ = gz.Write([]byte(form.Encode()))
	_ = gz.Close()

	// ----------------- Connect (proxy + TLS) -----------------
	host := "graph.facebook.com"
	//address := host + ":443"

	// ----------------- Build request -----------------
	req, _ := http.NewRequest("POST", "https://"+host+"/graphql", &compressed)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "ActorSubscribeCoreMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Navigation-Chain", generateNavigationChainfollow())
	req.Header.Set("x-fb-privacy-context", "2368177546817046")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-ta-logging-ids", "graphql:"+generateUUIDfollow())
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)                                        // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                        // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthfollow()) //เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

	// ----------------- Send -----------------
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

	_, err = db.Exec("DELETE FROM subscribee_id_table WHERE subscribee_id = ?", subscribee) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ subscribee_id_table ออกจากฐานข้อมูลแล้ว:", subscribee)
	}

}

func generateUUIDfollow() string {
	return uuid.New().String()
}

func shortIDfollow() string {
	id := uuid.New().String()
	return id[:6]
}

func generateAttributionIDfollow() string {
	timestamp := float64(time.Now().UnixNano()) / 1e9
	return fmt.Sprintf("ProfileFragment,profile_vnext_tab_posts,,%.3f,248631074,,,;ProfileFragment,timeline,,%.3f,248631074,,,", timestamp, timestamp-0.5)
}

func generateNavigationChainfollow() string {
	now := float64(time.Now().UnixNano()) / 1e9
	return fmt.Sprintf("ProfileFragment,profile_vnext_tab_posts,,%.3f,248631074,,,;ProfileFragment,timeline,,%.3f,248631074,,,", now, now-0.5)
}
