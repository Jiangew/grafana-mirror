package notifiers

import (
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
)

const defaultWeChatMsgType = "link"

func init() {
	alerting.RegisterNotifier(&alerting.NotifierPlugin{
		Type:        "wechat",
		Name:        "WeChat",
		Description: "Sends HTTP POST request to WeChat",
		Heading:     "WeChat settings",
		Factory:     newWeChatNotifier,
		Options: []alerting.NotifierOption{
			{
				Label:        "Url",
				Element:      alerting.ElementTypeInput,
				InputType:    alerting.InputTypeText,
				Placeholder:  "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxxx",
				PropertyName: "url",
				Required:     true,
			},
			{
				Label:        "Message Type",
				Element:      alerting.ElementTypeSelect,
				PropertyName: "msgType",
				SelectOptions: []alerting.SelectOption{
					{
						Value: "text",
						Label: "Text"},
					{
						Value: "markdown",
						Label: "Markdown",
					},
				},
			},
		},
	})
}

func newWeChatNotifier(model *models.AlertNotification) (alerting.Notifier, error) {
	url := model.Settings.Get("url").MustString()
	if url == "" {
		return nil, alerting.ValidationError{Reason: "Could not find url property in settings"}
	}

	msgType := model.Settings.Get("msgType").MustString(defaultWeChatMsgType)

	return &WeChatNotifier{
		NotifierBase: NewNotifierBase(model),
		MsgType:      msgType,
		URL:          url,
		log:          log.New("alerting.notifier.wechat"),
	}, nil
}

// WeChatNotifier is responsible for sending alert notifications to ding ding.
type WeChatNotifier struct {
	NotifierBase
	MsgType string
	URL     string
	log     log.Logger
}

// Notify sends the alert notification to wechat.
func (wc *WeChatNotifier) Notify(evalContext *alerting.EvalContext) error {
	wc.log.Info("Sending wechat")

	messageURL, err := evalContext.GetRuleURL()
	if err != nil {
		wc.log.Error("Failed to get messageUrl", "error", err, "wechat", wc.Name)
		messageURL = ""
	}

	body, err := wc.genBody(evalContext, messageURL)
	if err != nil {
		return err
	}

	cmd := &models.SendWebhookSync{
		Url:  wc.URL,
		Body: string(body),
	}

	if err := bus.DispatchCtx(evalContext.Ctx, cmd); err != nil {
		wc.log.Error("Failed to send WeChat", "error", err, "wechat", wc.Name)
		return err
	}

	return nil
}

func (wc *WeChatNotifier) genBody(evalContext *alerting.EvalContext, messageURL string) ([]byte, error) {
	message := evalContext.Rule.Message
	title := evalContext.GetNotificationTitle()
	if message == "" {
		message = title
	}

	for i, match := range evalContext.EvalMatches {
		message += fmt.Sprintf("\n%2d. %s: %s", i+1, match.Metric, match.Value)
	}

	var bodyMsg map[string]interface{}
	if wc.MsgType == "markdown" {
		bodyMsg = map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"content": message,
			},
		}
	} else {
		bodyMsg = map[string]interface{}{
			"msgtype": "text",
			"text": map[string]string{
				"content": message,
			},
		}
	}
	return json.Marshal(bodyMsg)
}
