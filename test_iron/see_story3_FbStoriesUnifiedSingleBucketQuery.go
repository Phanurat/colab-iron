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
	hostsee_story3_FbStoriesUnifiedSingleBucketQuery     = "graph.facebook.com"
	endpointsee_story3_FbStoriesUnifiedSingleBucketQuery = "https://graph.facebook.com/graphql"
)

func randomExcellentBandwidthsee_story3_FbStoriesUnifiedSingleBucketQuery() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runsee_story3_FbStoriesUnifiedSingleBucketQuery(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	//	address := hostsee_story3_FbStoriesUnifiedSingleBucketQuery + ":443"
	viewerSessionID := genUUIDsee_story3_FbStoriesUnifiedSingleBucketQuery()

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

	// utlsConn := utls.UClient(conn, &utls.Config{ServerName: hostsee_story3_FbStoriesUnifiedSingleBucketQuery}, utls.HelloAndroid_11_OkHttp)
	// if err := utlsConn.Handshake(); err != nil {
	// 	panic("❌ TLS handshake fail: " + err.Error())
	// }

	variables := map[string]interface{}{
		"height":                           1920,
		"restore_all_removed_fields_wave1": true,
		"video_thumbnail_height":           672,
		"comment_previews_order":           []string{"admin_first"},
		"bloks_version":                    "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		"fbstory_tray_preview_height":      565,
		"msqrd_supported_capabilities": []map[string]string{
			{"value": "multiplane_disabled", "name": "multiplane"},
			{"value": "world_tracker_disabled", "name": "world_tracker"},
			{"value": "xray_disabled", "name": "xray"},
			{"value": "world_tracking_disabled", "name": "world_tracking"},
			{"value": "half_float_render_pass_enabled", "name": "half_float_render_pass"},
			{"value": "multiple_render_targets_enabled", "name": "multiple_render_targets"},
			{"value": "vertex_texture_fetch_enabled", "name": "vertex_texture_fetch"},
			{"value": "render_settings_high_enabled", "name": "render_settings_high"},
			{"value": "body_tracking_disabled", "name": "body_tracking"},
			{"value": "gyroscope_enabled", "name": "gyroscope"},
			{"value": "geoanchor_disabled", "name": "geoanchor"},
			{"value": "scene_depth_disabled", "name": "scene_depth"},
			{"value": "segmentation_disabled", "name": "segmentation"},
			{"value": "hand_tracking_disabled", "name": "hand_tracking"},
			{"value": "real_scale_estimation_disabled", "name": "real_scale_estimation"},
			{"value": "hair_segmentation_disabled", "name": "hair_segmentation"},
			{"value": "depth_shader_read_enabled", "name": "depth_shader_read"},
			{"value": "etc2_compression", "name": "compression"},
			{"value": "0", "name": "face_tracker_version"},
			{"value": "133.0,134.0,135.0,136.0,137.0,138.0,139.0,140.0,141.0,142.0,143.0,144.0,145.0,146.0,147.0,148.0,149.0,150.0,151.0,152.0,153.0,154.0,155.0,156.0,157.0,158.0,159.0,160.0,161.0,162.0,163.0,164.0,165.0,166.0,167.0,168.0,169.0,170.0,171.0,172.0,173.0,174.0,175.0", "name": "supported_sdk_versions"},
		},
		"comment_previews_include_attachments": true,
		"reaction_image_size":                  63,
		"enable_cix_screen_rollout":            true,
		"restore_bounds_to_location_sticker":   true,
		"fbstory_tray_sizing_type":             "cover-fill-cropped",
		"intro_card_preview_width":             318,
		"ar_effect_capabilities": []map[string]string{
			{"value": "multiplane_disabled", "name": "multiplane"},
			{"value": "world_tracker_disabled", "name": "world_tracker"},
			{"value": "xray_disabled", "name": "xray"},
			{"value": "world_tracking_disabled", "name": "world_tracking"},
			{"value": "half_float_render_pass_enabled", "name": "half_float_render_pass"},
			{"value": "multiple_render_targets_enabled", "name": "multiple_render_targets"},
			{"value": "vertex_texture_fetch_enabled", "name": "vertex_texture_fetch"},
			{"value": "render_settings_high_enabled", "name": "render_settings_high"},
			{"value": "body_tracking_disabled", "name": "body_tracking"},
			{"value": "gyroscope_enabled", "name": "gyroscope"},
			{"value": "geoanchor_disabled", "name": "geoanchor"},
			{"value": "scene_depth_disabled", "name": "scene_depth"},
			{"value": "segmentation_disabled", "name": "segmentation"},
			{"value": "hand_tracking_disabled", "name": "hand_tracking"},
			{"value": "real_scale_estimation_disabled", "name": "real_scale_estimation"},
			{"value": "hair_segmentation_disabled", "name": "hair_segmentation"},
			{"value": "depth_shader_read_enabled", "name": "depth_shader_read"},
			{"value": "etc2_compression", "name": "compression"},
			{"value": "0", "name": "face_tracker_version"},
			{"value": "133.0,134.0,135.0,136.0,137.0", "name": "supported_sdk_versions"},
		},
		"fbstory_tray_preview_width": 318,
		"nt_surface":                 "STORIES_VIEWER_SHEET",
		"reaction_image_scale":       2.5,
		"bucket_id":                  "173012030849267",
		"scale":                      "3",
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
		"comment_previews_count":             3,
		"page_profile_image_size_experiment": 105,
		"profile_image_size":                 105,
		"should_include_live_ring_fields":    true,
		"width":                              1080,
		"video_thumbnail_width":              672,
		"viewer_session_id":                  viewerSessionID,
	}

	varBuf, _ := json.Marshal(variables)
	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("purpose", "fetch")
	form.Set("fb_api_req_friendly_name", "FbStoriesUnifiedSingleBucketQuery")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "3636326671221641623534994213")
	form.Set("variables", string(varBuf))
	form.Set("fb_api_analytics_tags", `["trigger=forward_prefetch","surface=story_viewer","GraphServices","type=prefetch"]`)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(form.Encode()))
	gz.Close()

	req, _ := http.NewRequest("POST", endpointsee_story3_FbStoriesUnifiedSingleBucketQuery, &buf)
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostsee_story3_FbStoriesUnifiedSingleBucketQuery)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "FbStoriesUnifiedSingleBucketQuery")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "1326330710893128")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"fetch","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-graphql-request-purpose", "fetch")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthsee_story3_FbStoriesUnifiedSingleBucketQuery())
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

func genUUIDsee_story3_FbStoriesUnifiedSingleBucketQuery() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
