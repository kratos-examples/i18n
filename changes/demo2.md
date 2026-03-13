# Changes

Code differences compared to source project.

## internal/pkg/middleware/localize/i18n_message/active.en-US.yaml (+3 -0)

```diff
@@ -0,0 +1,3 @@
+greeting: "Hello, {{.Name}}"
+welcome: "Welcome to our service, {{.User}}!"
+errorMsg: "Operation failed: {{.ErrorMessage}}"
```

## internal/pkg/middleware/localize/i18n_message/active.zh-CN.yaml (+3 -0)

```diff
@@ -0,0 +1,3 @@
+greeting: "你好，{{.Name}}"
+welcome: "欢迎使用我们的服务，{{.User}}！"
+errorMsg: "操作失败：{{.ErrorMessage}}"
```

## internal/pkg/middleware/localize/i18n_message/i18n.gen.go (+68 -0)

```diff
@@ -0,0 +1,68 @@
+package i18n_message
+
+import (
+	"github.com/nicksnyder/go-i18n/v2/i18n"
+)
+
+type ErrorMsgParam struct {
+	ErrorMessage any
+}
+
+func (p *ErrorMsgParam) GetTemplateValues() map[string]any {
+	res := make(map[string]any)
+	if p.ErrorMessage != nil {
+		res["ErrorMessage"] = p.ErrorMessage
+	}
+	return res
+}
+
+func I18nErrorMsg(data *ErrorMsgParam) *i18n.LocalizeConfig {
+	const messageID = "errorMsg"
+	var valuesMap = data.GetTemplateValues()
+	return &i18n.LocalizeConfig{
+		MessageID:    messageID,
+		TemplateData: valuesMap,
+	}
+}
+
+type GreetingParam struct {
+	Name any
+}
+
+func (p *GreetingParam) GetTemplateValues() map[string]any {
+	res := make(map[string]any)
+	if p.Name != nil {
+		res["Name"] = p.Name
+	}
+	return res
+}
+
+func I18nGreeting(data *GreetingParam) *i18n.LocalizeConfig {
+	const messageID = "greeting"
+	var valuesMap = data.GetTemplateValues()
+	return &i18n.LocalizeConfig{
+		MessageID:    messageID,
+		TemplateData: valuesMap,
+	}
+}
+
+type WelcomeParam struct {
+	User any
+}
+
+func (p *WelcomeParam) GetTemplateValues() map[string]any {
+	res := make(map[string]any)
+	if p.User != nil {
+		res["User"] = p.User
+	}
+	return res
+}
+
+func I18nWelcome(data *WelcomeParam) *i18n.LocalizeConfig {
+	const messageID = "welcome"
+	var valuesMap = data.GetTemplateValues()
+	return &i18n.LocalizeConfig{
+		MessageID:    messageID,
+		TemplateData: valuesMap,
+	}
+}
```

## internal/pkg/middleware/localize/i18n_message/i18n.gen_test.go (+23 -0)

```diff
@@ -0,0 +1,23 @@
+package i18n_message_test
+
+import (
+	"testing"
+
+	"github.com/yylego/goi18n"
+	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize/i18n_message"
+	"github.com/yylego/neatjson/neatjsons"
+	"github.com/yylego/osexistpath/osmustexist"
+	"github.com/yylego/runpath/runtestpath"
+	"github.com/yylego/zaplog"
+)
+
+//go:generate go test -v -run ^TestGenerate$
+func TestGenerate(t *testing.T) {
+	bundle, messageFiles := i18n_message.LoadI18nFiles(true)
+	zaplog.SUG.Debugln(neatjsons.S(bundle.LanguageTags()))
+
+	outputPath := osmustexist.FILE(runtestpath.SrcPath(t))
+	options := goi18n.NewOptions().WithOutputPathWithPkgName(outputPath)
+	t.Log(neatjsons.S(options))
+	goi18n.Generate(messageFiles, options)
+}
```

## internal/pkg/middleware/localize/i18n_message/i18n.go (+52 -0)

