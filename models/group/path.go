// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package group

import (
	"context"
	"fmt"

	"gitea.dev/models/db"
	user_model "gitea.dev/models/user"
	giturl "gitea.dev/modules/git/url"
	"gitea.dev/modules/log"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/util"
	"xorm.io/builder"
)

func PathByID(gid int64, ctxs ...context.Context) (string, error) {
	if gid <= 0 {
		return "", nil
	}
	ctx := util.OptionalArg(ctxs, context.TODO())
	var strs []string
	err := db.GetEngine(ctx).SQL(groupPathCTEBuilder()+`
select path from group_hierarchy where id = ?`, gid).Find(&strs)
	if err != nil {
		log.Error("unable to find group path: %w", err)
		return "", err
	}
	if len(strs) < 1 {
		return "", nil
	}
	return strs[0], nil
}

func PathsByIDs(gids []int64, ctxs ...context.Context) (m map[int64]string, err error) {
	type idPath struct {
		ID   int64
		Path string
	}
	m = map[int64]string{}
	if len(gids) == 0 {
		return m, err
	}
	ctx := util.OptionalArg(ctxs, context.TODO())
	selPart := builder.Select("id", "path").
		From("group_hierarchy").
		Where(builder.In("id", gids))

	selSQL, args, err := builder.ToSQL(selPart)
	if err != nil {
		return m, err
	}
	r := new(idPath)
	rows, err := db.GetEngine(ctx).SQL(groupPathCTEBuilder()+"\n"+selSQL, args...).Rows(r)
	if err != nil {
		return m, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(r); err != nil {
			return map[int64]string{}, err
		}
		m[r.ID] = r.Path
	}
	return m, err
}

func groupPathCTEBuilder() string {
	var recursiveKeyword string
	if !setting.Database.Type.IsMSSQL() {
		recursiveKeyword = " RECURSIVE"
	}
	return fmt.Sprintf(`WITH%s %s`, recursiveKeyword, groupHierarchyCTEBuilder(nil))
}

func GetGroupByPathname(ctx context.Context, owner, pathname string) (*Group, error) {
	pathname = giturl.NormalizeGroupPath(pathname)
	rawSQL := groupPathCTEBuilder() + `
SELECT *
FROM group_hierarchy
WHERE owner_name = ? and path = ?;`
	g := new(Group)
	has, err := db.GetEngine(ctx).SQL(rawSQL, owner, pathname).Get(g)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrGroupNotExist{Path: pathname}
	}

	return g, nil
}

func IDByPathname(ctx context.Context, ownerID int64, pathname string) (int64, error) {
	pathname = giturl.NormalizeGroupPath(pathname)
	if pathname == "" {
		return 0, nil
	}
	owner, err := user_model.GetUserByID(ctx, ownerID)
	if err != nil {
		return 0, nil
	}

	rg, err := GetGroupByPathname(ctx, owner.LowerName, pathname)
	if err != nil {
		return 0, err
	}
	if rg == nil {
		return 0, nil
	}
	return rg.ID, nil
}
