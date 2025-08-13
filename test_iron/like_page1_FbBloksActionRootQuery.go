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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/klauspost/compress/zstd"
	_ "github.com/mattn/go-sqlite3"
)

func decodeZstdRunlike_page1_FbBloksActionRootQuery(data []byte) ([]byte, error) {
	dec, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}
	defer dec.Close()
	return dec.DecodeAll(data, nil)
}
func extractFacebookIDsRunlike_page1_FbBloksActionRootQuery(rawurl string) (string, string, error) {
	var postID, ownerID string
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", "", err
	}
	query := u.Query()

	reStory := regexp.MustCompile(`story_fbid=(\d+)`)
	rePath := regexp.MustCompile(`facebook\.com/(\d+)/(?:videos|posts)/(\d+)`)
	reUsername := regexp.MustCompile(`facebook\.com/([^/?&]+)`)

	if match := reStory.FindStringSubmatch(rawurl); len(match) > 1 {
		postID = match[1]
	}
	if match := rePath.FindStringSubmatch(rawurl); len(match) > 2 {
		ownerID = match[1]
		postID = match[2]
	}
	if postID == "" {
		re := regexp.MustCompile(`/posts/(\d+)|/videos/(\d+)`)
		match := re.FindStringSubmatch(u.Path)
		if len(match) > 1 {
			if match[1] != "" {
				postID = match[1]
			} else {
				postID = match[2]
			}
		}
	}
	if id := query.Get("id"); id != "" {
		ownerID = id
	}
	if ownerID == "" {
		if match := reUsername.FindStringSubmatch(rawurl); len(match) > 1 {
			username := match[1]
			if isNumericRunlike_page1_FbBloksActionRootQuery(username) {
				ownerID = username
			} else {
				fbid, err := getFBIDFromUsernameRunlike_page1_FbBloksActionRootQuery(username)
				if err != nil {
					return "", "", err
				}
				ownerID = fbid
			}
		}
	}
	return ownerID, postID, nil
}

func isNumericRunlike_page1_FbBloksActionRootQuery(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func getFBIDFromUsernameRunlike_page1_FbBloksActionRootQuery(username string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest("GET", "https://mbasic.facebook.com/"+username, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if strings.HasPrefix(location, "intent://profile/") {
		re := regexp.MustCompile(`intent://profile/(\d+)`)
		match := re.FindStringSubmatch(location)
		if len(match) > 1 {
			return match[1], nil
		}
	}

	resp, err = http.Get("https://mbasic.facebook.com/" + username)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	text := string(body)

	re := regexp.MustCompile(`owner_id=(\d+)`)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1], nil
	}
	re = regexp.MustCompile(`profile\.php\?id=(\d+)`)
	match = re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("❌ ไม่พบ owner_id จาก username")
}

func buildEncodedPayloadRunlike_page1_FbBloksActionRootQuery(ownerID string) string {
	serverParams := map[string]interface{}{
		"INTERNAL__latency_qpl_marker_id":   36707139,
		"INTERNAL__latency_qpl_instance_id": generateQPLIDRunlike_page1_FbBloksActionRootQuery(),
		"render_location":                   12,
		"profile_id":                        ownerID,
	}
	innerParams := map[string]interface{}{
		"client_input_params": map[string]interface{}{},
		"server_params":       serverParams,
	}
	innerParamsJSON, _ := json.Marshal(innerParams)
	outerParams := map[string]interface{}{
		"params": string(innerParamsJSON),
	}
	outerParamsJSON, _ := json.Marshal(outerParams)

	variables := map[string]interface{}{
		"params": map[string]interface{}{
			"params":              string(outerParamsJSON),
			"bloks_versioning_id": "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
			"app_id":              "com.bloks.www.fb.profile.action_bar.action.like_button_tap",
		},
		"scale": "3",
		"nt_context": map[string]interface{}{
			"using_white_navbar": true,
			"pixel_ratio":        3,
			"is_push_on":         true,
			"styles_id":          "196702b4d5dfb9dbf1ded6d58ee42767",
			"bloks_version":      "c459b951c037ad3fbe67f94342f309a73154e66c326b3cd823682078d9eeb722",
		},
	}
	jsonVars, _ := json.Marshal(variables)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("purpose", "fetch")
	form.Set("fb_api_req_friendly_name", "FbBloksActionRootQuery-com.bloks.www.fb.profile.action_bar.action.like_button_tap")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "11994080424240083948543644217")
	form.Set("variables", string(jsonVars))
	form.Set("fb_api_analytics_tags", `["GraphServices"]`)
	form.Set("client_trace_id", uuid.New().String())
	return form.Encode()
}

func generateQPLIDRunlike_page1_FbBloksActionRootQuery() float64 {
	rand.Seed(time.Now().UnixNano())
	return float64(rand.Int63n(899999999999999) + 100000000000000)
}

func generateSessionIDRunlike_page1_FbBloksActionRootQuery() string {
	u := uuid.New().String()
	tid := strconv.Itoa(rand.Intn(900) + 100)
	return fmt.Sprintf("nid=%s;tid=%s;nc=1;fc=2;bc=1;cid=%s", u, tid, uuid.New().String())
}

func gzipCompressRunlike_page1_FbBloksActionRootQuery(data []byte) []byte {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(data)
	writer.Close()
	return buf.Bytes()
}