```diff
@@ -0,0 +1,52 @@
+package i18n_message
+
+import (
+	"embed"
+	"io/fs"
+	"path/filepath"
+
+	"github.com/nicksnyder/go-i18n/v2/i18n"
+	"github.com/yylego/must"
+	"github.com/yylego/neatjson/neatjsons"
+	"github.com/yylego/rese"
+	"github.com/yylego/zaplog"
+	"golang.org/x/text/language"
+	"gopkg.in/yaml.v3"
+)
+
+// DefaultLanguage 配置默认语言
+var DefaultLanguage = language.AmericanEnglish
+
+//go:embed active.en-US.yaml active.zh-CN.yaml
+var files embed.FS
+
+func LoadI18nFiles(debugModeOpen bool) (*i18n.Bundle, []*i18n.MessageFile) {
+	bundle := i18n.NewBundle(DefaultLanguage)
+	const format = "yaml"
+	bundle.RegisterUnmarshalFunc(format, yaml.Unmarshal)
+
+	var messageFiles []*i18n.MessageFile
+	must.Done(fs.WalkDir(files, ".", func(fileName string, stat fs.DirEntry, err error) error {
+		if err != nil {
+			return err
+		}
+		if stat.IsDir() || filepath.Ext(stat.Name()) != "."+format {
+			return nil
+		}
+		content := rese.A1(files.ReadFile(fileName))
+		//这里文件名 file-name 写 "active.en-US.toml" 或者 "en-US.toml" 都行，内部会通过这个解析出语言标签名称
+		messageFile := rese.P1(bundle.ParseMessageFileBytes(content, fileName))
+		if debugModeOpen {
+			zaplog.SUG.Debugln(neatjsons.S(messageFile)) //安利下我的俩工具包
+		}
+		messageFiles = append(messageFiles, messageFile)
+		return nil
+	}))
+	must.Have(messageFiles)
+	must.Have(bundle.LanguageTags())
+
+	if debugModeOpen {
+		zaplog.SUG.Debugln(neatjsons.S(bundle.LanguageTags()))
+	}
+	return bundle, messageFiles
+}
```

## internal/pkg/middleware/localize/i18n_message/i18n_test.go (+13 -0)

```diff
@@ -0,0 +1,13 @@
+package i18n_message_test
+
+import (
+	"testing"
+
+	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize/i18n_message"
+)
+
+func TestLoadI18nFiles(t *testing.T) {
+	bundle, messageFiles := i18n_message.LoadI18nFiles(true)
+	t.Log(len(messageFiles))
+	t.Log(len(bundle.LanguageTags()))
+}
```

## internal/pkg/middleware/localize/localize.go (+74 -0)

```diff
@@ -0,0 +1,74 @@
+package localize
+
+import (
+	"context"
+
+	"github.com/go-kratos/kratos/v2/middleware"
+	"github.com/go-kratos/kratos/v2/transport"
+	"github.com/nicksnyder/go-i18n/v2/i18n"
+	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize/i18n_message"
+	"github.com/yylego/tern/zerotern"
+	"github.com/yylego/zaplog"
+	"go.uber.org/zap"
+)
+
+type localizerKey struct{}
+
+// I18N 组件用于翻译
+// cp from: https://github.com/go-kratos/examples/blob/3a46aa32f7dbecbb01f2e3ecb28af187b2d9b53c/i18n/internal/pkg/middleware/localize/localize.go#L16
+func I18N() middleware.Middleware {
+	bundle, _ := i18n_message.LoadI18nFiles(false)
+
+	return func(handler middleware.Handler) middleware.Handler {
+		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
+			var acceptLanguage string
+			// parse accept-language from transport context
+			if tsp, ok := transport.FromServerContext(ctx); ok {
+				acceptLanguage = tsp.RequestHeader().Get("accept-language") //这里有可能得到个空字符串
+			}
+			// when accept == "" set accept = defaultValue
+			acceptLanguage = zerotern.VV(acceptLanguage, i18n_message.DefaultLanguage.String()) //在这里设置默认语言
+
+			localizer := i18n.NewLocalizer(bundle, acceptLanguage)
+			ctx = context.WithValue(ctx, localizerKey{}, localizer)
+			return handler(ctx, req)
+		}
+	}
+}
+
+func FromContext(ctx context.Context) *LocalizerI18n {
+	return NewLocalizerI18n(ctx.Value(localizerKey{}).(*i18n.Localizer))
+}
+
+// LocalizerI18n 把原始的翻译逻辑封装起来，增加些装饰逻辑，让翻译"永不报错"，因为更多的时候服务优先是逻辑正确，语言翻译不应该影响业务
+type LocalizerI18n struct {
+	localizer *i18n.Localizer
+}
+
+func NewLocalizerI18n(localizer *i18n.Localizer) *LocalizerI18n {
+	return &LocalizerI18n{localizer: localizer}
+}
+
+// Localize 封装 localizer 的 Localize 函数，但"永不报错"而不抛出异常，单返回值只返回消息而不返回错误
+func (T *LocalizerI18n) Localize(config *i18n.LocalizeConfig) string {
+	if T.localizer == nil {
+		zaplog.LOG.Warn("错误的翻译器-没有完成初始化逻辑")
+		return "[unknown]" + config.MessageID
+	}
+	if config.MessageID == "" {
+		zaplog.LOG.Warn("翻译参数错误-没有填写消息的ID")
+		return "[unknown]"
+	}
+	if config.DefaultMessage == nil {
+		config.DefaultMessage = &i18n.Message{ // 这里会改参数的配置，但认为没事，因为外部用完配置就会丢掉的，而且默认值不会造成负面影响
+			ID:    config.MessageID,
+			Other: "[message]" + config.MessageID, // 配置默认消息，当某种语言缺翻译时会返回这种格式的
+		}
+	}
+	msg, err := T.localizer.Localize(config)
+	if err != nil {
+		zaplog.LOG.Warn("翻译遇到问题-理论上这里不该出错", zap.Error(err))
+		return "[unknown]" + config.MessageID
+	}
+	return msg //这就是翻译成功的信息
+}
```

