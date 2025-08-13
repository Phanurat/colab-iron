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
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	hostsee_story2_FBStoriesAdsPaginatingQuery_At_Connection     = "graph.facebook.com"
	endpointsee_story2_FBStoriesAdsPaginatingQuery_At_Connection = "https://graph.facebook.com/graphql"
)

type Capability struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func genUUIDsee_story2_FBStoriesAdsPaginatingQuery_At_Connection() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func randomExcellentBandwidthsee_story2_FBStoriesAdsPaginatingQuery_At_Connection() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func buildVariablessee_story2_FBStoriesAdsPaginatingQuery_At_Connection() map[string]interface{} {
	viewerSessionID := genUUIDsee_story2_FBStoriesAdsPaginatingQuery_At_Connection()

	capabilities := []Capability{
		{"multiplane", "multiplane_disabled"},
		{"world_tracker", "world_tracker_disabled"},
		{"xray", "xray_disabled"},
		{"world_tracking", "world_tracking_disabled"},
		{"half_float_render_pass", "half_float_render_pass_enabled"},
		{"multiple_render_targets", "multiple_render_targets_enabled"},
		{"vertex_texture_fetch", "vertex_texture_fetch_enabled"},
		{"render_settings_high", "render_settings_high_enabled"},
		{"body_tracking", "body_tracking_disabled"},
		{"gyroscope", "gyroscope_enabled"},
		{"geoanchor", "geoanchor_disabled"},
		{"scene_depth", "scene_depth_disabled"},
		{"segmentation", "segmentation_disabled"},
		{"hand_tracking", "hand_tracking_disabled"},
		{"real_scale_estimation", "real_scale_estimation_disabled"},
		{"hair_segmentation", "hair_segmentation_disabled"},
		{"depth_shader_read", "depth_shader_read_enabled"},
		{"compression", "etc2_compression"},
		{"face_tracker_version", "0"},
		{"supported_sdk_versions", "133.0,134.0,135.0,136.0,137.0,138.0,139.0,140.0,141.0,142.0,143.0,144.0,145.0,146.0,147.0,148.0,149.0,150.0,151.0,152.0,153.0,154.0,155.0,156.0,157.0,158.0,159.0,160.0,161.0,162.0,163.0,164.0,165.0,166.0,167.0,168.0,169.0,170.0,171.0,172.0,173.0,174.0,175.0"},
	}

	return map[string]interface{}{
		"entry_point":                          "STORIES_VIEWER_SHEET",
		"num_of_organic_stories":               4,
		"scale":                                "3",
		"viewer_session_id":                    viewerSessionID,
		"bucket_id":                            "173012030849267",
		"profile_image_size":                   105,
		"reaction_image_size":                  63,
		"reaction_image_scale":                 2.5,
		"nt_surface":                           "STORIES_VIEWER_SHEET",
		"fbstory_tray_preview_width":           318,
		"fbstory_tray_preview_height":          565,
		"fbstory_tray_sizing_type":             "cover-fill-cropped",
		"comment_previews_order":               []string{"admin_first"},
		"comment_previews_count":               3,
		"comment_previews_include_attachments": true,
		"page_profile_image_size_experiment":   105,
		"should_include_live_ring_fields":      true,
		"intro_card_preview_width":             318,
		"video_thumbnail_width":                672,
		"video_thumbnail_height":               672,
		"height":                               1920,
		"width":                                1080,
		"enable_cix_screen_rollout":            true,
		"restore_bounds_to_location_sticker":   true,
		"restore_all_removed_fields_wave1":     true,
		"bloks_version":                        "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
		"msqrd_supported_capabilities": capabilities,
		"ar_effect_capabilities":       capabilities,
	}
}

func Runsee_story2_FBStoriesAdsPaginatingQuery_At_Connection(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

	variables := buildVariablessee_story2_FBStoriesAdsPaginatingQuery_At_Connection()
	varBuf, _ := json.Marshal(variables)
	//	address := hostsee_story2_FBStoriesAdsPaginatingQuery_At_Connection + ":443"

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

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("purpose", "fetch")
	form.Set("fb_api_req_friendly_name", "FBStoriesAdsPaginatingQuery_At_Connection_Pagination_Viewer_facebook_story_ads_paginating")
	form.Set("fb_api_caller_class", "AtConnection")
	form.Set("client_doc_id", "267184416012693006368362197886")
	form.Set("fb_api_analytics_tags", `["At_Connection","GraphServices"]`)
	form.Set("variables", string(varBuf))

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(form.Encode()))
	gz.Close()

	req, _ := http.NewRequest("POST", endpointsee_story2_FBStoriesAdsPaginatingQuery_At_Connection, &buf)
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostsee_story2_FBStoriesAdsPaginatingQuery_At_Connection)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "FBStoriesAdsPaginatingQuery_At_Connection_Pagination_Viewer_facebook_story_ads_paginating")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "674480340001932")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"fetch","request_category":"graphql","retry_attempt":"0"},"application_tags":"AtConnection"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-graphql-request-purpose", "fetch")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthsee_story2_FBStoriesAdsPaginatingQuery_At_Connection())
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
