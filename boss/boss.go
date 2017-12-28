package boss

import (
	"fmt"
	"math/rand"
	"study-spider/spider"
	"study-spider/storage"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BossCell struct {
	Href        string `json:"href"`
	Img         string `json:"img"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Salary      string `json:"salary"`
	City        string `json:"city"`
	Time        string `json:"time"`
	Educational string `json:"educational"`
}

type TimerFunc func()

type Boss struct {
	sp spider.ISpider
}

func NewBoss() *Boss {
	return &Boss{
		sp: spider.NewBossSpider(),
	}
}

func (b *Boss) Run() {
	go b.exec(10*time.Minute, 20*time.Minute, b.spider)
	fmt.Println("boss spider running")
	select {}
}

func (b *Boss) spider() {
	// 抓取5页 golang 数据
	for i := 0; i < 5; i++ {
		if err := b.sp.Spider(fmt.Sprintf(`http://www.zhipin.com/mobile/jobs.json?page=%d&city=101270100&query=golang`, i+1), nil, b.filter); err != nil {
			fmt.Println(err)
			break
		}
	}
}

func (b *Boss) filter(doc *goquery.Document) error {
	cells := []storage.StorageCell{}
	doc.Find("li a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists && href != "" {
			title := ""
			img := ""
			company := ""
			salary := ""
			city := ""
			t := ""
			edu := ""
			if ss := s.Find("img"); ss != nil {
				img, _ = ss.Attr("img")
			}
			if ss := s.Find("h4"); ss != nil {
				title, _ = ss.Html()
			}
			if ss := s.Find(".salary"); ss != nil {
				salary, _ = ss.Html()
			}
			if ss := s.Find(".name"); ss != nil {
				company, _ = ss.Html()
			}
			s.Find("em").Each(func(i int, ss *goquery.Selection) {
				if i == 0 {
					city, _ = ss.Html()
				} else if i == 1 {
					t, _ = ss.Html()
				} else if i == 2 {
					edu, _ = ss.Html()
				}
			})
			cells = append(cells, storage.StorageCell{
				Key: "boss-" + title,
				Value: BossCell{
					Title:       title,
					Href:        "http://www.zhipin.com" + href,
					Img:         img,
					Company:     company,
					Salary:      salary,
					City:        city,
					Time:        t,
					Educational: edu,
				},
			})
		}
	})
	storage.Storage("boss", cells)
	return nil

}
func (b *Boss) exec(minTime time.Duration, maxTime time.Duration, callback TimerFunc) {
	for {
		callback()
		v := time.Duration(rand.Int63n(int64(maxTime - minTime)))
		t := time.NewTicker(minTime + v)
		<-t.C
	}
}
