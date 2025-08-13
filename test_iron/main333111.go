package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	utls "github.com/refraction-networking/utls"
)

var fullProxy string

type TLSConnections struct {
	RWGraphConn   io.Closer
	RWBGraphConn  io.Closer
	RWGatewayConn io.Closer // "static.xx.fbcdn.net"
	RWWebConn     io.Closer // ✅ เปลี่ยนชื่อให้ไม่ซ้ำ rupload.facebook.com
	RWruploadConn io.Closer
	RWstaticConn  io.Closer

	RWGraph   *bufio.ReadWriter
	RWBGraph  *bufio.ReadWriter
	RWGateway *bufio.ReadWriter
	RWWeb     *bufio.ReadWriter // ✅ เปลี่ยนชื่อให้ไม่ซ้ำ
	RWrupload *bufio.ReadWriter
	RWstatic  *bufio.ReadWriter
}

func dialTLS(host string) (io.Closer, *bufio.ReadWriter, error) {
	rawConn, err := net.Dial("tcp", host+":443")
	if err != nil {
		return nil, nil, err
	}
	config := &utls.Config{ServerName: host}
	conn := utls.UClient(rawConn, config, utls.HelloAndroid_11_OkHttp)

	if err := conn.Handshake(); err != nil {
		return nil, nil, err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	return conn, rw, nil
}

func initTLSConns() (*TLSConnections, error) {
	gConn, gRW, err := dialTLS("graph.facebook.com")
	if err != nil {
		return nil, fmt.Errorf("graph handshake failed: %v", err)
	}
	bgConn, bgRW, err := dialTLS("graph.facebook.com")
	if err != nil {
		return nil, fmt.Errorf("b-graph handshake failed: %v", err)
	}
	// ✅ เพิ่ม handshake กับ gateway.facebook.com
	gatewayConn, gatewayRW, err := dialTLS("gateway.facebook.com")
	if err != nil {
		return nil, fmt.Errorf("gateway handshake failed: %v", err)
	}

	webConn, webRW, err := dialTLS("web.facebook.com")
	if err != nil {
		return nil, fmt.Errorf("web handshake failed: %v", err)
	}

	ruploadConn, ruploadRW, err := dialTLS("rupload.facebook.com")
	if err != nil {
		return nil, fmt.Errorf("web handshake failed: %v", err)
	}

	staticConn, staticRW, err := dialTLS("static.xx.fbcdn.net")
	if err != nil {
		return nil, fmt.Errorf("web handshake failed: %v", err)
	}

	return &TLSConnections{
		RWGraphConn:   gConn,
		RWBGraphConn:  bgConn,
		RWGatewayConn: gatewayConn, // ✅ เพิ่ม
		RWWebConn:     webConn,     // ✅ ชื่อตรงกับ struct
		RWruploadConn: ruploadConn,
		RWstaticConn:  staticConn,

		RWGraph:   gRW,
		RWBGraph:  bgRW,
		RWGateway: gatewayRW, // ✅ เพิ่ม
		RWWeb:     webRW,     // ✅ ชื่อตรงกับ struct
		RWrupload: ruploadRW,
		RWstatic:  staticRW,
	}, nil
}

func closeConns(tlsConns *TLSConnections) {
	tlsConns.RWGraphConn.Close()
	tlsConns.RWBGraphConn.Close()
	tlsConns.RWGatewayConn.Close() // ✅ ปิด connection ที่เพิ่ม
	tlsConns.RWWebConn.Close()
	tlsConns.RWruploadConn.Close()
	tlsConns.RWstaticConn.Close()
}

func init() {
	// ดึง DB path จาก ENV
	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}
	dbPath := filepath.Join(folder, "fb_comment_system.db")

	// เปิดฐานข้อมูล
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("❌ เปิด DB ไม่ได้: %v", err)
	}
	defer db.Close()

	// ดึง proxy info จาก proxy_table
	var ip string

	err = db.QueryRow(`SELECT proxy_key FROM proxy_table LIMIT 1`).Scan(
		&ip)
	if err != nil {
		log.Fatalf("❌ ดึง proxy ไม่สำเร็จ: %v", err)
	}

	fullProxy = fmt.Sprintf(ip)
	fmt.Println("✅ fullProxy:", fullProxy)
}

// legit
func parseProxy(full string) (addr string, auth string) {
	parts := strings.Split(full, "@")
	if len(parts) != 2 {
		panic("❌ รูปแบบ proxy ต้องเป็น user:pass@ip:port")
	}
	auth = base64.StdEncoding.EncodeToString([]byte(parts[0]))
	addr = parts[1]
	return
}

// legit
//func runScript(path, proxyAddr, proxyAuth string) {
//	fmt.Printf("▶️ รัน: %s\n", path)
//	cmd := exec.Command("./" + path) // <== เพิ่ม ./ นี่แหละจุดสำคัญ
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//	cmd.Env = append(os.Environ(),
//		"USE_PROXY="+proxyAddr,
//		"USE_PROXY_AUTH="+proxyAuth,
//	)
//	if err := cmd.Run(); err != nil {
//		fmt.Printf("❌ พังที่ %s : %v\n", path, err)
//	}
//}

// auto
func loopJewel(tlsConns *TLSConnections, proxyAddr, proxyAuth string) {
	for {
		RunJewel(tlsConns, proxyAddr, proxyAuth)

		// สุ่มเวลา 1 - 13 วิ
		delay := time.Duration(rand.Intn(13)+1) * time.Second
		time.Sleep(delay)
	}
}

// // auto
// //func sendAnalyticsLog(proxyAddr, proxyAuth string) {
// //	for {
// //		runScript("sendAnalyticsLog.exe", proxyAddr, proxyAuth)
// //		delay := time.Duration(rand.Intn(9)+6) * time.Second
// //		time.Sleep(delay)
// //		fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ...\n", delay.Seconds())
// //	}
// //}

// auto
func loopSendPing(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	for {
		if *active {
			fmt.Println("▶️ [1]ส่งรีเควสสสสสสสสสสส SendPing")
			RunsendPing(tlsConns, proxyAddr, proxyAuth)
		}
		time.Sleep(60 * time.Second) // ยิงทุก 1 นาที เฉพาะตอนแอพเปิด
	}
}

//func runSequence(files []string, proxyAddr, proxyAuth string) {
//	for _, exe := range files {
//		fmt.Println("▶️ รัน:", exe)
//		cmd := exec.Command("./" + exe) // ← สำคัญตรงนี้
//		cmd.Stdout = os.Stdout
//		cmd.Stderr = os.Stderr
//		cmd.Env = append(os.Environ(),
//			"USE_PROXY="+proxyAddr,
//			"USE_PROXY_AUTH="+proxyAuth,
//		)
//		err := cmd.Run()
//		if err != nil {
//			fmt.Println("❌ พังที่", exe, ":", err)
//		} else {
//			fmt.Println("✅ สำเร็จ:", exe)
//		}
//	}
//}

// เปิด
func OpenApp(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	fmt.Println("🚀 เปิดแอพ")

	Runopenapp1(tlsConns, proxyAddr, proxyAuth)
	Runopenapp2(tlsConns, proxyAddr, proxyAuth)
	Runopenapp3(tlsConns, proxyAddr, proxyAuth)
	Runopenapp4(tlsConns, proxyAddr, proxyAuth)
	Runopenapp5_sendAnalyticsLog(tlsConns, proxyAddr, proxyAuth)
	Runopenapp6(tlsConns, proxyAddr, proxyAuth)
	Runopenapp7(tlsConns, proxyAddr, proxyAuth)
	Runopenapp8(tlsConns, proxyAddr, proxyAuth)
	Runopenapp9(tlsConns, proxyAddr, proxyAuth)
	Runopenapp10(tlsConns, proxyAddr, proxyAuth)
	Runopenapp11(tlsConns, proxyAddr, proxyAuth)
	Runopenapp12(tlsConns, proxyAddr, proxyAuth)
	Runopenapp14(tlsConns, proxyAddr, proxyAuth)

	delay := time.Duration(rand.Intn(9)+6) * time.Second
	time.Sleep(delay)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ...\n", delay.Seconds())
}

// // เปิด
// // ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func friend_accept(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {

	fmt.Println("👥 รับเพื่อน")

	Runfriend_accept1(tlsConns, proxyAddr, proxyAuth)
	Runfriend_accept2(tlsConns, proxyAddr, proxyAuth)
	Runfriend_accept3(tlsConns, proxyAddr, proxyAuth)

	delay2 := time.Duration(rand.Intn(5)+1) * time.Second
	time.Sleep(delay2)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1)\n", delay2.Seconds())
}

