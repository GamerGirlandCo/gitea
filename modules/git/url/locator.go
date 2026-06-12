// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package url

import (
	"net/url"
	"strings"

	"gitea.dev/modules/util"
)

// Locator holds information needed to build various repository paths
type Locator struct {
	Owner     string
	Repo      string
	GroupPath string
}

func (l Locator) groupSegment() string {
	return strings.Trim(l.GroupPath, "/")
}

func (l Locator) groupSegmentWithTrailingSlash() string {
	return util.Iif(l.GroupPath != "", l.groupSegment()+"/", "")
}

func (l Locator) urlEscapedGroupSegment() string {
	return strings.Join(util.SliceMap(strings.Split(l.groupSegment(), "/"), url.PathEscape), "/")
}

func (l Locator) urlEscapedGroupSegmentWithTrailingSlash() string {
	return util.Iif(l.GroupPath != "", l.urlEscapedGroupSegment()+"/", "")
}

func (l Locator) cloneGroupSegment() string {
	if l.GroupPath == "" {
		return ""
	}
	return FormatExplicitGroupPath(l.urlEscapedGroupSegment()) + "/"
}

func (l Locator) ClonePath() string {
	return url.PathEscape(l.Owner) + "/" + l.cloneGroupSegment() + url.PathEscape(l.Repo)
}

func (l Locator) StoragePath() string {
	return strings.ToLower(l.Owner) + "/" + strings.ToLower(l.groupSegmentWithTrailingSlash()) + strings.ToLower(l.Repo)
}

func (l Locator) WebPath() string {
	return url.PathEscape(l.Owner) + "/" + l.urlEscapedGroupSegmentWithTrailingSlash() + url.PathEscape(l.Repo)
}

func (l Locator) FullName() string {
	return l.Owner + "/" + l.groupSegmentWithTrailingSlash() + l.Repo
}

func NewLocator(ownerName, repoName, groupPath string) Locator {
	return Locator{
		Owner:     ownerName,
		Repo:      repoName,
		GroupPath: groupPath,
	}
}
