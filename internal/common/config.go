package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func initDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "release")

	viper.SetDefault("user.recheck_interval", "1h")

	viper.SetDefault("agent.enable", false)
	viper.SetDefault("agent.user.concurrency", 3)
	viper.SetDefault("agent.item.concurrency", 3)
	viper.SetDefault("agent.item.max", 3000)
	viper.SetDefault("agent.discover.level", 0)

	viper.SetDefault("crawl.enable", true)
	viper.SetDefault("storage.enable", true)

	viper.SetDefault("datasource.driver", "mysql")
	viper.SetDefault("datasource.port", "3306")
	viper.SetDefault("datasource.charset", "utf8mb4")
	viper.SetDefault("datasource.loc", "Asia/Shanghai")

	viper.SetDefault("http.timeout", 10000)
	viper.SetDefault("http.retry_max", 20)
	viper.SetDefault("http.interval.user", 4000)
	viper.SetDefault("http.interval.item", 4000)
	viper.SetDefault("http.interval.discover", 4000)
}

// LoadDotEnv loads environment variables from the first existing .env file.
// Search order: current working directory, parent, grandparent.
func LoadDotEnv() {
	workDir, _ := os.Getwd()
	candidates := []string{
		filepath.Join(workDir, ".env"),
		filepath.Join(workDir, "../.env"),
		filepath.Join(workDir, "../../.env"),
	}

	for _, filePath := range candidates {
		if _, err := os.Stat(filePath); err == nil {
			_ = parseAndSetEnv(filePath)
			return
		}
	}
}

func parseAndSetEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		value = strings.Trim(value, "\"")
		value = strings.Trim(value, "'")

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// BindEnv binds MOUBAN_* environment variables to viper dot-keys.
func BindEnv() {
	viper.SetEnvPrefix("MOUBAN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func toEnvKey(dotKey string) string {
	return "MOUBAN_" + strings.ToUpper(strings.ReplaceAll(dotKey, ".", "_"))
}

func mustNotEmpty(missing *[]string, key string) {
	if strings.TrimSpace(viper.GetString(key)) == "" {
		*missing = append(*missing, toEnvKey(key))
	}
}

func ValidateConfig() error {
	var missing []string

	mustNotEmpty(&missing, "datasource.host")
	mustNotEmpty(&missing, "datasource.database")
	mustNotEmpty(&missing, "datasource.username")
	mustNotEmpty(&missing, "datasource.password")

	if viper.GetBool("crawl.enable") {
		mustNotEmpty(&missing, "http.auth")
	}

	if viper.GetBool("storage.enable") {
		mustNotEmpty(&missing, "s3.endpoint")
		mustNotEmpty(&missing, "s3.region")
		mustNotEmpty(&missing, "s3.bucket")
		mustNotEmpty(&missing, "s3.access_key")
		mustNotEmpty(&missing, "s3.secret_key")
	}

	if viper.GetBool("agent.enable") && !viper.GetBool("crawl.enable") {
		missing = append(missing, "MOUBAN_CRAWL_ENABLE must be true when MOUBAN_AGENT_ENABLE=true")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}
	return nil
}

// InitConfigBase loads defaults + .env + environment variables without validation.
// Useful for tests that only need partial config domains.
func InitConfigBase() {
	initDefaults()
	LoadDotEnv()
	BindEnv()
}

func InitConfig() {
	InitConfigBase()

	if err := ValidateConfig(); err != nil {
		panic(err)
	}
}
