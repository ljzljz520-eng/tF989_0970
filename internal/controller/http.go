package controller

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"text/template"

	"followuplabel/internal/model"
	"followuplabel/internal/repository"
	"followuplabel/internal/service"
	"followuplabel/internal/validation"
	frontend "followuplabel/web"
)

type HTTP struct {
	tags       *service.Tags
	adminView  *template.Template
	labelsView *template.Template
}

func NewHTTP(tags *service.Tags) http.Handler {
	assets, err := fs.Sub(frontend.Files, "assets")
	if err != nil {
		panic(err)
	}
	controller := &HTTP{
		tags:       tags,
		adminView:  template.Must(template.ParseFS(frontend.Files, "views/admin.html")),
		labelsView: template.Must(template.ParseFS(frontend.Files, "views/labels.html")),
	}
	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets))))
	mux.HandleFunc("GET /admin/tags", controller.adminPage)
	mux.HandleFunc("GET /tags", controller.labelsPage)
	mux.HandleFunc("GET /api/tags", controller.listTags)
	mux.HandleFunc("GET /api/admin/state", controller.adminState)
	mux.HandleFunc("POST /api/admin/tags", controller.createTag)
	mux.HandleFunc("PUT /api/admin/tags/{id}", controller.updateTag)
	mux.HandleFunc("DELETE /api/admin/tags/{id}", controller.deleteTag)
	mux.HandleFunc("GET /", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/admin/tags", http.StatusFound)
	})
	return mux
}

func (c *HTTP) adminPage(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.adminView.Execute(response, nil); err != nil {
		http.Error(response, "页面渲染失败", http.StatusInternalServerError)
	}
}

func (c *HTTP) labelsPage(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.labelsView.Execute(response, nil); err != nil {
		http.Error(response, "页面渲染失败", http.StatusInternalServerError)
	}
}

func (c *HTTP) listTags(response http.ResponseWriter, request *http.Request) {
	tags, err := c.tags.List(request.Context())
	if err != nil {
		writeError(response, http.StatusInternalServerError, err)
		return
	}
	writeJSON(response, http.StatusOK, tags)
}

func (c *HTTP) adminState(response http.ResponseWriter, request *http.Request) {
	state, err := c.tags.AdminState(request.Context())
	if err != nil {
		writeError(response, http.StatusInternalServerError, err)
		return
	}
	writeJSON(response, http.StatusOK, state)
}

func (c *HTTP) createTag(response http.ResponseWriter, request *http.Request) {
	input, ok := decodeInput(response, request)
	if !ok {
		return
	}
	tag, err := c.tags.Create(request.Context(), input)
	if err != nil {
		writeServiceError(response, err)
		return
	}
	writeJSON(response, http.StatusCreated, tag)
}

func (c *HTTP) updateTag(response http.ResponseWriter, request *http.Request) {
	input, ok := decodeInput(response, request)
	if !ok {
		return
	}
	tag, err := c.tags.Update(request.Context(), request.PathValue("id"), input)
	if err != nil {
		writeServiceError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, tag)
}

func (c *HTTP) deleteTag(response http.ResponseWriter, request *http.Request) {
	if err := c.tags.Delete(request.Context(), request.PathValue("id")); err != nil {
		writeServiceError(response, err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func decodeInput(response http.ResponseWriter, request *http.Request) (model.TagInput, bool) {
	defer request.Body.Close()
	decoder := json.NewDecoder(http.MaxBytesReader(response, request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	var input model.TagInput
	if err := decoder.Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, errors.New("请求内容不是有效的标签数据"))
		return model.TagInput{}, false
	}
	return input, true
}

func writeServiceError(response http.ResponseWriter, err error) {
	var fieldErrors validation.Errors
	switch {
	case errors.As(err, &fieldErrors):
		writeJSON(response, http.StatusUnprocessableEntity, map[string]any{"error": "标签数据校验失败", "fields": fieldErrors})
	case errors.Is(err, repository.ErrNotFound):
		writeError(response, http.StatusNotFound, err)
	default:
		writeError(response, http.StatusInternalServerError, err)
	}
}

func writeError(response http.ResponseWriter, status int, err error) {
	writeJSON(response, status, map[string]string{"error": err.Error()})
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}
