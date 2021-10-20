package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin:     func(r *http.Request) bool { return true },
}

// type IndexData struct {
// 	Title   string
// 	Content string
// }

// func test(c *gin.Context) {
// 	data := new(IndexData)
// 	data.Title = "WebSocket Test"
// 	data.Content = "Gin Test"
// 	c.HTML(http.StatusOK, "index.html", data)
// }

func main() {
	fmt.Println("Go WebSocket")
	rou()
}

func rou() {
	server := gin.New()
	// server := gin.Default()
	//指定靜態資料夾路徑
	// server.LoadHTMLGlob("./public/*")  //搭配 c.HTML 指定文件位置使用
	server.GET("/ws", wsEndpoint)

	server.NoRoute(gin.WrapH(http.FileServer(http.Dir("./public"))), func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		fmt.Println(path)
		fmt.Println(method)
		//檢查path的開頭使是否為"/"
		if strings.HasPrefix(path, "/") {
			fmt.Println("ok")
		}
	})
	err := server.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func wsEndpoint(c *gin.Context) {
	// 透過http請求程序調用upgrader.Upgrade，來獲取*Conn (代表WebSocket連接)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("使用者已連線")

	var newImg []byte

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Printf("連線WebSocket失敗！ %s", err.Error())
		}

		if string(p) == "run" {

			func() {
				//設定視訊鏡頭，0 = 預設鏡頭
				webcam, err := gocv.VideoCaptureDevice(0)
				if err != nil {
					log.Println(err)
				}

				time.Sleep(time.Second)
				img := gocv.NewMat()
				defer img.Close()

				webcam.Read(&img)
				defer webcam.Close()
				//設定副檔名、來源
				buf, err := gocv.IMEncode(".jpg", img)
				if err != nil {
					log.Fatal(err)
				}
				defer buf.Close() //nolint
				//設定變數取得暫存檔案
				newImg = buf.GetBytes()
				// d, _ := os.ReadFile(a)
				//轉換成base64的字串型別
				data := base64.StdEncoding.EncodeToString(newImg)
				//把轉換好的字串傳送到前端，前端接收在轉換回圖片
				if err := ws.WriteMessage(messageType, []byte(data)); err != nil {
					log.Println(err)
					return
				}
			}()

		}
		if string(p) == "save" {
			//用來生成新文件使用(檔名、來源、)
			err := os.WriteFile("demo.jpg", newImg, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
		log.Println("使用者訊息: " + string(p))
	}
}
