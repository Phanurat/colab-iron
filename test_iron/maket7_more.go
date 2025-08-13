package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func randomExcellentBandwidthmaket7_more() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runmaket7_more(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	host := "graph.facebook.com"
	//address := host + ":443"

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

	variables := genVariablesWithPaginationmaket7_more()
	fmt.Println("🧾 Variables:")
	fmt.Println(prettyPrintJSONmaket7_more([]byte(variables)))

	form := url.Values{}
	form.Set("access_token", accessToken)
	form.Set("fb_api_caller_class", "RelayModern")
	form.Set("fb_api_req_friendly_name", "MarketplaceHomeFeedPaginationQuery")
	form.Set("variables", variables)
	form.Set("server_timestamps", "true")
	form.Set("doc_id", "7634017873281297")

	var zipped bytes.Buffer
	gw := gzip.NewWriter(&zipped)
	_, _ = gw.Write([]byte(form.Encode()))
	gw.Close()

	// 🧾 HTTP req
	req, err := http.NewRequest("POST", "https://"+host+"/graphql?locale=en_US", &zipped)
	if err != nil {
		log.Fatalf("❌ build req failed: %v", err)
	}
	req.Host = host
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "RelayFBNetwork_MarketplaceHomeFeedPaginationQuery")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","retry_attempt":"0"},"application_tags":"unknown"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-zero-eh", "2,,AS_f4gKEkvTH3SfgCYJVzirQUTh0TLneLfitj0JQoYbj90OG3tBV9erihGNv-2EK4YE")
	req.Header.Set("Zero-Rated", "0")
	req.Header.Set("x-fb-net-hni", netHni)
	req.Header.Set("x-fb-sim-hni", simHni)
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthmaket7_more())
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

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