// เปิด ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func friend_requester(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	fmt.Println("👥 แอดเพื่อน")

	Runfriend_request(tlsConns, proxyAddr, proxyAuth)

	delay3 := time.Duration(rand.Intn(5)+1) * time.Second
	time.Sleep(delay3)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1)\n", delay3.Seconds())

}

// เอาไปสุ่ม /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func maket(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string, duration time.Duration) {
	start := time.Now()
	fmt.Println("🛒 เริ่มโหมด market")

	for time.Since(start) < duration {
		// เช็กเวลาว่าเหลือพอไหมก่อนเริ่มรอบใหม่
		remaining := duration - time.Since(start)
		if remaining < 10*time.Second {
			fmt.Println("⏹️ เวลาน้อยกว่า 10 วิแล้ว ออกจาก maket")
			break
		}

		fmt.Println("🛒 เข้า market รอบใหม่")

		Runmaket1(tlsConns, proxyAddr, proxyAuth)
		Runmaket2(tlsConns, proxyAddr, proxyAuth)
		Runmaket3(tlsConns, proxyAddr, proxyAuth)
		Runmaket4(tlsConns, proxyAddr, proxyAuth)
		Runmaket5(tlsConns, proxyAddr, proxyAuth)
		Runmaket6(tlsConns, proxyAddr, proxyAuth)
		Runmaket7_more(tlsConns, proxyAddr, proxyAuth)

		// ✅ สุ่มรอรอบแรก
		delay := time.Duration(rand.Intn(51)+10) * time.Second
		if time.Since(start)+delay >= duration {
			fmt.Println("⏹️ จะหมดเวลาก่อนเริ่ม maket7_more.exe loop หยุดก่อน")
			break
		}
		time.Sleep(delay)

		loopCount := rand.Intn(6) + 3
		for i := 1; i <= loopCount; i++ {
			if time.Since(start) >= duration {
				fmt.Println("⏹️ เวลาหมดในรอบ maket7_more.exe ออก")
				return
			}
			Runmaket7_more(tlsConns, proxyAddr, proxyAuth)

			delay := time.Duration(rand.Intn(51)+8) * time.Second
			if time.Since(start)+delay >= duration {
				fmt.Printf("⏹️ จะหมดเวลาในรอบ %d / %d หยุดก่อน\n", i, loopCount)
				return
			}
			time.Sleep(delay)
		}
	}
	fmt.Println("✅ ออกจากโหมด market แล้ว")
}

// // เอาไปสุ่ม/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func see_story(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	// 	//	for {
	fmt.Println("🛒 เข้า see_story (9%)")

	Runsee_story1_FBStoriesAdsPaginatingQuery(tlsConns, proxyAddr, proxyAuth)
	Runsee_story2_FBStoriesAdsPaginatingQuery_At_Connection(tlsConns, proxyAddr, proxyAuth)
	Runsee_story3_FbStoriesUnifiedSingleBucketQuery(tlsConns, proxyAddr, proxyAuth)

	// ✅ สุ่มรอรอบแรก (10–60 วิ)
	delay := time.Duration(rand.Intn(51)+10) * time.Second
	time.Sleep(delay)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนเริ่ม loop ดู story เพิ่มเติม\n", delay.Seconds())

	// ✅ สุ่มจำนวนรอบที่ยิง see_story3 (เช่น 3–10 ครั้ง)
	loopCount := rand.Intn(6) + 2 // 3 ถึง 10

	for i := 1; i <= loopCount; i++ {

		Runsee_story3_FbStoriesUnifiedSingleBucketQuery(tlsConns, proxyAddr, proxyAuth)

		delayEach := time.Duration(rand.Intn(61)+7) * time.Second // 10–100 วิ
		time.Sleep(delayEach)
		fmt.Printf("⏱️ รอบ %d / %d: รอ %.0f วินาทีก่อนดู story ต่อไป\n", i, loopCount, delayEach.Seconds())
		//		}
	}
}

// // ก้อนเม้น/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func see_watch_start(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	//	for {
	fmt.Println("🛒 เข้า see_watch_start")

	Runsee_watch1_VideoHomeFeedVisitMutation(tlsConns, proxyAddr, proxyAuth)
	Runsee_watch2_SurveyIntegrationPointQuery(tlsConns, proxyAddr, proxyAuth)

	delay := time.Duration(rand.Intn(100)+10) * time.Second
	time.Sleep(delay)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1/4)\n", delay.Seconds())

	Runsee_watch3_FbShortsImpressionMutation(tlsConns, proxyAddr, proxyAuth)

	delay2 := time.Duration(rand.Intn(30)+10) * time.Second
	time.Sleep(delay2)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2/4)\n", delay2.Seconds())

	Runsee_watch3_FbShortsImpressionMutation(tlsConns, proxyAddr, proxyAuth)

	delay3 := time.Duration(rand.Intn(20)+10) * time.Second
	time.Sleep(delay3)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 3/4)\n", delay3.Seconds())

	Runsee_watch4_VideoHomeSectionQuery(tlsConns, proxyAddr, proxyAuth)

	delay4 := time.Duration(rand.Intn(50)+10) * time.Second
	time.Sleep(delay4)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 4/4)\n", delay4.Seconds())

	// 	//	}
}

// // ก้อนเม้น
func see_watch_continue(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	//	for {
	fmt.Println("🛒 เข้า see_watch_continue")

	Runsee_watch4_VideoHomeSectionQuery(tlsConns, proxyAddr, proxyAuth)

	delay5 := time.Duration(rand.Intn(20)+10) * time.Second
	time.Sleep(delay5)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1/2)\n", delay5.Seconds())

	Runsee_watch3_FbShortsImpressionMutation(tlsConns, proxyAddr, proxyAuth)

	delay6 := time.Duration(rand.Intn(100)+10) * time.Second
	time.Sleep(delay6)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2/2)\n", delay6.Seconds())

	// 	//	}
}

// // /ก้อนเม้น/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func fetch_feed_start(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	// 	//	for {
	fmt.Println("🛒 เข้า fetch_feed 1 more")

	Runfetch_feed1(tlsConns, proxyAddr, proxyAuth)
	Runfetch_feed2_lightspeed(tlsConns, proxyAddr, proxyAuth)
	Runfetch_feed2_more(tlsConns, proxyAddr, proxyAuth)

	// 	// ✅ สุ่มรอ 10-60 วินาที
	delay := time.Duration(rand.Intn(51)+10) * time.Second
	time.Sleep(delay)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2/6)\n", delay.Seconds())

	//RunsendAnalyticsLog(tlsConns, proxyAddr, proxyAuth)
	Runfetch_feed3_scroll_past_group(tlsConns, proxyAddr, proxyAuth)

	delay2 := time.Duration(rand.Intn(31)+10) * time.Second
	time.Sleep(delay2)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 3/6)\n", delay2.Seconds())

	Runfetch_feed2_more(tlsConns, proxyAddr, proxyAuth)

	delay3 := time.Duration(rand.Intn(120)+10) * time.Second
	time.Sleep(delay3)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 4/6)\n", delay3.Seconds())

	Runfetch_feed3_scroll_past_group(tlsConns, proxyAddr, proxyAuth)

	delay4 := time.Duration(rand.Intn(50)+10) * time.Second
	time.Sleep(delay4)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 5/6)\n", delay4.Seconds())

	Runfetch_feed2_more(tlsConns, proxyAddr, proxyAuth)

	delay5 := time.Duration(rand.Intn(10)+10) * time.Second
	time.Sleep(delay5)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบสุดท้าย)\n", delay5.Seconds())

	// 	//	}
}

// // /ก้อนเม้น
func fetch_feed_continue(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	// 	//	for {
	fmt.Println("🛒 เข้า fetch_feed_continue")

	Runfetch_feed3_scroll_past_group(tlsConns, proxyAddr, proxyAuth)

	delay6 := time.Duration(rand.Intn(10)+3) * time.Second
	time.Sleep(delay6)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1/4)\n", delay6.Seconds())

	Runfetch_feed2_more(tlsConns, proxyAddr, proxyAuth)

	delay7 := time.Duration(rand.Intn(20)+5) * time.Second
	time.Sleep(delay7)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2/4)\n", delay7.Seconds())

	Runfetch_feed2_more(tlsConns, proxyAddr, proxyAuth)

	delay8 := time.Duration(rand.Intn(80)+10) * time.Second
	time.Sleep(delay8)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 3/4)\n", delay8.Seconds())

	Runfetch_feed3_scroll_past_group(tlsConns, proxyAddr, proxyAuth)
	delay9 := time.Duration(rand.Intn(20)+5) * time.Second
	time.Sleep(delay9)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 4/4)\n", delay9.Seconds())

	// 	//	}
}

