package i18n_message

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/yylego/must"
	"github.com/yylego/neatjson/neatjsons"
	"github.com/yylego/rese"
	"github.com/yylego/zaplog"
	"golang.org/x/text/language"
)

// DefaultLanguage 配置默认语言
var DefaultLanguage = language.English

//go:embed active.en.toml active.zh.toml
var files embed.FS

func LoadI18nFiles(debugModeOpen bool) (*i18n.Bundle, []*i18n.MessageFile) {
	bundle := i18n.NewBundle(DefaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	var messageFiles []*i18n.MessageFile
	for _, fileName := range []string{"active.en.toml", "active.zh.toml"} {
		content := rese.A1(files.ReadFile(fileName))
		//这里文件名 file-name 写 "active.en.toml" 或者 "en.toml" 都行，内部会通过这个解析出语言标签名称
		messageFile := rese.P1(bundle.ParseMessageFileBytes(content, fileName))
		if debugModeOpen {
			zaplog.SUG.Debugln(neatjsons.S(messageFile)) //安利下我的俩工具包
		}
		messageFiles = append(messageFiles, messageFile)
	}
	must.Have(messageFiles)
	must.Have(bundle.LanguageTags())

	if debugModeOpen {
		zaplog.SUG.Debugln(neatjsons.S(bundle.LanguageTags()))
	}
	return bundle, messageFiles
}
