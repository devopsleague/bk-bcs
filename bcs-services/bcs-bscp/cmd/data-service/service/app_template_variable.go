/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"bytes"
	"context"

	"bscp.io/pkg/criteria/constant"
	"bscp.io/pkg/dal/table"
	"bscp.io/pkg/kit"
	"bscp.io/pkg/logs"
	pbatv "bscp.io/pkg/protocol/core/app-template-variable"
	pbbase "bscp.io/pkg/protocol/core/base"
	pbcommit "bscp.io/pkg/protocol/core/commit"
	pbci "bscp.io/pkg/protocol/core/config-item"
	pbcontent "bscp.io/pkg/protocol/core/content"
	pbtv "bscp.io/pkg/protocol/core/template-variable"
	pbds "bscp.io/pkg/protocol/data-service"
	"bscp.io/pkg/tools"
	"bscp.io/pkg/types"
)

// ExtractAppTmplVariables extract app template variables.
// the variables come from template and non-template config items
func (s *Service) ExtractAppTmplVariables(ctx context.Context, req *pbds.ExtractAppTmplVariablesReq) (
	*pbds.ExtractAppTmplVariablesResp, error) {
	kt := kit.FromGrpcContext(ctx)

	tmplRevisions, cis, err := s.getAllAppCIs(kt)
	if err != nil {
		logs.Errorf("get all app config items failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	contents, err := s.downloadTmplContent(kt, tmplRevisions)
	if err != nil {
		logs.Errorf("download template content failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}
	ciContents, err := s.downloadCIContent(kt, cis)
	if err != nil {
		logs.Errorf("download config item content failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}
	contents = append(contents, ciContents...)

	// merge all template content
	allContent := bytes.Join(contents, []byte(" "))
	// extract all template variables
	variables := s.tmplProc.ExtractVariables(allContent)

	return &pbds.ExtractAppTmplVariablesResp{
		Details: variables,
	}, nil
}

// GetAppTmplVariableRefs get app template variable references.
// the variables come from template and non-template config items
func (s *Service) GetAppTmplVariableRefs(ctx context.Context, req *pbds.GetAppTmplVariableRefsReq) (
	*pbds.GetAppTmplVariableRefsResp, error) {
	kt := kit.FromGrpcContext(ctx)

	tmplRevisions, cis, err := s.getAllAppCIs(kt)
	if err != nil {
		logs.Errorf("get all app config items failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	refs, err := s.getVariableReferences(kt, tmplRevisions, cis)
	if err != nil {
		logs.Errorf("get variable references failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	return &pbds.GetAppTmplVariableRefsResp{
		Details: refs,
	}, nil
}

// GetReleasedAppTmplVariableRefs get released app template variable references.
// the variables come from template and non-template config items
func (s *Service) GetReleasedAppTmplVariableRefs(ctx context.Context,
	req *pbds.GetReleasedAppTmplVariableRefsReq) (
	*pbds.GetReleasedAppTmplVariableRefsResp, error) {
	kt := kit.FromGrpcContext(ctx)

	releasedTmpls, _, err := s.dao.ReleasedAppTemplate().List(kt, req.BizId, req.AppId, req.ReleaseId, nil,
		&types.BasePage{All: true})
	if err != nil {
		logs.Errorf("list released app templates failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}
	tmplRevisions := getTmplRevisionsFromReleased(releasedTmpls)
	tmplRevisions = filterSizeForTmplRevisions(tmplRevisions)

	releasedCIs, _, err := s.dao.ReleasedCI().List(kt, req.BizId, req.AppId, req.ReleaseId, nil,
		&types.BasePage{All: true})
	if err != nil {
		logs.Errorf("list released config items failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}
	cis := getPbConfigItemsFromReleased(releasedCIs)
	cis = filterSizeForConfigItems(cis)

	refs, err := s.getVariableReferences(kt, tmplRevisions, cis)
	if err != nil {
		logs.Errorf("get variable references failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	return &pbds.GetReleasedAppTmplVariableRefsResp{
		Details: refs,
	}, nil
}

func getTmplRevisionsFromReleased(releasedTmpls []*table.ReleasedAppTemplate) []*table.TemplateRevision {
	tmplRevisions := make([]*table.TemplateRevision, len(releasedTmpls))
	for idx, r := range releasedTmpls {
		tmplRevisions[idx] = &table.TemplateRevision{
			ID: r.ID,
			Spec: &table.TemplateRevisionSpec{
				RevisionName: r.Spec.TemplateRevisionName,
				RevisionMemo: r.Spec.TemplateRevisionMemo,
				Name:         r.Spec.Name,
				Path:         r.Spec.Path,
				FileType:     table.FileFormat(r.Spec.FileType),
				FileMode:     table.FileMode(r.Spec.FileMode),
				Permission: &table.FilePermission{
					User:      r.Spec.User,
					UserGroup: r.Spec.UserGroup,
					Privilege: r.Spec.Privilege,
				},
				ContentSpec: &table.ContentSpec{
					Signature: r.Spec.OriginSignature,
					ByteSize:  r.Spec.OriginByteSize,
				},
			},
			Attachment: &table.TemplateRevisionAttachment{
				BizID:           r.Attachment.BizID,
				TemplateSpaceID: r.Spec.TemplateSpaceID,
				TemplateID:      r.Spec.TemplateID,
			},
			Revision: &table.CreatedRevision{
				Creator:   r.Revision.Creator,
				CreatedAt: r.Revision.CreatedAt,
			},
		}
	}

	return tmplRevisions
}

func getPbConfigItemsFromReleased(releasedCIs []*table.ReleasedConfigItem) []*pbci.ConfigItem {
	cis := make([]*pbci.ConfigItem, len(releasedCIs))
	for idx, r := range releasedCIs {
		cis[idx] = &pbci.ConfigItem{
			Id: r.ID,
			Spec: &pbci.ConfigItemSpec{
				Name:     r.ConfigItemSpec.Name,
				Path:     r.ConfigItemSpec.Path,
				FileType: string(r.ConfigItemSpec.FileType),
				FileMode: string(r.ConfigItemSpec.FileMode),
				Permission: &pbci.FilePermission{
					User:      r.ConfigItemSpec.Permission.User,
					UserGroup: r.ConfigItemSpec.Permission.UserGroup,
					Privilege: r.ConfigItemSpec.Permission.Privilege,
				},
			},
			CommitSpec: &pbcommit.CommitSpec{
				ContentId: r.CommitSpec.ContentID,
				Content: &pbcontent.ContentSpec{
					Signature: r.CommitSpec.Content.OriginSignature,
					ByteSize:  r.CommitSpec.Content.OriginByteSize,
				},
				Memo: r.CommitSpec.Memo,
			},
			Attachment: &pbci.ConfigItemAttachment{
				BizId: r.Attachment.BizID,
				AppId: r.Attachment.AppID,
			},
			Revision: &pbbase.Revision{
				Creator:  r.Revision.Creator,
				Reviser:  r.Revision.Creator,
				CreateAt: r.Revision.CreatedAt.Format(constant.TimeStdFormat),
				UpdateAt: r.Revision.CreatedAt.Format(constant.TimeStdFormat),
			},
		}
	}

	return cis
}

// GetAppTmplVariableRefs get app template variable references.
func (s *Service) getVariableReferences(kt *kit.Kit, tmplRevisions []*table.TemplateRevision, cis []*pbci.ConfigItem) (
	[]*pbatv.AppTemplateVariableReference, error) {
	contents, err := s.downloadTmplContent(kt, tmplRevisions)
	if err != nil {
		logs.Errorf("download template content failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}
	ciContents, err := s.downloadCIContent(kt, cis)
	if err != nil {
		logs.Errorf("download config item content failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	allVariables := make([]string, 0)
	revisionVariableMap := make(map[uint32]map[string]struct{}, len(tmplRevisions))
	revisionMap := make(map[uint32]*table.TemplateRevision, len(tmplRevisions))
	for idx, r := range tmplRevisions {
		// extract template variables for one template config item
		variables := s.tmplProc.ExtractVariables(contents[idx])
		allVariables = append(allVariables, variables...)
		revisionMap[r.ID] = r

		revisionVariableMap[r.ID] = map[string]struct{}{}
		for _, v := range variables {
			revisionVariableMap[r.ID][v] = struct{}{}
		}
	}

	ciVariableMap := make(map[uint32]map[string]struct{}, len(cis))
	ciMap := make(map[uint32]*pbci.ConfigItem, len(cis))
	for idx, ci := range cis {
		// extract config item variables for one config item
		variables := s.tmplProc.ExtractVariables(ciContents[idx])
		allVariables = append(allVariables, variables...)
		ciMap[ci.Id] = ci

		ciVariableMap[ci.Id] = map[string]struct{}{}
		for _, v := range variables {
			ciVariableMap[ci.Id][v] = struct{}{}
		}
	}

	allVariables = tools.RemoveDuplicateStrings(allVariables)
	refs := make([]*pbatv.AppTemplateVariableReference, len(allVariables))
	for idx, v := range allVariables {
		ref := &pbatv.AppTemplateVariableReference{
			VariableName: v,
		}
		for rID, variables := range revisionVariableMap {
			if _, ok := variables[v]; ok {
				ref.References = append(ref.References, &pbatv.AppTemplateVariableReferenceReference{
					Id:                 revisionMap[rID].Attachment.TemplateID,
					TemplateRevisionId: rID,
					Name:               revisionMap[rID].Spec.Name,
					Path:               revisionMap[rID].Spec.Path,
				})
			}
		}
		for cID, variables := range ciVariableMap {
			if _, ok := variables[v]; ok {
				ref.References = append(ref.References, &pbatv.AppTemplateVariableReferenceReference{
					Id:                 ciMap[cID].Id,
					TemplateRevisionId: 0,
					Name:               ciMap[cID].Spec.Name,
					Path:               ciMap[cID].Spec.Path,
				})
			}
		}
		refs[idx] = ref
	}

	return refs, nil
}

// ListAppTmplVariables list app template variables.
func (s *Service) ListAppTmplVariables(ctx context.Context, req *pbds.ListAppTmplVariablesReq) (
	*pbds.ListAppTmplVariablesResp, error) {
	kt := kit.FromGrpcContext(ctx)

	// extract all variables for current app
	extractRep, err := s.ExtractAppTmplVariables(ctx, &pbds.ExtractAppTmplVariablesReq{
		BizId: req.BizId,
		AppId: req.AppId,
	})
	if err != nil {
		logs.Errorf("extract app template variables failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}
	allVariables := extractRep.Details
	if len(allVariables) == 0 {
		return &pbds.ListAppTmplVariablesResp{
			Details: []*pbtv.TemplateVariableSpec{},
		}, nil
	}

	// get app template variables
	appVars, err := s.dao.AppTemplateVariable().ListVariables(kt, req.BizId, req.AppId)
	if err != nil {
		logs.Errorf("list app template variables failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}
	appVarMap := make(map[string]*table.TemplateVariableSpec, len(appVars))
	for _, v := range appVars {
		appVarMap[v.Name] = v
	}

	// get biz template variables
	bizVars, _, err := s.dao.TemplateVariable().List(kt, req.BizId, nil, &types.BasePage{All: true})
	if err != nil {
		logs.Errorf("list template variables failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}
	bizVarMap := make(map[string]*table.TemplateVariableSpec, len(bizVars))
	for _, v := range bizVars {
		bizVarMap[v.Spec.Name] = v.Spec
	}

	// get final app template variables
	// use app variables first, then use biz variables
	finalVar := make([]*pbtv.TemplateVariableSpec, 0)
	for _, name := range allVariables {
		if v, ok := appVarMap[name]; ok {
			finalVar = append(finalVar, pbtv.PbTemplateVariableSpec(v))
			continue
		}
		if v, ok := bizVarMap[name]; ok {
			finalVar = append(finalVar, pbtv.PbTemplateVariableSpec(v))
			continue
		}
		// for unset variable, just return its name, other fields keep empty
		finalVar = append(finalVar, &pbtv.TemplateVariableSpec{Name: name})
	}

	return &pbds.ListAppTmplVariablesResp{
		Details: finalVar,
	}, nil
}

// ListReleasedAppTmplVariables get app template variable references.
func (s *Service) ListReleasedAppTmplVariables(ctx context.Context, req *pbds.ListReleasedAppTmplVariablesReq) (
	*pbds.ListReleasedAppTmplVariablesResp, error) {
	kt := kit.FromGrpcContext(ctx)

	details, err := s.dao.ReleasedAppTemplateVariable().ListVariables(kt, req.BizId, req.AppId, req.ReleaseId)
	if err != nil {
		logs.Errorf("list released app template variables failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	return &pbds.ListReleasedAppTmplVariablesResp{
		Details: pbtv.PbTemplateVariableSpecs(details),
	}, nil
}

// UpdateAppTmplVariables update app template variables.
func (s *Service) UpdateAppTmplVariables(ctx context.Context, req *pbds.UpdateAppTmplVariablesReq) (
	*pbbase.EmptyResp, error) {
	kt := kit.FromGrpcContext(ctx)
	// set for empty slice to ensure the data in db is not `null` but `[]`
	if len(req.Spec.Variables) == 0 {
		req.Spec.Variables = []*pbtv.TemplateVariableSpec{}
	}

	appVar := &table.AppTemplateVariable{
		Spec:       req.Spec.AppTemplateVariableSpec(),
		Attachment: req.Attachment.AppTemplateVariableAttachment(),
		Revision: &table.Revision{
			Creator: kt.User,
			Reviser: kt.User,
		},
	}
	if err := s.dao.AppTemplateVariable().Upsert(kt, appVar); err != nil {
		logs.Errorf("update app template variables failed, err: %v, rid: %s", err, kt.Rid)
		return nil, err
	}

	return new(pbbase.EmptyResp), nil
}
