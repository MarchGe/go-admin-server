package email

type HtmlMailer struct {
	Mailer
}

func CreateHtmlMailer(c *MailConfig) *HtmlMailer {
	tm := &HtmlMailer{
		Mailer{
			cfg: c,
		},
	}
	return tm
}

func (hm *HtmlMailer) MsgBuilder() *HtmlMsgBuilder {
	return &HtmlMsgBuilder{
		from: hm.cfg.From,
	}
}
