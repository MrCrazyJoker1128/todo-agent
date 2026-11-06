package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"service/internal/dao"
	"service/internal/middleware"
	"service/internal/model/entity"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
)

// SummarizeTasksArgs 工作总结维度参数：全是 AI 的维度建议（可空），系统先弹确认卡用户可改后才真执行。
type SummarizeTasksArgs struct {
	Period    string `json:"period" jsonschema:"时间维度建议：day|week|month|custom，可空"`
	StartDate string `json:"startDate" jsonschema:"开始时间 YYYY-MM-DD HH:mm:ss，可空"`
	EndDate   string `json:"endDate" jsonschema:"结束时间 YYYY-MM-DD HH:mm:ss，可空"`
	TaskType  string `json:"taskType" jsonschema:"任务类型建议：work|personal|project|adhoc，可空"`
	Priority  string `json:"priority" jsonschema:"紧急程度建议：high|medium|low，可空"`
}

// SummaryStats 统计结果。
type SummaryStats struct {
	Total      int            `json:"total"`
	Done       int            `json:"done"`
	Doing      int            `json:"doing"`
	Todo       int            `json:"todo"`
	Cancelled  int            `json:"cancelled"`
	Rate       int            `json:"rate"`
	ByType     map[string]int `json:"byType"`
	ByPriority map[string]int `json:"byPriority"`
}

