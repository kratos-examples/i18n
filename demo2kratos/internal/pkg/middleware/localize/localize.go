package localize

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize/i18n_message"
	"github.com/yylego/tern/zerotern"
	"github.com/yylego/zaplog"
	"go.uber.org/zap"
)

type localizerKey struct{}

// I18N 组件用于翻译
// cp from: https://github.com/go-kratos/examples/blob/3a46aa32f7dbecbb01f2e3ecb28af187b2d9b53c/i18n/internal/pkg/middleware/localize/localize.go#L16
func I18N() middleware.Middleware {
	bundle, _ := i18n_message.LoadI18nFiles(false)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var acceptLanguage string
			// parse accept-language from transport context
			if tsp, ok := transport.FromServerContext(ctx); ok {
				acceptLanguage = tsp.RequestHeader().Get("accept-language") //这里有可能得到个空字符串
			}
			// when accept == "" set accept = defaultValue
			acceptLanguage = zerotern.VV(acceptLanguage, i18n_message.DefaultLanguage.String()) //在这里设置默认语言

			localizer := i18n.NewLocalizer(bundle, acceptLanguage)
			ctx = context.WithValue(ctx, localizerKey{}, localizer)
			return handler(ctx, req)
		}
	}
}

func FromContext(ctx context.Context) *LocalizerI18n {
	return NewLocalizerI18n(ctx.Value(localizerKey{}).(*i18n.Localizer))
}

// LocalizerI18n 把原始的翻译逻辑封装起来，增加些装饰逻辑，让翻译"永不报错"，因为更多的时候服务优先是逻辑正确，语言翻译不应该影响业务
type LocalizerI18n struct {
	localizer *i18n.Localizer
}

func NewLocalizerI18n(localizer *i18n.Localizer) *LocalizerI18n {
	return &LocalizerI18n{localizer: localizer}
}

// Localize 封装 localizer 的 Localize 函数，但"永不报错"而不抛出异常，单返回值只返回消息而不返回错误
func (T *LocalizerI18n) Localize(config *i18n.LocalizeConfig) string {
	if T.localizer == nil {
		zaplog.LOG.Warn("错误的翻译器-没有完成初始化逻辑")
		return "[unknown]" + config.MessageID
	}
	if config.MessageID == "" {
		zaplog.LOG.Warn("翻译参数错误-没有填写消息的ID")
		return "[unknown]"
	}
	if config.DefaultMessage == nil {
		config.DefaultMessage = &i18n.Message{ // 这里会改参数的配置，但认为没事，因为外部用完配置就会丢掉的，而且默认值不会造成负面影响
			ID:    config.MessageID,
			Other: "[message]" + config.MessageID, // 配置默认消息，当某种语言缺翻译时会返回这种格式的
		}
	}
	msg, err := T.localizer.Localize(config)
	if err != nil {
		zaplog.LOG.Warn("翻译遇到问题-理论上这里不该出错", zap.Error(err))
		return "[unknown]" + config.MessageID
	}
	return msg //这就是翻译成功的信息
}
