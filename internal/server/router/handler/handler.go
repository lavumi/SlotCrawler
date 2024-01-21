package handler

import (
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"io"
	"net/http"
	"slot-crawler/internal/crawler"
	"slot-crawler/internal/database"
	"strconv"
	"time"
)

func GetProgress() gin.HandlerFunc {
	return func(c *gin.Context) {
		//database.GetSpinDataByUser()
	}
}

func StartCrawling(c *gin.Context) {
	userQuery := c.Query("user")
	countQuery := c.Query("count")
	slotQuery := c.Query("slot")

	//fmt.Printf("%s, %s", userQuery, countQuery)

	user, _ := strconv.ParseInt(userQuery, 10, 32)
	count, _ := strconv.ParseInt(countQuery, 10, 32)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	uuid := uuid2.New().String()
	c.SSEvent("uuid", uuid)
	slot, err := crawler.Initialize(slotQuery, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "initialize slot fail [" + slotQuery + "] " + err.Error()})
		return
	}

	finished := make(chan bool)
	counter := make(chan int)
	quit := make(chan bool)

	defer func() {
		close(counter)
		close(finished)
		close(quit)
	}()

	go func() {
		for {
			select {
			case <-quit:
				finished <- true
				return
			default:
				for i := 0; i < int(user); i++ {
					err = slot.StartCrawling(int(count), time.Millisecond*300, counter)
					if err != nil {
						c.SSEvent("event", err.Error())
						return
					}
				}
				finished <- true
				return
			}
		}

	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			quit <- true
			return false
		case cnt := <-counter:
			c.SSEvent("count", cnt)
			return true
		case <-finished:
			c.SSEvent("close", "close")
			return false
		}
	})
}

type UserRequest struct {
	Slot string `uri:"slot"`
}

func GetUserList(c *gin.Context) {

	var req UserRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	//vs4096bufking
	users := database.GetUserList(req.Slot)
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

type DataRequest struct {
	Slot string `json:"slot"`
	User string `json:"user"`
}

func GetSpinData(c *gin.Context) {
	userQuery := c.Query("user")
	slotQuery := c.Query("slot")

	spinResult, _ := database.GetSpinDataByUser(userQuery, slotQuery)
	c.JSON(http.StatusOK, spinResult)
}
