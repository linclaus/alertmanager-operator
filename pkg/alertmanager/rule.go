package alertmanager

import (
	"net/url"

	alertmangerconfig "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
)

var (
	ALERTMANAGER_CONFIG_PATH     string
	ALERTMANAGER_CONFIG_NAME     string
	ALERTMANAGER_HOST            = "http://localhost:9093"
	ALERTMANAGER_GROUP_INTERVAL  = "5m"
	ALERTMANAGER_REPEAT_INTERVAL = "1d"
	WEBHOOK_URL                  string
)

func updateReceivers(rvs []*Receiver, strategyId string, contactValues []string) []*Receiver {
	//TODO multiEmailConfig if multiContactValues
	var rv *Receiver
	index := -1
	for i, receive := range rvs {
		if receive.Name == strategyId {
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
	rawurl, _ := url.Parse(WEBHOOK_URL)
	wc := &alertmangerconfig.WebhookConfig{
		URL: &alertmangerconfig.URL{
			URL: rawurl,
		},
	}
	rv = &Receiver{
		Name:           strategyId,
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

func deleteReceivers(rvs []*Receiver, strategyId string, contactValues []string, all bool) ([]*Receiver, bool) {
	index := -1
	for i, receive := range rvs {
		if receive.Name == strategyId {
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
				Name:         strategyId,
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

func updateRoutes(rts []*alertmangerconfig.Route, strategyId string) []*alertmangerconfig.Route {
	var rt *alertmangerconfig.Route
	index := -1
	for i, route := range rts {
		if route.Receiver == strategyId {
			index = i
			break
		}
	}
	ri, _ := model.ParseDuration(ALERTMANAGER_REPEAT_INTERVAL)
	gi, _ := model.ParseDuration(ALERTMANAGER_GROUP_INTERVAL)
	rt = &alertmangerconfig.Route{
		Match:          map[string]string{"strategy_id": strategyId},
		Receiver:       strategyId,
		RepeatInterval: &ri,
		GroupInterval:  &gi,
		GroupByStr:     []string{"strategy_id"},
	}
	if index == -1 {
		return append(rts, rt)
	} else {
		rts[index] = rt
		return rts
	}
}

func deleteRoutes(rts []*alertmangerconfig.Route, strategyId string) []*alertmangerconfig.Route {
	index := -1
	for i, route := range rts {
		if route.Receiver == strategyId {
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