// ก้อนเม้น
func like_only_only(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {

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

	var reactionType string
	err = db.QueryRow("SELECT reaction_type FROM like_only_table LIMIT 1").Scan(&reactionType)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ reaction_type ในตาราง like_only_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล reaction_type ล้มเหลว:", err)
		return
	}

	///////////////////////////////////////////////////////////////////
	// ===== 1) สุ่มลำดับ "พร้อมรีเจน id/rowid" =====
	// - ดึงสคีมาเพื่อหา PK INTEGER
	type colInfo struct {
		name, ctype string
		pk          int
	}
	var cols []colInfo
	rows, err := db.Query(`PRAGMA table_info(like_and_comment_table)`)
	if err != nil {
		fmt.Println("❌ PRAGMA table_info ล้มเหลว:", err)
		return
	}
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			fmt.Println("❌ อ่าน table_info ล้มเหลว:", err)
			rows.Close()
			return
		}
		cols = append(cols, colInfo{name: name, ctype: strings.ToUpper(ctype), pk: pk})
	}
	rows.Close()

	// ตรวจว่ามี INTEGER PRIMARY KEY ไหม
	hasIntPK := false
	var pkName string
	for _, c := range cols {
		if c.pk == 1 && strings.Contains(c.ctype, "INT") { // พอสำหรับตรวจ INTEGER
			hasIntPK = true
			pkName = c.name
			break
		}
	}

	// เริ่มทรานแซคชัน เพื่อสลับลำดับจริง ๆ
	if _, err := db.Exec(`PRAGMA foreign_keys=OFF`); err != nil {
		fmt.Println("❌ ปิด foreign_keys ไม่ได้:", err)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("❌ เปิดทรานแซคชันไม่ได้:", err)
		return
	}

	// temp table เรียงแบบสุ่ม
	if _, err := tx.Exec(`
		CREATE TEMP TABLE __tmp_like AS
		SELECT * FROM like_and_comment_table
		ORDER BY RANDOM()
	`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ สร้าง temp table ไม่สำเร็จ:", err)
		return
	}

	// ล้างของเดิม
	if _, err := tx.Exec(`DELETE FROM like_and_comment_table`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ ล้างตารางเดิมไม่สำเร็จ:", err)
		return
	}

	if hasIntPK {
		// เตรียมคอลัมน์ที่ไม่ใช่ PK เพื่อให้ SQLite สร้าง id ใหม่ตามลำดับสุ่ม
		var nonPK []string
		for _, c := range cols {
			if c.name != pkName {
				nonPK = append(nonPK, c.name)
			}
		}
		if len(nonPK) == 0 {
			_ = tx.Rollback()
			fmt.Println("❌ ไม่พบคอลัมน์สำหรับ insert (นอกจาก PK)")
			return
		}
		colList := strings.Join(nonPK, ",")
		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO like_and_comment_table (%s)
			SELECT %s FROM __tmp_like
		`, colList, colList))
		if err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับ (รีเจน PK) ไม่สำเร็จ:", err)
			return
		}
	} else {
		// ไม่มี INTEGER PK -> แค่ลบแล้วใส่กลับก็ได้ rowid ใหม่ตามลำดับสุ่ม
		if _, err := tx.Exec(`INSERT INTO like_and_comment_table SELECT * FROM __tmp_like`); err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับไม่สำเร็จ:", err)
			return
		}
	}

	_, _ = tx.Exec(`DROP TABLE __tmp_like`)
	if err := tx.Commit(); err != nil {
		fmt.Println("❌ คอมมิตไม่สำเร็จ:", err)
		return
	}
	_, _ = db.Exec(`PRAGMA foreign_keys=ON`)
	fmt.Println("✅ สุ่มลำดับแถวใน DB + รีเจน rowid/id แล้ว")

	// ===== 2) จากนี้ LIMIT 1 จะได้ "แถวแรกแบบใหม่" จริง ๆ =====

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	if strings.ToUpper(reactionType) == "LIKE" {

		Runlike_only_befor_reactiontype_like_DelightsMLEAnimationQuery(tlsConns, proxyAddr, proxyAuth)
		Runlike_only(tlsConns, proxyAddr, proxyAuth)

	} else {

		Runlike_only(tlsConns, proxyAddr, proxyAuth)

		delay6 := time.Duration(rand.Intn(11)+2) * time.Second
		time.Sleep(delay6)
		fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1/1)\n", delay6.Seconds())

	}
}

// /ก้อนเม้น/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func like_and_comment(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {

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

	var reactionType string
	err = db.QueryRow("SELECT reaction_type FROM like_and_comment_table LIMIT 1").Scan(&reactionType)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ reaction_type ในตาราง like_and_comment_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล reaction_type ล้มเหลว:", err)
		return

	}
	///////////////////////////////////////////////////////////////////
	// ===== 1) สุ่มลำดับ "พร้อมรีเจน id/rowid" =====
	// - ดึงสคีมาเพื่อหา PK INTEGER
	type colInfo struct {
		name, ctype string
		pk          int
	}
	var cols []colInfo
	rows, err := db.Query(`PRAGMA table_info(like_and_comment_table)`)
	if err != nil {
		fmt.Println("❌ PRAGMA table_info ล้มเหลว:", err)
		return
	}
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			fmt.Println("❌ อ่าน table_info ล้มเหลว:", err)
			rows.Close()
			return
		}
		cols = append(cols, colInfo{name: name, ctype: strings.ToUpper(ctype), pk: pk})
	}
	rows.Close()

	// ตรวจว่ามี INTEGER PRIMARY KEY ไหม
	hasIntPK := false
	var pkName string
	for _, c := range cols {
		if c.pk == 1 && strings.Contains(c.ctype, "INT") { // พอสำหรับตรวจ INTEGER
			hasIntPK = true
			pkName = c.name
			break
		}
	}

	// เริ่มทรานแซคชัน เพื่อสลับลำดับจริง ๆ
	if _, err := db.Exec(`PRAGMA foreign_keys=OFF`); err != nil {
		fmt.Println("❌ ปิด foreign_keys ไม่ได้:", err)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("❌ เปิดทรานแซคชันไม่ได้:", err)
		return
	}

	// temp table เรียงแบบสุ่ม
	if _, err := tx.Exec(`
		CREATE TEMP TABLE __tmp_like AS
		SELECT * FROM like_and_comment_table
		ORDER BY RANDOM()
	`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ สร้าง temp table ไม่สำเร็จ:", err)
		return
	}

	// ล้างของเดิม
	if _, err := tx.Exec(`DELETE FROM like_and_comment_table`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ ล้างตารางเดิมไม่สำเร็จ:", err)
		return
	}

	if hasIntPK {
		// เตรียมคอลัมน์ที่ไม่ใช่ PK เพื่อให้ SQLite สร้าง id ใหม่ตามลำดับสุ่ม
		var nonPK []string
		for _, c := range cols {
			if c.name != pkName {
				nonPK = append(nonPK, c.name)
			}
		}
		if len(nonPK) == 0 {
			_ = tx.Rollback()
			fmt.Println("❌ ไม่พบคอลัมน์สำหรับ insert (นอกจาก PK)")
			return
		}
		colList := strings.Join(nonPK, ",")
		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO like_and_comment_table (%s)
			SELECT %s FROM __tmp_like
		`, colList, colList))
		if err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับ (รีเจน PK) ไม่สำเร็จ:", err)
			return
		}
	} else {
		// ไม่มี INTEGER PK -> แค่ลบแล้วใส่กลับก็ได้ rowid ใหม่ตามลำดับสุ่ม
		if _, err := tx.Exec(`INSERT INTO like_and_comment_table SELECT * FROM __tmp_like`); err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับไม่สำเร็จ:", err)
			return
		}
	}

	_, _ = tx.Exec(`DROP TABLE __tmp_like`)
	if err := tx.Commit(); err != nil {
		fmt.Println("❌ คอมมิตไม่สำเร็จ:", err)
		return
	}
	_, _ = db.Exec(`PRAGMA foreign_keys=ON`)
	fmt.Println("✅ สุ่มลำดับแถวใน DB + รีเจน rowid/id แล้ว")

	// ===== 2) จากนี้ LIMIT 1 จะได้ "แถวแรกแบบใหม่" จริง ๆ =====

	fmt.Println("🛒 เข้า fetch_feed (75%)")
	//////////////////////////////////////////////////////////////////
	if strings.ToUpper(reactionType) == "LIKE" {

		Runbefor_reactiontype_like_DelightsMLEAnimationQuery(tlsConns, proxyAddr, proxyAuth)
		Runlike_before_comment(tlsConns, proxyAddr, proxyAuth)

		delay61 := time.Duration(rand.Intn(10)+2) * time.Second
		time.Sleep(delay61)
		fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2)\n", delay61.Seconds())

		RunpreC0mment1_CommentHidingTransparencyNUXTooltipTextQuery(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment2_FetchPredictiveTextSuggestions(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment3_FeedbackStartTypingCoreMutation(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment4_MentionsSuggestionQuery(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment5_FetchPredictiveTextSuggestions(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment6_FamilyNonUserMemberTagQuery(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment7_FetchMentionsBootstrapEntities(tlsConns, proxyAddr, proxyAuth)

		delay62 := time.Duration(rand.Intn(51)+10) * time.Second
		time.Sleep(delay62)
		fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2)\n", delay62.Seconds())

		RunpreC0mment8_FeedbackStopTypingCoreMutation(tlsConns, proxyAddr, proxyAuth)
		Runcomment(tlsConns, proxyAddr, proxyAuth)

	} else {

		Runlike_before_comment(tlsConns, proxyAddr, proxyAuth)

		delay63 := time.Duration(rand.Intn(10)+2) * time.Second
		time.Sleep(delay63)
		fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2)\n", delay63.Seconds())

		RunpreC0mment1_CommentHidingTransparencyNUXTooltipTextQuery(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment2_FetchPredictiveTextSuggestions(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment3_FeedbackStartTypingCoreMutation(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment4_MentionsSuggestionQuery(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment5_FetchPredictiveTextSuggestions(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment6_FamilyNonUserMemberTagQuery(tlsConns, proxyAddr, proxyAuth)
		RunpreC0mment7_FetchMentionsBootstrapEntities(tlsConns, proxyAddr, proxyAuth)

		delay64 := time.Duration(rand.Intn(51)+10) * time.Second
		time.Sleep(delay64)
		fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2)\n", delay64.Seconds())

		RunpreC0mment8_FeedbackStopTypingCoreMutation(tlsConns, proxyAddr, proxyAuth)
		Runcomment(tlsConns, proxyAddr, proxyAuth)

		delay6 := time.Duration(rand.Intn(11)+2) * time.Second
		time.Sleep(delay6)
		fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2)\n", delay6.Seconds())
	}
}

// // /ก้อนเม้น////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func like_reel_only_only(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var reactionType string
	err = db.QueryRow("SELECT reaction_type FROM like_reel_only_table LIMIT 1").Scan(&reactionType)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ reaction_type ในตาราง like_reel_only_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล reaction_type ล้มเหลว:", err)
		return
	}

	///////////////////////////////////////////////////////////////////
	// ===== 1) สุ่มลำดับ "พร้อมรีเจน id/rowid" =====
	// - ดึงสคีมาเพื่อหา PK INTEGER
	type colInfo struct {
		name, ctype string
		pk          int
	}
	var cols []colInfo
	rows, err := db.Query(`PRAGMA table_info(like_and_comment_table)`)
	if err != nil {
		fmt.Println("❌ PRAGMA table_info ล้มเหลว:", err)
		return
	}
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			fmt.Println("❌ อ่าน table_info ล้มเหลว:", err)
			rows.Close()
			return
		}
		cols = append(cols, colInfo{name: name, ctype: strings.ToUpper(ctype), pk: pk})
	}
	rows.Close()

	// ตรวจว่ามี INTEGER PRIMARY KEY ไหม
	hasIntPK := false
	var pkName string
	for _, c := range cols {
		if c.pk == 1 && strings.Contains(c.ctype, "INT") { // พอสำหรับตรวจ INTEGER
			hasIntPK = true
			pkName = c.name
			break
		}
	}

	// เริ่มทรานแซคชัน เพื่อสลับลำดับจริง ๆ
	if _, err := db.Exec(`PRAGMA foreign_keys=OFF`); err != nil {
		fmt.Println("❌ ปิด foreign_keys ไม่ได้:", err)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("❌ เปิดทรานแซคชันไม่ได้:", err)
		return
	}

	// temp table เรียงแบบสุ่ม
	if _, err := tx.Exec(`
		CREATE TEMP TABLE __tmp_like AS
		SELECT * FROM like_and_comment_table
		ORDER BY RANDOM()
	`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ สร้าง temp table ไม่สำเร็จ:", err)
		return
	}

	// ล้างของเดิม
	if _, err := tx.Exec(`DELETE FROM like_and_comment_table`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ ล้างตารางเดิมไม่สำเร็จ:", err)
		return
	}

	if hasIntPK {
		// เตรียมคอลัมน์ที่ไม่ใช่ PK เพื่อให้ SQLite สร้าง id ใหม่ตามลำดับสุ่ม
		var nonPK []string
		for _, c := range cols {
			if c.name != pkName {
				nonPK = append(nonPK, c.name)
			}
		}
		if len(nonPK) == 0 {
			_ = tx.Rollback()
			fmt.Println("❌ ไม่พบคอลัมน์สำหรับ insert (นอกจาก PK)")
			return
		}
		colList := strings.Join(nonPK, ",")
		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO like_and_comment_table (%s)
			SELECT %s FROM __tmp_like
		`, colList, colList))
		if err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับ (รีเจน PK) ไม่สำเร็จ:", err)
			return
		}
	} else {
		// ไม่มี INTEGER PK -> แค่ลบแล้วใส่กลับก็ได้ rowid ใหม่ตามลำดับสุ่ม
		if _, err := tx.Exec(`INSERT INTO like_and_comment_table SELECT * FROM __tmp_like`); err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับไม่สำเร็จ:", err)
			return
		}
	}

	_, _ = tx.Exec(`DROP TABLE __tmp_like`)
	if err := tx.Commit(); err != nil {
		fmt.Println("❌ คอมมิตไม่สำเร็จ:", err)
		return
	}
	_, _ = db.Exec(`PRAGMA foreign_keys=ON`)
	fmt.Println("✅ สุ่มลำดับแถวใน DB + รีเจน rowid/id แล้ว")

	// ===== 2) จากนี้ LIMIT 1 จะได้ "แถวแรกแบบใหม่" จริง ๆ =====

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runlike_reel_only(tlsConns, proxyAddr, proxyAuth)

	delay6 := time.Duration(rand.Intn(11)+2) * time.Second
	time.Sleep(delay6)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1)\n", delay6.Seconds())
}

// // ก้อนเม้น
func like_reel_and_comment_reel(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var reactionType string
	err = db.QueryRow("SELECT reaction_type FROM like_reel_and_comment_reel_table LIMIT 1").Scan(&reactionType)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ reaction_type ในตาราง like_reel_and_comment_reel_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล reaction_type ล้มเหลว:", err)
		return
	}

	///////////////////////////////////////////////////////////////////
	// ===== 1) สุ่มลำดับ "พร้อมรีเจน id/rowid" =====
	// - ดึงสคีมาเพื่อหา PK INTEGER
	type colInfo struct {
		name, ctype string
		pk          int
	}
	var cols []colInfo
	rows, err := db.Query(`PRAGMA table_info(like_and_comment_table)`)
	if err != nil {
		fmt.Println("❌ PRAGMA table_info ล้มเหลว:", err)
		return
	}
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			fmt.Println("❌ อ่าน table_info ล้มเหลว:", err)
			rows.Close()
			return
		}
		cols = append(cols, colInfo{name: name, ctype: strings.ToUpper(ctype), pk: pk})
	}
	rows.Close()

	// ตรวจว่ามี INTEGER PRIMARY KEY ไหม
	hasIntPK := false
	var pkName string
	for _, c := range cols {
		if c.pk == 1 && strings.Contains(c.ctype, "INT") { // พอสำหรับตรวจ INTEGER
			hasIntPK = true
			pkName = c.name
			break
		}
	}

	// เริ่มทรานแซคชัน เพื่อสลับลำดับจริง ๆ
	if _, err := db.Exec(`PRAGMA foreign_keys=OFF`); err != nil {
		fmt.Println("❌ ปิด foreign_keys ไม่ได้:", err)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("❌ เปิดทรานแซคชันไม่ได้:", err)
		return
	}

	// temp table เรียงแบบสุ่ม
	if _, err := tx.Exec(`
		CREATE TEMP TABLE __tmp_like AS
		SELECT * FROM like_and_comment_table
		ORDER BY RANDOM()
	`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ สร้าง temp table ไม่สำเร็จ:", err)
		return
	}

	// ล้างของเดิม
	if _, err := tx.Exec(`DELETE FROM like_and_comment_table`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ ล้างตารางเดิมไม่สำเร็จ:", err)
		return
	}

	if hasIntPK {
		// เตรียมคอลัมน์ที่ไม่ใช่ PK เพื่อให้ SQLite สร้าง id ใหม่ตามลำดับสุ่ม
		var nonPK []string
		for _, c := range cols {
			if c.name != pkName {
				nonPK = append(nonPK, c.name)
			}
		}
		if len(nonPK) == 0 {
			_ = tx.Rollback()
			fmt.Println("❌ ไม่พบคอลัมน์สำหรับ insert (นอกจาก PK)")
			return
		}
		colList := strings.Join(nonPK, ",")
		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO like_and_comment_table (%s)
			SELECT %s FROM __tmp_like
		`, colList, colList))
		if err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับ (รีเจน PK) ไม่สำเร็จ:", err)
			return
		}
	} else {
		// ไม่มี INTEGER PK -> แค่ลบแล้วใส่กลับก็ได้ rowid ใหม่ตามลำดับสุ่ม
		if _, err := tx.Exec(`INSERT INTO like_and_comment_table SELECT * FROM __tmp_like`); err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับไม่สำเร็จ:", err)
			return
		}
	}

	_, _ = tx.Exec(`DROP TABLE __tmp_like`)
	if err := tx.Commit(); err != nil {
		fmt.Println("❌ คอมมิตไม่สำเร็จ:", err)
		return
	}
	_, _ = db.Exec(`PRAGMA foreign_keys=ON`)
	fmt.Println("✅ สุ่มลำดับแถวใน DB + รีเจน rowid/id แล้ว")

	// ===== 2) จากนี้ LIMIT 1 จะได้ "แถวแรกแบบใหม่" จริง ๆ =====

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runlike_reel_before_comment_reel(tlsConns, proxyAddr, proxyAuth)

	delay66 := time.Duration(rand.Intn(10)+2) * time.Second
	time.Sleep(delay66)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1)\n", delay66.Seconds())

	RunpreC0mment1_for_reel(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment2_for_reel(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment3_for_reel(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment4_for_reel(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment5_for_reel(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment6_for_reel(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment7_for_reel(tlsConns, proxyAddr, proxyAuth)

	delay68 := time.Duration(rand.Intn(51)+10) * time.Second
	time.Sleep(delay68)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2)\n", delay68.Seconds())

	RunpreC0mment8_for_reel(tlsConns, proxyAddr, proxyAuth)
	Runcomment_reel(tlsConns, proxyAddr, proxyAuth)

	delay67 := time.Duration(rand.Intn(11)+2) * time.Second
	time.Sleep(delay67)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1)\n", delay67.Seconds())

}

// // /ก้อนเม้น////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func like_comment_only_only(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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
	var reactionType string
	err = db.QueryRow("SELECT reaction_type FROM like_comment_only_table LIMIT 1").Scan(&reactionType)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ reaction_type ในตาราง like_comment_only_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล reaction_type ล้มเหลว:", err)
		return
	}

	///////////////////////////////////////////////////////////////////
	// ===== 1) สุ่มลำดับ "พร้อมรีเจน id/rowid" =====
	// - ดึงสคีมาเพื่อหา PK INTEGER
	type colInfo struct {
		name, ctype string
		pk          int
	}
	var cols []colInfo
	rows, err := db.Query(`PRAGMA table_info(like_and_comment_table)`)
	if err != nil {
		fmt.Println("❌ PRAGMA table_info ล้มเหลว:", err)
		return
	}
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			fmt.Println("❌ อ่าน table_info ล้มเหลว:", err)
			rows.Close()
			return
		}
		cols = append(cols, colInfo{name: name, ctype: strings.ToUpper(ctype), pk: pk})
	}
	rows.Close()

	// ตรวจว่ามี INTEGER PRIMARY KEY ไหม
	hasIntPK := false
	var pkName string
	for _, c := range cols {
		if c.pk == 1 && strings.Contains(c.ctype, "INT") { // พอสำหรับตรวจ INTEGER
			hasIntPK = true
			pkName = c.name
			break
		}
	}

	// เริ่มทรานแซคชัน เพื่อสลับลำดับจริง ๆ
	if _, err := db.Exec(`PRAGMA foreign_keys=OFF`); err != nil {
		fmt.Println("❌ ปิด foreign_keys ไม่ได้:", err)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("❌ เปิดทรานแซคชันไม่ได้:", err)
		return
	}

	// temp table เรียงแบบสุ่ม
	if _, err := tx.Exec(`
		CREATE TEMP TABLE __tmp_like AS
		SELECT * FROM like_and_comment_table
		ORDER BY RANDOM()
	`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ สร้าง temp table ไม่สำเร็จ:", err)
		return
	}

	// ล้างของเดิม
	if _, err := tx.Exec(`DELETE FROM like_and_comment_table`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ ล้างตารางเดิมไม่สำเร็จ:", err)
		return
	}

	if hasIntPK {
		// เตรียมคอลัมน์ที่ไม่ใช่ PK เพื่อให้ SQLite สร้าง id ใหม่ตามลำดับสุ่ม
		var nonPK []string
		for _, c := range cols {
			if c.name != pkName {
				nonPK = append(nonPK, c.name)
			}
		}
		if len(nonPK) == 0 {
			_ = tx.Rollback()
			fmt.Println("❌ ไม่พบคอลัมน์สำหรับ insert (นอกจาก PK)")
			return
		}
		colList := strings.Join(nonPK, ",")
		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO like_and_comment_table (%s)
			SELECT %s FROM __tmp_like
		`, colList, colList))
		if err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับ (รีเจน PK) ไม่สำเร็จ:", err)
			return
		}
	} else {
		// ไม่มี INTEGER PK -> แค่ลบแล้วใส่กลับก็ได้ rowid ใหม่ตามลำดับสุ่ม
		if _, err := tx.Exec(`INSERT INTO like_and_comment_table SELECT * FROM __tmp_like`); err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับไม่สำเร็จ:", err)
			return
		}
	}

	_, _ = tx.Exec(`DROP TABLE __tmp_like`)
	if err := tx.Commit(); err != nil {
		fmt.Println("❌ คอมมิตไม่สำเร็จ:", err)
		return
	}
	_, _ = db.Exec(`PRAGMA foreign_keys=ON`)
	fmt.Println("✅ สุ่มลำดับแถวใน DB + รีเจน rowid/id แล้ว")

	// ===== 2) จากนี้ LIMIT 1 จะได้ "แถวแรกแบบใหม่" จริง ๆ =====

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runlike_comment_only(tlsConns, proxyAddr, proxyAuth)

	delay6 := time.Duration(rand.Intn(11)+2) * time.Second
	time.Sleep(delay6)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2)\n", delay6.Seconds())
}

// // ก้อนเม้น
func like_comment_and_reply_comment_table(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var reactionType string
	err = db.QueryRow("SELECT reaction_type FROM like_comment_and_reply_comment_table LIMIT 1").Scan(&reactionType)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ reaction_type ในตาราง like_comment_and_reply_comment_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล reaction_type ล้มเหลว:", err)
		return
	}

	///////////////////////////////////////////////////////////////////
	// ===== 1) สุ่มลำดับ "พร้อมรีเจน id/rowid" =====
	// - ดึงสคีมาเพื่อหา PK INTEGER
	type colInfo struct {
		name, ctype string
		pk          int
	}
	var cols []colInfo
	rows, err := db.Query(`PRAGMA table_info(like_and_comment_table)`)
	if err != nil {
		fmt.Println("❌ PRAGMA table_info ล้มเหลว:", err)
		return
	}
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			fmt.Println("❌ อ่าน table_info ล้มเหลว:", err)
			rows.Close()
			return
		}
		cols = append(cols, colInfo{name: name, ctype: strings.ToUpper(ctype), pk: pk})
	}
	rows.Close()

	// ตรวจว่ามี INTEGER PRIMARY KEY ไหม
	hasIntPK := false
	var pkName string
	for _, c := range cols {
		if c.pk == 1 && strings.Contains(c.ctype, "INT") { // พอสำหรับตรวจ INTEGER
			hasIntPK = true
			pkName = c.name
			break
		}
	}

	// เริ่มทรานแซคชัน เพื่อสลับลำดับจริง ๆ
	if _, err := db.Exec(`PRAGMA foreign_keys=OFF`); err != nil {
		fmt.Println("❌ ปิด foreign_keys ไม่ได้:", err)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("❌ เปิดทรานแซคชันไม่ได้:", err)
		return
	}

	// temp table เรียงแบบสุ่ม
	if _, err := tx.Exec(`
		CREATE TEMP TABLE __tmp_like AS
		SELECT * FROM like_and_comment_table
		ORDER BY RANDOM()
	`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ สร้าง temp table ไม่สำเร็จ:", err)
		return
	}

	// ล้างของเดิม
	if _, err := tx.Exec(`DELETE FROM like_and_comment_table`); err != nil {
		_ = tx.Rollback()
		fmt.Println("❌ ล้างตารางเดิมไม่สำเร็จ:", err)
		return
	}

	if hasIntPK {
		// เตรียมคอลัมน์ที่ไม่ใช่ PK เพื่อให้ SQLite สร้าง id ใหม่ตามลำดับสุ่ม
		var nonPK []string
		for _, c := range cols {
			if c.name != pkName {
				nonPK = append(nonPK, c.name)
			}
		}
		if len(nonPK) == 0 {
			_ = tx.Rollback()
			fmt.Println("❌ ไม่พบคอลัมน์สำหรับ insert (นอกจาก PK)")
			return
		}
		colList := strings.Join(nonPK, ",")
		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO like_and_comment_table (%s)
			SELECT %s FROM __tmp_like
		`, colList, colList))
		if err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับ (รีเจน PK) ไม่สำเร็จ:", err)
			return
		}
	} else {
		// ไม่มี INTEGER PK -> แค่ลบแล้วใส่กลับก็ได้ rowid ใหม่ตามลำดับสุ่ม
		if _, err := tx.Exec(`INSERT INTO like_and_comment_table SELECT * FROM __tmp_like`); err != nil {
			_ = tx.Rollback()
			fmt.Println("❌ ใส่กลับไม่สำเร็จ:", err)
			return
		}
	}

	_, _ = tx.Exec(`DROP TABLE __tmp_like`)
	if err := tx.Commit(); err != nil {
		fmt.Println("❌ คอมมิตไม่สำเร็จ:", err)
		return
	}
	_, _ = db.Exec(`PRAGMA foreign_keys=ON`)
	fmt.Println("✅ สุ่มลำดับแถวใน DB + รีเจน rowid/id แล้ว")

	// ===== 2) จากนี้ LIMIT 1 จะได้ "แถวแรกแบบใหม่" จริง ๆ =====

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runlike_comment_before_reply_comment(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment1_for_comment_comment(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment2_for_comment_comment(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment3_for_comment_comment(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment4_for_comment_comment(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment5_for_comment_comment(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment6_for_comment_comment(tlsConns, proxyAddr, proxyAuth)
	RunpreC0mment7_for_comment_comment(tlsConns, proxyAddr, proxyAuth)

	delay68 := time.Duration(rand.Intn(51)+10) * time.Second
	time.Sleep(delay68)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 2)\n", delay68.Seconds())

	RunpreC0mment8_for_comment_comment(tlsConns, proxyAddr, proxyAuth)
	Runcomment_comment(tlsConns, proxyAddr, proxyAuth)

	delay6 := time.Duration(rand.Intn(11)+2) * time.Second
	time.Sleep(delay6)
	fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำงานต่อ... (รอบที่ 1)\n", delay6.Seconds())

}

// // //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// // switch
func bio_bio(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var status string
	err = db.QueryRow("SELECT status_id FROM switch_for_bio_profile_table LIMIT 1").Scan(&status)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ status_id ในตาราง switch_for_bio_profile_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล switch_for_bio_profile_table ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runchange_profile1_details(tlsConns, proxyAddr, proxyAuth)
	Runbio(tlsConns, proxyAddr, proxyAuth)
	Runcity(tlsConns, proxyAddr, proxyAuth)
	Runchange_name1(tlsConns, proxyAddr, proxyAuth)
	Runchange_name2(tlsConns, proxyAddr, proxyAuth)
	Runschool1(tlsConns, proxyAddr, proxyAuth)
	Runschool2_real_change(tlsConns, proxyAddr, proxyAuth)

	_, err = db.Exec("DELETE FROM switch_for_bio_profile_table WHERE status_id = ?", status) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ switch_for_bio_profile_table ออกจากฐานข้อมูลแล้ว:", status)
	}

}

// // สุ่ม///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// // folder
func cover_pic(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}
	coverPath := filepath.Join(folder, "cover_photo")

	files, err := ioutil.ReadDir(coverPath)
	if err != nil {
		fmt.Println("❌ เปิดโฟลเดอร์ cover_photo ไม่ได้:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("📂 ไม่มีไฟล์ในโฟลเดอร์ cover_photo จบการทำงาน")
		return
	}

	fmt.Printf("📸 เจอไฟล์ %d รายการใน cover_photo เริ่มทำงานต่อ...\n", len(files))

	Runcover_pic1up(tlsConns, proxyAddr, proxyAuth)
	Runcover_pic2(tlsConns, proxyAddr, proxyAuth)
	Runcover_pic3(tlsConns, proxyAddr, proxyAuth)
	Runcover_pic4(tlsConns, proxyAddr, proxyAuth)

}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// // สุ่ม//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// // folder
func profile_pic(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {

	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}
	profilePhotoPath := filepath.Join(folder, "profile_photo")

	files, err := ioutil.ReadDir(profilePhotoPath)
	if err != nil {
		fmt.Println("❌ เปิดโฟลเดอร์ profile_photo ไม่ได้:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("📂 ไม่มีไฟล์ในโฟลเดอร์ profile_photo จบการทำงาน")
		return
	}

	fmt.Printf("📸 เจอไฟล์ %d รายการใน profile_photo เริ่มทำงานต่อ...\n", len(files))

	Runprofile_pic1(tlsConns, proxyAddr, proxyAuth)
	Runprofile_pic2(tlsConns, proxyAddr, proxyAuth)
	Runprofile_pic4(tlsConns, proxyAddr, proxyAuth)
	Runprofile_pic5(tlsConns, proxyAddr, proxyAuth)
	Runprofile_pic6_up(tlsConns, proxyAddr, proxyAuth)
	Runprofile_pic7_2_set(tlsConns, proxyAddr, proxyAuth)
	Runprofile_pic8(tlsConns, proxyAddr, proxyAuth)

}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// // /สุ่ม/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func unfollow_unfollow(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var unsubscribee_id string
	err = db.QueryRow("SELECT unsubscribee_id FROM unsubscribee_id_table LIMIT 1").Scan(&unsubscribee_id)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ unsubscribee_id ในตาราง unsubscribee_id_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล unsubscribee_id ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Rununfollow(tlsConns, proxyAddr, proxyAuth)

}

// // /สุ่ม/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func follow_follow(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var subscribee_id string
	err = db.QueryRow("SELECT subscribee_id FROM subscribee_id_table LIMIT 1").Scan(&subscribee_id)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ subscribee_id ในตาราง subscribee_id_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล subscribee_id ล้มเหลว:", err)
		return
	}
	// 	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runfollow(tlsConns, proxyAddr, proxyAuth)

}

// // /สุ่ม/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func joint_group_joint_group(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var group_id string
	err = db.QueryRow("SELECT group_id FROM group_id_table LIMIT 1").Scan(&group_id)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ group_id ในตาราง group_id_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล group_id ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runjoint_group(tlsConns, proxyAddr, proxyAuth)

}

