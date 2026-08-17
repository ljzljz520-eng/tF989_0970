package service

import (
	"context"

	"golang.org/x/sync/errgroup"

	"followuplabel/internal/model"
	"followuplabel/internal/validation"
)

type TagStore interface {
	ListTags(context.Context) ([]model.Tag, error)
	CreateTag(context.Context, model.TagInput) (model.Tag, error)
	UpdateTag(context.Context, string, model.TagInput) (model.Tag, error)
	DeleteTag(context.Context, string) error
}

type MenuReader interface {
	ListMenus(context.Context) ([]model.Menu, error)
}

type Tags struct {
	tags  TagStore
	menus MenuReader
}

func NewTags(tags TagStore, menus MenuReader) *Tags {
	return &Tags{tags: tags, menus: menus}
}

func (s *Tags) List(ctx context.Context) ([]model.Tag, error) {
	return s.tags.ListTags(ctx)
}

func (s *Tags) Create(ctx context.Context, input model.TagInput) (model.Tag, error) {
	normalized, err := validation.NormalizeTag(input)
	if err != nil {
		return model.Tag{}, err
	}
	return s.tags.CreateTag(ctx, normalized)
}

func (s *Tags) Update(ctx context.Context, id string, input model.TagInput) (model.Tag, error) {
	normalized, err := validation.NormalizeTag(input)
	if err != nil {
		return model.Tag{}, err
	}
	return s.tags.UpdateTag(ctx, id, normalized)
}

func (s *Tags) Delete(ctx context.Context, id string) error {
	return s.tags.DeleteTag(ctx, id)
}

func (s *Tags) AdminState(ctx context.Context) (model.AdminState, error) {
	var state model.AdminState
	group, groupContext := errgroup.WithContext(ctx)
	group.Go(func() error {
		tags, err := s.tags.ListTags(groupContext)
		state.Tags = tags
		return err
	})
	group.Go(func() error {
		menus, err := s.menus.ListMenus(groupContext)
		state.Menus = menus
		return err
	})
	if err := group.Wait(); err != nil {
		return model.AdminState{}, err
	}
	return state, nil
}
