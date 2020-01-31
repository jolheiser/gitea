// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notify

import (
	"net/http"

	"github.com/jolheiser/gitea/models"
	"github.com/jolheiser/gitea/modules/context"
	api "github.com/jolheiser/gitea/sdk"
)

// NewAvailable check if unread notifications exist
func NewAvailable(ctx *context.APIContext) {
	// swagger:operation GET /notifications/new notification notifyNewAvailable
	// ---
	// summary: Check if unread notifications exist
	// responses:
	//   "200":
	//    "$ref": "#/responses/NotificationCount"
	//   "204":
	//     description: No unread notification

	count := models.CountUnread(ctx.User)

	if count > 0 {
		ctx.JSON(http.StatusOK, api.NotificationCount{New: count})
	} else {
		ctx.Status(http.StatusNoContent)
	}
}
