package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/linebot"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type lineCred struct {
	Key   string `datastore:"LINE_CHANNEL_SECRET"`
	Token string `datastore:"LINE_CHANNEL_TOKEN"`
}

const (
	appURL string = "https://godtone.df.r.appspot.com/"

	imGoingTosleep         string = appURL + "m4a/imGoingTosleep.m4a"
	youHaveToSaidItFirst   string = appURL + "m4a/youHaveToSaidItFirst.m4a"
	sodaDrinkingFeelSoGood string = appURL + "m4a/sodaDrinkingFeelSoGood.m4a"
	sevenInTheChat         string = appURL + "m4a/sevenInTheChat.m4a"
	comeHereBitch          string = appURL + "m4a/comeHereBitch.m4a"
	penetrateMy88          string = appURL + "m4a/penetrateMy88.m4a"
	goDie                  string = appURL + "m4a/godie.m4a"
	carry                  string = appURL + "m4a/carry.m4a"
	doYouHaveAJob          string = appURL + "m4a/doYouHaveAJob.m4a"
)

var (
	client *secretmanager.Client
	bot    *linebot.Client
	err    error

	secret, token string
)

func main() {
	server := os.Getenv("APP_SERVER_ENV")
	if server == "app_engine" {
		client, err = secretmanager.NewClient(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		secretResp, err := client.AccessSecretVersion(context.TODO(), &secretmanagerpb.AccessSecretVersionRequest{
			Name: "projects/209245468577/secrets/LINE_CHANNEL_SECRET/versions/latest",
		})
		if err != nil {
			log.Fatal(err)
		}
		tokenResp, err := client.AccessSecretVersion(context.TODO(), &secretmanagerpb.AccessSecretVersionRequest{
			Name: "projects/209245468577/secrets/LINE_CHANNEL_TOKEN/versions/latest",
		})
		if err != nil {
			log.Fatal(err)
		}
		secret = string(secretResp.Payload.Data)
		token = string(tokenResp.Payload.Data)
	} else {
		secret = os.Getenv("LINE_CHANNEL_SECRET")
		token = os.Getenv("LINE_CHANNEL_TOKEN")
	}

	bot, err = linebot.New(
		secret,
		token,
	)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "80"
	}

	router.Static("/m4a", "./m4a")

	router.POST("/callback", callback)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "pong",
		})
	})

	router.Run(":" + port)
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
			groupID := event.Source.GroupID

			switch message := event.Message.(type) {

			case *linebot.StickerMessage:
				fmt.Println(message.StickerID)
				switch message.StickerID {
				case "277504782":
					log.Println(imGoingTosleep)
					aud := linebot.NewAudioMessage(imGoingTosleep, 1000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				case "277504783":
					log.Print(comeHereBitch)
					aud := linebot.NewAudioMessage(comeHereBitch, 3000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				case "277504795":
					log.Println(youHaveToSaidItFirst)
					aud := linebot.NewAudioMessage(youHaveToSaidItFirst, 1000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				case "335919878":
					log.Println(sodaDrinkingFeelSoGood)
					aud := linebot.NewAudioMessage(sodaDrinkingFeelSoGood, 2000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				case "277504788":
					log.Println(doYouHaveAJob)
					aud := linebot.NewAudioMessage(doYouHaveAJob, 2000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				case "335919905":
					log.Println(carry)
					aud := linebot.NewAudioMessage(carry, 3000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}

			case *linebot.TextMessage:
				if strings.Contains(message.Text, "那我也要睡") {
					log.Println(imGoingTosleep)
					aud := linebot.NewAudioMessage(imGoingTosleep, 1000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}
				if message.Text == "你要先講" {
					log.Println(youHaveToSaidItFirst)
					aud := linebot.NewAudioMessage(youHaveToSaidItFirst, 1000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}
				if message.Text == "爽阿刺阿" || message.Text == "爽啊刺啊" {
					log.Println(sodaDrinkingFeelSoGood)
					aud := linebot.NewAudioMessage(sodaDrinkingFeelSoGood, 2000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}
				if strings.Contains(message.Text, "777") ||
					strings.Contains(message.Text, "聊天室") {
					log.Println(sevenInTheChat)
					aud := linebot.NewAudioMessage(sevenInTheChat, 2000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}
				if strings.Contains(message.Text, "過來一下") {
					log.Print(comeHereBitch)
					aud := linebot.NewAudioMessage(comeHereBitch, 3000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}
				if strings.Contains(message.Text, "穿過我的") {
					log.Println(penetrateMy88)
					aud := linebot.NewAudioMessage(penetrateMy88, 3000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}
				if strings.Contains(message.Text, "去死一死") {
					log.Println(goDie)
					aud := linebot.NewAudioMessage(goDie, 1000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}
				if strings.Contains(message.Text, "太神拉") || message.Text == carry {
					log.Println(carry)
					aud := linebot.NewAudioMessage(carry, 3000)
					if _, err := bot.PushMessage(groupID, aud).Do(); err != nil {
						log.Printf(" [linebot] error: %v\n", err.Error())
					}
					return
				}
			}
		}
	}
}
