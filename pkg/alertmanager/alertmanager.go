package alertmanager

import (
	"fmt"
	"net/url"

	"gopkg.in/yaml.v2"

	alertmanagerv1 "github.com/linclaus/alertmanager-operator/api/v1"
	alertmangerconfig "github.com/prometheus/alertmanager/config"
	commoncfg "github.com/prometheus/common/config"
	prometheuscfg "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/spf13/viper"
)

var (
	ALERTMANAGER_HOST            = "http://localhost:9093"
	ALERTMANAGER_GROUP_INTERVAL  = "5m"
	ALERTMANAGER_REPEAT_INTERVAL = "1d"
)

// Config is the top-level configuration for Alertmanager's config files.
type Config struct {
	Global       *GlobalConfig                    `yaml:"global,omitempty" json:"global,omitempty"`
	Route        *alertmangerconfig.Route         `yaml:"route,omitempty" json:"route,omitempty"`
	InhibitRules []*alertmangerconfig.InhibitRule `yaml:"inhibit_rules,omitempty" json:"inhibit_rules,omitempty"`
	Receivers    []*Receiver                      `yaml:"receivers,omitempty" json:"receivers,omitempty"`
	Templates    []string                         `yaml:"templates" json:"templates"`

	// original is the input from which the config was parsed.
	original string
}