// generate randomized variables
func genVariablesWithPaginationmaket7_more() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	postID := fmt.Sprintf("7%d", r.Int63n(1e16))
	qid := fmt.Sprintf("-%d", r.Int63n(1e16))
	rankSig := fmt.Sprintf("%d", r.Int63n(1e16))

	rankObj := map[string]interface{}{
		"target_id":         postID,
		"target_type":       6,
		"primary_position":  0,
		"ranking_signature": rankSig,
		"commerce_channel":  501,
		"value":             0,
		"upsell_type":       21,
		"candidate_retrieval_source_map": map[string]int{
			"24233918622879705": 805,
			"9725946894191496":  3001,
		},
		"grouping_info": nil,
	}

	seenAds := []string{}
	for i := 0; i < 3; i++ {
		seenAds = append(seenAds, fmt.Sprintf("%d", r.Int63n(1e15)))
	}

	cursorObj := map[string]interface{}{
		"basic": map[string]interface{}{
			"item_index": r.Intn(10),
		},
		"ads": map[string]interface{}{
			"items_since_last_ad":            r.Intn(10),
			"items_retrieved":                6,
			"ad_index":                       r.Intn(10),
			"ad_slot":                        r.Intn(5),
			"dynamic_gap_rule":               0,
			"counted_organic_items":          0,
			"average_organic_score":          0,
			"is_dynamic_gap_rule_set":        false,
			"first_organic_score":            0,
			"is_dynamic_initial_gap_set":     false,
			"iterated_organic_items":         1,
			"top_organic_score":              0,
			"feed_slice_number":              1,
			"feed_retrieved_items":           6,
			"ad_req_id":                      r.Intn(1e9),
			"refresh_ts":                     0,
			"cursor_id":                      r.Intn(50000),
			"mc_id":                          0,
			"ad_index_e2e":                   0,
			"seen_ads":                       map[string]interface{}{"ad_ids": seenAds, "page_ids": seenAds, "campaign_ids": []string{}},
			"has_ad_index_been_reset":        false,
			"is_reconsideration_ads_dropped": false,
		},
		"boosted_ads": map[string]interface{}{
			"items_since_last_ad":            r.Intn(10),
			"items_retrieved":                6,
			"ad_index":                       r.Intn(10),
			"ad_slot":                        r.Intn(5),
			"dynamic_gap_rule":               0,
			"counted_organic_items":          0,
			"average_organic_score":          0,
			"is_dynamic_gap_rule_set":        false,
			"first_organic_score":            0,
			"is_dynamic_initial_gap_set":     false,
			"iterated_organic_items":         0,
			"top_organic_score":              0,
			"feed_slice_number":              0,
			"feed_retrieved_items":           0,
			"ad_req_id":                      0,
			"refresh_ts":                     0,
			"cursor_id":                      r.Intn(50000),
			"mc_id":                          0,
			"ad_index_e2e":                   0,
			"seen_ads":                       map[string]interface{}{"ad_ids": []string{}, "page_ids": []string{}},
			"has_ad_index_been_reset":        false,
			"is_reconsideration_ads_dropped": false,
		},
		"lightning": map[string]interface{}{
			"initial_request":   false,
			"top_unit_item_ids": []string{postID},
			"ranking_signature": rankSig,
			"qid":               qid,
		},
		"delivered_ids": map[string]interface{}{
			"delivered_product_item_ids": []string{},
			"delivered_module_ids":       []string{},
		},
	}

	vars := map[string]interface{}{
		"FBFallbackAttachment_LARGE_HEIGHT":         472,
		"FBFallbackAttachment_LARGE_WIDTH":          1080,
		"FBFallbackAttachment_PORTRAIT_HEIGHT":      472,
		"FBFallbackAttachment_PORTRAIT_WIDTH":       315,
		"FBFallbackAttachment_SMALL_HEIGHT":         315,
		"FBFallbackAttachment_SMALL_WIDTH":          315,
		"FBPhoto_LEGACY_FULL_HEIGHT":                731,
		"FBPhoto_LEGACY_FULL_WIDTH":                 411,
		"FBPhoto_LEGACY_LARGE_HEIGHT":               200,
		"FBPhoto_LEGACY_LARGE_WIDTH":                411,
		"FBPhoto_LEGACY_MEDIUM_HEIGHT":              120,
		"FBPhoto_LEGACY_MEDIUM_WIDTH":               205,
		"FBPhoto_LEGACY_SMALL_HEIGHT":               90,
		"FBPhoto_LEGACY_SMALL_WIDTH":                137,
		"FBPhoto_ZOOMED_HEIGHT":                     365,
		"FBPhoto_ZOOMED_WIDTH":                      205,
		"FBStickerAttachment_SIZE":                  157,
		"FBVideoAttachment_HEIGHT":                  341,
		"FBVideoAttachment_WIDTH":                   1080,
		"MarketplaceBrowseFeedScrollView_SIZE":      405,
		"MarketplaceExploreView_IMAGE_CONTEXT_FEED": "",
		"MarketplaceExploreView_IMAGE_MEDIA_TYPE":   "image/jpeg",
		"ShopsUIUtils_FEED_IMAGE_WIDTH":             202,
		"ad_cursor_override":                        nil,
		"cappedScale":                               2,
		"count":                                     6,
		"cursor":                                    cursorObj,
		"inSessionIntentProfile": map[string]interface{}{
			"referral_surface": "MARKETPLACE_UNKNOWN",
			"seed_product_id":  "",
		},
		"injected_unit_config": nil,
		"isPullToRefresh":      false,
		"localOnly":            nil,
		"marketplaceID":        nil,
		"queryType":            "ADS_ORGANIC_JOINT_QUERY",
		"real_time_signal_store": map[string]interface{}{
			"clickData":          []interface{}{},
			"dwellData":          []interface{}{},
			"homeFeedFcrt":       false,
			"impressionData":     []interface{}{},
			"latestAdEvent":      nil,
			"latestOrganicEvent": nil,
			"revisitData":        []interface{}{},
			"searchData":         []interface{}{},
			"skimData":           []interface{}{},
			"surfaceEnterData":   []interface{}{},
			"surfaceExitData":    []interface{}{},
			"swipeData":          []interface{}{},
		},
		"scale":       2.625,
		"seenItemIDs": []string{postID},
		"seenItemTracking": map[string]interface{}{
			"qid":               qid,
			"mf_story_key":      postID,
			"top_level_post_id": postID,
			"commerce_rank_obj": rankObj,
		},
		"shippedOnly":                                nil,
		"shouldURIPrefetchSmallSingleAdImages":       false,
		"shouldUseNativeBSGPostRenderingForMegamall": false,
		"should_skip_ad_request":                     false,
		"skipLargeAdFormat":                          false,
		"streaming":                                  false,
		"topOfFeedFetchReason":                       "UNKNOWN",
		"__relay_internal__pv__enablePrefetchFeedUnitImagesrelayprovider": false,
		"__relay_internal__pv__enableDeferVideoFragmentrelayprovider":     false,
		"__relay_internal__pv__enableListingMediaAPIrelayprovider":        false,
	}

	buff, _ := json.Marshal(vars)
	return string(buff)
}

// pretty JSON
func prettyPrintJSONmaket7_more(input []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, input, "", "  ")
	if err != nil {
		return string(input)
	}
	return out.String()
}
