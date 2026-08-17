package service_test

import (
	"context"
	"testing"

	"followuplabel/internal/model"
	"followuplabel/internal/repository"
	"followuplabel/internal/service"
)

func TestCreatedTagIsNormalizedAndOrdered(t *testing.T) {
	store := repository.NewMemory(repository.FixtureTags())
	tags := service.NewTags(store, store)

	created, err := tags.Create(context.Background(), model.TagInput{
		Name:             "  等待客户确认  ",
		Color:            "#1a6fba",
		ApplicableScenes: []string{"  报价跟进 ", "报价跟进", ""},
		SortOrder:        5,
	})
	if err != nil {
		t.Fatalf("创建标签返回错误: %v", err)
	}
	if created.ID != "tag-004" {
		t.Fatalf("标签编号为 %q", created.ID)
	}
	if created.Name != "等待客户确认" {
		t.Fatalf("标签名称为 %q", created.Name)
	}
	if created.Color != "#1A6FBA" {
		t.Fatalf("标签颜色为 %q", created.Color)
	}
	if len(created.ApplicableScenes) != 1 || created.ApplicableScenes[0] != "报价跟进" {
		t.Fatalf("适用场景为 %v", created.ApplicableScenes)
	}

	listed, err := tags.List(context.Background())
	if err != nil {
		t.Fatalf("读取标签返回错误: %v", err)
	}
	if listed[0].ID != created.ID {
		t.Fatalf("首个标签编号为 %q", listed[0].ID)
	}
}

func TestInvalidTagIsRejectedWithFieldResults(t *testing.T) {
	store := repository.NewMemory(nil)
	tags := service.NewTags(store, store)

	_, err := tags.Create(context.Background(), model.TagInput{Color: "green", SortOrder: -1})
	if err == nil {
		t.Fatal("无效标签被创建")
	}

	listed, listErr := tags.List(context.Background())
	if listErr != nil {
		t.Fatalf("读取标签返回错误: %v", listErr)
	}
	if len(listed) != 0 {
		t.Fatalf("标签数量为 %d", len(listed))
	}
}
