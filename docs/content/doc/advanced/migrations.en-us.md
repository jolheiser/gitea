---
date: "2019-04-15T17:29:00+08:00"
title: "Advanced: Migrations Interfaces"
slug: "migrations-interfaces"
weight: 30
toc: true
draft: false
menu:
  sidebar:
    parent: "advanced"
    name: "Migrations Interfaces"
    weight: 55
    identifier: "migrations-interfaces"
---

# Migration Features

The new migration feature was introduced in Gitea 1.9.0.  
It defines two interfaces to support migrating repository data from other git hosting platforms to Gitea. Currently, only migrations from GitHub via APIv3 to Gitea is implemented.

Gitea defines some standard objects in packages `modules/migrations/base`. They are
 `Repository`, `Milestone`, `Release`, `Label`, `Issue`, `Comment`, `PullRequest`.

## Downloader Interfaces

To add support for migrating from a new git hosting platform:

- You should implement a `Downloader` which will get all repository information.
- You should implement a `DownloaderFactory` which is used to detect if the URL matches and 
creates a `Downloader`.
- You'll need to register the `DownloaderFactory` via `RegisterDownloaderFactory` on init.

```Go
type Downloader interface {
	SetContext(context.Context)
	GetRepoInfo() (*Repository, error)
	GetTopics() ([]string, error)
	GetMilestones() ([]*Milestone, error)
	GetReleases() ([]*Release, error)
	GetLabels() ([]*Label, error)
	GetIssues(page, perPage int) ([]*Issue, bool, error)
	GetComments(issueNumber int64) ([]*Comment, error)
	GetPullRequests(page, perPage int) ([]*PullRequest, error)
}
```

```Go
type DownloaderFactory interface {
	Match(opts MigrateOptions) (bool, error)
	New(opts MigrateOptions) (Downloader, error)
	GitServiceType() structs.GitServiceType
}
```

## Uploader Interfaces

Currently, only a `GiteaLocalUploader` is implemented, so we only save downloaded 
data via this `Uploader` on the local Gitea instance. Other uploaders are not currently supported
and will be implemented in the future.

```Go
type Uploader interface {
	MaxBatchInsertSize(tp string) int
	CreateRepo(repo *Repository, opts MigrateOptions) error
	CreateTopics(topic ...string) error
	CreateMilestones(milestones ...*Milestone) error
	CreateReleases(releases ...*Release) error
	SyncTags() error
	CreateLabels(labels ...*Label) error
	CreateIssues(issues ...*Issue) error
	CreateComments(comments ...*Comment) error
	CreatePullRequests(prs ...*PullRequest) error
	Rollback() error
	Close()
}

```
