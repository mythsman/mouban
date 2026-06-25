package crawl

import "github.com/spf13/viper"

func isCrawlEnabled() bool {
	if !viper.IsSet("crawl.enable") {
		return true
	}
	return viper.GetBool("crawl.enable")
}

func isStorageEnabled() bool {
	if !viper.IsSet("storage.enable") {
		return true
	}
	return viper.GetBool("storage.enable")
}
