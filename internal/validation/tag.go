package validation

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"followuplabel/internal/model"
)

var hexColor = regexp.MustCompile(`^#[0-9A-F]{6}$`)

type Errors map[string]string

func (e Errors) Error() string {
	keys := make([]string, 0, len(e))
	for key := range e {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s", key, e[key]))
	}
	return strings.Join(parts, "; ")
}

func NormalizeTag(input model.TagInput) (model.TagInput, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Color = strings.ToUpper(strings.TrimSpace(input.Color))

	seen := make(map[string]struct{}, len(input.ApplicableScenes))
	scenes := make([]string, 0, len(input.ApplicableScenes))
	for _, scene := range input.ApplicableScenes {
		scene = strings.TrimSpace(scene)
		if scene == "" {
			continue
		}
		if _, exists := seen[scene]; exists {
			continue
		}
		seen[scene] = struct{}{}
		scenes = append(scenes, scene)
	}
	input.ApplicableScenes = scenes

	errs := Errors{}
	if input.Name == "" {
		errs["name"] = "标签名称不能为空"
	} else if len([]rune(input.Name)) > 40 {
		errs["name"] = "标签名称不能超过40个字符"
	}
	if !hexColor.MatchString(input.Color) {
		errs["color"] = "颜色必须是六位十六进制值"
	}
	if len(input.ApplicableScenes) == 0 {
		errs["applicableScenes"] = "至少选择一个适用场景"
	}
	if input.SortOrder < 0 || input.SortOrder > 9999 {
		errs["sortOrder"] = "排序必须在0到9999之间"
	}
	if len(errs) > 0 {
		return model.TagInput{}, errs
	}
	return input, nil
}
