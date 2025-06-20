package api

import (
	"math/rand"
	"net/http"
	"strings"

	"mbti-site/pkg/cache"
	"mbti-site/pkg/content"
	"mbti-site/pkg/i18n"

	"github.com/gin-gonic/gin"
)

// Handler API处理器
type Handler struct {
	cache *cache.Manager
	i18n  *i18n.Manager
}

// NewHandler 创建新的API处理器
func NewHandler(cacheManager *cache.Manager, i18nManager *i18n.Manager) *Handler {
	return &Handler{
		cache: cacheManager,
		i18n:  i18nManager,
	}
}

// GetPersonalityTypes 获取人格类型列表
func (h *Handler) GetPersonalityTypes(c *gin.Context) {
	lang := c.Param("lang")
	if !h.i18n.IsValidLanguage(lang) {
		c.JSON(http.StatusBadRequest, content.APIResponse{
			Success: false,
			Error:   "Invalid language",
		})
		return
	}

	// 解析查询参数
	var params content.TypesSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, content.APIResponse{
			Success: false,
			Error:   "Invalid query parameters",
		})
		return
	}

	// 获取所有人格类型
	types, err := h.cache.GetAllPersonalityTypes(lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, content.APIResponse{
			Success: false,
			Error:   "Failed to get personality types",
		})
		return
	}

	// 应用筛选
	filteredTypes := h.filterPersonalityTypes(types, &params)

	// 应用分页
	total := len(filteredTypes)
	start := (params.Page - 1) * params.Limit
	end := start + params.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedTypes := filteredTypes[start:end]

	c.JSON(http.StatusOK, content.PaginatedResponse{
		Success: true,
		Data:    paginatedTypes,
		Total:   total,
		Page:    params.Page,
		Limit:   params.Limit,
	})
}

// GetPersonalityType 获取单个人格类型
func (h *Handler) GetPersonalityType(c *gin.Context) {
	lang := c.Param("lang")
	code := strings.ToUpper(c.Param("code"))
	gender := c.Param("gender")

	if !h.i18n.IsValidLanguage(lang) {
		c.JSON(http.StatusBadRequest, content.APIResponse{
			Success: false,
			Error:   "Invalid language",
		})
		return
	}

	if gender != "m" && gender != "f" {
		c.JSON(http.StatusBadRequest, content.APIResponse{
			Success: false,
			Error:   "Invalid gender, must be 'm' or 'f'",
		})
		return
	}

	personalityType, err := h.cache.GetPersonalityType(lang, code, gender)
	if err != nil {
		c.JSON(http.StatusNotFound, content.APIResponse{
			Success: false,
			Error:   "Personality type not found",
		})
		return
	}

	c.JSON(http.StatusOK, content.APIResponse{
		Success: true,
		Data:    personalityType,
	})
}

// GetPersonalityPairings 获取人格配对列表
func (h *Handler) GetPersonalityPairings(c *gin.Context) {
	lang := c.Param("lang")
	if !h.i18n.IsValidLanguage(lang) {
		c.JSON(http.StatusBadRequest, content.APIResponse{
			Success: false,
			Error:   "Invalid language",
		})
		return
	}

	// 解析查询参数
	var params content.PairingsSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, content.APIResponse{
			Success: false,
			Error:   "Invalid query parameters",
		})
		return
	}

	// 获取所有配对
	pairings, err := h.cache.GetAllPersonalityPairings(lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, content.APIResponse{
			Success: false,
			Error:   "Failed to get personality pairings",
		})
		return
	}

	// 应用筛选
	filteredPairings := h.filterPersonalityPairings(pairings, &params)

	// 应用分页
	total := len(filteredPairings)
	start := (params.Page - 1) * params.Limit
	end := start + params.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedPairings := filteredPairings[start:end]

	c.JSON(http.StatusOK, content.PaginatedResponse{
		Success: true,
		Data:    paginatedPairings,
		Total:   total,
		Page:    params.Page,
		Limit:   params.Limit,
	})
}

// GetPersonalityPairing 获取单个人格配对
func (h *Handler) GetPersonalityPairing(c *gin.Context) {
	lang := c.Param("lang")
	male := strings.ToUpper(c.Param("male"))
	female := strings.ToUpper(c.Param("female"))

	if !h.i18n.IsValidLanguage(lang) {
		c.JSON(http.StatusBadRequest, content.APIResponse{
			Success: false,
			Error:   "Invalid language",
		})
		return
	}

	pairing, err := h.cache.GetPersonalityPairing(lang, male, female)
	if err != nil {
		c.JSON(http.StatusNotFound, content.APIResponse{
			Success: false,
			Error:   "Personality pairing not found",
		})
		return
	}

	c.JSON(http.StatusOK, content.APIResponse{
		Success: true,
		Data:    pairing,
	})
}

