package cron

import (
	iLog "log"
	"moscode-cdn-fiber/configs"
	"moscode-cdn-fiber/internal/pkg/httpclient"
)

func UpdateStaticJob() {
	if configs.GetBaseURL() != "" {
		if err := httpclient.UpdateStatic(configs.GetBaseURL() + "/update"); err != nil {
			iLog.Println("Error handling update:", err)
		}
	}
}

func UpdateIndexJob() {
	if configs.GetBaseURL() != "" {
		err := httpclient.RefreshFile(configs.GetBaseURL(), "index.html", configs.GetTempDir(), configs.GetStaticDir())
		if err != nil {
			iLog.Println("Error refreshing file:", err)
		}
	}
}
