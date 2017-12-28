package spider

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
)

type BossSpider struct {
}

func NewBossSpider() ISpider {
	return &BossSpider{}
}

type JobResponse struct {
	HasMore bool   `json:"hasMore"`
	ResMsg  string `json:"resmsg"`
	ResCode int    `json:"rescode"`
	Html    string `json:"html"`
}

func (bs *BossSpider) Spider(url string, header http.Header, next SpiderFunc) error {
	query := gorequest.New().Get(url)
	for k, v := range header {
		if len(v) > 0 {
			query = query.Set(k, v[0])
		}
	}
	job := JobResponse{}
	if _, _, errs := query.EndStruct(&job); len(errs) > 0 {
		return errs[0]
	}
	if job.ResCode == 1 {
		// 获取数据成功
		buf := bytes.NewBuffer([]byte(job.Html))
		if doc, err := goquery.NewDocumentFromReader(buf); err != nil {
			return err
		} else if err := next(doc); err != nil {
			return err
		}
	} else {
		return errors.New(job.ResMsg)
	}
	if !job.HasMore {
		// 没有更多数据了
		return errors.New("hasn't more data")
	}
	return nil
}
