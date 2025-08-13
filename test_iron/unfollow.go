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

func randomExcellentBandwidthunfollow() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Rununfollow(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่

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

	var unsubscribee string
	err = db.QueryRow("SELECT unsubscribee_id FROM unsubscribee_id_table LIMIT 1").Scan(
		&unsubscribee)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	mutationID := uuid.New().String()
	traceID := uuid.New().String()
	visitID := shortIDunfollow()
	attribution := generateAttributionunfollow()

	variables := fmt.Sprintf(`{
	"input":{
		"tracking":[],
		"subscribe_location":"PROFILE",
		"story_id":null,
		"unsubscribee_id":"%s",
		"client_mutation_id":"%s",
		"is_tracking_encrypted":false,
		"attribution_id_v2":"%s",
		"actor_id":"%s"
	},
	"fetch_profile_context_row":true,
	"nt_context":{}
	}`, unsubscribee, mutationID, attribution, userID)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "ActorUnsubscribeCoreMutation")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "431523233186746802998151697")
	form.Set("variables", variables)
	form.Set("fb_api_analytics_tags", fmt.Sprintf(`["visitation_id=391724414624676:%s:0:%d","GraphServices"]`, visitID, time.Now().Unix()))
	form.Set("client_trace_id", traceID)

	req, _ := http.NewRequest("POST", "https://graph.facebook.com/graphql", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "ActorUnsubscribeCoreMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-privacy-context", "2368177546817046")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-ta-logging-ids", "graphql:"+traceID)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)                                          // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                          // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthunfollow()) //เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

	// ---------- SEND ----------
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

	_, err = db.Exec("DELETE FROM unsubscribee_id_table WHERE unsubscribee_id = ?", unsubscribee) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ unsubscribee_id_table ออกจากฐานข้อมูลแล้ว:", unsubscribee)
	}

}

func shortIDunfollow() string {
	id := uuid.New().String()
	return strings.ReplaceAll(id[:6], "-", "")
}

func generateAttributionunfollow() string {
	type frag struct {
		Name   string
		Source string
		Action string
	}

	fragments := []frag{
		{"GraphSearchFragment", "search_typeahead", "tap_search_bar"},
		{"SearchResultsFragment", "graph_search_results_page_blended", "tap_search_result"},
		{"ProfileFragment", "timeline", ""},
		{"ProfileFragment", "profile_vnext_tab_posts", ""},
	}

	var output []string
	now := time.Now().Unix()
	baseTime := float64(now) - rand.Float64()*300
	baseSession := rand.Intn(999999999)

	for i := len(fragments) - 1; i >= 0; i-- {
		f := fragments[i]
		timestamp := baseTime + float64(i*2) + rand.Float64()
		sessionID := baseSession + i*10000 + rand.Intn(9999)

		entry := []string{
			f.Name,
			f.Source,
			f.Action,
			fmt.Sprintf("%.3f", timestamp),
			fmt.Sprintf("%d", sessionID),
			shortIDunfollow() + shortIDunfollow(),
			"",
		}
		output = append(output, strings.Join(entry, ","))
	}

	return strings.Join(output, ";")
}
