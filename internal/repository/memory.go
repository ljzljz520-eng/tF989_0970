package repository

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"sync"

	"followuplabel/internal/model"
)

var ErrNotFound = errors.New("记录不存在")

type Memory struct {
	mu     sync.RWMutex
	tags   map[string]model.Tag
	menus  map[string]model.Menu
	nextID int
}

func NewMemory(tags []model.Tag) *Memory {
	store := &Memory{
		tags:   make(map[string]model.Tag, len(tags)),
		menus:  make(map[string]model.Menu),
		nextID: len(tags) + 1,
	}
	for _, tag := range tags {
		store.tags[tag.ID] = cloneTag(tag)
	}
	return store
}

func FixtureTags() []model.Tag {
	return []model.Tag{
		{ID: "tag-001", Name: "需要二次沟通", Color: "#D97706", ApplicableScenes: []string{"售前咨询", "报价跟进"}, SortOrder: 10},
		{ID: "tag-002", Name: "问题已解决", Color: "#16855B", ApplicableScenes: []string{"售后回访"}, SortOrder: 20},
		{ID: "tag-003", Name: "高价值客户", Color: "#B42355", ApplicableScenes: []string{"客户关怀", "续约提醒"}, SortOrder: 30},
	}
}

func (m *Memory) ListTags(ctx context.Context) ([]model.Tag, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	tags := make([]model.Tag, 0, len(m.tags))
	for _, tag := range m.tags {
		tags = append(tags, cloneTag(tag))
	}
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].SortOrder != tags[j].SortOrder {
			return tags[i].SortOrder < tags[j].SortOrder
		}
		if tags[i].Name != tags[j].Name {
			return tags[i].Name < tags[j].Name
		}
		return tags[i].ID < tags[j].ID
	})
	return tags, nil
}

func (m *Memory) CreateTag(ctx context.Context, input model.TagInput) (model.Tag, error) {
	if err := ctx.Err(); err != nil {
		return model.Tag{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("tag-%03d", m.nextID)
	m.nextID++
	tag := model.Tag{
		ID:               id,
		Name:             input.Name,
		Color:            input.Color,
		ApplicableScenes: slices.Clone(input.ApplicableScenes),
		SortOrder:        input.SortOrder,
	}
	m.tags[id] = tag
	return cloneTag(tag), nil
}

func (m *Memory) UpdateTag(ctx context.Context, id string, input model.TagInput) (model.Tag, error) {
	if err := ctx.Err(); err != nil {
		return model.Tag{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.tags[id]; !exists {
		return model.Tag{}, ErrNotFound
	}
	tag := model.Tag{
		ID:               id,
		Name:             input.Name,
		Color:            input.Color,
		ApplicableScenes: slices.Clone(input.ApplicableScenes),
		SortOrder:        input.SortOrder,
	}
	m.tags[id] = tag
	return cloneTag(tag), nil
}

func (m *Memory) DeleteTag(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.tags[id]; !exists {
		return ErrNotFound
	}
	delete(m.tags, id)
	return nil
}

func (m *Memory) DeleteAllTags(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.tags)
	return nil
}

func (m *Memory) PutMenu(ctx context.Context, menu model.Menu) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.menus[menu.Key] = menu
	return nil
}

func (m *Memory) RemoveMenu(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.menus, key)
	return nil
}

func (m *Memory) ListMenus(ctx context.Context) ([]model.Menu, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	menus := make([]model.Menu, 0, len(m.menus))
	for _, menu := range m.menus {
		menus = append(menus, menu)
	}
	sort.Slice(menus, func(i, j int) bool {
		if menus[i].SortOrder != menus[j].SortOrder {
			return menus[i].SortOrder < menus[j].SortOrder
		}
		return menus[i].Key < menus[j].Key
	})
	return menus, nil
}

func cloneTag(tag model.Tag) model.Tag {
	tag.ApplicableScenes = slices.Clone(tag.ApplicableScenes)
	return tag
}
