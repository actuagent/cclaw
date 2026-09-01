package config

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	mu                sync.RWMutex
	currentConfigPath string
)

// SetConfigPath stores the active config file path so that runtime directories
// can be derived from its parent (supports multi-instance deployments).
func SetConfigPath(path string) {
	mu.Lock()
	defer mu.Unlock()
	currentConfigPath = expandHome(path)
}

// GetConfigPath 返回当前生效的配置文件路径；未显式设置时返回 `~/.cclaw/config.yaml`。
func GetConfigPath() string {
	mu.RLock()
	defer mu.RUnlock()
	if currentConfigPath != "" {
		return currentConfigPath
	}
	return DefaultConfigPath()
}

// GetDataDir 返回实例级运行数据目录，该目录由当前配置文件的父目录确定。
// 此规则与上游 nanobot 的 get_data_dir 行为保持一致：
//
//	~/.cclaw/config.yaml  →  ~/.cclaw/
func GetDataDir() string {
	return ensureDir(filepath.Dir(GetConfigPath()))
}

// GetRuntimeSubdir 返回数据目录下指定名称的子目录，并在目录不存在时创建该目录。
// 例如 GetRuntimeSubdir("cron") 返回 `~/.cclaw/cron`。
func GetRuntimeSubdir(name string) string {
	return ensureDir(filepath.Join(GetDataDir(), name))
}

// GetWorkspacePath 解析并创建 Agent 工作区目录；workspace 为空时使用 `~/.cclaw/workspace`。
func GetWorkspacePath(workspace string) string {
	if workspace != "" {
		return ensureDir(expandHome(workspace))
	}
	return ensureDir(filepath.Join(GetDataDir(), "workspace"))
}

// GetSessionsDir 返回并创建会话数据目录 `~/.cclaw/sessions`。
func GetSessionsDir() string {
	return GetRuntimeSubdir("sessions")
}

// GetMemoryDir 返回并创建记忆数据目录 `~/.cclaw/memory`。
func GetMemoryDir() string {
	return GetRuntimeSubdir("memory")
}

// GetMediaDir 返回并创建媒体目录；channel 非空时返回对应渠道子目录。
func GetMediaDir(channel string) string {
	base := GetRuntimeSubdir("media")
	if channel != "" {
		return ensureDir(filepath.Join(base, channel))
	}
	return base
}

// GetCronDir 返回并创建定时任务数据目录 `~/.cclaw/cron`。
func GetCronDir() string {
	return GetRuntimeSubdir("cron")
}

// GetLogsDir 返回并创建日志目录 `~/.cclaw/logs`。
func GetLogsDir() string {
	return GetRuntimeSubdir("logs")
}

// GetPromptsDir 返回并创建 Prompt 模板目录 `~/.cclaw/prompts`。
func GetPromptsDir() string {
	return GetRuntimeSubdir("prompts")
}

// GetSkillsDir 返回并创建内建 Skill 目录 `~/.cclaw/skills`。
func GetSkillsDir() string {
	return GetRuntimeSubdir("skills")
}

// GetCLIHistoryPath 返回 CLI 历史记录文件路径 `~/.cclaw/history/cli_history`。
func GetCLIHistoryPath() string {
	dir := GetRuntimeSubdir("history")
	return filepath.Join(dir, "cli_history")
}

// GetCronStorePath 返回定时任务持久化文件路径 `~/.cclaw/cron/jobs.json`。
func GetCronStorePath() string {
	return filepath.Join(GetCronDir(), "jobs.json")
}

func ensureDir(path string) string {
	_ = os.MkdirAll(path, 0755)
	return path
}