// // //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// // switch
func lock_profile(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var status string
	err = db.QueryRow("SELECT status_id FROM switch_for_lock_profile_table LIMIT 1").Scan(&status)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ status_id ในตาราง switch_for_lock_profile_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล switch_for_lock_profile_table ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runlock_profile1(tlsConns, proxyAddr, proxyAuth)
	Runlock_profile2(tlsConns, proxyAddr, proxyAuth)
	RunJewel(tlsConns, proxyAddr, proxyAuth)
	Runlock_profile4_truelock(tlsConns, proxyAddr, proxyAuth)
	Runlock_profile5(tlsConns, proxyAddr, proxyAuth)
	Runlock_profile6(tlsConns, proxyAddr, proxyAuth)

	_, err = db.Exec("DELETE FROM switch_for_lock_profile_table WHERE status_id = ?", status) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ switch_for_lock_profile_table ออกจากฐานข้อมูลแล้ว:", status)
	}
}

// // //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// // switch
func unlock_profile(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var status string
	err = db.QueryRow("SELECT status_id FROM switch_for_unlock_profile_table LIMIT 1").Scan(&status)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ status_id ในตาราง switch_for_unlock_profile_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล switch_for_unlock_profile_table ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Rununlock_profile1(tlsConns, proxyAddr, proxyAuth)
	Rununlock_profile2(tlsConns, proxyAddr, proxyAuth)
	Rununlock_profile3(tlsConns, proxyAddr, proxyAuth)
	Rununlock_profile4(tlsConns, proxyAddr, proxyAuth)
	Rununlock_profile5(tlsConns, proxyAddr, proxyAuth)

	_, err = db.Exec("DELETE FROM switch_for_unlock_profile_table WHERE status_id = ?", status) // commentText, postLink
	if err != nil {
		fmt.Println("❌ ลบไม่สำเร็จ:", err)
	} else {
		fmt.Println("🧹 ลบ switch_for_unlock_profile_table ออกจากฐานข้อมูลแล้ว:", status)
	}
}

