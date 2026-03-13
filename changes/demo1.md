# Changes

Code differences compared to source project.

## internal/pkg/middleware/localize/i18n_message/active.en.toml (+11 -0)

```diff
@@ -0,0 +1,11 @@
+[greeting]
+description = "A greeting message with the user's name"
+other = "Hello, {{.Name}}"
+
+[welcome]
+description = "A welcome message for users accessing the service"
+other = "Welcome to our service, {{.User}}!"
+
+[errorMsg]
+description = "A generic error message for failed operations"
+other = "Operation failed: {{.ErrorMessage}}"
```

## internal/pkg/middleware/localize/i18n_message/active.zh.toml (+11 -0)

```diff
@@ -0,0 +1,11 @@
+[greeting]
+description = "A greeting message with the user's name"
+other = "你好，{{.Name}}"
+
+[welcome]
+description = "A welcome message for users accessing the service"
+other = "欢迎使用我们的服务，{{.User}}！"
+
+[errorMsg]
+description = "A generic error message for failed operations"
+other = "操作失败：{{.ErrorMessage}}"
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
+	"github.com/yylego/kratos-examples/demo1kratos/internal/pkg/middleware/localize/i18n_message"
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

## internal/pkg/middleware/localize/i18n_message/i18n.go (+42 -0)

```diff
@@ -0,0 +1,42 @@
+package i18n_message
+
+import (
+	"embed"
+
+	"github.com/BurntSushi/toml"
+	"github.com/nicksnyder/go-i18n/v2/i18n"
+	"github.com/yylego/must"
+	"github.com/yylego/neatjson/neatjsons"
+	"github.com/yylego/rese"
+	"github.com/yylego/zaplog"
+	"golang.org/x/text/language"
+)
+
+// DefaultLanguage 配置默认语言
+var DefaultLanguage = language.English
+
+//go:embed active.en.toml active.zh.toml
+var files embed.FS
+
+func LoadI18nFiles(debugModeOpen bool) (*i18n.Bundle, []*i18n.MessageFile) {
+	bundle := i18n.NewBundle(DefaultLanguage)
+	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
+
+	var messageFiles []*i18n.MessageFile
+	for _, fileName := range []string{"active.en.toml", "active.zh.toml"} {
+		content := rese.A1(files.ReadFile(fileName))
+		//这里文件名 file-name 写 "active.en.toml" 或者 "en.toml" 都行，内部会通过这个解析出语言标签名称
+		messageFile := rese.P1(bundle.ParseMessageFileBytes(content, fileName))
+		if debugModeOpen {
+			zaplog.SUG.Debugln(neatjsons.S(messageFile)) //安利下我的俩工具包
+		}
+		messageFiles = append(messageFiles, messageFile)
+	}
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
+	"github.com/yylego/kratos-examples/demo1kratos/internal/pkg/middleware/localize/i18n_message"
+)
+
+func TestLoadI18nFiles(t *testing.T) {
+	bundle, messageFiles := i18n_message.LoadI18nFiles(true)
+	t.Log(len(messageFiles))
+	t.Log(len(bundle.LanguageTags()))
+}
```

## internal/pkg/middleware/localize/localize.go (+39 -0)

```diff
@@ -0,0 +1,39 @@
+package localize
+
+import (
+	"context"
+
+	"github.com/go-kratos/kratos/v2/middleware"
+	"github.com/go-kratos/kratos/v2/transport"
+	"github.com/nicksnyder/go-i18n/v2/i18n"
+	"github.com/yylego/kratos-examples/demo1kratos/internal/pkg/middleware/localize/i18n_message"
+	"github.com/yylego/tern/zerotern"
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
+func FromContext(ctx context.Context) *i18n.Localizer {
+	return ctx.Value(localizerKey{}).(*i18n.Localizer)
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
+	"github.com/yylego/kratos-examples/demo1kratos/internal/pkg/middleware/localize"
+	"github.com/yylego/kratos-examples/demo1kratos/internal/pkg/middleware/localize/i18n_message"
+	"github.com/yylego/neatjson/neatjsons"
+	"google.golang.org/protobuf/types/known/emptypb"
+	"google.golang.org/protobuf/types/known/wrapperspb"
+)
+
+// 其实也没必要对 middleware 进行单元测试的，这里只是自己尝试写写，掌握点新技能
+func TestI18N(t *testing.T) {
+	var rawFunc middleware.Handler = func(ctx context.Context, req any) (any, error) {
+		msg := localize.FromContext(ctx).MustLocalize(i18n_message.I18nWelcome(&i18n_message.WelcomeParam{
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
 	pb "github.com/yylego/kratos-examples/demo1kratos/api/student"
 	"github.com/yylego/kratos-examples/demo1kratos/internal/conf"
+	"github.com/yylego/kratos-examples/demo1kratos/internal/pkg/middleware/localize"
 	"github.com/yylego/kratos-examples/demo1kratos/internal/service"
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

## internal/service/student.go (+8 -1)

```diff
@@ -5,6 +5,8 @@
 
 	pb "github.com/yylego/kratos-examples/demo1kratos/api/student"
 	"github.com/yylego/kratos-examples/demo1kratos/internal/biz"
+	"github.com/yylego/kratos-examples/demo1kratos/internal/pkg/middleware/localize"
+	"github.com/yylego/kratos-examples/demo1kratos/internal/pkg/middleware/localize/i18n_message"
 )
 
 type StudentService struct {
@@ -22,7 +24,12 @@
 	if ebz != nil {
 		return nil, ebz.Erk
 	}
-	return &pb.CreateStudentReply{Student: &pb.StudentInfo{Id: v.ID, Name: v.Name, Age: v.Age, ClassName: v.ClassName}}, nil
+
+	message := localize.FromContext(ctx).MustLocalize(i18n_message.I18nGreeting(&i18n_message.GreetingParam{
+		Name: v.Name,
+	}))
+
+	return &pb.CreateStudentReply{Student: &pb.StudentInfo{Id: v.ID, Name: message, Age: v.Age, ClassName: v.ClassName}}, nil
 }
 
 func (s *StudentService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentReply, error) {
```

