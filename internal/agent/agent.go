package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"newsplatform/internal/database"
	"newsplatform/internal/models"
)

type AgentOpts struct {
	Endpoint     string
	Token        string
	Q            string
	TimeInterval time.Duration
	PageSie      int
	MaxPage      int
}

type Agent struct {
	Endpoint     string
	Token        string
	Q            string
	TimeInterval time.Duration
	PageSize     int
	MaxPage      int
}

func NewAgent(opts AgentOpts) Agent {
	return Agent{
		Endpoint:     opts.Endpoint,
		Token:        opts.Token,
		Q:            opts.Q,
		TimeInterval: opts.TimeInterval,
		PageSize:     opts.PageSie,
		MaxPage:      opts.MaxPage,
	}
}

func (a *Agent) Run(ctx context.Context) {
	ticker := time.NewTicker(a.TimeInterval)
	for {
		select {
		case <-ticker.C:
			log.Debug("Running agent...")
			for p := 1; p <= a.MaxPage; p++ {
				url := fmt.Sprintf("https://%s/?q=%s&apiKey=%s&pagesize=%d&page=%d", a.Endpoint, a.Q, a.Token, a.PageSize, p)
				articles, err := request(url)
				if err != nil {
					panic(err)
				}
				// TODO: add sync.WaitGroup to parse full content
				// and save to database concurrently
				saveToDB(articles)
			}
		case <-ctx.Done():
			log.Info("agent is exiting")
			break
		}
	}
}

func request(url string) ([]models.News, error) {
	log.Debugf("pulling content from: %s", url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("API call error", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var resp models.Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("can not parse response body into news")
	}
	return resp.Articles, nil
}

func getFullContent(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Errorf("http Request error: %v, %s", err, url)
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("get full-content error: %v", err)
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Errorf("read response body error: %v", err)
		return "", err
	}
	var buffer bytes.Buffer
	tags := "h1, h2, h3, p, title, table, tr, td, meta"
	doc.Find(tags).Each(func(i int, s *goquery.Selection) {
		tagName := s.Get(0).Data
		if tagName == "meta" {
			a, exists := s.Attr("content")
			if exists {
				buffer.WriteString(a + "\n\n")
			}
		} else {
			buffer.WriteString(s.Text() + "\n\n")
		}
	})
	log.Debug("Get full content:", buffer.String())
	return buffer.String(), nil
}

func saveToDB(news []models.News) {
	for _, item := range news {
		var n models.News
		if err := database.DB.Where("title = ?", item.Title).First(&n).Error; err != nil {
			log.Debugf("this article is not in db, save it,: %s", item.Title)
			fullContent, err := getFullContent(item.URL)
			if err != nil {
				log.Errorf("get full content error for article: %s, %v", n.Title, err)
			}
			item.FullContent = fullContent
			database.DB.Create(&item)
		}
		log.Warnf("this article is already in database, skip saving: %s", item.Title)
	}
}
