package email

import (
	"github.com/wneessen/go-mail"
	"html/template"
)

type HtmlMsgBuilder struct {
	from         string
	to           []string
	cc           []string
	bcc          []string
	subject      string
	htmlTemplate string
	data         any
}

func (hmb *HtmlMsgBuilder) To(to []string) *HtmlMsgBuilder {
	hmb.to = to
	return hmb
}
func (hmb *HtmlMsgBuilder) Cc(cc []string) *HtmlMsgBuilder {
	hmb.cc = cc
	return hmb
}
func (hmb *HtmlMsgBuilder) Bcc(bcc []string) *HtmlMsgBuilder {
	hmb.bcc = bcc
	return hmb
}
func (hmb *HtmlMsgBuilder) Subject(subject string) *HtmlMsgBuilder {
	hmb.subject = subject
	return hmb
}

func (hmb *HtmlMsgBuilder) HtmlTemplate(template string) *HtmlMsgBuilder {
	hmb.htmlTemplate = template
	return hmb
}

func (hmb *HtmlMsgBuilder) Data(data any) *HtmlMsgBuilder {
	hmb.data = data
	return hmb
}

func (hmb *HtmlMsgBuilder) Build() (*mail.Msg, error) {
	msg := mail.NewMsg()
	if err := msg.From(hmb.from); err != nil {
		return nil, err
	}
	if err := msg.To(hmb.to...); err != nil {
		return nil, err
	}
	if err := msg.Cc(hmb.cc...); err != nil {
		return nil, err
	}
	if err := msg.Bcc(hmb.bcc...); err != nil {
		return nil, err
	}
	msg.Subject(hmb.subject)
	tmpl, err := template.New("mailTemplate").Parse(hmb.htmlTemplate)
	if err != nil {
		return nil, err
	}
	if err = msg.SetBodyHTMLTemplate(tmpl, hmb.data); err != nil {
		return nil, err
	}
	return msg, nil
}
