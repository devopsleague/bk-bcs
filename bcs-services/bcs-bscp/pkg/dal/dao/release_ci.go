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

package dao

import (
	"errors"
	"fmt"

	rawgen "gorm.io/gen"

	"bscp.io/pkg/criteria/errf"
	"bscp.io/pkg/dal/gen"
	"bscp.io/pkg/dal/table"
	"bscp.io/pkg/kit"
	"bscp.io/pkg/search"
	"bscp.io/pkg/types"
)

// ReleasedCI supplies all the released config item related operations.
type ReleasedCI interface {
	// BulkCreateWithTx bulk create released config items with tx.
	BulkCreateWithTx(kit *kit.Kit, tx *gen.QueryTx, items []*table.ReleasedConfigItem) error
	// Get released config item by id and released id
	Get(kit *kit.Kit, id, bizID, releasedID uint32) (*table.ReleasedConfigItem, error)
	// GetReleasedLately released config item by app id and biz id
	GetReleasedLately(kit *kit.Kit, appId, bizID uint32, searchKey string) ([]*table.ReleasedConfigItem, error)
	// List released config items with options.
	List(kit *kit.Kit, bizID, appID, releaseID uint32, s search.Searcher, opt *types.BasePage) (
		[]*table.ReleasedConfigItem, int64, error)
	// ListAll list all released config items in biz.
	ListAll(kit *kit.Kit, bizID uint32) ([]*table.ReleasedConfigItem, error)
	// ListAllByAppID list all released config items by appID.
	ListAllByAppID(kit *kit.Kit, appID, bizID uint32) ([]*table.ReleasedConfigItem, error)
	// ListAllByAppIDs batch list released config items by appIDs.
	ListAllByAppIDs(kit *kit.Kit, appIDs []uint32, bizID uint32) ([]*table.ReleasedConfigItem, error)
	// ListAllByReleaseIDs batch list released config items by releaseIDs.
	ListAllByReleaseIDs(kit *kit.Kit, releasedIDs []uint32, bizID uint32) ([]*table.ReleasedConfigItem, error)
}

var _ ReleasedCI = new(releasedCIDao)

type releasedCIDao struct {
	genQ     *gen.Query
	idGen    IDGenInterface
	auditDao AuditDao
}

