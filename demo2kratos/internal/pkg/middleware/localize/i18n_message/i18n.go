package i18n_message

import (
	"embed"
	"io/fs"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/yylego/must"
	"github.com/yylego/neatjson/neatjsons"
	"github.com/yylego/rese"
	"github.com/yylego/zaplog"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// DefaultLanguage 配置默认语言
var DefaultLanguage = language.AmericanEnglish

//go:embed active.en-US.yaml active.zh-CN.yaml
var files embed.FS

func LoadI18nFiles(debugModeOpen bool) (*i18n.Bundle, []*i18n.MessageFile) {
	bundle := i18n.NewBundle(DefaultLanguage)
	const format = "yaml"
	bundle.RegisterUnmarshalFunc(format, yaml.Unmarshal)

	var messageFiles []*i18n.MessageFile
	must.Done(fs.WalkDir(files, ".", func(fileName string, stat fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if stat.IsDir() || filepath.Ext(stat.Name()) != "."+format {
			return nil
		}
		content := rese.A1(files.ReadFile(fileName))
		//这里文件名 file-name 写 "active.en-US.toml" 或者 "en-US.toml" 都行，内部会通过这个解析出语言标签名称
		messageFile := rese.P1(bundle.ParseMessageFileBytes(content, fileName))
		if debugModeOpen {
			zaplog.SUG.Debugln(neatjsons.S(messageFile)) //安利下我的俩工具包
		}
		messageFiles = append(messageFiles, messageFile)
		return nil
	}))
	must.Have(messageFiles)
	must.Have(bundle.LanguageTags())

	if debugModeOpen {
		zaplog.SUG.Debugln(neatjsons.S(bundle.LanguageTags()))
	}
	return bundle, messageFiles
}
