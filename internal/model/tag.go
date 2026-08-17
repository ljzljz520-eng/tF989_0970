package model

type Tag struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Color            string   `json:"color"`
	ApplicableScenes []string `json:"applicableScenes"`
	SortOrder        int      `json:"sortOrder"`
}

type TagInput struct {
	Name             string   `json:"name"`
	Color            string   `json:"color"`
	ApplicableScenes []string `json:"applicableScenes"`
	SortOrder        int      `json:"sortOrder"`
}

type Menu struct {
	Key       string `json:"key"`
	Title     string `json:"title"`
	Path      string `json:"path"`
	SortOrder int    `json:"sortOrder"`
}

type AdminState struct {
	Tags  []Tag  `json:"tags"`
	Menus []Menu `json:"menus"`
}
