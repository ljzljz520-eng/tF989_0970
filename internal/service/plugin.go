package service

import (
	"context"

	"followuplabel/internal/model"
)

const MenuKey = "customer-followup-labels"

type TagLifecycleStore interface {
	DeleteAllTags(context.Context) error
}

type MenuStore interface {
	PutMenu(context.Context, model.Menu) error
	RemoveMenu(context.Context, string) error
}

type RetainedDataCleanup func(context.Context) error

type Plugin struct {
	tags                TagLifecycleStore
	menus               MenuStore
	retainedDataCleanup RetainedDataCleanup
}

func NewPlugin(tags TagLifecycleStore, menus MenuStore, cleanup RetainedDataCleanup) *Plugin {
	return &Plugin{tags: tags, menus: menus, retainedDataCleanup: cleanup}
}

func (p *Plugin) Install(ctx context.Context) error {
	return p.menus.PutMenu(ctx, model.Menu{
		Key:       MenuKey,
		Title:     "客户回访标签",
		Path:      "/admin/tags",
		SortOrder: 60,
	})
}

func (p *Plugin) Uninstall(ctx context.Context, retainBusinessData bool) error {
	if err := p.menus.RemoveMenu(ctx, MenuKey); err != nil {
		return err
	}
	if retainBusinessData {
		// 未配置清理回调时，保留业务数据卸载仅移除插件菜单，标签数据完整保留，
		// 不应调用空回调而引发运行时崩溃。
		if p.retainedDataCleanup == nil {
			return nil
		}
		return p.retainedDataCleanup(ctx)
	}
	return p.tags.DeleteAllTags(ctx)
}
