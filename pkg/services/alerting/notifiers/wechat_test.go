package notifiers

import (
	"context"
	"testing"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/alerting"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWeChatNotifier(t *testing.T) {
	Convey("WeChat notifier tests", t, func() {
		Convey("empty settings should return error", func() {
			json := `{ }`

			settingsJSON, _ := simplejson.NewJson([]byte(json))
			model := &models.AlertNotification{
				Name:     "wechat_testing",
				Type:     "wechat",
				Settings: settingsJSON,
			}

			_, err := newWeChatNotifier(model)
			So(err, ShouldNotBeNil)
		})
		Convey("settings should trigger incident", func() {
			json := `{ "url": "https://www.google.com" }`

			settingsJSON, _ := simplejson.NewJson([]byte(json))
			model := &models.AlertNotification{
				Name:     "wechat_testing",
				Type:     "wechat",
				Settings: settingsJSON,
			}

			not, err := newWeChatNotifier(model)
			notifier := not.(*WeChatNotifier)

			So(err, ShouldBeNil)
			So(notifier.Name, ShouldEqual, "wechat_testing")
			So(notifier.Type, ShouldEqual, "wechat")
			So(notifier.URL, ShouldEqual, "https://www.google.com")

			Convey("genBody should not panic", func() {
				evalContext := alerting.NewEvalContext(context.Background(),
					&alerting.Rule{
						State:   models.AlertStateAlerting,
						Message: `{host="localhost"}`,
					})
				_, err = notifier.genBody(evalContext, "")
				So(err, ShouldBeNil)
			})
		})
	})
}
