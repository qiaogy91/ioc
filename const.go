package ioc

import "github.com/qiaogy91/ioc/config"

func OtlpMustEnabled() {
	if Config().Get(config.OtlpName) == nil {
		panic("Current module depends on Otlp, please add Otlp to ioc first")
	}
}
