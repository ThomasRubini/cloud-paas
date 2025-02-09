package paas_cli

import (
	"github.com/ThomasRubini/cloud-paas/internal/paas_cli/config"

	"gopkg.in/resty.v1"
)

func getAPIClient() *resty.Client {
	r := resty.New()
	c := config.Get()
	if c.REFRESH_TOKEN != "" {
		r.SetHeader("Authorization", "Bearer "+c.REFRESH_TOKEN)
	}

	return r
}
