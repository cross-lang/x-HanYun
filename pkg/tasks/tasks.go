package tasks

import (
	"strconv"
	"strings"
	"time"
	"x-HanYun/pkg/core/log"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// StartPeriodicTask 启动周期任务，支持指定首次执行时间
func StartPeriodicTask(interval string, startTime string, cmd func()) {
	// 解析间隔时间
	duration, err := time.ParseDuration(interval)
	if err != nil {
		// 默认使用小时为单位
		duration, _ = time.ParseDuration(interval + "h")
	}
	log.Info(">>>>>>>> Periodic Task Duration", zap.String("duration", duration.String()))

	// 计算距离开始时间的延迟
	var startDelay time.Duration
	if startTime != "" {
		now := time.Now()
		// 解析开始时间（格式：HH:MM）
		startTimeParts := strings.Split(startTime, ":")
		if len(startTimeParts) == 2 {
			hour, _ := strconv.Atoi(startTimeParts[0])
			minute, _ := strconv.Atoi(startTimeParts[1])

			// 计算今天的开始时间
			startToday := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

			// 如果今天的开始时间已过，则计算明天的开始时间
			if startToday.Before(now) {
				startToday = startToday.Add(24 * time.Hour)
			}

			// 计算距离开始时间的延迟
			startDelay = startToday.Sub(now)
		}
	}

	log.Info(">>>>>>>> Periodic Task Configuration Follows",
		zap.String("interval", duration.String()),
		zap.String("startTime", startTime),
		zap.Duration("firstExecutionDelay", startDelay))

	// 创建并启动定时任务
	c := cron.New(cron.WithSeconds())

	// 使用 time.AfterFunc 处理延迟启动
	if startDelay > 0 {
		log.Info(">>>>>>>> Periodic Task Will been Scheduled After Delay", zap.Duration("delay", startDelay))
		time.AfterFunc(startDelay, func() {
			log.Info(">>>>>>>> Periodic Task First Execution", zap.String("nowTime", time.Now().Format("2006-01-02 15:04:05")))
			cmd()

			// 添加周期性任务
			_, err = c.AddFunc("@every "+duration.String(), func() {
				log.Info(">>>>>>>> Periodic Task Execution Func", zap.String("nowTime", time.Now().Format("2006-01-02 15:04:05")))

				// 执行周期任务的业务逻辑
				cmd()
			})
			if err != nil {
				log.Error("<<<<<<<< Failed to Add Periodic Task", zap.Error(err))
				return
			}
			log.Info(">>>>>>>> Periodic Task Has been Scheduled", zap.String("interval", duration.String()))
		})
	} else {
		log.Info(">>>>>>>> Periodic Task Will been Scheduled Immediately", zap.String("interval", duration.String()))
		// 如果没有延迟，立即添加周期性任务
		_, err = c.AddFunc("@every "+duration.String(), func() {
			log.Info(">>>>>>>> Periodic Task Execution Func", zap.String("nowTime", time.Now().Format("2006-01-02 15:04:05")))

			// 执行周期任务的业务逻辑
			cmd()
		})
		if err != nil {
			log.Error("<<<<<<<< Failed to Add Periodic Task", zap.Error(err))
			return
		}
		log.Info(">>>>>>>> Periodic Task Has been Scheduled Immediately", zap.String("interval", duration.String()))
	}

	// 启动 cron 任务
	c.Start()
	log.Info(">>>>>>>> Periodic Task Has been Started Successfully")

	// 阻塞主函数，防止程序退出
	select {}
}

// StartScheduledTask 启动定时任务（支持秒级精度）
// spec
// 格式：
// 秒 分 时 日 月 星期
// 例子：
// "每天0点0分30秒": "30 0 0 * * *"
// "每天9点0分0秒":  "0 0 9 * * *"
// "每天12点30分0秒": "0 30 12 * * *"
// "每天18点0分0秒":  "0 0 18 * * *"
// "每分钟的第15秒":  "15 * * * * *"
// "每小时的第30分钟0秒": "0 30 * * * *"
// "工作日9点0分0秒":  "0 0 9 * * 1-5"
// "周末10点0分0秒":  "0 0 10 * * 6,0"
// "每月1号0点0分0秒": "0 0 0 1 * *"
func StartScheduledTask(spec string, cmd func()) {
	// 创建支持秒级的cron实例
	c := cron.New(cron.WithSeconds())

	entryID, err := c.AddFunc(spec, func() {
		log.Info(">>>>>>>> Scheduled Task Exec Func", zap.String("timeFormat", time.Now().Format("2006-01-02 15:04:05.000")))
		// 执行定时任务的业务逻辑
		cmd()
	})
	if err != nil {
		log.Error("<<<<<<<< Scheduled Task Created Error", zap.Error(err), zap.Int("entryID", int(entryID)))
		return
	}
	log.Info(">>>>>>>> Scheduled Task Created Successfully", zap.Int("entryID", int(entryID)), zap.String("spec", spec))

	// 启动 cron 任务
	c.Start()
	log.Info(">>>>>>>> Scheduled Task Started Successfully", zap.Int("entryID", int(entryID)))

	// 阻塞主函数，防止程序退出
	select {}
}
