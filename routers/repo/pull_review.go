// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"fmt"

	"github.com/jolheiser/gitea/models"
	"github.com/jolheiser/gitea/modules/auth"
	"github.com/jolheiser/gitea/modules/context"
	"github.com/jolheiser/gitea/modules/log"
	pull_service "github.com/jolheiser/gitea/services/pull"
)

// CreateCodeComment will create a code comment including an pending review if required
func CreateCodeComment(ctx *context.Context, form auth.CodeCommentForm) {
	issue := GetActionIssue(ctx)
	if !issue.IsPull {
		return
	}
	if ctx.Written() {
		return
	}

	if ctx.HasError() {
		ctx.Flash.Error(ctx.Data["ErrorMsg"].(string))
		ctx.Redirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
		return
	}

	signedLine := form.Line
	if form.Side == "previous" {
		signedLine *= -1
	}

	comment, err := pull_service.CreateCodeComment(
		ctx.User,
		ctx.Repo.GitRepo,
		issue,
		signedLine,
		form.Content,
		form.TreePath,
		form.IsReview,
		form.Reply,
		form.LatestCommitID,
	)
	if err != nil {
		ctx.ServerError("CreateCodeComment", err)
		return
	}

	log.Trace("Comment created: %d/%d/%d", ctx.Repo.Repository.ID, issue.ID, comment.ID)

	if comment != nil {
		ctx.Redirect(comment.HTMLURL())
	} else {
		ctx.Redirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
	}
}

// SubmitReview creates a review out of the existing pending review or creates a new one if no pending review exist
func SubmitReview(ctx *context.Context, form auth.SubmitReviewForm) {
	issue := GetActionIssue(ctx)
	if !issue.IsPull {
		return
	}
	if ctx.Written() {
		return
	}
	if ctx.HasError() {
		ctx.Flash.Error(ctx.Data["ErrorMsg"].(string))
		ctx.Redirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
		return
	}

	reviewType := form.ReviewType()
	switch reviewType {
	case models.ReviewTypeUnknown:
		ctx.ServerError("ReviewType", fmt.Errorf("unknown ReviewType: %s", form.Type))
		return

	// can not approve/reject your own PR
	case models.ReviewTypeApprove, models.ReviewTypeReject:
		if issue.IsPoster(ctx.User.ID) {
			var translated string
			if reviewType == models.ReviewTypeApprove {
				translated = ctx.Tr("repo.issues.review.self.approval")
			} else {
				translated = ctx.Tr("repo.issues.review.self.rejection")
			}

			ctx.Flash.Error(translated)
			ctx.Redirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
			return
		}
	}

	_, comm, err := pull_service.SubmitReview(ctx.User, ctx.Repo.GitRepo, issue, reviewType, form.Content, form.CommitID)
	if err != nil {
		if models.IsContentEmptyErr(err) {
			ctx.Flash.Error(ctx.Tr("repo.issues.review.content.empty"))
			ctx.Redirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
		} else {
			ctx.ServerError("SubmitReview", err)
		}
		return
	}

	ctx.Redirect(fmt.Sprintf("%s/pulls/%d#%s", ctx.Repo.RepoLink, issue.Index, comm.HashTag()))
}
