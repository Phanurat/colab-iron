package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/base64"
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

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// ------------------------------
// ปรับได้แค่ตรงนี้
var (
	//	deviceGroup = "6301"

	hostopenapp4         = "graph.facebook.com"
	friendlyNameopenapp4 = "FirstNotificationsPageQueryNewApi"
	clientDocIDopenapp4  = "307348724111638261155333340906"
	privacyCtxopenapp4   = "138965567254360"
)

// ------------------------------

func randomExcellentBandwidthopenapp4() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000 // 20 Mbps
	max := 35000000 // 35 Mbps
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runopenapp4(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	deviceID := uuid.New().String()
	cacheTokens := generateCacheTokens(10)

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

	variables := map[string]interface{}{
		"enable_comment_shares":                   true,
		"reading_attachment_profile_image_height": 354,
		"enable_interesting_replies":              true,
		"default_image_scale":                     3,
		"thumbnail_height":                        158,
		"bloks_version":                           "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
		"size_style":               "contain-fit",
		"image_large_aspect_width": 1080,
		"inspiration_capabilities": []map[string]interface{}{
			{"version": 175, "type": "MSQRD_MASK", "capabilities": []interface{}{}},
			{"version": 1, "type": "FRAME"},
			{"version": 1, "type": "SHADER_FILTER", "capabilities": []map[string]string{{"value": "true", "name": "multipass"}}},
		},
		"device_type":                            "",
		"should_fetch_ats_nt_view":               true,
		"angora_attachment_cover_image_size":     1260,
		"angora_attachment_profile_image_size":   105,
		"widget_image_scaling":                   2.625,
		"device_id":                              deviceID,
		"thumbnail_width":                        158,
		"profile_image_size":                     168,
		"enable_feedback_animation_config":       true,
		"msqrd_instruction_image_width":          263,
		"image_low_height":                       2048,
		"image_high_width":                       1080,
		"cache_tokens":                           cacheTokens,
		"image_medium_width":                     540,
		"environment":                            "MAIN_SURFACE",
		"icon_scale":                             3,
		"frame_scale":                            "3",
		"image_medium_height":                    2048,
		"widget_image_size":                      378,
		"fbstory_image_width":                    1080,
		"fbstory_image_height":                   1920,
		"image_preview_size":                     105,
		"media_type":                             "image/jpeg",
		"num_full_relevant_comments":             1,
		"feedback_reactions_floating_effect":     true,
		"notification_request_source":            "prefetch",
		"reading_attachment_profile_image_width": 236,
		"overlapping_glyph_enabled":              true,
		"should_fetch_full_relevant_comments":    false,
		"skip_sample_entities_fields":            true,
		"first_notification_stories":             7,
		"in_channel_eligibility_experiment":      false,
		"enable_comment_identity_badge":          true,
		"scale":                                  "3",
		"reaction_context": map[string]interface{}{
			"unit_styles":      []string{"VERTICAL_COMPONENTS"},
			"surface":          "ANDROID_NOTIFICATIONS_FRIENDING",
			"request_type":     "normal",
			"component_styles": []string{"FRIEND_REQUEST_ACTION_LIST"},
			"action_styles":    []string{"CONFIRM_FRIEND_REQUEST", "DELETE_FRIEND_REQUEST"},
		},
		"image_large_aspect_height":         565,
		"notif_query_flags":                 []string{"MESSENGER_NOT_INSTALLED", "CAMERA_ROLL_PERMISSION_NOT_GRANTED"},
		"image_high_height":                 2048,
		"supported_model_compression_types": []string{"TAR_BROTLI", "NONE"},
		"fb_shorts_location":                "fb_shorts_notification",
		"supported_compression_types":       []string{"ZIP", "TAR_BROTLI"},
		"first_n_feedback_reactions":        1,
		"msqrd_instruction_image_height":    263,
		"notif_option_set_context": map[string]interface{}{
			"supported_display_styles": []map[string]interface{}{
				{
					"option_set_display_style": "LONGPRESS_MENU",
					"option_display_styles":    []string{"POPUP_MENU_OPTION", "CHEVRON_MENU_OPTION", "HEADER_OPTION"},
				},
			},
			"client_action_types": []string{
				"HIDE", "MODSUB", "MARK_AS_READ", "MARK_AS_UNREAD", "OPEN_ACTION_SHEET",
				"NT_ACTION", "SHOW_MORE", "UNSUB", "SERVER_ACTION", "SAVE_ITEM", "UNSAVE_ITEM",
				"PIN_TO_NEW", "SNOOZE", "OPEN_NOTIF_SETTINGS", "TOGGLE_ACTIVE", "USEFUL_SURVEY",
				"REPORT_BUG", "REPORT_USER", "OPEN_IN_DEBUG_VIEWER",
			},
		},
		"profile_pic_media_type": "image/x-auto",
		"feedback_referrer":      "UNKNOWN_TAP_SOURCE",
		"image_low_width":        360,
		"feedback_source":        "NOTIFICATIONS_PREFETCH",
	}

	varBuf, _ := json.Marshal(variables)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", friendlyNameopenapp4)
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", clientDocIDopenapp4)
	form.Set("variables", string(varBuf))
	form.Set("fb_api_analytics_tags", `["At_Connection","fetch_location:0","GraphServices"]`)

	// address := host + ":443"
	// var conn net.Conn

	// proxy := os.Getenv("USE_PROXY")
	// auth := os.Getenv("USE_PROXY_AUTH")

	// conn, err = net.DialTimeout("tcp", proxy, 10*time.Second)
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

	req, _ := http.NewRequest("POST", "https://"+hostopenapp4+"/graphql", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostopenapp4)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	// req.Header.Set("X-FB-Connection-Type", "WIFI")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA") //
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", friendlyNameopenapp4)
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", privacyCtxopenapp4)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")

	req.Header.Set("x-fb-net-hni", netHni)                                          // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                          // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthopenapp4()) //เพิ่มเข้าไป
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

func generateCacheTokens(n int) []string {
	tokens := make([]string, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		b := make([]byte, 24)
		rand.Read(b)
		tokens[i] = base64.StdEncoding.EncodeToString(b)
	}
	return tokens
}
