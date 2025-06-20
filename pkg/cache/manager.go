package cache

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"mbti-site/pkg/content"

	"github.com/fsnotify/fsnotify"
)

// Manager 缓存管理器
type Manager struct {
	parser           *content.Parser
	personalityTypes map[string]*content.PersonalityType    // key: lang:code:gender
	pairings         map[string]*content.PersonalityPairing // key: lang:male:female
	mutex            sync.RWMutex
	watcher          *fsnotify.Watcher
}

// NewManager 创建新的缓存管理器
func NewManager(parser *content.Parser) *Manager {
	manager := &Manager{
		parser:           parser,
		personalityTypes: make(map[string]*content.PersonalityType),
		pairings:         make(map[string]*content.PersonalityPairing),
	}

	// 初始化文件监听器
	if err := manager.initWatcher(); err != nil {
		log.Printf("Warning: failed to initialize file watcher: %v", err)
	}

	// 预加载缓存
	manager.preloadCache()

	return manager
}

// GetPersonalityType 获取人格类型（带缓存）
func (m *Manager) GetPersonalityType(lang, code, gender string) (*content.PersonalityType, error) {
	key := fmt.Sprintf("%s:%s:%s", lang, strings.ToUpper(code), gender)

	m.mutex.RLock()
	if cached, exists := m.personalityTypes[key]; exists {
		m.mutex.RUnlock()
		return cached, nil
	}
	m.mutex.RUnlock()

	// 缓存中不存在，从文件加载
	personalityType, err := m.parser.GetPersonalityTypeByCodeAndGender(lang, code, gender)
	if err != nil {
		return nil, err
	}

	// 存入缓存
	m.mutex.Lock()
	m.personalityTypes[key] = personalityType
	m.mutex.Unlock()

	return personalityType, nil
}

// GetPersonalityPairing 获取人格配对（带缓存）
func (m *Manager) GetPersonalityPairing(lang, maleCode, femaleCode string) (*content.PersonalityPairing, error) {
	key := fmt.Sprintf("%s:%s:%s", lang, strings.ToUpper(maleCode), strings.ToUpper(femaleCode))

	m.mutex.RLock()
	if cached, exists := m.pairings[key]; exists {
		m.mutex.RUnlock()
		return cached, nil
	}
	m.mutex.RUnlock()

	// 缓存中不存在，从文件加载
	pairing, err := m.parser.GetPersonalityPairingByMaleAndFemale(lang, maleCode, femaleCode)
	if err != nil {
		return nil, err
	}

	// 存入缓存
	m.mutex.Lock()
	m.pairings[key] = pairing
	m.mutex.Unlock()

	return pairing, nil
}

// GetAllPersonalityTypes 获取所有人格类型
func (m *Manager) GetAllPersonalityTypes(lang string) ([]*content.PersonalityType, error) {
	// 尝试从缓存获取
	m.mutex.RLock()
	var cached []*content.PersonalityType
	for key, personalityType := range m.personalityTypes {
		if strings.HasPrefix(key, lang+":") {
			cached = append(cached, personalityType)
		}
	}
	m.mutex.RUnlock()

	// 如果缓存中有数据，返回缓存
	if len(cached) > 0 {
		return cached, nil
	}

	// 缓存中没有数据，从文件加载
	types, err := m.parser.ScanPersonalityTypes(lang)
	if err != nil {
		return nil, err
	}

	// 存入缓存
	m.mutex.Lock()
	for _, personalityType := range types {
		key := fmt.Sprintf("%s:%s:%s", lang, personalityType.Code, personalityType.Gender)
		m.personalityTypes[key] = personalityType
	}
	m.mutex.Unlock()

	return types, nil
}

// GetAllPersonalityPairings 获取所有人格配对
func (m *Manager) GetAllPersonalityPairings(lang string) ([]*content.PersonalityPairing, error) {
	// 尝试从缓存获取
	m.mutex.RLock()
	var cached []*content.PersonalityPairing
	for key, pairing := range m.pairings {
		if strings.HasPrefix(key, lang+":") {
			cached = append(cached, pairing)
		}
	}
	m.mutex.RUnlock()

	// 如果缓存中有数据，返回缓存
	if len(cached) > 0 {
		return cached, nil
	}

	// 缓存中没有数据，从文件加载
	pairings, err := m.parser.ScanPersonalityPairings(lang)
	if err != nil {
		return nil, err
	}

	// 存入缓存
	m.mutex.Lock()
	for _, pairing := range pairings {
		key := fmt.Sprintf("%s:%s:%s", lang, pairing.MaleCode, pairing.FemaleCode)
		m.pairings[key] = pairing
	}
	m.mutex.Unlock()

	return pairings, nil
}

// ClearCache 清空缓存
func (m *Manager) ClearCache() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.personalityTypes = make(map[string]*content.PersonalityType)
	m.pairings = make(map[string]*content.PersonalityPairing)
}

// preloadCache 预加载缓存
func (m *Manager) preloadCache() {
	languages := []string{"zh", "en"}

	for _, lang := range languages {
		// 预加载人格类型
		if _, err := m.GetAllPersonalityTypes(lang); err != nil {
			log.Printf("Warning: failed to preload personality types for %s: %v", lang, err)
		}

		// 预加载配对（如果存在）
		if _, err := m.GetAllPersonalityPairings(lang); err != nil {
			log.Printf("Warning: failed to preload pairings for %s: %v", lang, err)
		}
	}

	log.Println("Cache preloaded successfully")
}

// initWatcher 初始化文件监听器
func (m *Manager) initWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	m.watcher = watcher

	// 监听content目录
	contentDirs := []string{
		"./content/zh/types",
		"./content/zh/pairings",
		"./content/en/types",
		"./content/en/pairings",
	}

	for _, dir := range contentDirs {
		if err := watcher.Add(dir); err != nil {
			log.Printf("Warning: failed to watch directory %s: %v", dir, err)
		}
	}

	// 启动监听协程
	go m.watchFiles()

	return nil
}

// watchFiles 监听文件变化
func (m *Manager) watchFiles() {
	for {
		select {
		case event, ok := <-m.watcher.Events:
			if !ok {
				return
			}

			// 只处理.md文件的写入和删除事件
			if !strings.HasSuffix(event.Name, ".md") {
				continue
			}

			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Printf("File changed: %s, clearing cache", event.Name)
				m.invalidateCache(event.Name)
			}

		case err, ok := <-m.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("File watcher error: %v", err)
		}
	}
}

// invalidateCache 使缓存失效
func (m *Manager) invalidateCache(filePath string) {
	// 简单实现：清空所有缓存
	// 更精细的实现可以只清空相关的缓存项
	m.ClearCache()

	// 延迟重新加载缓存
	go func() {
		time.Sleep(100 * time.Millisecond)
		m.preloadCache()
	}()
}

// Close 关闭缓存管理器
func (m *Manager) Close() error {
	if m.watcher != nil {
		return m.watcher.Close()
	}
	return nil
}
