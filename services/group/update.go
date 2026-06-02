// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package group

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	group_model "gitea.dev/models/group"
	"gitea.dev/modules/optional"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/structs"
	"gitea.dev/modules/util"
)

type UpdateOptions struct {
	Name        optional.Option[string]
	Description optional.Option[string]
	Visibility  optional.Option[structs.VisibleType]
}

func UpdateGroup(ctx context.Context, g *group_model.Group, opts *UpdateOptions) error {
	var nameChanged bool
	var oldName string
	if opts.Name.Has() {
		oldName = g.Name
		nameChanged = !strings.EqualFold(opts.Name.Value(), g.Name)
		g.Name = opts.Name.Value()
		g.LowerName = strings.ToLower(g.Name)
	}
	if opts.Description.Has() {
		g.Description = opts.Description.Value()
	}
	if opts.Visibility.Has() {
		g.Visibility = opts.Visibility.Value()
	}
	if nameChanged {
		parentDir := filepath.Dir(filepath.Join(setting.RepoRootPath, filepath.FromSlash(g.FullPath(ctx))))
		oldDir := filepath.Join(parentDir, oldName)
		ndir := filepath.Join(setting.RepoRootPath, filepath.FromSlash(g.FullPath(ctx)))
		if err := util.Rename(oldDir, ndir); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}
	}
	return group_model.UpdateGroup(ctx, g)
}