func Runlike_page1_FbBloksActionRootQuery(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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
	var token, userAgent, netHni, simHni, deviceGroup string
	err = db.QueryRow(`SELECT access_token, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1`).Scan(
		&token, &userAgent, &netHni, &simHni, &deviceGroup)
	if err != nil {
		fmt.Println("❌ ดึง app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	var link string
	err = db.QueryRow(`SELECT link_page FROM link_page_for_like_table LIMIT 1`).Scan(&link)
	if err != nil {
		fmt.Println("❌ ดึง link_page ไม่สำเร็จ: " + err.Error())
		return
	}

	ownerID, _, err := extractFacebookIDsRunlike_page1_FbBloksActionRootQuery(link)
	if err != nil {
		fmt.Println("❌ ขุด post_id ไม่สำเร็จ: " + err.Error())
		return
	}

	payload := buildEncodedPayloadRunlike_page1_FbBloksActionRootQuery(ownerID)
	host := "graph.facebook.com"

	req, err := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBuffer(gzipCompressRunlike_page1_FbBloksActionRootQuery([]byte(payload))))
	if err != nil {
		fmt.Println("❌ NewRequest fail: " + err.Error())
		return
	}

	traceID := uuid.New().String()
	sessionID := generateSessionIDRunlike_page1_FbBloksActionRootQuery()
	connectionToken := uuid.New().String()

	req.Header = map[string][]string{
		"Authorization":               {"OAuth " + token},
		"Accept-Encoding":             {"zstd, gzip, deflate"},
		"Connection":                  {"keep-alive"},
		"Content-Encoding":            {"gzip"},
		"Content-Type":                {"application/x-www-form-urlencoded"},
		"Host":                        {host},
		"User-Agent":                  {userAgent},
		"X-FB-Friendly-Name":          {"FbBloksActionRootQuery-com.bloks.www.fb.profile.action_bar.action.like_button_tap"},
		"x-fb-client-ip":              {"True"},
		"x-fb-connection-token":       {connectionToken},
		"X-FB-Connection-Type":        {"MOBILE.HSDPA"},
		"x-fb-device-group":           {deviceGroup},
		"X-FB-HTTP-Engine":            {"Liger"},
		"x-fb-privacy-context":        {"3643298472347298"},
		"x-fb-qpl-active-flows-json":  {`{"schema_version":"v2","inprogress_qpls":[],"snapshot_attributes":{}}`},
		"X-FB-Request-Analytics-Tags": {`{"network_tags":{"product":"350685531728","purpose":"fetch","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`},
		"x-fb-rmd":                    {"state=URL_ELIGIBLE"},
		"x-fb-server-cluster":         {"True"},
		"x-fb-session-id":             {sessionID},
		"x-fb-ta-logging-ids":         {"graphql:" + traceID},
		"x-graphql-client-library":    {"graphservice"},
		"x-graphql-request-purpose":   {"fetch"},
		"x-tigon-is-retry":            {"False"},
	}

	// ---------- SEND ----------
	bw := tlsConns.RWGraph.Writer
	br := tlsConns.RWGraph.Reader

	// ✅ เขียน request ผ่าน buffered writer
	err = req.Write(bw)
	if err != nil {
		fmt.Println("❌ Write fail: " + err.Error())
		return
	}
	bw.Flush() // ✅ ต้อง flush เพื่อให้ข้อมูลถูกส่งออกจริง ๆ

	// ✅ อ่าน response ผ่าน buffered reader
	resp, err := http.ReadResponse(br, req)
	if err != nil {
		fmt.Println("❌ Read fail: " + err.Error())
		return
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	var body []byte

	// ✅ จัดการ Content-Encoding เหมือนเดิม
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("❌ GZIP decompress fail: " + err.Error())
			return
		}
		defer reader.Close()
		body, err = io.ReadAll(reader)
		if err != nil {
			fmt.Println("❌ Body read fail: " + err.Error())
			return
		}
	case "zstd":
		reader = io.NopCloser(resp.Body)
		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			fmt.Println("❌ Read zstd body fail: " + err.Error())
			return
		}
		decoded, err := decodeZstdRunlike_page1_FbBloksActionRootQuery(bodyBytes)
		if err != nil {
			fmt.Println("❌ ZSTD decode fail: " + err.Error())
			return
		}
		body = decoded
	default:
		reader = resp.Body
		body, err = io.ReadAll(reader)
		if err != nil {
			fmt.Println("❌ Body read fail: " + err.Error())
			return
		}
	}

	fmt.Println("✅ Status:", resp.Status)
	fmt.Println("📦 Response:", string(body))

	// ✅ ดึงค่า delegate_page.id จาก response แล้วเก็บลง DB
	decoder := json.NewDecoder(bytes.NewReader(body))
	var jsonResponse map[string]interface{}
	if err := decoder.Decode(&jsonResponse); err != nil {
		fmt.Println("❌ JSON decode fail:", err)
		return
	}

	if data, ok := jsonResponse["data"].(map[string]interface{}); ok {
		if fbAction, ok := data["fb_bloks_action"].(map[string]interface{}); ok {
			if gqlVars, ok := fbAction["gql_variables"].([]interface{}); ok && len(gqlVars) > 0 {
				if profileAction, ok := gqlVars[0].(map[string]interface{})["profile_action"].(map[string]interface{}); ok {
					if entity, ok := profileAction["entity"].(map[string]interface{}); ok {
						if delegatePage, ok := entity["delegate_page"].(map[string]interface{}); ok {
							if pageIDStr, ok := delegatePage["id"].(string); ok {
								_, err := db.Exec("INSERT OR REPLACE INTO link_page_for_id_page_table (page_id) VALUES (?)", pageIDStr)
								if err != nil {
									fmt.Println("❌ บันทึก page_id ไม่สำเร็จ:", err)
								} else {
									fmt.Println("✅ บันทึก page_id สำเร็จ:", pageIDStr)
								}
							}
						}
					}
				}
			}
		}

	} else {
		fmt.Println("❌ JSON decode fail:", err)
	}
}
