package content

import "time"

// PersonalityType 人格类型结构
type PersonalityType struct {
	Code     string `yaml:"code" json:"code"`
	Gender   string `yaml:"gender" json:"gender"`
	Name     string `yaml:"name" json:"name"`
	Tagline  string `yaml:"tagline" json:"tagline"`
	Summary  string `yaml:"summary" json:"summary"`
	Body     string `yaml:"-" json:"body"`
	Lang     string `yaml:"lang" json:"lang"`
	FilePath string `yaml:"-" json:"-"`
	ModTime  time.Time `yaml:"-" json:"-"`
}

// PersonalityPairing 人格配对结构
type PersonalityPairing struct {
	ID                string `yaml:"id" json:"id"`
	MaleCode          string `yaml:"male_code" json:"male_code"`
	FemaleCode        string `yaml:"female_code" json:"female_code"`
	CompatibilityScore int    `yaml:"compatibility_score" json:"compatibility_score"`
	Body              string `yaml:"-" json:"body"`
	Lang              string `yaml:"lang" json:"lang"`
	FilePath          string `yaml:"-" json:"-"`
	ModTime           time.Time `yaml:"-" json:"-"`
}

// FrontMatter 通用的front-matter结构
type FrontMatter struct {
	Lang               string `yaml:"lang"`
	Code               string `yaml:"code,omitempty"`
	Gender             string `yaml:"gender,omitempty"`
	Name               string `yaml:"name,omitempty"`
	Tagline            string `yaml:"tagline,omitempty"`
	Summary            string `yaml:"summary,omitempty"`
	ID                 string `yaml:"id,omitempty"`
	MaleCode           string `yaml:"male_code,omitempty"`
	FemaleCode         string `yaml:"female_code,omitempty"`
	CompatibilityScore int    `yaml:"compatibility_score,omitempty"`
}

// TypesSearchParams 人格类型搜索参数
type TypesSearchParams struct {
	Gender      string `form:"gender"`
	Search      string `form:"search"`
	Temperament string `form:"temperament"`
	Page        int    `form:"page,default=1"`
	Limit       int    `form:"limit,default=20"`
}

// PairingsSearchParams 配对搜索参数
type PairingsSearchParams struct {
	Male   string `form:"male"`
	Female string `form:"female"`
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=20"`
}

// APIResponse 通用API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse 分页响应结构
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}
