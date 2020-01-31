// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	"github.com/jolheiser/gitea/models"
	"github.com/jolheiser/gitea/modules/base"
	"github.com/jolheiser/gitea/modules/context"
	"github.com/jolheiser/gitea/modules/setting"
	structs "github.com/jolheiser/gitea/sdk"
	"github.com/jolheiser/gitea/routers"
)

const (
	tplOrgs base.TplName = "admin/org/list"
)

// Organizations show all the organizations
func Organizations(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("admin.organizations")
	ctx.Data["PageIsAdmin"] = true
	ctx.Data["PageIsAdminOrganizations"] = true

	routers.RenderUserSearch(ctx, &models.SearchUserOptions{
		Type: models.UserTypeOrganization,
		ListOptions: models.ListOptions{
			PageSize: setting.UI.Admin.OrgPagingNum,
		},
		Visible: []structs.VisibleType{structs.VisibleTypePublic, structs.VisibleTypeLimited, structs.VisibleTypePrivate},
	}, tplOrgs)
}