// SummarizeTasksResult 总结执行结果（含落库 recordId）。
type SummarizeTasksResult struct {
	RecordID  string         `json:"recordId"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Stats     SummaryStats   `json:"stats"`
	Period    string         `json:"period"`
	StartDate string         `json:"startDate"`
	EndDate   string         `json:"endDate"`
	TaskType  string         `json:"taskType"`
	Priority  string         `json:"priority"`
	Tasks     []SummaryItem `json:"tasks"`
}

// SummaryItem 总结清单行。
type SummaryItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	DueDate string `json:"dueDate"`
	Status string `json:"status"`
}

func SummarizeTasks() tool.InvokableTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "summarize_tasks",
			Desc: "按时间维度（日/周/月）×任务类型×紧急程度查询任务并生成工作总结落库。所有参数都只是 AI 的维度建议（全可空，空=AI 不确定），系统一定先弹确认卡（用户可改），AI 禁止假设维度已确认。",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"period": {
						Desc:     "时间维度建议：day|week|month|custom，可空",
						Type:     "string",
						Required: false,
					},
					"startDate": {
						Desc:     "开始时间 YYYY-MM-DD HH:mm:ss，可空",
						Type:     "string",
						Required: false,
					},
					"endDate": {
						Desc:     "结束时间 YYYY-MM-DD HH:mm:ss，可空",
						Type:     "string",
						Required: false,
					},
					"taskType": {
						Desc:     "任务类型建议：work|personal|project|adhoc，可空",
						Type:     "string",
						Required: false,
					},
					"priority": {
						Desc:     "紧急程度建议：high|medium|low，可空",
						Type:     "string",
						Required: false,
					},
				}),
		},
		summarizeTasks,
	)
}

// shanghaiLoc Asia/Shanghai 时区，load 失败回退 UTC+8 固定偏移。
func shanghaiLoc() *time.Location {
	if loc, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		return loc
	}
	return time.FixedZone("UTC+8", 8*3600)
}

// normalizeSummaryPeriod 归一化维度：非法 period→空；空 period+齐起止→custom；day|week|month 按当天算窗。返回 period/start/end。
func normalizeSummaryPeriod(period, start, end string, now time.Time) (string, string, string) {
	p := strings.ToLower(strings.TrimSpace(period))
	if p != "day" && p != "week" && p != "month" && p != "custom" {
		p = ""
	}
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)
	if p == "" && start != "" && end != "" {
		return "custom", start, end
	}
	loc := shanghaiLoc()
	now = now.In(loc)
	y, m, d := now.Date()
	switch p {
	case "day":
		s := time.Date(y, m, d, 0, 0, 0, 0, loc)
		e := time.Date(y, m, d, 23, 59, 59, 0, loc)
		return p, s.Format("2006-01-02 15:04:05"), e.Format("2006-01-02 15:04:05")
	case "week":
		wd := int(now.Weekday())
		if wd == 0 {
			wd = 7
		}
		monday := time.Date(y, m, d, 0, 0, 0, 0, loc).AddDate(0, 0, -(wd - 1))
		sunday := monday.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		return p, monday.Format("2006-01-02 15:04:05"), sunday.Format("2006-01-02 15:04:05")
	case "month":
		first := time.Date(y, m, 1, 0, 0, 0, 0, loc)
		last := first.AddDate(0, 1, 0).Add(-time.Second)
		return p, first.Format("2006-01-02 15:04:05"), last.Format("2006-01-02 15:04:05")
	default:
		return p, start, end
	}
}

func summarizeTasks(ctx context.Context, input *SummarizeTasksArgs) (string, error) {
	userId := middleware.GetUserIDFromContext(ctx)
	if userId == "" {
		return "", errors.New("user_id not found in context")
	}
	uid, err := uuid.Parse(userId)
	if err != nil {
		return "", errors.New("invalid userId")
	}
	period, start, end := normalizeSummaryPeriod(input.Period, input.StartDate, input.EndDate, time.Now())
	taskType := strings.TrimSpace(input.TaskType)
	priority := strings.ToLower(strings.TrimSpace(input.Priority))
	if priority != "high" && priority != "medium" && priority != "low" {
		priority = ""
	}

	tasks, err := dao.DbTasks.QueryTasksForSummary(ctx, uid, start, end, taskType, priority)
	if err != nil {
		glog.Error(ctx, "[summarizeTasks] query failed: %v", err)
		return "", err
	}

	stats := SummaryStats{ByType: map[string]int{}, ByPriority: map[string]int{}, Total: len(tasks)}
	var doneTitles, todoTitles []SummaryItem
	for _, t := range tasks {
		stats.ByType[t.TaskType]++
		stats.ByPriority[t.Priority]++
		switch t.Status {
		case "done":
			stats.Done++
		case "doing":
			stats.Doing++
		case "cancelled":
			stats.Cancelled++
		default:
			stats.Todo++
		}
		item := SummaryItem{ID: t.ID, Title: t.Title, DueDate: t.DueDate, Status: t.Status}
		if t.Status == "done" {
			doneTitles = append(doneTitles, item)
		} else {
			todoTitles = append(todoTitles, item)
		}
		_ = doneTitles
		_ = todoTitles
	}
	if stats.Total > 0 {
		stats.Rate = stats.Done * 100 / stats.Total
	}

	title := "全部任务总结"
	rangeText := "全部时间"
	if period != "" {
		switch period {
		case "day":
			title = "日报 " + start[:10]
		case "week":
			title = fmt.Sprintf("周总结 %s~%s", start[5:10], end[5:10])
		case "month":
			title = fmt.Sprintf("月总结 %s~%s", start[:7], end[:7])
		case "custom":
			title = fmt.Sprintf("时段总结 %s~%s", start[:10], end[:10])
		}
		if start != "" {
			rangeText = start + " ~ " + end
		}
	}

	var sb strings.Builder
	sb.WriteString("# " + title + "\n")
	sb.WriteString("范围：" + rangeText + "\n")
	if taskType != "" {
		sb.WriteString("类型：" + taskType + "\n")
	}
	if priority != "" {
		sb.WriteString("紧急程度：" + priority + "\n")
	}
	sb.WriteString(fmt.Sprintf("总数：%d 完成：%d 进行中：%d 待办：%d", stats.Total, stats.Done, stats.Doing, stats.Todo) + "\n")
	if stats.Cancelled > 0 {
		sb.WriteString(fmt.Sprintf("已取消：%d\n", stats.Cancelled))
	}
	sb.WriteString(fmt.Sprintf("完成率：%d%%\n", stats.Rate))
	for typ, n := range stats.ByType {
		sb.WriteString(fmt.Sprintf("- %s：%d\n", typ, n))
	}
	sb.WriteString("任务清单：\n")
	shown := tasks
	extra := 0
	if len(shown) > 100 {
		extra = len(shown) - 100
		shown = shown[:100]
	}
	for _, t := range shown {
		mark := "○"
		if t.Status == "done" {
			mark = "✓"
		}
		sb.WriteString(fmt.Sprintf("%s %s（%s）\n", mark, t.Title, t.DueDate))
	}
	if extra > 0 {
		sb.WriteString(fmt.Sprintf("…等 %d 条\n", extra))
	}
	content := sb.String()

	statsJSON, _ := json.Marshal(stats)
	id := uuid.New()
	var startT, endT *gtime.Time
	if start != "" {
		if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", start, shanghaiLoc()); err == nil {
			startT = gtime.NewFromTime(parsed)
		}
	}
	if end != "" {
		if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", end, shanghaiLoc()); err == nil {
			endT = gtime.NewFromTime(parsed)
		}
	}
	if err := dao.DbTaskSummaries.InsertSummary(ctx, &entity.DbTaskSummaries{
		Id: id, UserId: uid, Period: period, StartDate: startT, EndDate: endT,
		TaskType: taskType, Priority: priority, Title: title, Content: content, Stats: string(statsJSON),
	}); err != nil {
		glog.Error(ctx, "[summarizeTasks] insert failed: %v", err)
		return "", err
	}

	items := make([]SummaryItem, 0, len(tasks))
	for _, t := range tasks {
		items = append(items, SummaryItem{ID: t.ID, Title: t.Title, DueDate: t.DueDate, Status: t.Status})
	}
	resultJSON, err := json.Marshal(SummarizeTasksResult{
		RecordID: id.String(), Title: title, Content: content, Stats: stats,
		Period: period, StartDate: start, EndDate: end, TaskType: taskType, Priority: priority, Tasks: items,
	})
	if err != nil {
		return "", err
	}
	return string(resultJSON), nil
}
