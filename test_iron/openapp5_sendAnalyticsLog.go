package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ปรับค่าตรงนี้
var (
	//	deviceGroup = "6301"
	hostopenapp5     = "graph.facebook.com"
	boundaryopenapp5 = "a_2r0-OU6lNDmY8oQ4UAX1HcrxDbdxqJ"
)

func randomExcellentBandwidthopenapp5() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000 // 20 Mbps
	max := 35000000 // 35 Mbps
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func Runopenapp5_sendAnalyticsLog(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	// mock ตัวอย่าง cmsg (gzip/binary); แทนที่ด้วย []byte ที่ต้องการ
	cmsg := mockCmsg()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.SetBoundary(boundaryopenapp5)

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

	// ฟิลด์ text
	formFields := map[string]string{
		"compressed":               "0",
		"cmethod":                  "gzip",
		"multi_batch":              "1",
		"locale":                   "en_US",
		"client_country_code":      "TH",
		"fb_api_req_friendly_name": "sendAnalyticsLog",
		"fb_api_caller_class":      "FbHttpUploader",
		"access_token":             "350685531728|62f8ce9f74b12f84c123cc23437a4a32",
	}

	for k, v := range formFields {
		w, _ := writer.CreateFormField(k)
		w.Write([]byte(v))
	}

	// ฟิลด์ไฟล์
	cmsgHeader := make(textproto.MIMEHeader)
	cmsgHeader.Set("Content-Disposition", `form-data; name="cmsg"; filename="message"`)
	cmsgHeader.Set("Content-Type", "application/octet-stream")
	cmsgHeader.Set("Content-Transfer-Encoding", "binary")
	part, _ := writer.CreatePart(cmsgHeader)
	part.Write(cmsg)

	writer.Close()

	//	address := host + ":443"
	//	var conn net.Conn

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

	req, _ := http.NewRequest("POST", "https://"+hostopenapp5+"/logging_client_events", bytes.NewBuffer(body.Bytes()))
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundaryopenapp5)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", hostopenapp5)
	req.Header.Set("X-FB-Background-State", "1")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("X-FB-Connection-Bandwidth", "34131200")
	req.Header.Set("X-FB-Connection-Quality", "EXCELLENT")
	// req.Header.Set("X-FB-Connection-Type", "WIFI")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "sendAnalyticsLog")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","request_category":"analytics","retry_attempt":"0"},"application_tags":"unknown"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-zero-eh", "2,,AS_f4gKEkvTH3SfgCYJVzirQUTh0TLneLfitj0JQoYbj90OG3tBV9erihGNv-2EK4YE")
	req.Header.Set("X-ZERO-STATE", "unknown")

	req.Header.Set("x-fb-net-hni", netHni)                                          // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)                                          // เพิ่มเข้าไป
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthopenapp5()) //เพิ่มเข้าไป
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

// ใช้แทน payload จริง ถ้าจะส่ง gzip หรือ payload จริงใส่ []byte ตรงนี้เอง
func mockCmsg() []byte {
	b := make([]byte, 80)
	rand.Read(b)
	return b
}