## internal/pkg/middleware/localize/localize_test.go (+37 -0)

```diff
@@ -0,0 +1,37 @@
+package localize_test
+
+import (
+	"context"
+	"testing"
+
+	"github.com/go-kratos/kratos/v2/middleware"
+	"github.com/stretchr/testify/require"
+	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize"
+	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize/i18n_message"
+	"github.com/yylego/neatjson/neatjsons"
+	"google.golang.org/protobuf/types/known/emptypb"
+	"google.golang.org/protobuf/types/known/wrapperspb"
+)
+
+// 其实也没必要对 middleware 进行单元测试的，这里只是自己尝试写写，掌握点新技能
+func TestI18N(t *testing.T) {
+	var rawFunc middleware.Handler = func(ctx context.Context, req any) (any, error) {
+		msg := localize.FromContext(ctx).Localize(i18n_message.I18nWelcome(&i18n_message.WelcomeParam{
+			User: "abc",
+		}))
+		return wrapperspb.String(msg), nil
+	}
+
+	runFunc := localize.I18N()(rawFunc)
+
+	// Unable to create a context with transport metadata (e.g., accept-language header) without simulating a web request.
+	ctx := context.Background()
+	// This unit test only tests the default English language due to the absence of transport metadata.
+	res, err := runFunc(ctx, &emptypb.Empty{})
+	require.NoError(t, err)
+	t.Log(neatjsons.S(res))
+
+	stringValue, ok := res.(*wrapperspb.StringValue)
+	require.True(t, ok)
+	require.Equal(t, "Welcome to our service, abc!", stringValue.Value)
+}
```

## internal/server/http.go (+2 -0)

```diff
@@ -6,6 +6,7 @@
 	"github.com/go-kratos/kratos/v2/transport/http"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
+	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
 )
 
@@ -13,6 +14,7 @@
 	var opts = []http.ServerOption{
 		http.Middleware(
 			recovery.Recovery(),
+			localize.I18N(),
 		),
 	}
 	if c.Http.Network != "" {
```

## internal/service/article.go (+8 -1)

```diff
@@ -5,6 +5,8 @@
 
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/biz"
+	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize"
+	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize/i18n_message"
 )
 
 type ArticleService struct {
@@ -22,7 +24,12 @@
 	if ebz != nil {
 		return nil, ebz.Erk
 	}
-	return &pb.CreateArticleReply{Article: &pb.ArticleInfo{Id: v.ID, Title: v.Title, Content: v.Content, StudentId: v.StudentID}}, nil
+
+	message := localize.FromContext(ctx).Localize(i18n_message.I18nGreeting(&i18n_message.GreetingParam{
+		Name: v.Title,
+	}))
+
+	return &pb.CreateArticleReply{Article: &pb.ArticleInfo{Id: v.ID, Title: message, Content: v.Content, StudentId: v.StudentID}}, nil
 }
 
 func (s *ArticleService) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleReply, error) {
```

