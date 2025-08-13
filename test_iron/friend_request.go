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

func genVisitationIDfriend_request(userID string) string {
	hexPart := fmt.Sprintf("%x", rand.Intn(0xfffff))
	flags := rand.Intn(3) + 1
	ts := float64(time.Now().UnixNano()) / 1e9
	return fmt.Sprintf("%s:%s:%d:%.3f", userID, hexPart, flags, ts)
}

func getRandomUIDsfriend_request(n int) ([]string, error) {

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "uid_for_add_friend.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())

	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/uid_for_add_friend.db")

	rows, err := db.Query("SELECT user_id FROM uid_table ORDER BY RANDOM() LIMIT ?", n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uids []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			continue
		}
		uids = append(uids, uid)
	}
	return uids, nil
}

func deleteUIDfriend_request(uid string) {
	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "uid_for_add_friend.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/uid_for_add_friend.db")

	db.Exec("DELETE FROM uid_table WHERE user_id = ?", uid)
}

// Edit
func randomExcellentBandwidthfriend_request() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000 // 20 Mbps
	max := 35000000 // 35 Mbps
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runfriend_request(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	rand.Seed(time.Now().UnixNano())

	// สุ่มจำนวน UID (2–8) แล้วดึง
	limit := rand.Intn(7) + 2
	uids, err := getRandomUIDsfriend_request(limit)
	if err != nil || len(uids) == 0 {
		fmt.Println("❌ ไม่พบ UID ในฐานข้อมูล")
		return
	}
	fmt.Println("🎯 ดึงมา", len(uids), "ID:", uids)

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

	host := "graph.facebook.com"

	for _, targetID := range uids {
		delay := rand.Intn(5) + 1
		time.Sleep(time.Duration(delay) * time.Second)

		now := float64(time.Now().UnixNano()) / 1e9
		visitationID := genVisitationIDfriend_request(userID)
		attribution := strings.Join([]string{
			fmt.Sprintf("ProfileFragment,profile_vnext_tab_posts,,%.3f,235476155,,,", now),
			fmt.Sprintf("ProfileFragment,timeline,,%.3f,235476155,,,", now-0.6),
			fmt.Sprintf("NewsFeedFragment,native_newsfeed,tap_top_jewel_bar,%.3f,42846365,%s,,", now-3, userID),
		}, ";")

		variables := map[string]interface{}{
			"input": map[string]interface{}{
				"source":               "profile_button",
				"friend_requestee_ids": []string{targetID},
				"actor_id":             userID,
				"refs":                 []string{"pymk_feed"},
				"client_mutation_id":   uuid.New().String(),
				"attribution_id_v2":    attribution,
			},
		}
		varJSON, _ := json.Marshal(variables)

		analytics := []string{
			fmt.Sprintf("visitation_id=%s", visitationID),
			"GraphServices",
		}
		analyticsJSON, _ := json.Marshal(analytics)

		form := url.Values{}
		form.Set("method", "post")
		form.Set("pretty", "false")
		form.Set("format", "json")
		form.Set("server_timestamps", "true")
		form.Set("locale", "en_US")
		form.Set("fb_api_req_friendly_name", "FriendRequestSendCoreMutation")
		form.Set("fb_api_caller_class", "graphservice")
		form.Set("client_doc_id", "8268251071582849202978527632")
		form.Set("variables", string(varJSON))
		form.Set("fb_api_analytics_tags", string(analyticsJSON))

		req, err := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBufferString(form.Encode()))
		if err != nil {
			fmt.Println("❌ build req failed:", err)
			continue
		}

		req.Host = host
		req.Header.Set("Authorization", "OAuth "+accessToken)
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Host", host)
		req.Header.Set("X-FB-Background-State", "1")
		req.Header.Set("x-fb-client-ip", "True")
		req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
		req.Header.Set("x-fb-device-group", devicegroup)
		req.Header.Set("X-FB-Friendly-Name", "FriendingJewelContentQuery")
		req.Header.Set("X-FB-HTTP-Engine", "Liger")
		req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","purpose":"fetch","request_category":"graphql","retry_attempt":"0"},"application_tags":"graphservice"}`)
		req.Header.Set("x-fb-server-cluster", "True")
		req.Header.Set("x-graphql-client-library", "graphservice")
		req.Header.Set("x-graphql-request-purpose", "fetch")
		req.Header.Set("x-tigon-is-retry", "False")
		req.Header.Set("x-fb-net-hni", netHni)
		req.Header.Set("x-fb-sim-hni", simHni)
		req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthfriend_request())
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

		fmt.Println("🎯", targetID, "=>", resp.Status)

		deleteUIDfriend_request(targetID)
	}
}
