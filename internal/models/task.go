package models

import (
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"time"
)

// Task 定时任务表
type Task struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TaskID     string `json:"task_id" gorm:"column:task_id;type:varchar(64);uniqueIndex"` // 任务唯一标识
	Name       string `json:"name" gorm:"column:name;type:varchar(256)"`                  // 任务名称
	Type       string `json:"type" gorm:"column:type;type:varchar(64)"`                   // 任务类型
	HandlerKey string `json:"handler_key" gorm:"column:handler_key;type:varchar(64)"`     // 处理器标识
	*TaskValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
}

type TaskValues struct {
	Cron       *string          `json:"cron" gorm:"column:cron;type:varchar(32)"`               // cron表达式
	Status     *string          `json:"status" gorm:"column:status;type:varchar(32)"`           // 任务状态
	LastTime   *int64           `json:"last_time" gorm:"column:last_time"`                      // 上次执行时间
	NextTime   *int64           `json:"next_time" gorm:"column:next_time"`                      // 下次执行时间
	LastResult *string          `json:"last_result" gorm:"column:last_result;type:varchar(64)"` // 上次执行结果
	RetryCount *int             `json:"retry_count" gorm:"column:retry_count"`                  // 重试次数
	MaxRetries *int             `json:"max_retries" gorm:"column:max_retries"`                  // 最大重试次数
	Timeout    *int             `json:"timeout" gorm:"column:timeout"`                          // 超时时间(秒)
	Params     protocol.MapData `json:"params" gorm:"column:params;type:json;serializer:json"`  // 任务参数
	Remark     *string          `json:"remark" gorm:"column:remark;type:varchar(512)"`          // 备注说明
}

// TableName 表名
func (Task) TableName() string {
	return "t_tasks"
}

// ScanTask 扫描可执行的任务
func ScanTask() ([]Task, error) {
	var tasks []Task
	err := ReadDB.Where("status = ? AND next_time>=0 and next_time <= ?",
		protocol.StatusEnabled,
		time.Now().UnixMilli(),
	).Find(&tasks).Error

	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// UpdateTaskExecution 更新任务执行信息
func UpdateTaskExecution(task *Task, result string, nextTime time.Time) error {
	task.SetLastTime(utils.TimeNowMilli())
	task.SetNextTime(nextTime.UnixMilli())
	task.SetLastResult(result)
	return WriteDB.Save(task).Error
}

func CheckTaskExist(taskID string) (isExist bool) {
	err := ReadDB.Model(&Task{}).Where("task_id = ?", taskID).Select("task_id").First(&taskID).Error
	return err == nil
}

// GetTaskByTaskID 根据TaskID获取任务
func GetTaskByTaskID(taskID string) (*Task, error) {
	var task Task
	err := ReadDB.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateTask 创建任务
func CreateTask(task *Task) error {
	return WriteDB.Create(task).Error
}

// UpdateTask 更新任务
func UpdateTask(task *Task) error {
	return WriteDB.Save(task).Error
}

// TaskValues Getter Methods
// GetCron returns the Cron value
func (tv *TaskValues) GetCron() string {
	if tv.Cron == nil {
		return ""
	}
	return *tv.Cron
}

// GetStatus returns the Status value
func (tv *TaskValues) GetStatus() string {
	if tv.Status == nil {
		return ""
	}
	return *tv.Status
}

// GetLastTime returns the LastTime value
func (tv *TaskValues) GetLastTime() int64 {
	if tv.LastTime == nil {
		return 0
	}
	return *tv.LastTime
}

// GetNextTime returns the NextTime value
func (tv *TaskValues) GetNextTime() int64 {
	if tv.NextTime == nil {
		return 0
	}
	return *tv.NextTime
}

// GetLastResult returns the LastResult value
func (tv *TaskValues) GetLastResult() string {
	if tv.LastResult == nil {
		return ""
	}
	return *tv.LastResult
}

// GetRetryCount returns the RetryCount value
func (tv *TaskValues) GetRetryCount() int {
	if tv.RetryCount == nil {
		return 0
	}
	return *tv.RetryCount
}

// GetMaxRetries returns the MaxRetries value
func (tv *TaskValues) GetMaxRetries() int {
	if tv.MaxRetries == nil {
		return 0
	}
	return *tv.MaxRetries
}

// GetTimeout returns the Timeout value
func (tv *TaskValues) GetTimeout() int {
	if tv.Timeout == nil {
		return 0
	}
	return *tv.Timeout
}

// GetParams returns the Params value
func (tv *TaskValues) GetParams() protocol.MapData {
	return tv.Params
}

// GetRemark returns the Remark value
func (tv *TaskValues) GetRemark() string {
	if tv.Remark == nil {
		return ""
	}
	return *tv.Remark
}

// TaskValues Setter Methods (support method chaining)
// SetCron sets the Cron value
func (tv *TaskValues) SetCron(value string) *TaskValues {
	tv.Cron = &value
	return tv
}

// SetStatus sets the Status value
func (tv *TaskValues) SetStatus(value string) *TaskValues {
	tv.Status = &value
	return tv
}

// SetLastTime sets the LastTime value
func (tv *TaskValues) SetLastTime(value int64) *TaskValues {
	tv.LastTime = &value
	return tv
}

// SetNextTime sets the NextTime value
func (tv *TaskValues) SetNextTime(value int64) *TaskValues {
	tv.NextTime = &value
	return tv
}

// SetLastResult sets the LastResult value
func (tv *TaskValues) SetLastResult(value string) *TaskValues {
	tv.LastResult = &value
	return tv
}

// SetRetryCount sets the RetryCount value
func (tv *TaskValues) SetRetryCount(value int) *TaskValues {
	tv.RetryCount = &value
	return tv
}

// SetMaxRetries sets the MaxRetries value
func (tv *TaskValues) SetMaxRetries(value int) *TaskValues {
	tv.MaxRetries = &value
	return tv
}

// SetTimeout sets the Timeout value
func (tv *TaskValues) SetTimeout(value int) *TaskValues {
	tv.Timeout = &value
	return tv
}

// SetParams sets the Params value
func (tv *TaskValues) SetParams(value protocol.MapData) *TaskValues {
	tv.Params = value
	return tv
}

// SetRemark sets the Remark value
func (tv *TaskValues) SetRemark(value string) *TaskValues {
	tv.Remark = &value
	return tv
}

// SetValues sets multiple TaskValues fields at once
func (t *Task) SetValues(values *TaskValues) *Task {
	if values == nil {
		return t
	}

	if t.TaskValues == nil {
		t.TaskValues = &TaskValues{}
	}

	// Set all fields from the provided values
	if values.Cron != nil {
		t.TaskValues.SetCron(*values.Cron)
	}
	if values.Status != nil {
		t.TaskValues.SetStatus(*values.Status)
	}
	if values.LastTime != nil {
		t.TaskValues.SetLastTime(*values.LastTime)
	}
	if values.NextTime != nil {
		t.TaskValues.SetNextTime(*values.NextTime)
	}
	if values.LastResult != nil {
		t.TaskValues.SetLastResult(*values.LastResult)
	}
	if values.RetryCount != nil {
		t.TaskValues.SetRetryCount(*values.RetryCount)
	}
	if values.MaxRetries != nil {
		t.TaskValues.SetMaxRetries(*values.MaxRetries)
	}
	if values.Timeout != nil {
		t.TaskValues.SetTimeout(*values.Timeout)
	}
	// Params is not a pointer, so we always set it
	t.TaskValues.SetParams(values.Params)
	if values.Remark != nil {
		t.TaskValues.SetRemark(*values.Remark)
	}

	return t
}