// BulkCreateWithTx bulk create released config items.
func (dao *releasedCIDao) BulkCreateWithTx(kit *kit.Kit, tx *gen.QueryTx, items []*table.ReleasedConfigItem) error {
	if len(items) == 0 {
		return errors.New("released config items is empty")
	}

	// validate released config item field.
	for _, item := range items {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	// generate released config items id.
	ids, err := dao.idGen.Batch(kit, table.ReleasedConfigItemTable, len(items))
	if err != nil {
		return err
	}

	start := 0
	for _, item := range items {
		item.ID = ids[start]
		start++
	}
	batchSize := 100

	q := tx.ReleasedConfigItem.WithContext(kit.Ctx)
	if err := q.CreateInBatches(items, batchSize); err != nil {
		return fmt.Errorf("create released config items in batch failed, err: %v", err)
	}

	ad := dao.auditDao.DecoratorV2(kit, items[0].Attachment.BizID).PrepareCreate(table.RciList(items))
	if err := ad.Do(tx.Query); err != nil {
		return err
	}

	return nil
}

// Get released config item by ID and config item id and release id.
func (dao *releasedCIDao) Get(kit *kit.Kit, configItemID, bizID, releaseID uint32) (*table.ReleasedConfigItem, error) {

	if configItemID == 0 {
		return nil, errf.New(errf.InvalidParameter, "config item id can not be 0")
	}

	if releaseID == 0 {
		return nil, errf.New(errf.InvalidParameter, "release id can not be 0")
	}

	m := dao.genQ.ReleasedConfigItem
	return m.WithContext(kit.Ctx).Where(
		m.ConfigItemID.Eq(configItemID), m.ReleaseID.Eq(releaseID), m.BizID.Eq(bizID)).Take()
}

// List released config items with options.
func (dao *releasedCIDao) List(kit *kit.Kit, bizID, appID, releaseID uint32, s search.Searcher, opt *types.BasePage) (
	[]*table.ReleasedConfigItem, int64, error) {
	m := dao.genQ.ReleasedConfigItem
	q := dao.genQ.ReleasedConfigItem.WithContext(kit.Ctx)

	var conds []rawgen.Condition
	// add search condition
	if s != nil {
		exprs := s.SearchExprs(dao.genQ)
		if len(exprs) > 0 {
			var do gen.IReleasedConfigItemDo
			for i := range exprs {
				if i == 0 {
					do = q.Where(exprs[i])
				}
				do = do.Or(exprs[i])
			}
			conds = append(conds, do)
		}
	}

	d := q.Where(m.BizID.Eq(bizID), m.AppID.Eq(appID), m.ReleaseID.Eq(releaseID), m.ConfigItemID.Neq(0)).
		Where(conds...)
	if opt.All {
		result, err := d.Find()
		if err != nil {
			return nil, 0, err
		}
		return result, int64(len(result)), err
	}

	return d.FindByPage(opt.Offset(), opt.LimitInt())
}

// ListAll list all released config items in biz.
func (dao *releasedCIDao) ListAll(kit *kit.Kit, bizID uint32) ([]*table.ReleasedConfigItem, error) {
	if bizID == 0 {
		return nil, errf.New(errf.InvalidParameter, "biz_id can not be 0")
	}

	m := dao.genQ.ReleasedConfigItem
	return m.WithContext(kit.Ctx).Where(m.BizID.Eq(bizID)).Find()
}

// ListAll list all released config items in biz.
func (dao *releasedCIDao) ListAllByAppID(kit *kit.Kit, appID, bizID uint32) ([]*table.ReleasedConfigItem, error) {
	if bizID == 0 {
		return nil, errf.New(errf.InvalidParameter, "biz_id can not be 0")
	}

	m := dao.genQ.ReleasedConfigItem
	return m.WithContext(kit.Ctx).Where(m.AppID.Eq(appID), m.BizID.Eq(bizID)).Find()
}

// ListAllByAppIDs list all released config items by appIDs.
func (dao *releasedCIDao) ListAllByAppIDs(kit *kit.Kit,
	appIDs []uint32, bizID uint32) ([]*table.ReleasedConfigItem, error) {
	if bizID == 0 {
		return nil, errf.New(errf.InvalidParameter, "biz_id can not be 0")
	}

	m := dao.genQ.ReleasedConfigItem
	return m.WithContext(kit.Ctx).Where(m.AppID.In(appIDs...), m.BizID.Eq(bizID)).Find()
}

// ListAllByReleaseIDs list all released config items by releaseIDs.
func (dao *releasedCIDao) ListAllByReleaseIDs(kit *kit.Kit,
	releaseIDs []uint32, bizID uint32) ([]*table.ReleasedConfigItem, error) {
	if bizID == 0 {
		return nil, errf.New(errf.InvalidParameter, "biz_id can not be 0")
	}

	m := dao.genQ.ReleasedConfigItem
	return m.WithContext(kit.Ctx).Where(m.ReleaseID.In(releaseIDs...), m.BizID.Eq(bizID)).Find()
}

// GetReleasedLately
func (dao *releasedCIDao) GetReleasedLately(kit *kit.Kit, appId, bizID uint32, searchKey string) (
	[]*table.ReleasedConfigItem, error) {
	if bizID == 0 {
		return nil, errf.New(errf.InvalidParameter, "biz_id can not be 0")
	}

	m := dao.genQ.ReleasedConfigItem
	q := dao.genQ.ReleasedConfigItem.WithContext(kit.Ctx)
	// m.ConfigItemID.Neq(0) means not to match template config items
	query := q.Where(m.BizID.Eq(bizID), m.AppID.Eq(appId), m.ConfigItemID.Neq(0))
	if searchKey != "" {
		param := "%" + searchKey + "%"
		query = q.Where(query, q.Where(m.Name.Like(param)).Or(m.Creator.Like(param)).Or(m.Reviser.Like(param)))
	}
	subQuery := q.Where(m.BizID.Eq(bizID), m.AppID.Eq(appId)).Order(m.ReleaseID.Desc()).Limit(1).Select(m.ReleaseID)
	return query.Where(q.Columns(m.ReleaseID).Eq(subQuery)).Find()
}
