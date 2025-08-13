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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func genVisitationIDfriend_accept3(actorID string) string {
	hashPart := make([]rune, 5)
	chars := []rune("abcdef0123456789")
	for i := range hashPart {
		hashPart[i] = chars[rand.Intn(len(chars))]
	}
	flags := rand.Intn(3) + 1
	ts := float64(time.Now().UnixNano()) / 1e9
	return fmt.Sprintf("%s:%s:%d:%.3f", actorID, string(hashPart), flags, ts)
}

func genAttributionIDfriend_accept3() string {
	parts := []string{
		"FriendingJewelFragment",
		"friend_requests",
		"tap_top_jewel_bar",
		fmt.Sprintf("%.3f", float64(time.Now().UnixNano())/1e9),
		fmt.Sprintf("%d", rand.Intn(899999999)+100000000),
		fmt.Sprintf("%d", rand.Int63n(899999999999999)+100000000000000),
		"",
	}
	return strings.Join(parts, ",")
}

// Edit
func randomExcellentBandwidthfriend_accept3() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000 // 20 Mbps
	max := 35000000 // 35 Mbps
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runfriend_accept3(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	host := "graph.facebook.com"
	//	address := host + ":443"

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "friend.db") + "?_busy_timeout=5000&_journal_mode=WAL"
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
		return
	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/friend.db")

	//	dir, _ := os.Getwd()
	//	dbPath := filepath.Join(dir, "friend.db") + "?_busy_timeout=5000&_journal_mode=WAL"
	//	db, err := sql.Open("sqlite3", dbPath)
	//	if err != nil {
	//		panic("❌ เปิดฐานข้อมูลไม่ได้: " + err.Error())
	//	}
	//	defer db.Close()

	rows, err := db.Query(`SELECT friend_requester_id FROM friend_info`)
	if err != nil {
		fmt.Println("❌ อ่าน friend_requester_id ไม่ได้: " + err.Error())
		return
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		fmt.Println("📭 ไม่มี friend_requester_id ให้ยิง")
		return
	}

	dbApp := filepath.Join(folder, "fb_comment_system.db")
	fmt.Println("📂 DB PATH:", dbApp)

	dbacc, err := sql.Open("sqlite3", dbApp)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
		return
	}
	defer dbacc.Close()

	fmt.Println("📂 DB PATH:", folder+"/fb_comment_system.db")

	var accessToken, userID, userAgent, netHni, simHni, devicegroup string
	err = dbacc.QueryRow("SELECT access_token, actor_id, user_agent, net_hni, sim_hni, device_group FROM app_profiles LIMIT 1").Scan(
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

	for _, friendRequesterID := range ids {
		visitationID := genVisitationIDfriend_accept3(userID)
		attributionID := genAttributionIDfriend_accept3()
		idempotenceToken := uuid.New().String()
		clientMutationID := uuid.New().String()

		variables := map[string]interface{}{
			"input": map[string]interface{}{
				"attribution_id_v2":   attributionID,
				"friend_requester_id": friendRequesterID,
				"source":              "friends_home_main",
				"idempotence_token":   idempotenceToken,
				"client_mutation_id":  clientMutationID,
				"origin":              "FRIENDING_TAB_OPEN",
				"actor_id":            userID,
			},
		}

		variablesJSON, _ := json.Marshal(variables)
		analyticsTags := []string{fmt.Sprintf("visitation_id=%s", visitationID), "GraphServices"}
		analyticsJSON, _ := json.Marshal(analyticsTags)

		payload := fmt.Sprintf(
			"method=post&pretty=false&format=json&server_timestamps=true&locale=en_US"+
				"&fb_api_req_friendly_name=FriendRequestAcceptCoreMutation"+
				"&fb_api_caller_class=graphservice"+
				"&client_doc_id=38817391810048484601801151473"+
				"&variables=%s"+
				"&fb_api_analytics_tags=%s",
			string(variablesJSON), string(analyticsJSON))

		req, err := http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBufferString(payload))
		if err != nil {
			fmt.Println("❌ new request fail:", err)
			//return
			continue
		}

		//	req, _ = http.NewRequest("POST", "https://"+host+"/graphql", bytes.NewBufferString(payload)) // ✅ ใช้ = แทน :=
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
		req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthfriend_accept3())
		req.Header.Set("x-fb-connection-quality", "EXCELLENT")

		// ---------- SEND ----------
		bw := tlsConns.RWGraph.Writer
		br := tlsConns.RWGraph.Reader

		err = req.Write(bw)
		if err != nil {
			fmt.Println("❌ Write fail: " + err.Error())
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

		_, err = db.Exec(`DELETE FROM friend_info WHERE friend_requester_id = ?`, friendRequesterID)
		if err != nil {
			fmt.Println("❌ ลบไม่สำเร็จ:", err)
		} else {
			fmt.Println("🧹 ลบ friend_id ออกจากฐานข้อมูลแล้ว:", friendRequesterID)
		}

		delay := time.Duration(rand.Intn(5)+1) * time.Second
		fmt.Printf("⏱️ รอ %.0f วินาทีก่อนยิงคนถัดไป...\n", delay.Seconds())
		time.Sleep(delay)
	}

}
