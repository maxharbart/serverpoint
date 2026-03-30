package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"serverpoint/config"
	"serverpoint/models"
)

type CollectorResult struct {
	Metric    models.Metric
	Processes []models.Process
	OS        string
	Hostname  string
	Uptime    string
	Online    bool
}

func CollectMetrics(server models.Server, cfg *config.Config) CollectorResult {
	result := CollectorResult{
		Metric: models.Metric{ServerID: server.ID, CreatedAt: time.Now()},
	}

	// Decrypt credentials
	var password, privateKey string
	if server.AuthType == "password" && server.EncryptedPassword != "" {
		var err error
		password, err = Decrypt(server.EncryptedPassword, cfg.AESKey)
		if err != nil {
			log.Printf("Failed to decrypt password for server %d: %v", server.ID, err)
			return result
		}
	}
	if server.AuthType == "key" && server.EncryptedKey != "" {
		var err error
		privateKey, err = Decrypt(server.EncryptedKey, cfg.AESKey)
		if err != nil {
			log.Printf("Failed to decrypt key for server %d: %v", server.ID, err)
			return result
		}
	}

	client, err := NewSSHClient(server.Host, server.Port, server.Username, password, privateKey, server.AuthType)
	if err != nil {
		log.Printf("SSH connection failed for server %d (%s): %v", server.ID, server.Host, err)
		return result
	}
	defer client.Close()

	result.Online = true

	// Collect hostname
	if out, err := client.Run("hostname"); err == nil {
		result.Hostname = strings.TrimSpace(out)
	}

	// Collect OS info
	if out, err := client.Run("cat /etc/os-release 2>/dev/null | grep PRETTY_NAME | cut -d'\"' -f2"); err == nil {
		result.OS = strings.TrimSpace(out)
	}
	if result.OS == "" {
		if out, err := client.Run("uname -o"); err == nil {
			result.OS = strings.TrimSpace(out)
		}
	}

	// Collect uptime
	if out, err := client.Run("uptime -p 2>/dev/null || uptime"); err == nil {
		result.Uptime = strings.TrimSpace(out)
	}

	// Collect CPU usage
	cpuCmd := `top -bn2 -d0.5 | grep "Cpu(s)" | tail -1 | awk '{print $2+$4}'`
	if out, err := client.Run(cpuCmd); err == nil {
		if v, err := strconv.ParseFloat(strings.TrimSpace(out), 64); err == nil {
			result.Metric.CPUPercent = v
		}
	}

	// Collect RAM
	ramCmd := `free -b | awk '/^Mem:/{printf "%d %d", $2, $3}'`
	if out, err := client.Run(ramCmd); err == nil {
		parts := strings.Fields(strings.TrimSpace(out))
		if len(parts) == 2 {
			if total, err := strconv.ParseUint(parts[0], 10, 64); err == nil {
				result.Metric.RAMTotal = total
			}
			if used, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
				result.Metric.RAMUsed = used
			}
			if result.Metric.RAMTotal > 0 {
				result.Metric.RAMPercent = float64(result.Metric.RAMUsed) / float64(result.Metric.RAMTotal) * 100
			}
		}
	}

	// Collect Disk
	diskCmd := `df -B1 / | awk 'NR==2{printf "%d %d", $2, $3}'`
	if out, err := client.Run(diskCmd); err == nil {
		parts := strings.Fields(strings.TrimSpace(out))
		if len(parts) == 2 {
			if total, err := strconv.ParseUint(parts[0], 10, 64); err == nil {
				result.Metric.DiskTotal = total
			}
			if used, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
				result.Metric.DiskUsed = used
			}
			if result.Metric.DiskTotal > 0 {
				result.Metric.DiskPercent = float64(result.Metric.DiskUsed) / float64(result.Metric.DiskTotal) * 100
			}
		}
	}

	// Collect top processes
	procCmd := `ps aux --sort=-%cpu | head -11 | tail -10 | awk '{printf "%s %s %s %s %s\n", $2, $1, $3, $4, $11}'`
	if out, err := client.Run(procCmd); err == nil {
		lines := strings.Split(strings.TrimSpace(out), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}
			pid, _ := strconv.Atoi(fields[0])
			cpu, _ := strconv.ParseFloat(fields[2], 64)
			ram, _ := strconv.ParseFloat(fields[3], 64)
			result.Processes = append(result.Processes, models.Process{
				ServerID:  server.ID,
				PID:       pid,
				User:      fields[1],
				CPU:       cpu,
				RAM:       ram,
				Command:   fields[4],
				CreatedAt: time.Now(),
			})
		}
	}

	return result
}

func FormatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.1f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
