package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/weeee9/godtone/config"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/linebot"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var (
	client *secretmanager.Client
	cfg    config.Config
	bot    *linebot.Client
	err    error
)

var userUseCina = map[string]time.Time{}

func init() {
	cfg, err = config.Environ()
	if err != nil {
		log.Fatal(err)
	}
	if !cfg.Server.Debug {
		log.Println(" [godtone] Load Cred From CSM")
		if err := loadCredFromGSM(&cfg); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(cfg.LineCred)

}

func init() {
	bot, err = linebot.New(
		cfg.LineCred.Secret,
		cfg.LineCred.Token,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := gin.Default()

	router.Static("/m4a", "./m4a")
	router.Static("/mp4", "./mp4")
	router.Static("/img", "./img")

	router.POST("/callback", callback)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "pong",
		})
	})

	router.Run(":" + cfg.Server.Port)
}

func callback(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {

			// 10 分鐘內不給使用 bot
			if t, ok := userUseCina[event.Source.UserID]; ok {
				if time.Since(t) < time.Minute*10 {
					return
				}
			}

			groupID := event.Source.GroupID

			switch message := event.Message.(type) {
			case *linebot.StickerMessage:
				audioMsg := &linebot.AudioMessage{}
				switch message.StickerID {
				case "277504782":
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "imGoingToSleep"), 1000)
				case "277504783":
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "comeHereBitch"), 3000)
				case "277504795":
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "youHaveToSaidItFirst"), 1000)
				case "335919878":
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "sodaDrinkingFeelSoGood"), 2000)
				case "277504788":
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "doYouHaveAJob"), 2000)
				case "335919905":
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "carry"), 3000)
				}
				if _, err := bot.PushMessage(groupID, audioMsg).Do(); err != nil {
					log.Printf(" [linebot] error: %v\n", err.Error())
					return
				}
			case *linebot.TextMessage:

				if png := checkCinaLanguage(message.Text); png != "" {
					userUseCina[event.Source.UserID] = time.Now()
					if _, err := bot.PushMessage(groupID, linebot.NewImageMessage(getImgURL(cfg, png), getImgURL(cfg, png))).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}

				audioMsg := &linebot.AudioMessage{}
				if strings.Contains(message.Text, "那我也要睡") {
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "imGoingToSleep"), 1000)
				}
				if message.Text == "你要先講" {
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "youHaveToSaidItFirst"), 1000)
				}
				if message.Text == "爽阿刺阿" || message.Text == "爽啊刺啊" {
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "sodaDrinkingFeelSoGood"), 2000)
				}
				if strings.Contains(message.Text, "777") ||
					strings.Contains(message.Text, "聊天室") {
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "carry"), 2000)
				}
				if strings.Contains(message.Text, "過來一下") {
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "comeHereBitch"), 3000)
				}
				if strings.Contains(message.Text, "穿過我的") {
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "penetrateMy88"), 3000)
				}
				if strings.Contains(message.Text, "去死一死") || message.Text == "7414" {
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "goDie"), 1000)
				}
				if strings.Contains(message.Text, "太神拉") || message.Text == "carry" {
					audioMsg = linebot.NewAudioMessage(getM4AURL(cfg, "carry"), 3000)
				}
				if _, err := bot.PushMessage(groupID, audioMsg).Do(); err != nil {
					log.Printf(" [linebot] error: %v\n", err.Error())
					return
				}
			}
		}
	}
}

var cinaLanguage = []string{
	"視頻", "質量", "黑客", "老鐵", "牛逼", "激活", "高清",
}

var pngs = []string{
	"agent", "bailiff", "canine", "captain", "informant", "judge", "police", "v",
}

func checkCinaLanguage(text string) string {
	rand.Seed(time.Now().Unix())
	for _, word := range cinaLanguage {
		if strings.Contains(text, word) {
			return pngs[rand.Intn(len(pngs))]
		}
	}
	return ""
}

func getM4AURL(c config.Config, filename string) string {
	return fmt.Sprintf("%s/m4a/%s.m4a", c.Server.Addr, filename)
}

func getMP4URL(c config.Config, filename string) string {
	return fmt.Sprintf("%s/mp4/%s.mp4", c.Server.Addr, filename)
}

func getImgURL(c config.Config, filename string) string {
	return fmt.Sprintf("%s/img/%s.png", c.Server.Addr, filename)
}

func loadCredFromGSM(cfg *config.Config) error {
	client, err = secretmanager.NewClient(context.TODO())
	if err != nil {
		return err
	}
	secretResp, err := client.AccessSecretVersion(context.TODO(), &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/209245468577/secrets/LINE_CHANNEL_SECRET/versions/latest",
	})
	if err != nil {
		return err
	}
	tokenResp, err := client.AccessSecretVersion(context.TODO(), &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/209245468577/secrets/LINE_CHANNEL_TOKEN/versions/latest",
	})
	if err != nil {
		return err
	}
	cfg.LineCred.Secret = string(secretResp.Payload.Data)
	cfg.LineCred.Token = string(tokenResp.Payload.Data)
	return nil
}
