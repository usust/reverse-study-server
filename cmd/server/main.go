package main

import (
	"reverse-study-server/internal/bootstrap"
	localhttp "reverse-study-server/internal/transport/http"
)

func main() {
	// _ = os.Setenv("DATA_SECURITY_LAB_CONFIG_FILE", "/Users/lyu/Code/sec-study/platform/config.yaml")

	if err := bootstrap.Init(); err != nil {
		return
	}
	localhttp.InitRouter().Run("0.0.0.0:10000")
}
