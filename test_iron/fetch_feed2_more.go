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

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// ================= CONFIG ====================
var (
	hostfetch_feed2_more = "graph.facebook.com"
)

// =============================================

func randomExcellentBandwidthfetch_feed2_more() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runfetch_feed2_more(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	rand.Seed(time.Now().UnixNano())

	deviceID := uuid.New().String()
	clientQueryID := fmt.Sprintf("%d_%s", time.Now().Unix(), uuid.New().String())
	feedSessionID := fmt.Sprintf("%d", time.Now().UnixNano())

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

	// recent_vpvs_v2
	type VPVEntry struct {
		FeedBackendDataSerializedPayloads string      `json:"feed_backend_data_serialized_payloads"`
		Vsid                              string      `json:"vsid"`
		Qid                               string      `json:"qid"`
		Objid                             string      `json:"objid"`
		StoryType                         interface{} `json:"story_type"`
		OriginalQid                       string      `json:"original_qid"`
		FetchTracking                     bool        `json:"fetch_tracking"`
		Vspos                             int         `json:"vspos"`
		FeedSessionID                     string      `json:"feed_session_id"`
		Vvt                               int         `json:"vvt"`
		CurrentVideoTimeSpent             int         `json:"current_video_time_spent"`
		Timestamp                         int64       `json:"timestamp"`
		ClientVpvToken                    string      `json:"client_vpv_token"`
	}
	recentVPVS := []VPVEntry{}
	for i := 0; i < 3; i++ {
		recentVPVS = append(recentVPVS, VPVEntry{
			FeedBackendDataSerializedPayloads: "Gw1GAvK6q6zjnpIE/AMCwgHuA7QDDiLw4pGqxMAtBKKg9JLjnpIEBg4Q8rqrrOOekgQc8rqrrOOekgQAtujv9JAMHo7liKTSl48B+gMG8gOw6gEbAkjAAQhmb2xsb3dlZQ4IZm9sbG93ZWUbA0m25gMG3uYDBoTnAwYA",
			Vsid:                              fmt.Sprintf("%d", rand.Int63()),
			Qid:                               fmt.Sprintf("%d", rand.Int63()),
			Objid:                             "",
			StoryType:                         nil,
			OriginalQid:                       "0",
			FetchTracking:                     false,
			Vspos:                             0,
			FeedSessionID:                     feedSessionID,
			Vvt:                               -1,
			CurrentVideoTimeSpent:             -1,
			Timestamp:                         time.Now().Unix(),
			ClientVpvToken:                    fmt.Sprintf("%d", rand.Intn(8999)+1000),
		})
	}

	// cached_story_data
	type CachedStory struct {
		StoryRankingTime int64  `json:"storyRankingTime"`
		StoryId          string `json:"storyId"`
	}
	cachedStories := []CachedStory{}
	for i := 0; i < 3; i++ {
		cachedStories = append(cachedStories, CachedStory{
			StoryRankingTime: time.Now().Unix(),
			StoryId:          fmt.Sprintf("%d", rand.Int63()),
		})
	}

	variables := map[string]interface{}{
		// (เหมือนของคุณทุกบรรทัด)
		"enable_download": true,
		"include_post_render_format_conversation_first_ufi_options": true,
		"device_id":                    deviceID,
		"in_feed_guide_caller_id":      "FB4A_IFG_DYNAMIC",
		"experimental_fields":          []interface{}{},
		"client_query_id":              clientQueryID,
		"gysj_cover_photo_width_param": 738,
		"saved_lists_enabled":          true,
		"question_poll_count":          100,
		"inspiration_capabilities": []interface{}{
			map[string]interface{}{
				"version": 175, "type": "MSQRD_MASK", "capabilities": []interface{}{
					map[string]interface{}{"value": "multiplane_disabled", "name": "multiplane"},
					map[string]interface{}{"value": "world_tracker_disabled", "name": "world_tracker"},
					map[string]interface{}{"value": "xray_disabled", "name": "xray"},
					map[string]interface{}{"value": "world_tracking_disabled", "name": "world_tracking"},
					map[string]interface{}{"value": "half_float_render_pass_enabled", "name": "half_float_render_pass"},
					map[string]interface{}{"value": "multiple_render_targets_enabled", "name": "multiple_render_targets"},
					map[string]interface{}{"value": "vertex_texture_fetch_enabled", "name": "vertex_texture_fetch"},
					map[string]interface{}{"value": "render_settings_high_enabled", "name": "render_settings_high"},
					map[string]interface{}{"value": "body_tracking_disabled", "name": "body_tracking"},
					map[string]interface{}{"value": "gyroscope_enabled", "name": "gyroscope"},
					map[string]interface{}{"value": "geoanchor_disabled", "name": "geoanchor"},
					map[string]interface{}{"value": "scene_depth_disabled", "name": "scene_depth"},
					map[string]interface{}{"value": "segmentation_disabled", "name": "segmentation"},
					map[string]interface{}{"value": "hand_tracking_disabled", "name": "hand_tracking"},
					map[string]interface{}{"value": "real_scale_estimation_disabled", "name": "real_scale_estimation"},
					map[string]interface{}{"value": "hair_segmentation_disabled", "name": "hair_segmentation"},
					map[string]interface{}{"value": "depth_shader_read_enabled", "name": "depth_shader_read"},
					map[string]interface{}{"value": "etc2_compression", "name": "compression"},
					map[string]interface{}{"value": "0", "name": "face_tracker_version"},
					map[string]interface{}{"value": "133.0,134.0,135.0,136.0,137.0,138.0,139.0,140.0,141.0,142.0,143.0,144.0,145.0,146.0,147.0,148.0,149.0,150.0,151.0,152.0,153.0,154.0,155.0,156.0,157.0,158.0,159.0,160.0,161.0,162.0,163.0,164.0,165.0,166.0,167.0,168.0,169.0,170.0,171.0,172.0,173.0,174.0,175.0", "name": "supported_sdk_versions"},
				},
			},
			map[string]interface{}{"version": 1, "type": "FRAME"},
			map[string]interface{}{"version": 1, "type": "SHADER_FILTER", "capabilities": []interface{}{
				map[string]interface{}{"value": "true", "name": "multipass"},
			}},
		},
		"image_large_aspect_width":                  1080,
		"edge_metadata_fetch_enabled":               true,
		"profile_entry_point":                       "NEWS_FEED_POST",
		"include_is_currently_live":                 true,
		"creative_high_img_size":                    1080,
		"image_medium_height":                       2048,
		"creative_med_img_size":                     540,
		"goodwill_small_accent_image":               108,
		"fetch_fbc_header":                          true,
		"gysj_facepile_size_param":                  246,
		"gysj_size_param":                           158,
		"pyml_first_fetch_size":                     8,
		"image_large_aspect_height":                 565,
		"media_question_photo_size":                 1080,
		"quick_promotion_large_image_size_param":    1080,
		"msqrd_instruction_image_width":             263,
		"recent_vpvs_v2":                            recentVPVS,
		"multi_share_item_image_size_param":         551,
		"greeting_card_image_size_large":            1080,
		"num_media_question_options":                15,
		"include_comment_markdown":                  true,
		"should_fetch_sponsored_bumpers":            true,
		"skip_sample_entities_fields":               true,
		"fetch_is_instant_feed_cached_story":        true,
		"thumbnail_height":                          158,
		"profile_pic_swipe_size_param":              494,
		"bloks_version":                             "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		"inline_text_delight_comment_enabled":       true,
		"refresh_mode_param":                        "manual",
		"angora_attachment_cover_image_size":        1260,
		"orderby_home_story_param":                  []string{"top_stories"},
		"quick_promotion_image_size_param":          210,
		"reading_attachment_profile_image_height":   354,
		"instant_article_server_control_prefetch":   true,
		"reading_attachment_profile_image_width":    236,
		"action_links_location":                     "feed_mobile",
		"cached_story_data":                         cachedStories,
		"home_story_first_page_total_count":         10,
		"storyset_first_fetch_size":                 8,
		"poll_voters_count":                         5,
		"media_type":                                "image/jpeg",
		"fetch_has_bump_comment_in_story":           true,
		"top_ad_position_enabled":                   true,
		"include_stars_ufi_metadata":                true,
		"allocation_gap_hint_fetch_enabled":         true,
		"image_size_px":                             158,
		"feed_story_render_location":                "feed_mobile",
		"frame_scale":                               "3",
		"creative_low_img_size":                     360,
		"enable_hd":                                 true,
		"should_fetch_birthday_avatar_nt_action":    true,
		"include_predicted_feed_topics":             true,
		"quick_promotion_branding_image_size_param": 63,
		"fbstory_tray_preview_height":               565,
		"image_high_height":                         2048,
		"poll_facepile_size":                        105,
		"fetch_is_end_of_feed_story":                true,
		"include_open_message_in_ufi":               true,
		"fetch_partial_feedback_ctr":                true,
		"profile_pic_media_type":                    "image/x-auto",
		"feed_clarity_config":                       map[string]interface{}{"samplers": []interface{}{}, "is_enabled": false},
		"real_time_engagements":                     []interface{}{},
		"default_image_scale":                       3,
		"recent_comment_vpvs":                       []interface{}{},
		"scale":                                     "3",
		"enable_comment_identity_badge":             true,
		"device_type":                               "SM-J730G",
		"should_fetch_adaptive_ufi":                 true,
		"pymk_size_param":                           525,
		"image_high_width":                          1080,
		"fbstory_tray_preview_width":                318,
		"image_medium_width":                        540,
		"fb_shorts_group_author_picture_size":       32,
		"action_location":                           "feed",
		"image_low_height":                          2048,
		"greeting_card_image_size_medium":           540,
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
		"size_style":                               "contain-fit",
		"fbstory_tray_sizing_type":                 "cover-fill-cropped",
		"inline_text_bolding_comment_enabled":      true,
		"image_low_width":                          360,
		"include_post_render_format":               true,
		"friends_locations_profile_pic_size_param": 494,
		"msqrd_supported_capabilities": []interface{}{
			map[string]interface{}{"value": "multiplane_disabled", "name": "multiplane"},
			map[string]interface{}{"value": "world_tracker_disabled", "name": "world_tracker"},
			map[string]interface{}{"value": "xray_disabled", "name": "xray"},
			map[string]interface{}{"value": "world_tracking_disabled", "name": "world_tracking"},
			map[string]interface{}{"value": "half_float_render_pass_enabled", "name": "half_float_render_pass"},
			map[string]interface{}{"value": "multiple_render_targets_enabled", "name": "multiple_render_targets"},
			map[string]interface{}{"value": "vertex_texture_fetch_enabled", "name": "vertex_texture_fetch"},
			map[string]interface{}{"value": "render_settings_high_enabled", "name": "render_settings_high"},
			map[string]interface{}{"value": "body_tracking_disabled", "name": "body_tracking"},
			map[string]interface{}{"value": "gyroscope_enabled", "name": "gyroscope"},
			map[string]interface{}{"value": "geoanchor_disabled", "name": "geoanchor"},
			map[string]interface{}{"value": "scene_depth_disabled", "name": "scene_depth"},
			map[string]interface{}{"value": "segmentation_disabled", "name": "segmentation"},
			map[string]interface{}{"value": "hand_tracking_disabled", "name": "hand_tracking"},
			map[string]interface{}{"value": "real_scale_estimation_disabled", "name": "real_scale_estimation"},
			map[string]interface{}{"value": "hair_segmentation_disabled", "name": "hair_segmentation"},
			map[string]interface{}{"value": "depth_shader_read_enabled", "name": "depth_shader_read"},
			map[string]interface{}{"value": "etc2_compression", "name": "compression"},
			map[string]interface{}{"value": "0", "name": "face_tracker_version"},
			map[string]interface{}{"value": "133.0,134.0,135.0,136.0,137.0,138.0,139.0,140.0,141.0,142.0,143.0,144.0,145.0,146.0,147.0,148.0,149.0,150.0,151.0,152.0,153.0,154.0,155.0,156.0,157.0,158.0,159.0,160.0,161.0,162.0,163.0,164.0,165.0,166.0,167.0,168.0,169.0,170.0,171.0,172.0,173.0,174.0,175.0", "name": "supported_sdk_versions"},
		},
		"supported_model_compression_types": []string{"TAR_BROTLI", "NONE"},
		"battery_context":                   "{is_charging:\"false\",battery_level:100}",
		"msqrd_instruction_image_height":    263,
		"fb_shorts_location":                "fb_shorts_video_deep_dive",
		"supported_compression_types":       []string{"ZIP", "TAR_BROTLI"},
		"pyml_size_param":                   105,
		"thumbnail_width":                   158,
		"profile_image_size":                105,
		"should_fetch_fallback_actions":     true,
		"news_feed_only":                    true,
		"include_image_ranges":              true,
		"use_server_thumbnail":              true,
		"feed_session_id":                   feedSessionID,
	}

	// ============ Marshal JSON & Encode ===============
	varBuf, _ := json.Marshal(variables)
	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "fresh_feed_more_data_fetch")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "104560688215380751667448633671") // เปลี่ยนเป็น doc_id ที่ต้องการ 104560688215380751667448633671
	form.Set("variables", string(varBuf))
	form.Set("fb_api_analytics_tags", `["fetch_cause:PULL_TO_REFRESH","client_query_id:`+clientQueryID+`","GraphServices"]`)

	req, _ := http.NewRequest("POST", "https://"+hostfetch_feed2_more+"/graphql", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostfetch_feed2_more)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "fresh_feed_new_data_fetch")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "3130154110338948")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthfetch_feed2_more())
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

	// ถ้าต้อง spoof header เพิ่มเติมก็เติมตรงนี้
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
