package i18n_message_test

import (
	"testing"

	"github.com/yylego/goi18n"
	"github.com/yylego/kratos-examples/demo1kratos/internal/pkg/middleware/localize/i18n_message"
	"github.com/yylego/neatjson/neatjsons"
	"github.com/yylego/osexistpath/osmustexist"
	"github.com/yylego/runpath/runtestpath"
	"github.com/yylego/zaplog"
)

//go:generate go test -v -run ^TestGenerate$
func TestGenerate(t *testing.T) {
	bundle, messageFiles := i18n_message.LoadI18nFiles(true)
	zaplog.SUG.Debugln(neatjsons.S(bundle.LanguageTags()))

	outputPath := osmustexist.FILE(runtestpath.SrcPath(t))
	options := goi18n.NewOptions().WithOutputPathWithPkgName(outputPath)
	t.Log(neatjsons.S(options))
	goi18n.Generate(messageFiles, options)
}