// GetRandomPersonalityType 获取随机人格类型
func (h *Handler) GetRandomPersonalityType(c *gin.Context) {
	lang := c.Param("lang")
	if !h.i18n.IsValidLanguage(lang) {
		c.JSON(http.StatusBadRequest, content.APIResponse{
			Success: false,
			Error:   "Invalid language",
		})
		return
	}

	types, err := h.cache.GetAllPersonalityTypes(lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, content.APIResponse{
			Success: false,
			Error:   "Failed to get personality types",
		})
		return
	}

	if len(types) == 0 {
		c.JSON(http.StatusNotFound, content.APIResponse{
			Success: false,
			Error:   "No personality types found",
		})
		return
	}

	randomType := types[rand.Intn(len(types))]

	c.JSON(http.StatusOK, content.APIResponse{
		Success: true,
		Data:    randomType,
	})
}

// GetRandomPersonalityPairing 获取随机人格配对
func (h *Handler) GetRandomPersonalityPairing(c *gin.Context) {
	lang := c.Param("lang")
	if !h.i18n.IsValidLanguage(lang) {
		c.JSON(http.StatusBadRequest, content.APIResponse{
			Success: false,
			Error:   "Invalid language",
		})
		return
	}

	pairings, err := h.cache.GetAllPersonalityPairings(lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, content.APIResponse{
			Success: false,
			Error:   "Failed to get personality pairings",
		})
		return
	}

	if len(pairings) == 0 {
		c.JSON(http.StatusNotFound, content.APIResponse{
			Success: false,
			Error:   "No personality pairings found",
		})
		return
	}

	randomPairing := pairings[rand.Intn(len(pairings))]

	c.JSON(http.StatusOK, content.APIResponse{
		Success: true,
		Data:    randomPairing,
	})
}

// filterPersonalityTypes 筛选人格类型
func (h *Handler) filterPersonalityTypes(types []*content.PersonalityType, params *content.TypesSearchParams) []*content.PersonalityType {
	var filtered []*content.PersonalityType

	for _, personalityType := range types {
		// 性别筛选
		if params.Gender != "" && personalityType.Gender != params.Gender {
			continue
		}

		// 搜索筛选
		if params.Search != "" {
			searchTerm := strings.ToLower(params.Search)
			if !strings.Contains(strings.ToLower(personalityType.Code), searchTerm) &&
				!strings.Contains(strings.ToLower(personalityType.Name), searchTerm) &&
				!strings.Contains(strings.ToLower(personalityType.Tagline), searchTerm) {
				continue
			}
		}

		// 气质筛选
		if params.Temperament != "" {
			temperament := getTemperament(personalityType.Code)
			if temperament != params.Temperament {
				continue
			}
		}

		filtered = append(filtered, personalityType)
	}

	return filtered
}

// filterPersonalityPairings 筛选人格配对
func (h *Handler) filterPersonalityPairings(pairings []*content.PersonalityPairing, params *content.PairingsSearchParams) []*content.PersonalityPairing {
	var filtered []*content.PersonalityPairing

	for _, pairing := range pairings {
		// 男性人格筛选
		if params.Male != "" && pairing.MaleCode != strings.ToUpper(params.Male) {
			continue
		}

		// 女性人格筛选
		if params.Female != "" && pairing.FemaleCode != strings.ToUpper(params.Female) {
			continue
		}

		filtered = append(filtered, pairing)
	}

	return filtered
}

// getTemperament 获取人格的气质类型
func getTemperament(code string) string {
	if len(code) != 4 {
		return ""
	}

	// NT - 理性者
	if code[1] == 'N' && code[2] == 'T' {
		return "NT"
	}
	// NF - 理想主义者
	if code[1] == 'N' && code[2] == 'F' {
		return "NF"
	}
	// SJ - 守护者
	if code[1] == 'S' && code[3] == 'J' {
		return "SJ"
	}
	// SP - 艺术家
	if code[1] == 'S' && code[3] == 'P' {
		return "SP"
	}

	return ""
}
