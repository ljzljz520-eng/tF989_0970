package controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"followuplabel/internal/controller"
	"followuplabel/internal/model"
	"followuplabel/internal/repository"
	"followuplabel/internal/service"
)

func TestAdminAPIExposesInstalledMenuAndMaintainsTags(t *testing.T) {
	store := repository.NewMemory(repository.FixtureTags())
	plugin := service.NewPlugin(store, store, nil)
	if err := plugin.Install(context.Background()); err != nil {
		t.Fatalf("安装返回错误: %v", err)
	}
	handler := controller.NewHTTP(service.NewTags(store, store))

	createRequest := httptest.NewRequest(http.MethodPost, "/api/admin/tags", strings.NewReader(`{"name":"待续约","color":"#2368a2","applicableScenes":["续约提醒"],"sortOrder":15}`))
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("创建响应状态为 %d，内容为 %s", createResponse.Code, createResponse.Body.String())
	}

	stateRequest := httptest.NewRequest(http.MethodGet, "/api/admin/state", nil)
	stateResponse := httptest.NewRecorder()
	handler.ServeHTTP(stateResponse, stateRequest)
	if stateResponse.Code != http.StatusOK {
		t.Fatalf("后台状态响应为 %d", stateResponse.Code)
	}
	var state model.AdminState
	if err := json.NewDecoder(stateResponse.Body).Decode(&state); err != nil {
		t.Fatalf("后台状态内容无法读取: %v", err)
	}
	if len(state.Tags) != 4 {
		t.Fatalf("标签数量为 %d", len(state.Tags))
	}
	if len(state.Menus) != 1 || state.Menus[0].Key != service.MenuKey {
		t.Fatalf("菜单内容为 %+v", state.Menus)
	}
	if state.Tags[1].Name != "待续约" || state.Tags[1].Color != "#2368A2" {
		t.Fatalf("新增标签内容为 %+v", state.Tags[1])
	}
}

func TestReadonlyPageAndPublicTagsAreAvailable(t *testing.T) {
	store := repository.NewMemory(repository.FixtureTags())
	handler := controller.NewHTTP(service.NewTags(store, store))

	pageRequest := httptest.NewRequest(http.MethodGet, "/tags", nil)
	pageResponse := httptest.NewRecorder()
	handler.ServeHTTP(pageResponse, pageRequest)
	if pageResponse.Code != http.StatusOK {
		t.Fatalf("说明页响应为 %d", pageResponse.Code)
	}
	if !strings.Contains(pageResponse.Body.String(), "标签说明") {
		t.Fatalf("说明页内容为 %q", pageResponse.Body.String())
	}

	apiRequest := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	apiResponse := httptest.NewRecorder()
	handler.ServeHTTP(apiResponse, apiRequest)
	if apiResponse.Code != http.StatusOK {
		t.Fatalf("公开标签响应为 %d", apiResponse.Code)
	}
	var tags []model.Tag
	if err := json.NewDecoder(apiResponse.Body).Decode(&tags); err != nil {
		t.Fatalf("公开标签内容无法读取: %v", err)
	}
	if len(tags) != len(repository.FixtureTags()) {
		t.Fatalf("公开标签数量为 %d", len(tags))
	}
}
