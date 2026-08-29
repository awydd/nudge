package console

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/awydd/nudge/conf"
	"github.com/awydd/nudge/internal/notify"
	"github.com/awydd/nudge/utils"
)

func TestConsoleIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "anniversaries.json")

	// 注册测试结束后的清理回调（选填，t.TempDir() 本身会自动删除）
	t.Cleanup(func() {
		t.Logf("正在自动清理临时测试目录: %s", tmpDir)
	})

	t.Run("初始化", func(t *testing.T) {
		var store conf.Store
		err := utils.ReadJSON(dbPath, &store)
		if err != nil && os.IsNotExist(err) {
			store = conf.Store{
				Anniversaries: []conf.Anniversary{
					{
						ID:               "1",
						Title:            "nudge",
						Description:      "https://github.com/awydd/nudge",
						OriginalDate:     "2026-08-29",
						AdvanceDays:      0,
						LastNotifiedYear: 0,
					},
				},
			}
			err = utils.WriteJSON(dbPath, &store)
			if err != nil {
				t.Fatalf("初始化写入 JSON 失败: %v", err)
			}
		}

		var verifyStore conf.Store
		if err := utils.ReadJSON(dbPath, &verifyStore); err != nil {
			t.Fatalf("读取回显 JSON 失败: %v", err)
		}

		if len(verifyStore.Anniversaries) != 1 {
			t.Errorf("期望得到 1 个纪念日，实际得到 %d 个", len(verifyStore.Anniversaries))
		}
	})

	t.Run("检查并通知", func(t *testing.T) {
		var store conf.Store
		if err := utils.ReadJSON(dbPath, &store); err != nil {
			t.Fatalf("读取存储失败: %v", err)
		}

		updated := false
		for i := range store.Anniversaries {
			a := &store.Anniversaries[i]

			triggeredCount := 0
			triggered, err := notify.CheckAndNotify(a, func(title, desc, date string, years int) notify.Notify {
				triggeredCount++
				return &ConsoleNotifier{
					Title:       title,
					Description: desc,
					Date:        date,
					Years:       years,
				}
			})

			if err != nil {
				t.Errorf("发生意外错误: %v", err)
			}

			if triggered {
				updated = true
			}
		}

		if !updated {
			t.Error("期望由于 LastNotifiedYear=0 且提前天数为0而触发通知")
		}

		if err := utils.WriteJSON(dbPath, &store); err != nil {
			t.Fatalf("写入更新后的存储失败: %v", err)
		}

		var secondStore conf.Store
		_ = utils.ReadJSON(dbPath, &secondStore)

		for i := range secondStore.Anniversaries {
			a := &secondStore.Anniversaries[i]
			triggered, _ := notify.CheckAndNotify(a, func(title, desc, date string, years int) notify.Notify {
				return &ConsoleNotifier{
					Title:       title,
					Description: desc,
					Date:        date,
					Years:       years,
				}
			})
			if triggered {
				t.Error("期望同一年内的第二次运行不触发通知")
			}
		}
	})
}
