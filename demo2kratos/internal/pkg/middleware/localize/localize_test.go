package localize_test

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/stretchr/testify/require"
	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize"
	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/middleware/localize/i18n_message"
	"github.com/yylego/neatjson/neatjsons"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// 其实也没必要对 middleware 进行单元测试的，这里只是自己尝试写写，掌握点新技能
func TestI18N(t *testing.T) {
	var rawFunc middleware.Handler = func(ctx context.Context, req any) (any, error) {
		msg := localize.FromContext(ctx).Localize(i18n_message.I18nWelcome(&i18n_message.WelcomeParam{
			User: "abc",
		}))
		return wrapperspb.String(msg), nil
	}

	runFunc := localize.I18N()(rawFunc)

	// Unable to create a context with transport metadata (e.g., accept-language header) without simulating a web request.
	ctx := context.Background()
	// This unit test only tests the default English language due to the absence of transport metadata.
	res, err := runFunc(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	t.Log(neatjsons.S(res))

	stringValue, ok := res.(*wrapperspb.StringValue)
	require.True(t, ok)
	require.Equal(t, "Welcome to our service, abc!", stringValue.Value)
}