// // สุ่ม//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// // Folder
func story_upload(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	folder := strings.TrimSpace(os.Getenv("DBFOLDER"))
	if folder == "" {
		folder = "."
	}
	storyPhotoPath := filepath.Join(folder, "story_photo")

	files, err := ioutil.ReadDir(storyPhotoPath)
	if err != nil {
		fmt.Println("❌ เปิดโฟลเดอร์ story_photo ไม่ได้:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("📂 ไม่มีไฟล์ในโฟลเดอร์ story_photo จบการทำงาน")
		return
	}

	fmt.Printf("📸 เจอไฟล์ %d รายการใน story_photo เริ่มทำงานต่อ...\n", len(files))

	Runstory1_InspirationMusicPicker(tlsConns, proxyAddr, proxyAuth)
	Runstory2_StoriesPrivacySettingsQuery(tlsConns, proxyAddr, proxyAuth)
	Runstory3_upload_photo(tlsConns, proxyAddr, proxyAuth)
	Runstory4_set_photo(tlsConns, proxyAddr, proxyAuth)

}

// // ก้อนเม้น///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func shared_link_text_table(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {

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

	var status_text string
	err = db.QueryRow("SELECT status_text FROM shared_link_text_table LIMIT 1").Scan(&status_text)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ status_text ในตาราง shared_link_text_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล status_text ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runset_status_link(tlsConns, proxyAddr, proxyAuth)

}

