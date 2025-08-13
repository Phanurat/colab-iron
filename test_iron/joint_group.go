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

func randomExcellentBandwidthjoint_group() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runjoint_group(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var groupID string
	err = db.QueryRow("SELECT group_id FROM group_id_table LIMIT 1").Scan(
		&groupID)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล group_id_table ไม่สำเร็จ: " + err.Error())
		return
	}

	mutationID := uuid.New().String()
	traceID := uuid.New().String()
	visitID := shortIDjoint_group()
	attribution := randomAttributionjoint_group()
	bloksVersion := generateBloksVersionjoint_group()
	stylesID := uuid.New().String()

	variables := fmt.Sprintf(`{"input":{"tracking_codes":[],"source":"search","group_id":"%s","client_mutation_id":"%s","attribution_id_v2":"%s","actor_id":"%s"},"nt_context":{"using_white_navbar":true,"pixel_ratio":3,"is_push_on":true,"styles_id":"%s","bloks_version":"%s"}}`, groupID, mutationID, attribution, userID, stylesID, bloksVersion)

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "GroupJoinForumMutation")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "18856738412634049533340475479")
	form.Set("variables", variables)
	form.Set("fb_api_analytics_tags", fmt.Sprintf(`["visitation_id=391724414624676:%s:1:%d","GraphServices"]`, visitID, time.Now().Unix()))
	form.Set("client_trace_id", traceID)

	req, _ := http.NewRequest("POST", "https://graph.facebook.com/graphql", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "GroupJoinForumMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Navigation-Chain", randomNavChainjoint_group())
	req.Header.Set("x-fb-privacy-context", "3516766875019643")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-ta-logging-ids", "graphql:"+traceID)
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-fb-net-hni", netHni)                                             // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                             // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthjoint_group()) //เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-quality", "EXCELLENT")

	// ---------- SEND ----------
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

	_, err = db.Exec("DELETE FROM group_id_table WHERE group_id = ?", groupID) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ group_id_table ออกจากฐานข้อมูลแล้ว:", groupID)
	}

}

func shortIDjoint_group() string {
	id := uuid.New().String()
	return strings.ReplaceAll(id[:6], "-", "")
}

func generateBloksVersionjoint_group() string {
	return uuidPartjoint_group() + uuidPartjoint_group()
}

func uuidPartjoint_group() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func randomAttributionjoint_group() string {
	return "SearchResultsFragment,graph_search_results_page_blended,tap_search_result,1749" + shortIDjoint_group() + ".918,58226679,,,;" +
		"SuggestionsFragment,search_typeahead,,1749" + shortIDjoint_group() + ".265,158688020,,,;" +
		"SearchResultsFragment,graph_search_results_page_blended,tap_search_result,1749" + shortIDjoint_group() + ".157,217218979,,,;" +
		"GraphSearchFragment,search_typeahead,tap_search_bar,1749" + shortIDjoint_group() + ".344,56804238,391724414624676,,"
}

func randomNavChainjoint_group() string {
	return randomAttributionjoint_group()
}
