package conf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/goccy/go-yaml"
)

const (
	DataDir     = "data"
	cfgFilename = "config.yaml"
)

type Email struct {
	SMTPHost string   `yaml:"smtp_host"`
	SMTPPort int      `yaml:"smtp_port"`
	From     string   `yaml:"email_from"`
	Password string   `yaml:"email_password"`
	To       []string `yaml:"default_to"`
}

type Config struct {
	Email Email `yaml:"email"`
}

var (
	once    sync.Once
	initErr error

	mu      sync.RWMutex
	current *Config
)

func cfgPath() string {
	return filepath.Join(DataDir, cfgFilename)
}

func defaultConfig() *Config {
	return &Config{
		Email: Email{
			SMTPHost: "smtp.qq.com",
			SMTPPort: 25,
			From:     "",
			Password: "",
			To:       []string{},
		},
	}
}

func ensure() error {
	p := cfgPath()
	_, err := os.Stat(p)
	if err == nil {
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("conf: 检查文件状态失败: %w", err)
	}

	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("conf: 创建配置目录 %s 失败: %w", dir, err)
	}

	data, err := yaml.Marshal(defaultConfig())
	if err != nil {
		return fmt.Errorf("conf: 序列化默认配置失败: %w", err)
	}

	if err := os.WriteFile(p, data, 0644); err != nil {
		return fmt.Errorf("conf: 写入默认配置文件 %s 失败: %w", p, err)
	}

	return nil
}

func validate(cfg *Config) error {
	var errs []error

	// 验证 SMTP 端口范围
	if cfg.Email.SMTPPort <= 0 || cfg.Email.SMTPPort > 65535 {
		errs = append(errs, fmt.Errorf("无效的 smtp_port: %d (必须在 1 到 65535 之间)", cfg.Email.SMTPPort))
	}

	// 如果配置了发件人，则服务器地址不能为空
	if cfg.Email.From != "" && cfg.Email.SMTPHost == "" {
		errs = append(errs, errors.New("指定了发件人时，smtp_host 不能为空"))
	}

	return errors.Join(errs...)
}

func load() (*Config, error) {
	if err := ensure(); err != nil {
		return nil, err
	}

	p := cfgPath()
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("conf: 读取文件 %s 失败: %w", p, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("conf: 解析 yaml 文件 %s 失败: %w", p, err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("conf: 配置校验失败: %w", err)
	}

	return &cfg, nil
}

func Init() error {
	once.Do(func() {
		if err := os.MkdirAll(DataDir, 0755); err != nil {
			initErr = fmt.Errorf("conf: 创建目录 %s 失败: %w", DataDir, err)
			return
		}

		cfg, err := load()
		if err != nil {
			initErr = err
			return
		}

		mu.Lock()
		current = cfg
		mu.Unlock()
	})
	return initErr
}

func Get() *Config {
	mu.RLock()
	defer mu.RUnlock()
	if current == nil {
		panic("conf: 配置未初始化，请先调用 Init()")
	}
	return current
}
