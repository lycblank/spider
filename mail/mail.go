package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"
	"study-spider/boss"
	"study-spider/cnblogs"
	"study-spider/config"
	"study-spider/storage"
	"time"
)

type EmailPackage struct {
	User     string
	Password string
	Host     string
	Body     string
	Type     string
	To       string
	Subject  string
}

func init() {
	go func() {
		t := time.NewTicker(30 * time.Second)
		for {
			<-t.C
			sendBlogToMail()
			t = time.NewTicker(24 * time.Hour)
		}
	}()
	go func() {
		t := time.NewTicker(30 * time.Second)
		for {
			<-t.C
			sendBossToMail()
			t = time.NewTicker(24 * time.Hour)
		}
	}()
}

func sendBlogToMail() {
	cells := []cnblogs.BlogCell{}
	storage.GetStorageContent(time.Now().Add(-24*time.Hour), "blog", func(raw []byte) bool {
		cell := cnblogs.BlogCell{}
		if err := json.Unmarshal(raw, &cell); err != nil {
			return false
		}
		exists := storage.Check("blog", cell.Title)
		if !exists {
			storage.Put("blog", cell.Title, "1")
			cells = append(cells, cell)
		}
		return exists
	})
	var buf bytes.Buffer
	if len(cells) > 0 {
		// 组装html body
		for _, cell := range cells {
			buf.WriteString(fmt.Sprintf(`<div><h3><a href="%s">%s</a></h3></div>`, cell.Href, cell.Title))
		}

		emailContent := EmailPackage{
			User:     config.Conf.Mail.User,
			Password: config.Conf.Mail.Password,
			Host:     config.Conf.Mail.Host,
			Body:     buf.String(),
			Type:     "html",
			To:       config.Conf.Mail.To,
			Subject:  "每日一学",
		}

		SendToMail(emailContent)
	}
}

func sendBossToMail() {
	fmt.Println("exec boss to mail")
	cells := []boss.BossCell{}
	//storage.GetStorageContent(time.Now().Add(-24*time.Hour), "boss", func(raw []byte) bool {
	storage.GetStorageContent(time.Now().Add(-24*time.Hour), "boss", func(raw []byte) bool {
		cell := boss.BossCell{}
		if err := json.Unmarshal(raw, &cell); err != nil {
			return false
		}
		exists := storage.Check("boss", cell.Company+"-"+cell.Title)
		if !exists {
			storage.Put("boss", cell.Company+"-"+cell.Title, "1")
			cells = append(cells, cell)
		}
		return exists
	})
	var buf bytes.Buffer
	if len(cells) > 0 {
		// 组装html body
		buf.WriteString(`<ul>`)
		for _, cell := range cells {
			buf.WriteString(`<li>`)
			buf.WriteString(fmt.Sprintf(`<a href="%s"><div><div>%s&nbsp&nbsp&nbsp&nbsp&nbsp&nbsp%s</div><div>%s</div><div><em>%s&nbsp&nbsp</em><em>%s&nbsp&nbsp</em><em>%s&nbsp&nbsp</em></div></div></a>`, cell.Href, cell.Title, cell.Salary, cell.Company, cell.City, cell.Time, cell.Educational))
			buf.WriteString(`</li>`)
		}
		buf.WriteString(`</ul>`)

		emailContent := EmailPackage{
			User:     config.Conf.Mail.User,
			Password: config.Conf.Mail.Password,
			Host:     config.Conf.Mail.Host,
			Body:     buf.String(),
			Type:     "html",
			To:       config.Conf.Mail.To,
			Subject:  "golang聘",
		}
		SendToMail(emailContent)
	}
}

func SendToMail(emailContent EmailPackage) error {
	hp := strings.Split(emailContent.Host, ":")
	auth := smtp.PlainAuth("", emailContent.User, emailContent.Password, hp[0])
	contentType := "Content-Type: text/plain; charset=UTF-8"
	if emailContent.Type == "html" {
		contentType = "Content-Type: text/html; charset=UTF-8"
	}
	// 组装邮件内容
	var buf bytes.Buffer
	buf.WriteString("To: ")
	buf.WriteString(emailContent.To)
	buf.WriteString("\r\nFrom: ")
	buf.WriteString(emailContent.User)
	buf.WriteString(">\r\nSubject: ")
	buf.WriteString(emailContent.Subject)
	buf.WriteString("\r\n")
	buf.WriteString(contentType)
	buf.WriteString("\r\n\r\n")
	buf.WriteString(emailContent.Body)
	return smtp.SendMail(emailContent.Host, auth, emailContent.User, strings.Split(emailContent.To, ";"), buf.Bytes())
}