// // ก้อนเม้น///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func shared_link_link(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var link_link string
	err = db.QueryRow("SELECT link_link FROM shared_link_table LIMIT 1").Scan(&link_link)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ link_link ในตาราง shared_link_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล link_link ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runshared_link(tlsConns, proxyAddr, proxyAuth)

}

// // ก้อนเม้น
func set_status_status(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var status_text string
	err = db.QueryRow("SELECT status_text FROM set_status_text_table LIMIT 1").Scan(&status_text)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ status_text ในตาราง set_status_text_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล status_text ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runset_status(tlsConns, proxyAddr, proxyAuth)

}

// // ก้อนเม้น
func up_pic_caption_caption(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var caption_text string
	err = db.QueryRow("SELECT caption_text FROM pic_caption_text_table LIMIT 1").Scan(&caption_text)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ caption_text ในตาราง pic_caption_text_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล caption_text ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runup_pic_caption(tlsConns, proxyAddr, proxyAuth)
	Runpic_caption(tlsConns, proxyAddr, proxyAuth)

}

// // ก้อนเม้น///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func like_page(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
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

	var link_page string
	err = db.QueryRow("SELECT link_page FROM link_page_for_like_table LIMIT 1").Scan(&link_page)
	if err == sql.ErrNoRows {
		fmt.Println("❌ ไม่พบ link_page ในตาราง link_page_for_like_table จบการทำงาน")
		return
	} else if err != nil {
		fmt.Println("❌ ดึงข้อมูล link_page ล้มเหลว:", err)
		return
	}

	fmt.Println("🛒 เข้า fetch_feed (75%)")

	Runlike_page1_FbBloksActionRootQuery(tlsConns, proxyAddr, proxyAuth)
	Runlike_page2_FbBloksActionRootQuery(tlsConns, proxyAddr, proxyAuth)
	Runlike_page3_PageLike(tlsConns, proxyAddr, proxyAuth)
	Runlike_page4_ProfilePlusLikeChainingNTViewQuery(tlsConns, proxyAddr, proxyAuth)

}

