package cnblogs

import (
	"fmt"
	"math/rand"
	"study-spider/spider"
	"study-spider/storage"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BlogCell struct {
	Href   string `json:"href"`
	Title  string `json:"title"`
	Source string `json:"source"`
}

type CNBlogs struct {
	sp spider.ISpider
}

type TimerFunc func()

func NewBlog() *CNBlogs {
	return &CNBlogs{
		sp: spider.NewBlogSpider(),
	}
}

func (b *CNBlogs) Run() {
	go b.exec(10*time.Minute, 20*time.Minute, b.spider)
	fmt.Println("blog spider running")
	select {}
}

func (b *CNBlogs) spider() {
	if err := b.sp.Spider(`https://www.cnblogs.com/`, nil, b.filter); err != nil {
		fmt.Println(err)
	}
}

func (b *CNBlogs) filter(doc *goquery.Document) error {
	cells := []storage.StorageCell{}
	doc.Find(".titlelnk").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists && href != "" {
			if title, err := s.Html(); err == nil && title != "" {
				cells = append(cells, storage.StorageCell{
					Key: "blog-" + title,
					Value: BlogCell{
						Href:   href,
						Title:  title,
						Source: `博客园`,
					},
				})
			}
		}
	})
	b.storage(cells)
	return nil
}

func (b *CNBlogs) storage(cells []storage.StorageCell) {
	storage.Storage("blog", cells)
}

func (b *CNBlogs) exec(minTime time.Duration, maxTime time.Duration, callback TimerFunc) {
	for {
		callback()
		v := time.Duration(rand.Int63n(int64(maxTime - minTime)))
		t := time.NewTicker(minTime + v)
		<-t.C
	}
}
