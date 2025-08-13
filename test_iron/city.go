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

func randomExcellentBandwidthcity() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runcity(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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

	var city string
	err = db.QueryRow("SELECT city_id FROM change_city_table LIMIT 1").Scan(
		&city)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล change_city_table ไม่สำเร็จ: " + err.Error())
		return
	}

	// ===== สร้าง body
	clientMutationID := uuid.New().String()
	//---------------------------------------------------------------------------------------https://www.facebook.com/places--------------------------------------------------------------
	variables := fmt.Sprintf(`{"input":{"client_mutation_id":"%s","actor_id":"%s","hometown_city_id":"%s","privacy":{"allow":[],"base_state":"EVERYONE","deny":[],"tag_expansion_state":"UNSPECIFIED"}}}`, clientMutationID, userID, city)
	analyticsTags := fmt.Sprintf(`["visitation_id=281710865595635:ef4bb:1:%d","session_id=UFS-%s","nav_attribution_id="]`, time.Now().Unix(), uuid.New().String())

	form := url.Values{}
	form.Set("access_token", accessToken)
	form.Set("fb_api_caller_class", "RelayModern")
	form.Set("fb_api_req_friendly_name", "ProfileEditHometownSaveMutation")
	form.Set("variables", variables)
	form.Set("server_timestamps", "true")
	form.Set("doc_id", "2186879678019893")
	form.Set("fb_api_analytics_tags", analyticsTags)

	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	_, _ = gz.Write([]byte(form.Encode()))
	_ = gz.Close()

	// ===== สร้าง HTTP request
	req, _ := http.NewRequest("POST", "https://"+host+"/graphql?locale=en_US", &compressed)
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("X-FB-Friendly-Name", "RelayFBNetwork_ProfileEditHometownSaveMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","retry_attempt":"0"},"application_tags":"unknown"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-zero-eh", "2,,ASAySH2oMJyaEZ8_gJqsBkU9XiRTbXAuIxwORpw9LUATP-gBzQrZQZq8gPwqoERuScM")
	req.Header.Set("X-ZERO-STATE", "unknown")
	req.Header.Set("x-fb-net-hni", netHni)                                      // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                      // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthcity()) //เพิ่มเข้าไป
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

	_, err = db.Exec("DELETE FROM change_city_table WHERE city_id = ?", city) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ change_city_table ออกจากฐานข้อมูลแล้ว:", city)
	}

}