// //////////////// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func about_story(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string, duration time.Duration) {
	start := time.Now()
	fmt.Println("🚀 เข้าสู่โหมด about_story")

	for time.Since(start) < duration {
		remaining := duration - time.Since(start)
		if remaining < 5*time.Second {
			break
		}

		//	fmt.Println("👀 ดูสตอรี่")
		//	see_story(tlsConns, active, proxyAddr, proxyAuth)
		//	time.Sleep(time.Second * time.Duration(rand.Intn(3)+2)) // 2–4 วิ

		fmt.Println("📤 อัปโหลดสตอรี่")
		story_upload(tlsConns, active, proxyAddr, proxyAuth)
		time.Sleep(time.Second * time.Duration(rand.Intn(4)+3)) // 3–6 วิ
	}

	fmt.Println("🛑 ออกจาก about_story (หมดเวลา)")
}

// //เก็บไว้เผื่อใช้
// func about_story(active *bool, proxyAddr, proxyAuth string, duration time.Duration) {
// 	start := time.Now()
// 	fmt.Println("🚀 เข้าสู่โหมด about_story")

// 	for time.Since(start) < duration {
// 		remaining := duration - time.Since(start)
// 		if remaining < 5*time.Second {
// 			break // ถ้าเหลือน้อยกว่า 5 วิ ไม่ทำละ
// 		}

// 		action := rand.Intn(2) // 0 หรือ 1
// 		switch action {
// 		case 0:
// 			fmt.Println("👀 ดูสตอรี่")
// 			see_story(active, proxyAddr, proxyAuth)
// 		case 1:
// 			fmt.Println("📤 อัปโหลดสตอรี่")
// 			story_upload(active, proxyAddr, proxyAuth)
// 		}

// 		delay := time.Duration(rand.Intn(5)+2) * time.Second // 2–6 วิ
// 		fmt.Printf("⏱️ รอ %.0f วินาทีก่อนทำรอบใหม่\n", delay.Seconds())
// 		time.Sleep(delay)
// 	}

// 	fmt.Println("🛑 ออกจาก about_story (หมดเวลา)")
// }

func about_watch(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string, duration time.Duration) {
	fmt.Println("🚀 เริ่ม see_watch_start")
	//	see_watch_start(tlsConns, active, proxyAddr, proxyAuth)

	rand.Seed(time.Now().UnixNano()) // แค่ครั้งเดียว

	start := time.Now()
	for time.Since(start) < duration {
		x := rand.Intn(3) + 1

		switch x {
		case 1:
			fmt.Println("🚀 เริ่ม like_reel_only")
			like_reel_only_only(tlsConns, active, proxyAddr, proxyAuth)

		case 2:
			fmt.Println("🚀 เริ่ม see_watch_continue")
			//see_watch_continue(tlsConns, active, proxyAddr, proxyAuth)

		case 3:
			fmt.Println("🚀 เริ่ม like_reel_and_comment_reel")
			like_reel_and_comment_reel(tlsConns, active, proxyAddr, proxyAuth)
			// 💤 พักบ้างไม่ให้รันรัวแบบบอทโง่
			delay := time.Duration(rand.Intn(11)+5) * time.Second // 5–15 วิ
			fmt.Printf("⏱️ รอ %.0f วินาที\n", delay.Seconds())
			time.Sleep(delay)
		}

		fmt.Println("🛑 ออกจาก about_watch แล้ว (หมดเวลา)")
	}
}

