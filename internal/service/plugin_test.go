package service_test

import (
	"context"
	"testing"

	"followuplabel/internal/repository"
	"followuplabel/internal/service"
)

func TestInstallAddsCustomerFollowupMenu(t *testing.T) {
	store := repository.NewMemory(repository.FixtureTags())
	plugin := service.NewPlugin(store, store, nil)

	if err := plugin.Install(context.Background()); err != nil {
		t.Fatalf("安装返回错误: %v", err)
	}
	menus, err := store.ListMenus(context.Background())
	if err != nil {
		t.Fatalf("读取菜单返回错误: %v", err)
	}
	if len(menus) != 1 {
		t.Fatalf("菜单数量为 %d", len(menus))
	}
	if menus[0].Key != service.MenuKey || menus[0].Path != "/admin/tags" {
		t.Fatalf("菜单内容为 %+v", menus[0])
	}
}

func TestUninstallWithoutRetentionRemovesMenuAndTags(t *testing.T) {
	store := repository.NewMemory(repository.FixtureTags())
	plugin := service.NewPlugin(store, store, nil)
	if err := plugin.Install(context.Background()); err != nil {
		t.Fatalf("安装返回错误: %v", err)
	}

	if err := plugin.Uninstall(context.Background(), false); err != nil {
		t.Fatalf("卸载返回错误: %v", err)
	}
	menus, err := store.ListMenus(context.Background())
	if err != nil {
		t.Fatalf("读取菜单返回错误: %v", err)
	}
	tags, err := store.ListTags(context.Background())
	if err != nil {
		t.Fatalf("读取标签返回错误: %v", err)
	}
	if len(menus) != 0 {
		t.Fatalf("菜单数量为 %d", len(menus))
	}
	if len(tags) != 0 {
		t.Fatalf("标签数量为 %d", len(tags))
	}
}

func TestUninstallWithRetentionRemovesMenuAndPreservesTags(t *testing.T) {
	store := repository.NewMemory(repository.FixtureTags())
	plugin := service.NewPlugin(store, store, nil)
	if err := plugin.Install(context.Background()); err != nil {
		t.Fatalf("安装返回错误: %v", err)
	}

	if err := plugin.Uninstall(context.Background(), true); err != nil {
		t.Fatalf("卸载返回错误: %v", err)
	}
	menus, err := store.ListMenus(context.Background())
	if err != nil {
		t.Fatalf("读取菜单返回错误: %v", err)
	}
	tags, err := store.ListTags(context.Background())
	if err != nil {
		t.Fatalf("读取标签返回错误: %v", err)
	}
	if len(menus) != 0 {
		t.Fatalf("菜单数量为 %d", len(menus))
	}
	if len(tags) != len(repository.FixtureTags()) {
		t.Fatalf("标签数量为 %d", len(tags))
	}
}
