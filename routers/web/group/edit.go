// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package group

import (
	"net/http"

	"gitea.dev/models/db"
	group_model "gitea.dev/models/group"
	repo_model "gitea.dev/models/repo"
	"gitea.dev/modules/json"
	"gitea.dev/services/context"
	"gitea.dev/services/forms"
	group_service "gitea.dev/services/group"
	"xorm.io/builder"
)

type movedSubItem struct {
	NewPath  string `json:"newPath"`
	ID       int64  `json:"id"`
	IsGroup  bool   `json:"isGroup"`
	FullName string `json:"fullName,omitempty"`
}
type moveResult struct {
	NewPath  string         `json:"newPath"`
	FullName string         `json:"fullName,omitempty"`
	Children []movedSubItem `json:"children,omitempty"`
}

func MoveGroupItem(ctx *context.Context) {
	form := &forms.MovedGroupItemForm{}
	if err := json.NewDecoder(ctx.Req.Body).Decode(form); err != nil {
		ctx.ServerError("DecodeMovedGroupItemForm", err)
		return
	}
	if err := group_service.MoveGroupItem(ctx, group_service.MoveGroupOptions{
		IsGroup:   form.IsGroup,
		ItemID:    form.ItemID,
		NewPos:    form.NewPos,
		NewParent: form.NewParent,
	}, ctx.Doer); err != nil {
		ctx.ServerError("MoveGroupItem", err)
		return
	}
	var newPath, fullName string
	children := make([]movedSubItem, 0)
	if form.IsGroup {
		grp, err := group_model.GetGroupByID(ctx, form.ItemID)
		if err != nil {
			ctx.ServerError("GetGroupByID", err)
			return
		}
		newPath = grp.GroupLink()

		childGroupIDs, err := group_model.ChildGroupCond(ctx, grp.ID, group_model.AccessibleGroupCondition(ctx.Doer))
		if err != nil {
			ctx.ServerError("ChildGroupCond", err)
			return
		}
		repoParentIDs := make([]int64, len(childGroupIDs)+1)
		copy(repoParentIDs[:len(childGroupIDs)], childGroupIDs)
		repoParentIDs[len(childGroupIDs)] = grp.ID
		pathsToIDs, err := group_model.PathsByIDs(childGroupIDs, ctx)
		if err != nil {
			ctx.ServerError("PathsByIDs", err)
			return
		}
		for id, path := range pathsToIDs {
			children = append(children, movedSubItem{
				NewPath: group_model.Link(grp.OwnerName, path),
				ID:      id,
				IsGroup: true,
			})
		}
		childRepos, _, err := repo_model.SearchRepositoryByCondition(ctx, repo_model.SearchRepoOptions{
			ListOptions: db.ListOptions{ListAll: true},
		}, builder.In("group_id", repoParentIDs), false)
		if err != nil {
			ctx.ServerError("SearchRepositoryByCondition", err)
			return
		}
		for _, r := range childRepos {
			children = append(children, movedSubItem{
				IsGroup:  false,
				ID:       r.ID,
				NewPath:  r.Link(),
				FullName: r.FullName(),
			})
		}
	} else {
		repo, err := repo_model.GetRepositoryByID(ctx, form.ItemID)
		if err != nil {
			ctx.ServerError("GetRepositoryByID", err)
			return
		}
		fullName = repo.FullName()
		newPath = repo.Link()
	}
	ctx.JSON(http.StatusOK, moveResult{
		NewPath:  newPath,
		Children: children,
		FullName: fullName,
	})
}
