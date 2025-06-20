package content

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

// Parser 内容解析器
type Parser struct {
	contentPath string
	markdown    goldmark.Markdown
}

// NewParser 创建新的内容解析器
func NewParser(contentPath string) *Parser {
	return &Parser{
		contentPath: contentPath,
		markdown:    goldmark.New(),
	}
}

// ParsePersonalityType 解析人格类型文件
func (p *Parser) ParsePersonalityType(filePath string) (*PersonalityType, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 获取文件修改时间
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	// 解析front-matter和内容
	frontMatter, body, err := p.parseFrontMatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse front-matter in %s: %w", filePath, err)
	}

	// 将Markdown转换为HTML
	var htmlBuf bytes.Buffer
	if err := p.markdown.Convert([]byte(body), &htmlBuf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown to HTML: %w", err)
	}

	// 构建PersonalityType结构
	personalityType := &PersonalityType{
		Code:     frontMatter.Code,
		Gender:   frontMatter.Gender,
		Name:     frontMatter.Name,
		Tagline:  frontMatter.Tagline,
		Summary:  frontMatter.Summary,
		Body:     htmlBuf.String(),
		Lang:     frontMatter.Lang,
		FilePath: filePath,
		ModTime:  stat.ModTime(),
	}

	return personalityType, nil
}

// ParsePersonalityPairing 解析人格配对文件
func (p *Parser) ParsePersonalityPairing(filePath string) (*PersonalityPairing, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 获取文件修改时间
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	// 解析front-matter和内容
	frontMatter, body, err := p.parseFrontMatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse front-matter in %s: %w", filePath, err)
	}

	// 将Markdown转换为HTML
	var htmlBuf bytes.Buffer
	if err := p.markdown.Convert([]byte(body), &htmlBuf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown to HTML: %w", err)
	}

	// 构建PersonalityPairing结构
	pairing := &PersonalityPairing{
		ID:                 frontMatter.ID,
		MaleCode:           frontMatter.MaleCode,
		FemaleCode:         frontMatter.FemaleCode,
		CompatibilityScore: frontMatter.CompatibilityScore,
		Body:               htmlBuf.String(),
		Lang:               frontMatter.Lang,
		FilePath:           filePath,
		ModTime:            stat.ModTime(),
	}

	return pairing, nil
}

// parseFrontMatter 解析front-matter
func (p *Parser) parseFrontMatter(content []byte) (*FrontMatter, string, error) {
	// 匹配front-matter的正则表达式
	frontMatterRegex := regexp.MustCompile(`^---\s*\n(.*?)\n---\s*\n(.*)$`)
	matches := frontMatterRegex.FindSubmatch(content)

	if len(matches) != 3 {
		return nil, "", fmt.Errorf("invalid front-matter format")
	}

	// 解析YAML front-matter
	var frontMatter FrontMatter
	if err := yaml.Unmarshal(matches[1], &frontMatter); err != nil {
		return nil, "", fmt.Errorf("failed to parse YAML front-matter: %w", err)
	}

	// 返回front-matter和正文内容
	body := string(matches[2])
	return &frontMatter, body, nil
}

// ScanPersonalityTypes 扫描所有人格类型文件
func (p *Parser) ScanPersonalityTypes(lang string) ([]*PersonalityType, error) {
	typesPath := filepath.Join(p.contentPath, lang, "types")
	var types []*PersonalityType

	err := filepath.WalkDir(typesPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理.md文件
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		personalityType, parseErr := p.ParsePersonalityType(path)
		if parseErr != nil {
			// 记录错误但继续处理其他文件
			fmt.Printf("Warning: failed to parse %s: %v\n", path, parseErr)
			return nil
		}

		types = append(types, personalityType)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan personality types: %w", err)
	}

	return types, nil
}

// ScanPersonalityPairings 扫描所有人格配对文件
func (p *Parser) ScanPersonalityPairings(lang string) ([]*PersonalityPairing, error) {
	pairingsPath := filepath.Join(p.contentPath, lang, "pairings")
	var pairings []*PersonalityPairing

	err := filepath.WalkDir(pairingsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理.md文件
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		pairing, parseErr := p.ParsePersonalityPairing(path)
		if parseErr != nil {
			// 记录错误但继续处理其他文件
			fmt.Printf("Warning: failed to parse %s: %v\n", path, parseErr)
			return nil
		}

		pairings = append(pairings, pairing)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan personality pairings: %w", err)
	}

	return pairings, nil
}

// GetPersonalityTypeByCodeAndGender 根据代码和性别获取人格类型
func (p *Parser) GetPersonalityTypeByCodeAndGender(lang, code, gender string) (*PersonalityType, error) {
	fileName := fmt.Sprintf("%s_%s.md", strings.ToLower(code), gender)
	filePath := filepath.Join(p.contentPath, lang, "types", fileName)

	return p.ParsePersonalityType(filePath)
}

// GetPersonalityPairingByMaleAndFemale 根据男女人格代码获取配对
func (p *Parser) GetPersonalityPairingByMaleAndFemale(lang, maleCode, femaleCode string) (*PersonalityPairing, error) {
	fileName := fmt.Sprintf("%s_m__%s_f.md", strings.ToLower(maleCode), strings.ToLower(femaleCode))
	filePath := filepath.Join(p.contentPath, lang, "pairings", fileName)

	return p.ParsePersonalityPairing(filePath)
}