func about_feed(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string, duration time.Duration) {
	//fmt.Println("🚀 เริ่ม fetch_feed_start")
	//fetch_feed_start(tlsConns, active, proxyAddr, proxyAuth)
	delay := time.Duration(rand.Intn(36)+5) * time.Second // 5–15 วิ
	fmt.Printf("⏱️ รอ %.0f วินาที\n", delay.Seconds())
	time.Sleep(delay)

	rand.Seed(time.Now().UnixNano()) // แค่ครั้งเดียว

	start := time.Now()
	for time.Since(start) < duration {
		//	x := rand.Intn(13) + 1 // ได้เลข 1 ถึง 13

		//	switch x {
		//	case 1:
		//	delay := time.Duration(rand.Intn(600)+10) * time.Second
		//	fmt.Printf("⏱️ รอ %.0f วินาที\n", delay.Seconds())
		//	time.Sleep(delay)
		//fmt.Println("🚀 เริ่ม fetch_feed_continue")
		//fetch_feed_continue(tlsConns, active, proxyAddr, proxyAuth)

		fmt.Println("🚀 เริ่ม like_and_comment")
		like_and_comment(tlsConns, active, proxyAddr, proxyAuth)

		//	case 2:
		fmt.Println("🚀 เริ่ม like_only")
		like_only_only(tlsConns, active, proxyAddr, proxyAuth)

		//	case 3:
		fmt.Println("🚀 เริ่ม like_and_comment")
		like_and_comment(tlsConns, active, proxyAddr, proxyAuth)

		//	case 4:
		fmt.Println("🚀 เริ่ม like_comment_only")
		like_comment_only_only(tlsConns, active, proxyAddr, proxyAuth)

		//	case 5:
		fmt.Println("🚀 เริ่ม like_comment_and_reply_comment_table")
		like_comment_and_reply_comment_table(tlsConns, active, proxyAddr, proxyAuth)

		//	case 6:
		fmt.Println("🚀 เริ่ม up_pic_caption")
		up_pic_caption_caption(tlsConns, active, proxyAddr, proxyAuth)

		//	case 7:
		fmt.Println("🚀 เริ่ม unfollow")
		unfollow_unfollow(tlsConns, active, proxyAddr, proxyAuth)

		//	case 8:
		fmt.Println("🚀 เริ่ม follow")
		follow_follow(tlsConns, active, proxyAddr, proxyAuth)

		//	case 9:
		fmt.Println("🚀 เริ่ม joint_group")
		joint_group_joint_group(tlsConns, active, proxyAddr, proxyAuth)

		//	case 10:
		fmt.Println("🚀 เริ่ม shared_link_text_table")
		shared_link_text_table(tlsConns, active, proxyAddr, proxyAuth)

		//	case 11:
		fmt.Println("🚀 เริ่ม shared_link")
		shared_link_link(tlsConns, active, proxyAddr, proxyAuth)

		//	case 12:
		fmt.Println("🚀 เริ่ม set_status")
		set_status_status(tlsConns, active, proxyAddr, proxyAuth)

		//	case 13:
		fmt.Println("🚀 เริ่ม like_page")
		like_page(tlsConns, active, proxyAddr, proxyAuth)

		delay := time.Duration(rand.Intn(600)+10) * time.Second
		fmt.Printf("⏱️ รอ %.0f วินาที\n", delay.Seconds())
		time.Sleep(delay)

	}

	// delay := time.Duration(rand.Intn(11)+5) * time.Second // 5–15 วิ
	// fmt.Printf("⏱️ รอ %.0f วินาที\n", delay.Seconds())
	// time.Sleep(delay)
	//	}

	delay8 := time.Duration(rand.Intn(11)+5) * time.Second // 5–15 วิ
	fmt.Printf("⏱️ รอ %.0f วินาที\n", delay8.Seconds())
	time.Sleep(delay)

	fmt.Println("🛑 ออกจาก about_feed แล้ว (หมดเวลา)")
}

func simulateAppBehavior(tlsConns *TLSConnections, active *bool, proxyAddr, proxyAuth string) {
	fmt.Println("🚀 เริ่ม simulateAppBehavior")
	totalSeconds := rand.Intn(301) + 300
	//totalSeconds := rand.Intn(6301) + 900 // 15 - 120 นาที totalSeconds := rand.Intn(6301) + 900 totalSeconds := rand.Intn(120) + 10
	fmt.Printf("📱 เปิดแอพรอบใหม่: %d วินาที\n", totalSeconds)

	start := time.Now()
	for time.Since(start) < time.Duration(totalSeconds)*time.Second {
		//r := rand.Intn(100) + 1 // สุ่ม 1–100 เพื่อคุมสัดส่วน
		//	r := rand.Intn(3) + 1 // สุ่ม 1–100 เพื่อคุมสัดส่วน
		var durationSec int

		//	switch r {
		// case r <= 9: // 9%
		//	case 3:
		fmt.Println("🛒 เริ่ม maket")
		// 	durationSec = int(float64(totalSeconds) * 0.1)
		// 	maket(tlsConns, active, proxyAddr, proxyAuth, time.Duration(durationSec)*time.Second)

		// case r <= 12: // +3% = 12%
		//	case 4:
		fmt.Println("📖 เริ่ม about_story")
		// 	durationSec = int(float64(totalSeconds) * 0.05)
		// 	about_story(tlsConns, active, proxyAddr, proxyAuth, time.Duration(durationSec)*time.Second)

		//default:
		//	case 2:
		fmt.Println("📰 เริ่ม about_feed")
		durationSec = int(float64(totalSeconds) * 0.6)
		about_feed(tlsConns, active, proxyAddr, proxyAuth, time.Duration(durationSec)*time.Second)
		//}

		//	case r <= 30: // +18% = 30%
		//	case 1: // +18% = 30%
		fmt.Println("🎥 เริ่ม about_watch")
		durationSec = int(float64(totalSeconds) * 0.25)
		about_watch(tlsConns, active, proxyAddr, proxyAuth, time.Duration(durationSec)*time.Second)

	}

	fmt.Println("🔕 จบ simulateAppBehavior แล้ว")
}

func waitForTLS() *TLSConnections {
	for {
		tlsConns, err := initTLSConns()
		if err != nil {
			fmt.Println("❌ TLS Handshake Fail:", err)
			time.Sleep(3 * time.Second) // รอ 3 วินาทีก่อนลองใหม่
			continue
		}
		fmt.Println("✅ TLS Handshake สำเร็จ")
		return tlsConns
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())
	proxyAddr, proxyAuth := parseProxy(fullProxy)
	active := true
	tlsConns := waitForTLS()

	defer closeConns(tlsConns) // ปิด connection ตอนจบ

	go loopJewel(tlsConns, proxyAddr, proxyAuth)

	for {
		active = true
		fmt.Println("▶️ [1] SendPing")
		go loopSendPing(tlsConns, &active, proxyAddr, proxyAuth)

		// ✅ TLS handshake เฉพาะตอนเปิดแอป
		tlsConns, err := initTLSConns()
		if err != nil {
			log.Fatal("❌ TLS handshake พัง:", err)
		}

		fmt.Println("▶️ [1] เปิดแอพ")
		OpenApp(tlsConns, &active, proxyAddr, proxyAuth)

		fmt.Println("▶️ [2] รับเพื่อน")
		friend_accept(tlsConns, &active, proxyAddr, proxyAuth)

		fmt.Println("▶️ [3] ขอเพื่อน")
		friend_requester(tlsConns, &active, proxyAddr, proxyAuth)

		fmt.Println("▶️ [6] ตั้ง bio")
		bio_bio(tlsConns, &active, proxyAddr, proxyAuth)

		fmt.Println("▶️ [4] ล็อกโปรไฟล์")
		lock_profile(tlsConns, &active, proxyAddr, proxyAuth)

		fmt.Println("▶️ [5] ปลดล็อกโปรไฟล์")
		unlock_profile(tlsConns, &active, proxyAddr, proxyAuth)

		fmt.Println("▶️ [6] ตั้ง bio")
		bio_bio(tlsConns, &active, proxyAddr, proxyAuth)

		fmt.Println("▶️ [7] ตั้ง cover")
		cover_pic(tlsConns, &active, proxyAddr, proxyAuth)

		fmt.Println("▶️ [8] ตั้ง profile pic")
		profile_pic(tlsConns, &active, proxyAddr, proxyAuth)

		fmt.Println("▶️ [9] เริ่ม simulateAppBehavior")
		simulateAppBehavior(tlsConns, &active, proxyAddr, proxyAuth)

		//closeConns(tlsConns)

		// 🔻 ปิดแอพ
		active = false
		sleepSeconds := rand.Intn(121) + 180 // สุ่ม 180–300 วินาที //sleepSeconds := rand.Intn(106*60) + (15 * 60) // สุ่ม 900–7200 วินาที (15–120 นาที) sleepSeconds := rand.Intn(106*60) + (15 * 60) sleepSeconds := rand.Intn(2*60) + (2 * 6)
		fmt.Printf("📴 ปิดแอพ พัก %d วินาที\n", sleepSeconds)
		time.Sleep(time.Duration(sleepSeconds) * time.Second)

		closeConns(tlsConns)

	}
}
