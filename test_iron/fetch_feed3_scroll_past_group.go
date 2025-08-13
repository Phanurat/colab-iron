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

// ===================== CONFIG ปรับตรงนี้ ===========================
var (
	groupIDfetch_feed3_scroll_past_group = "1096205827471144"
	//	deviceGroup  = "6301"
	hostfetch_feed3_scroll_past_group         = "graph.facebook.com"
	friendlyNamefetch_feed3_scroll_past_group = "FetchGroupInformation"
	clientDocIDfetch_feed3_scroll_past_group  = "325529145212539348845560773234"
	privacyCtxfetch_feed3_scroll_past_group   = "3379608338725370"
)

// ===================== END CONFIG ===========================

func randomExcellentBandwidthfetch_feed3_scroll_past_group() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runfetch_feed3_scroll_past_group(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	variables := fmt.Sprintf(`{"image_low_width":360,"image_large_aspect_width":1080,"image_medium_width":540,"group_id":"%s","image_low_height":2048,"cover_photo_height":565,"media_type":"image/jpeg","cover_photo_width":1080,"top_promo_nux_id":"7383","profile_pic_media_type":"image/x-auto","default_image_scale":3,"image_large_aspect_height":565,"image_high_height":2048,"remove_unused_graphql_fields_group_composer_traits":false,"should_use_top_of_home_server_control":true,"nt_context":{"using_white_navbar":true,"pixel_ratio":3,"is_push_on":true,"styles_id":"196702b4d5dfb9dbf1ded6d58ee42767","bloks_version":"c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722"},"scale":"3","size_style":"contain-fit","image_high_width":1080,"group_composer_render_location":"group_mall","image_medium_height":2048,"action_source":"GROUP_MALL","should_fetch_action_intervention":true,"action_intervention_source":"GROUP_MALL","cover_image_navbar_size":84,"should_defer_rooms_creation_nt_action":true}`, groupIDfetch_feed3_scroll_past_group)

	clientContext := `{"is_notification_unread":"false","rank_index":"-1","has_blue_badge":"false","request_source":"native_newsfeed"}`
	analyticsTags := `["At_Connection","native_newsfeed","GraphServices"]`

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

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("purpose", "prefetch")
	form.Set("fb_api_req_friendly_name", friendlyNamefetch_feed3_scroll_past_group)
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", clientDocIDfetch_feed3_scroll_past_group)
	form.Set("fb_api_client_context", clientContext)
	form.Set("variables", variables)
	form.Set("fb_api_analytics_tags", analyticsTags)

	req, _ := http.NewRequest("POST", "https://"+hostfetch_feed3_scroll_past_group+"/graphql", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostfetch_feed3_scroll_past_group)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", friendlyNamefetch_feed3_scroll_past_group)
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", privacyCtxfetch_feed3_scroll_past_group)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"prefetch","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-graphql-request-purpose", "prefetch")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthfetch_feed3_scroll_past_group())
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

	// เติม header spoof เพิ่มตรงนี้ตามที่ดักได้จริง

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
