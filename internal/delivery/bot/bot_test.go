package bot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg-bot/internal/domain"
	"github.com/mr-linch/go-tg-bot/internal/service/mocks"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	bot, err := New(&Deps{
		Service: mocks.NewService(t),
	})

	assert.NoError(t, err)

	assert.NotNil(t, bot)
}

type tgClientCall struct {
	Method   string
	Params   url.Values
	Response string
}

func newTgClient(t *testing.T, calls []tgClientCall) *tg.Client {
	t.Helper()

	call := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.True(t, strings.HasPrefix(r.URL.Path, "/bot1234:secret/"))

		if call >= len(calls) {
			t.Fatalf("unexpected call %d", call)
		}

		assert.Equal(t, calls[call].Method, strings.ReplaceAll(r.URL.Path, "/bot1234:secret/", ""))
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, calls[call].Params, r.Form)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(calls[call].Response))

		call++
	}))

	t.Cleanup(func() {
		if call < len(calls)-1 {
			t.Fatalf("excepted calls: %d, actual calls: %d", len(calls), call)
		}
		server.Close()
	})

	tgClient := tg.New("1234:secret",
		tg.WithClientDoer(server.Client()),
		tg.WithClientServerURL(server.URL),
	)

	return tgClient
}

func TestBot_onStart(t *testing.T) {
	ctx := context.Background()

	authService := mocks.NewAuth(t)

	authService.On("AuthViaBot", ctx, &tg.User{
		ID:        1234,
		FirstName: "John",
	}).Return(&domain.User{
		ID:         1,
		FirstName:  "John",
		TelegramID: tg.UserID(1234),
	}, nil)

	service := mocks.NewService(t)
	service.On("Auth").Return(authService)

	bot, err := New(&Deps{
		Service: service,
	})

	assert.NoError(t, err)

	err = bot.Handle(ctx, &tgb.Update{
		Client: newTgClient(t, []tgClientCall{
			{
				"getMe",
				url.Values{},
				`{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"John","username":"john"}}`,
			},
			{
				"sendMessage",
				url.Values{
					"chat_id":    []string{"1234"},
					"parse_mode": []string{"HTML"},
					"text":       []string{"Hey, John!\n\n<strong>Your Bot ID:</strong> <code>1</code>\n<strong>Your Telegram ID:</strong> <code>1234</code>"},
				},
				`{"ok":true, "result": {}}`,
			},
		}),
		Update: &tg.Update{
			Message: &tg.Message{
				From: &tg.User{
					ID:        1234,
					FirstName: "John",
				},
				Chat: tg.Chat{
					ID:   1234,
					Type: tg.ChatTypePrivate,
				},
				Text: "/start",
			},
		},
	})

	assert.NoError(t, err)

}
