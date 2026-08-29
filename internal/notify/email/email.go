package email

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailNotifier struct {
	SMTPHost    string
	SMTPPort    int
	From        string
	Password    string
	To          []string
	Title       string
	Description string
	Date        string
	Years       int
}

func (e *EmailNotifier) Send() error {
	if len(e.To) == 0 {
		return errors.New("email: 收件人列表不能为空")
	}
	if e.SMTPHost == "" {
		return errors.New("email: SMTP 服务器地址不能为空")
	}
	if e.From == "" {
		return errors.New("email: 发件人邮箱不能为空")
	}

	auth := smtp.PlainAuth("", e.From, e.Password, e.SMTPHost)

	subject := fmt.Sprintf("纪念日提醒: %s (%d周年)", e.Title, e.Years)
	body := fmt.Sprintf("标题: %s (%d周年)\n日期: %s\n描述: %s", e.Title, e.Years, e.Date, e.Description)

	// 将切片拼接为符合邮件标准的多收件人头部格式
	toHeader := strings.Join(e.To, ", ")
	msg := []byte(
		"To: " + toHeader + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
			body,
	)

	addr := fmt.Sprintf("%s:%d", e.SMTPHost, e.SMTPPort)
	if err := smtp.SendMail(addr, auth, e.From, e.To, msg); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}
