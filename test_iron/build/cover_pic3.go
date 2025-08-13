package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/fs"
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

func randomExcellentBandwidthcover_pic3() string {
	rand.Seed(time.Now().UnixNano())
	min := 20000000
	max := 35000000
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func generateUploadIDscover_pic3(fileSize int64) (string, string, string, string) {
	sessionID := uuid.New().String()
	timestamp := time.Now().UnixMilli()
	uuidPart := strings.ReplaceAll(uuid.New().String(), "-", "")
	entityName := fmt.Sprintf("%s-0-%d-%d-%d", uuidPart[:32], fileSize, timestamp, timestamp)
	idempotence := fmt.Sprintf("%s_%d_0", sessionID, time.Now().UnixNano()%1000000000)
	waterfallID := fmt.Sprintf("%s_996C60C8D20C_Mixed_0", sessionID)
	return sessionID, idempotence, waterfallID, entityName
}

func getLatestPhotoIDcover_pic3() string {
	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	dbPath := filepath.Join(folder, "cover_photo_id.db")
	fmt.Println("📂 DB PATH:", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("❌ ดึงฐานข้อมูลไม่สำเร็จ: " + err.Error())
	}
	defer db.Close()

	fmt.Println("📂 DB PATH:", folder+"/cover_photo_id.db")

	row := db.QueryRow("SELECT cover_pic_id FROM cover_photo_id_table ORDER BY id DESC LIMIT 1")
	var photoID string
	err = row.Scan(&photoID)
	if err != nil {
		fmt.Println("❌ ดึง pic_id ไม่สำเร็จ: " + err.Error())
	}
	return photoID
}

func findSmallestImagePathcover_pic3(folder string) string {

	var smallest string
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".jpg") || strings.HasSuffix(strings.ToLower(d.Name()), ".jpeg") {
			if smallest == "" || d.Name() < smallest {
				smallest = d.Name()
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("❌ WalkDir fail: " + err.Error())
	}
	if smallest == "" {
		fmt.Println("❌ ไม่มีไฟล์ .jpg ในโฟลเดอร์เลย")
	}
	return filepath.Join(folder, smallest)
}

func Runcover_pic3(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	fmt.Println("🌐 ใช้ Proxy:", proxyAddr) // ✅ เพิ่ม debug แสดง proxy ที่ใช้อยู่
	host := "rupload.facebook.com"

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}

	profileFolder := filepath.Join("cover_photo")
	imagePath := findSmallestImagePathcover_pic3(profileFolder)

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

	// profileFolder := filepath.Join("cover_photo")
	// imagePath := findSmallestImagePathcover_pic3(profileFolder)

	imgFile, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("❌ Open image fail: " + err.Error())
		return
	}

	// Decode → Encode
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("❌ Decode image fail: " + err.Error())
		return
	}

	// ✅ Close ทันทีหลังใช้เสร็จ
	imgFile.Close()

	var imgBuf bytes.Buffer
	if err := jpeg.Encode(&imgBuf, img, &jpeg.Options{Quality: 90}); err != nil {
		fmt.Println("❌ Encode JPEG fail: " + err.Error())
		return
	}

	photoID := getLatestPhotoIDcover_pic3()

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("❌ Open image fail: " + err.Error())
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fileBuffer := new(bytes.Buffer)
	io.Copy(fileBuffer, file)

	file.Close()

	composerSessionID, idempotenceToken, waterfallID, filename := generateUploadIDscover_pic3(fileSize)
	url := "https://" + host + "/fb_photo/" + filename

	req, err := http.NewRequest("POST", url, bytes.NewReader(fileBuffer.Bytes()))
	if err != nil {
		fmt.Println("❌ Request build fail: " + err.Error())
		return
	}

	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", "OAuth "+accessToken)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fileSize))
	req.Header.Set("Host", host)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("add_to_photo_id", photoID)
	req.Header.Set("audience_exp", "true")
	req.Header.Set("composer_session_id", composerSessionID)
	req.Header.Set("idempotence_token", idempotenceToken)
	req.Header.Set("published", "false")
	req.Header.Set("qn", composerSessionID)
	req.Header.Set("Offset", "0")
	req.Header.Set("Segment-Start-Offset", "0")
	req.Header.Set("Segment-Type", "3")
	req.Header.Set("target_id", userID)
	req.Header.Set("X-Entity-Length", fmt.Sprintf("%d", fileSize))
	req.Header.Set("X-Entity-Name", filename)
	req.Header.Set("X-Entity-Type", "image/jpeg")
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-connection-bandwidth", randomExcellentBandwidthcover_pic3())
	req.Header.Set("X-FB-Connection-Quality", "EXCELLENT")
	req.Header.Set("X-FB-Connection-Type", "MOBILE.HSDPA")
	req.Header.Set("x-fb-device-group", devicegroup)
	req.Header.Set("X-FB-Friendly-Name", "Resumable-Upload-Post")
	req.Header.Set("X-FB-HTTP-Engine", "Liger")
	req.Header.Set("X-FB-Request-Analytics-Tags", `{"network_tags":{"product":"350685531728","retry_attempt":"0"},"application_tags":"unknown"}`)
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("x-tigon-is-retry", "False")
	req.Header.Set("x-zero-eh", "2,,ASAySH2oMJyaEZ8_gJqsBkU9XiRTbXAuIxwORpw9LUATP-gBzQrZQZq8gPwqoERuScM")
	req.Header.Set("X_FB_PHOTO_WATERFALL_ID", waterfallID)
	req.Header.Set("Zero-Rated", "0")
	req.Header.Set("x-fb-net-hni", netHni) // เพิ่มเข้าไป
	req.Header.Set("x-fb-sim-hni", simHni)

	// ---------- SEND ----------
	bw := tlsConns.RWrupload.Writer
	br := tlsConns.RWrupload.Reader

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

	if err := os.Remove(imagePath); err != nil {
		fmt.Println("⚠️ ลบไฟล์ไม่สำเร็จ:", err)
	} else {
		fmt.Println("🗑️ ลบไฟล์ที่ใช้แล้ว:", imagePath)
	}
}
