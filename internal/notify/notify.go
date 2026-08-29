package notify

import (
	"time"

	"github.com/awydd/nudge/conf"
)

type Notify interface {
	Send() error
}

func CheckAndNotify(a *conf.Anniversary, factory func(title, desc, date string, years int) Notify) (bool, error) {
	origTime, err := time.Parse(time.DateOnly, a.OriginalDate)
	if err != nil {
		return false, err
	}

	now := time.Now()
	currentYear := now.Year()

	// 今年已经通知过了，直接跳过
	if a.LastNotifiedYear >= currentYear {
		return false, nil
	}

	// 今年对应的纪念日日期
	targetDate := time.Date(currentYear, origTime.Month(), origTime.Day(), 0, 0, 0, 0, now.Location())

	// 计算提前通知的起始窗口日期
	advanceDays := max(a.AdvanceDays, 0)
	notifyStartDate := targetDate.AddDate(0, 0, -advanceDays)

	// 判断当前时间是否已到达提前通知窗口期
	if !now.Before(notifyStartDate) {
		years := max(currentYear-origTime.Year(), 0)

		// 实例化具体的通知通道并发送
		notifier := factory(a.Title, a.Description, targetDate.Format(time.DateOnly), years)
		if notifier != nil {
			notifier.Send()
		}

		// 更新已通知年份，避免重复通知
		a.LastNotifiedYear = currentYear
		return true, nil
	}

	return false, nil
}
