// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package common

import (
	"testing"

	"code.gitea.io/gitea/modules/optional"

	"github.com/stretchr/testify/assert"
)

func TestCompareRouterReq(t *testing.T) {
	cases := []struct {
		input            string
		CompareRouterReq *CompareRouterReq
	}{
		{
			input:            "",
			CompareRouterReq: &CompareRouterReq{},
		},
		{
			input: "v1.0...v1.1",
			CompareRouterReq: &CompareRouterReq{
				BaseOriRef:       "v1.0",
				CompareSeparator: "...",
				HeadGroupPath:    optional.Some(""),
				HeadOriRef:       "v1.1",
			},
		},
		{
			input: "main..develop",
			CompareRouterReq: &CompareRouterReq{
				BaseOriRef:       "main",
				CompareSeparator: "..",
				HeadGroupPath:    optional.Some(""),
				HeadOriRef:       "develop",
			},
		},
		{
			input: "main^...develop",
			CompareRouterReq: &CompareRouterReq{
				BaseOriRef:       "main",
				BaseOriRefSuffix: "^",
				CompareSeparator: "...",
				HeadGroupPath:    optional.Some(""),
				HeadOriRef:       "develop",
			},
		},
		{
			input: "main^^^^^...develop",
			CompareRouterReq: &CompareRouterReq{
				BaseOriRef:       "main",
				BaseOriRefSuffix: "^^^^^",
				CompareSeparator: "...",
				HeadGroupPath:    optional.Some(""),
				HeadOriRef:       "develop",
			},
		},
		{
			input: "develop",
			CompareRouterReq: &CompareRouterReq{
				CompareSeparator: "...",
				HeadGroupPath:    optional.Some(""),
				HeadOriRef:       "develop",
			},
		},
		{
			input: "teabot:feature1",
			CompareRouterReq: &CompareRouterReq{
				CompareSeparator: "...",
				HeadOwner:        "teabot",
				HeadOriRef:       "feature1",
			},
		},
		{
			input: "lunny/forked_repo:develop",
			CompareRouterReq: &CompareRouterReq{
				CompareSeparator: "...",
				HeadGroupPath:    optional.Some(""),
				HeadOwner:        "lunny",
				HeadRepoName:     "forked_repo",
				HeadOriRef:       "develop",
			},
		},
		{
			input: "main...lunny/forked_repo:develop",
			CompareRouterReq: &CompareRouterReq{
				BaseOriRef:       "main",
				CompareSeparator: "...",
				HeadGroupPath:    optional.Some(""),
				HeadOwner:        "lunny",
				HeadRepoName:     "forked_repo",
				HeadOriRef:       "develop",
			},
		},
		{
			input: "main^...lunny/forked_repo:develop",
			CompareRouterReq: &CompareRouterReq{
				BaseOriRef:       "main",
				BaseOriRefSuffix: "^",
				CompareSeparator: "...",
				HeadGroupPath:    optional.Some(""),
				HeadOwner:        "lunny",
				HeadRepoName:     "forked_repo",
				HeadOriRef:       "develop",
			},
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.CompareRouterReq, ParseCompareRouterParam(c.input), "input: %s", c.input)
	}
}
