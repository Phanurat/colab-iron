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

func generateNavChainRunshared_link(actorID string) string {
	timestamp := float64(time.Now().UnixNano()) / 1e9
	return fmt.Sprintf(
		"ComposerActivity,composer,tap_status_button,%.3f,250413799,,,;ProfileFragment,profile_vnext_tab_posts,,%.3f,256422804,,,;ProfileFragment,timeline,tap_bookmark,%.3f,256422804,%s,,;BookmarkComponentsFragment,bookmarks,tap_top_jewel_bar,%.3f,249114123,281710865595635,,",
		timestamp,
		timestamp-10,
		timestamp-11,
		timestamp-12,
	)
}

func randomExcellentBandwidthRunshared_link() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runshared_link(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
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
	var accessToken, userID, userAgent, netHni, simHni, devicegroup string
	err = db.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
		&accessToken, &userID, &userAgent, &netHni, &simHni, &devicegroup)
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล app_profiles ไม่สำเร็จ: " + err.Error())
		return
	}

	var link string                                                                  //statusLink
	err = db.QueryRow("SELECT link_link FROM shared_link_table LIMIT 1").Scan(&link) //, &statusLink) //, status_link
	if err != nil {
		fmt.Println("❌ ดึงข้อมูล shared_link_table ไม่สำเร็จ: " + err.Error())
		return
	}

	variablesObj := map[string]interface{}{
		"input": map[string]interface{}{
			"message": map[string]string{
				"text": "",
			},
			"composer_session_id": uuidRunshared_link(),
			"idempotence_token":   "FEED_" + uuidRunshared_link(),
			"client_mutation_id":  uuidRunshared_link(),
			"actor_id":            userID,
			"audiences": []interface{}{
				map[string]interface{}{
					"undirected": map[string]interface{}{
						"privacy": map[string]interface{}{
							"tag_expansion_state": "UNSPECIFIED",
							"deny":                []string{},
							"base_state":          "EVERYONE", ///EVERYONE //FRIENDS
							"allow":               []string{},
						},
					},
				},
			},
			"source": "MOBILE",
			"attachments": []interface{}{
				map[string]interface{}{
					"link": map[string]interface{}{
						"external": map[string]string{
							"url": link,
						},
					},
				},
			},
			"action_timestamp": time.Now().Unix(),
		},
	}

	escapedVariables, err := json.Marshal(variablesObj)
	if err != nil {
		fmt.Println("❌ แปลง variables ไม่ได้: " + err.Error())
		return
	}

	form := url.Values{}
	form.Set("method", "post")
	form.Set("pretty", "false")
	form.Set("format", "json")
	form.Set("server_timestamps", "true")
	form.Set("locale", "en_US")
	form.Set("fb_api_req_friendly_name", "ComposerStoryCreateMutation")
	form.Set("fb_api_caller_class", "graphservice")
	form.Set("client_doc_id", "91093790612716748765152950249")
	form.Set("variables", string(escapedVariables))

	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	_, _ = gz.Write([]byte(form.Encode()))
	_ = gz.Close()

	host := "graph.facebook.com"

	req, _ := http.NewRequest("POST", "https://"+host+"/graphql", &compressed)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-FB-Friendly-Name", "ComposerStoryCreateMutation")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-graphql-client-library", "graphservice")
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-fb-ta-logging-ids", fmt.Sprintf("graphql:%s", uuidRunshared_link()))
	req.Header.Set("x-fb-privacy-context", "496463117678580")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"none","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
	req.Header.Set("X-FB-Navigation-Chain", generateNavChainRunshared_link(userID))
	req.Header.Set("x-fb-net-hni", netHni)                                                // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                                // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthRunshared_link()) //เพิ่มเข้าไป
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

	_, err = db.Exec("DELETE FROM shared_link_table WHERE link_link = ?", link) //db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid)
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ shared_link_table ออกจากฐานข้อมูลแล้ว:", link)
	}

	//	_, err = dbApp.Exec(`UPDATE comments SET status_link = NULL WHERE status_link = ?`, statusLink)
	//	if err != nil {
	//		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	//	} else {
	//		fmt.Println("🧹 ลบ status_link ออกจากฐานข้อมูลแล้ว:", statusLink)
	//	}

}

func uuidRunshared_link() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		rand32Runshared_link(), rand16Runshared_link(), 0x4000|rand16Runshared_link()&0x0fff, 0x8000|rand16Runshared_link()&0x3fff, rand48Runshared_link())
}
func rand16Runshared_link() int {
	return int(time.Now().UnixNano()>>16) & 0xffff
}
func rand32Runshared_link() uint32 {
	return uint32(time.Now().UnixNano()>>8) & 0xffffffff
}
func rand48Runshared_link() int64 {
	return int64(time.Now().UnixNano()) & 0xffffffffffff
}
