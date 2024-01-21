package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"slot-crawler/internal/database"
	"strings"
	"time"
)

type Crawler struct {
	index      int32
	counter    int32
	uuid       string
	gameSymbol string
	mgcKey     string
	header     http.Header
	client     *http.Client
	baseUrl    string
}

func Initialize(gameSymbol string, uuid string) (*Crawler, error) {
	httpClient := http.Client{}
	baseRequestHeader := http.Header{
		"Host":            {"demogamesfree.pragmaticplay.net"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/116.0"},
		"Accept":          {"*/*"},
		"Accept-Language": {"en-US"},
		"Content-type":    {"application/x-www-form-urlencoded"},
		"Origin":          {"https://demogamesfree.pragmaticplay.net"},
		"Connection":      {"keep-alive"},
	}

	mgcKey, err := getMgcKey(gameSymbol)
	if err != nil {
		return nil, err
	}

	//uuid := uuid2.New().String()

	crawler := Crawler{
		index:      1,
		counter:    1,
		uuid:       uuid,
		mgcKey:     mgcKey,
		gameSymbol: gameSymbol,
		header:     baseRequestHeader,
		client:     &httpClient,
		baseUrl:    "https://demogamesfree.pragmaticplay.net/gs2c/v3/gameService",
	}

	return &crawler, nil
}

func getMgcKey(gameSymbol string) (string, error) {
	url := "https://demogamesfree-asia.pragmaticplay.net/gs2c/openGame.do?" +
		"gameSymbol=" + gameSymbol +
		"&websiteUrl=https://demogamesfree.pragmaticplay.net" +
		"&jurisdiction=99" +
		"&lobby_url=https://www.pragmaticplay.com/" +
		"&lang=EN" +
		"&cur=USD"

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	resQuery := res.Request.URL.Query()
	mgckey := resQuery.Get("mgckey")
	if mgckey == "" {
		return "", fmt.Errorf("mgcKey not exist")
	}
	return mgckey, nil
}

func (c *Crawler) initSlotClient() {

	c.index = 1
	c.counter = 1
	initBodyString := fmt.Sprintf("action=doInit&symbol=%s&cver=140343&repeat=0&mgckey=%s&index=%d&counter=%d", c.gameSymbol, c.mgcKey, c.index, c.counter)

	req, err := http.NewRequest(http.MethodPost, c.baseUrl, strings.NewReader(initBodyString))
	if err != nil {
		panic(err)
	}
	for key, values := range c.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	//httpClient := http.Client{}
	resp, err := c.client.Do(req)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		log.Panicf("initialize slot error")
	}
	log.Printf("init slot finished")
}

func (c *Crawler) spin() (map[string]interface{}, error) {

	c.index += 1
	c.counter += 2
	spinBodyString := fmt.Sprintf("action=doSpin&symbol=%s&c=0.05&l=4096&repeat=0&mgckey=%s&index=%d&counter=%d", c.gameSymbol, c.mgcKey, c.index, c.counter)
	if c.gameSymbol == "vswayschilheat" {
		spinBodyString += "&bl=0"
	}

	req, err := http.NewRequest(http.MethodPost, c.baseUrl, strings.NewReader(spinBodyString))
	if err != nil {
		panic(err)
	}
	for key, values := range c.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Panicf("spin Failed : %s", err.Error())
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Error reading response body: %s", err.Error())
	}

	responseBodyString := string(responseBody)
	fmt.Println(responseBodyString)
	if responseBodyString == "undefined (action invalid)" {
		return nil, fmt.Errorf("requestError : %s", responseBodyString)
	}

	json := UnMarshaling(responseBodyString, c.uuid)

	_, frozen := json["frozen"]
	if frozen {
		//panic("request error")
		return nil, fmt.Errorf("requestError : %s", responseBodyString)
	}

	database.InsertSpinData(c.gameSymbol, json)

	return json, nil
}

func (c *Crawler) collect() string {
	c.index += 1
	c.counter += 2
	bodyString := fmt.Sprintf("action=doCollect&symbol=%s&repeat=0&mgckey=%s&index=%d&counter=%d", c.gameSymbol, c.mgcKey, c.index, c.counter)
	req, err := http.NewRequest(http.MethodPost, c.baseUrl, strings.NewReader(bodyString))
	if err != nil {
		panic(err)
	}
	for key, values := range c.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Panicf("spin Failed : %s", err.Error())
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Error reading response body: %s", err.Error())
	}

	responseBodyString := string(responseBody)

	log.Println(responseBodyString)

	return responseBodyString
}

func (c *Crawler) StartCrawling(repeat int, wait time.Duration, cntChan chan int) error {
	c.initSlotClient()
	for i := 0; i < repeat; i++ {
		time.Sleep(wait)
		spinResult, err := c.spin()
		if err != nil {
			return err
		}

		w, _ := spinResult["tw"]
		_, fsmul := spinResult["fsmul"]
		_, total := spinResult["fswin_total"]
		_, rs_p := spinResult["rs_p"]

		if rs_p {
			continue
		}

		cntChan <- i
		if (w != "0.00" && !fsmul) || total {
			c.collect()
		}
	}

	return nil
}
