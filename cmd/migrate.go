// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"

	"github.com/jolheiser/gitea/models"
	"github.com/jolheiser/gitea/models/migrations"
	"github.com/jolheiser/gitea/modules/log"
	"github.com/jolheiser/gitea/modules/setting"

	"github.com/urfave/cli"
)

// CmdMigrate represents the available migrate sub-command.
var CmdMigrate = cli.Command{
	Name:        "migrate",
	Usage:       "Migrate the database",
	Description: "This is a command for migrating the database, so that you can run gitea admin create-user before starting the server.",
	Action:      runMigrate,
}

func runMigrate(ctx *cli.Context) error {
	if err := initDB(); err != nil {
		return err
	}

	log.Trace("AppPath: %s", setting.AppPath)
	log.Trace("AppWorkPath: %s", setting.AppWorkPath)
	log.Trace("Custom path: %s", setting.CustomPath)
	log.Trace("Log path: %s", setting.LogRootPath)
	setting.InitDBConfig()

	if err := models.NewEngine(context.Background(), migrations.Migrate); err != nil {
		log.Fatal("Failed to initialize ORM engine: %v", err)
		return err
	}

	return nil
}
