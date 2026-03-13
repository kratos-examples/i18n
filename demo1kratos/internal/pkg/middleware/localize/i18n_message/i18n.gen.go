package i18n_message

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ErrorMsgParam struct {
	ErrorMessage any
}

func (p *ErrorMsgParam) GetTemplateValues() map[string]any {
	res := make(map[string]any)
	if p.ErrorMessage != nil {
		res["ErrorMessage"] = p.ErrorMessage
	}
	return res
}

func I18nErrorMsg(data *ErrorMsgParam) *i18n.LocalizeConfig {
	const messageID = "errorMsg"
	var valuesMap = data.GetTemplateValues()
	return &i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: valuesMap,
	}
}

type GreetingParam struct {
	Name any
}

func (p *GreetingParam) GetTemplateValues() map[string]any {
	res := make(map[string]any)
	if p.Name != nil {
		res["Name"] = p.Name
	}
	return res
}

func I18nGreeting(data *GreetingParam) *i18n.LocalizeConfig {
	const messageID = "greeting"
	var valuesMap = data.GetTemplateValues()
	return &i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: valuesMap,
	}
}

type WelcomeParam struct {
	User any
}

func (p *WelcomeParam) GetTemplateValues() map[string]any {
	res := make(map[string]any)
	if p.User != nil {
		res["User"] = p.User
	}
	return res
}

func I18nWelcome(data *WelcomeParam) *i18n.LocalizeConfig {
	const messageID = "welcome"
	var valuesMap = data.GetTemplateValues()
	return &i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: valuesMap,
	}
}
