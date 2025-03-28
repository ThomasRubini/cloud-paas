package utils

import (
	"github.com/ThomasRubini/cloud-paas/internal/paas_cli/config"
	"github.com/go-resty/resty/v2"
)

func GetAPIClient() *resty.Client {
	r := resty.New()
	r.SetBaseURL(config.Get().BACKEND_URL)
	c := config.Get()
	if c.REFRESH_TOKEN != "" {
		r.SetHeader("Authorization", "Bearer "+c.REFRESH_TOKEN)
	}

	return r
}