func (c Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

// GlobalConfig defines configuration parameters that are valid globally
// unless overwritten.
type GlobalConfig struct {
	// ResolveTimeout is the time after which an alert is declared resolved
	// if it has not been updated.
	ResolveTimeout model.Duration `yaml:"resolve_timeout" json:"resolve_timeout"`

	HTTPConfig *commoncfg.HTTPClientConfig `yaml:"http_config,omitempty" json:"http_config,omitempty"`

	SMTPFrom         string                     `yaml:"smtp_from,omitempty" json:"smtp_from,omitempty"`
	SMTPHello        string                     `yaml:"smtp_hello,omitempty" json:"smtp_hello,omitempty"`
	SMTPSmarthost    alertmangerconfig.HostPort `yaml:"smtp_smarthost,omitempty" json:"smtp_smarthost,omitempty"`
	SMTPAuthUsername string                     `yaml:"smtp_auth_username,omitempty" json:"smtp_auth_username,omitempty"`
	SMTPAuthPassword string                     `yaml:"smtp_auth_password,omitempty" json:"smtp_auth_password,omitempty"`
	SMTPAuthSecret   string                     `yaml:"smtp_auth_secret,omitempty" json:"smtp_auth_secret,omitempty"`
	SMTPAuthIdentity string                     `yaml:"smtp_auth_identity,omitempty" json:"smtp_auth_identity,omitempty"`
	SMTPRequireTLS   bool                       `yaml:"smtp_require_tls,omitempty" json:"smtp_require_tls,omitempty"`
	SlackAPIURL      *alertmangerconfig.URL     `yaml:"slack_api_url,omitempty" json:"slack_api_url,omitempty"`
	PagerdutyURL     *alertmangerconfig.URL     `yaml:"pagerduty_url,omitempty" json:"pagerduty_url,omitempty"`
	HipchatAPIURL    *alertmangerconfig.URL     `yaml:"hipchat_api_url,omitempty" json:"hipchat_api_url,omitempty"`
	HipchatAuthToken string                     `yaml:"hipchat_auth_token,omitempty" json:"hipchat_auth_token,omitempty"`
	OpsGenieAPIURL   *alertmangerconfig.URL     `yaml:"opsgenie_api_url,omitempty" json:"opsgenie_api_url,omitempty"`
	OpsGenieAPIKey   string                     `yaml:"opsgenie_api_key,omitempty" json:"opsgenie_api_key,omitempty"`
	WeChatAPIURL     *alertmangerconfig.URL     `yaml:"wechat_api_url,omitempty" json:"wechat_api_url,omitempty"`
	WeChatAPISecret  string                     `yaml:"wechat_api_secret,omitempty" json:"wechat_api_secret,omitempty"`
	WeChatAPICorpID  string                     `yaml:"wechat_api_corp_id,omitempty" json:"wechat_api_corp_id,omitempty"`
	VictorOpsAPIURL  *alertmangerconfig.URL     `yaml:"victorops_api_url,omitempty" json:"victorops_api_url,omitempty"`
	VictorOpsAPIKey  string                     `yaml:"victorops_api_key,omitempty" json:"victorops_api_key,omitempty"`
}

// Receiver configuration provides configuration on how to contact a receiver.
type Receiver struct {
	// A unique identifier for this receiver.
	Name string `yaml:"name" json:"name"`

	EmailConfigs     []*EmailConfig                       `yaml:"email_configs,omitempty" json:"email_configs,omitempty"`
	PagerdutyConfigs []*alertmangerconfig.PagerdutyConfig `yaml:"pagerduty_configs,omitempty" json:"pagerduty_configs,omitempty"`
	HipchatConfigs   []*alertmangerconfig.HipchatConfig   `yaml:"hipchat_configs,omitempty" json:"hipchat_configs,omitempty"`
	SlackConfigs     []*alertmangerconfig.SlackConfig     `yaml:"slack_configs,omitempty" json:"slack_configs,omitempty"`
	WebhookConfigs   []*alertmangerconfig.WebhookConfig   `yaml:"webhook_configs,omitempty" json:"webhook_configs,omitempty"`
	OpsGenieConfigs  []*alertmangerconfig.OpsGenieConfig  `yaml:"opsgenie_configs,omitempty" json:"opsgenie_configs,omitempty"`
	WechatConfigs    []*alertmangerconfig.WechatConfig    `yaml:"wechat_configs,omitempty" json:"wechat_configs,omitempty"`
	PushoverConfigs  []*alertmangerconfig.PushoverConfig  `yaml:"pushover_configs,omitempty" json:"pushover_configs,omitempty"`
	VictorOpsConfigs []*alertmangerconfig.VictorOpsConfig `yaml:"victorops_configs,omitempty" json:"victorops_configs,omitempty"`
}

// EmailConfig configures notifications via mail.
type EmailConfig struct {
	alertmangerconfig.NotifierConfig `yaml:",inline" json:",inline"`

	// Email address to notify.
	To           string                     `yaml:"to,omitempty" json:"to,omitempty"`
	From         string                     `yaml:"from,omitempty" json:"from,omitempty"`
	Hello        string                     `yaml:"hello,omitempty" json:"hello,omitempty"`
	Smarthost    alertmangerconfig.HostPort `yaml:"smarthost,omitempty" json:"smarthost,omitempty"`
	AuthUsername string                     `yaml:"auth_username,omitempty" json:"auth_username,omitempty"`
	AuthPassword string                     `yaml:"auth_password,omitempty" json:"auth_password,omitempty"`
	AuthSecret   string                     `yaml:"auth_secret,omitempty" json:"auth_secret,omitempty"`
	AuthIdentity string                     `yaml:"auth_identity,omitempty" json:"auth_identity,omitempty"`
	Headers      map[string]string          `yaml:"headers,omitempty" json:"headers,omitempty"`
	HTML         string                     `yaml:"html,omitempty" json:"html,omitempty"`
	Text         string                     `yaml:"text,omitempty" json:"text,omitempty"`
	RequireTLS   *bool                      `yaml:"require_tls,omitempty" json:"require_tls,omitempty"`
	TLSConfig    prometheuscfg.TLSConfig    `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
}

func updateReceivers(rvs []*Receiver, receiverName string, contactValues []string) []*Receiver {
	//TODO multiEmailConfig if multiContactValues
	var rv *Receiver
	index := -1
	for i, receive := range rvs {
		if receive.Name == receiverName {
			index = i
			break
		}
	}
	ecs := []*EmailConfig{}
	if index != -1 {
		ecs = append(ecs, rvs[index].EmailConfigs...)
	}
	for _, cv := range contactValues {
		ec := &EmailConfig{
			Headers: map[string]string{
				"Subject": "{{ template \"email.custom.subject\" . }}",
			},
			HTML: "{{ template \"email.custom.html\" . }}",
		}
		ec.To = cv
		ecs = append(ecs, ec)
	}
	ecs = removeDuplicatedEmailConfigs(ecs)
	rawurl, _ := url.Parse(viper.GetViper().GetString("webhook"))
	wc := &alertmangerconfig.WebhookConfig{
		URL: &alertmangerconfig.URL{
			URL: rawurl,
		},
	}
	rv = &Receiver{
		Name:           receiverName,
		EmailConfigs:   ecs,
		WebhookConfigs: []*alertmangerconfig.WebhookConfig{wc},
	}
	if index == -1 {
		return append(rvs, rv)
	} else {
		rvs[index] = rv
		return rvs
	}
}

func deleteReceivers(rvs []*Receiver, receiverName string, contactValues []string, all bool) ([]*Receiver, bool) {
	index := -1
	for i, receive := range rvs {
		if receive.Name == receiverName {
			index = i
			break
		}
	}

	if index != -1 {
		if all {
			return append(rvs[:index], rvs[index+1:]...), true
		} else {
			ecs := []*EmailConfig{}
			for _, ec := range rvs[index].EmailConfigs {
				exists := false
				for _, cv := range contactValues {
					if ec.To == cv {
						exists = true
						break
					}
				}
				if !exists {
					ecs = append(ecs, ec)
				}
			}
			rv := &Receiver{
				Name:         receiverName,
				EmailConfigs: ecs,
			}
			if len(ecs) == 0 {
				return append(rvs[:index], rvs[index+1:]...), true
			} else {
				rvs[index] = rv
				return rvs, false
			}
		}

	} else {
		return rvs, false
	}
}

func updateRoutes(rts []*alertmangerconfig.Route, newRt alertmanagerv1.Route) []*alertmangerconfig.Route {
	var rt *alertmangerconfig.Route
	index := -1
	for i, route := range rts {
		if route.Receiver == newRt.Receiver {
			index = i
			break
		}
	}
	ri, _ := model.ParseDuration(ALERTMANAGER_REPEAT_INTERVAL)
	gi, _ := model.ParseDuration(ALERTMANAGER_GROUP_INTERVAL)
	rt = &alertmangerconfig.Route{
		Match:          newRt.Match,
		Receiver:       newRt.Receiver,
		RepeatInterval: &ri,
		GroupInterval:  &gi,
		GroupByStr:     newRt.GroupBy,
	}
	if index == -1 {
		return append(rts, rt)
	} else {
		rts[index] = rt
		return rts
	}
}

func deleteRoutes(rts []*alertmangerconfig.Route, receiverName string) []*alertmangerconfig.Route {
	index := -1
	for i, route := range rts {
		if route.Receiver == receiverName {
			index = i
			break
		}
	}
	if index != -1 {
		return append(rts[:index], rts[index+1:]...)
	} else {
		return rts
	}
}

func removeDuplicatedEmailConfigs(ems []*EmailConfig) []*EmailConfig {
	m := map[string]*EmailConfig{}
	for _, em := range ems {
		m[em.To] = em
	}
	rst := []*EmailConfig{}
	for _, v := range m {
		rst = append(rst, v)
	}
	return rst
}
