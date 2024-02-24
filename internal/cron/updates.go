package cron

import (
	iLog "log"
	"moscode-cdn-fiber/configs"
	"moscode-cdn-fiber/internal/pkg/httpclient"
)

func UpdateStaticJob() {
	appConfig := configs.GetConfig()
	if appConfig.BaseURL != "" {
		if err := httpclient.UpdateStatic(appConfig.BaseURL + "/update"); err != nil {
			iLog.Println("Error handling update:", err)
		}
	}
}

func UpdateIndexJob() {
	appConfig := configs.GetConfig()
	if appConfig.BaseURL != "" {
		err := httpclient.RefreshFile(appConfig.BaseURL, "index.html", appConfig.TempDir, appConfig.StaticDir)
		if err != nil {
			iLog.Println("Error refreshing file:", err)
		}
	}
}
