package controller

import (
	"github.com/G-Research/fasttrackml/pkg/api/aim2/api/response"
	"github.com/G-Research/fasttrackml/pkg/common/api"
	"github.com/G-Research/fasttrackml/pkg/common/middleware/namespace"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func (c Controller) GetTags(ctx *fiber.Ctx) error {
	ns, err := namespace.GetNamespaceFromContext(ctx.Context())
	if err != nil {
		return api.NewInternalError("error getting namespace from context")
	}
	log.Debugf("getTags namespace: %s", ns.Code)

	apps, err := c.tagService.
	if err != nil {
		return err
	}

	resp := response.NewGetAppsResponse(apps)
	log.Debugf("getApps response: %#v", resp)

	return ctx.JSON(resp)
}
