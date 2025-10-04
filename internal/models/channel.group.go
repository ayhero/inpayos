package models

import "inpayos/internal/protocol"

type ChannelGroup struct {
	ID   int64  `json:"id" gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	Code string `json:"code" gorm:"column:code;type:varchar(50);uniqueIndex"`
	*ChannelGroupValues
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"` // 更新时间 (毫秒时间戳)
}

type ChannelGroupValues struct {
	Name    *string       `json:"name" gorm:"column:name;type:varchar(255)"`
	Status  *string       `json:"status" gorm:"column:status;type:varchar(50);default:0"`
	Setting *GroupSetting `json:"setting" gorm:"column:setting;serializer:json"`
	Members GroupMembers  `json:"members" gorm:"column:members;serializer:json"`
}

type GroupSetting struct {
	Strategy  string `json:"strategy" gorm:"column:strategy"`
	Weight    string `json:"weight" gorm:"column:weight"`
	RankType  string `json:"rank_type" gorm:"column:rank_type"`
	TimeIndex string `json:"time_index" gorm:"column:time_index"`
	DataIndex string `json:"data_index" gorm:"column:data_index"`
	Timezone  string `json:"timezone" gorm:"column:timezone"`
}

func (t *ChannelGroup) TableName() string {
	return "t_channel_groups"
}

func GetActiveChannelGroupByCode(code string) *ChannelGroup {
	group := &ChannelGroup{}
	err := ReadDB.Where("code = ?", code).First(group).Error
	if err != nil {
		return nil
	}
	return group
}

// GetActiveChannelGroups returns all active channel groups
func GetActiveChannelGroups() ([]*ChannelGroup, error) {
	var groups []*ChannelGroup
	err := ReadDB.Where("status = ?", protocol.StatusActive).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

type GroupMembers []*GroupMember

// 组成员权重
type GroupMember struct {
	Member   string  `json:"member"`
	Target   float64 `json:"target"`
	Current  float64 `json:"current"`
	Distance float64 `json:"distance"`
	Weight   float64 `json:"weight"`
}

func (t *GroupMember) ComputeDistance() {
	t.Distance = t.Target - t.Current
}

func NewGroupMember(member string, target float64) *GroupMember {
	return &GroupMember{
		Member: member,
		Target: target,
	}
}

// ChannelGroupValues Getter Methods
// GetName returns the Name value
func (cgv *ChannelGroupValues) GetName() string {
	if cgv.Name == nil {
		return ""
	}
	return *cgv.Name
}

// GetStatus returns the Status value
func (cgv *ChannelGroupValues) GetStatus() string {
	if cgv.Status == nil {
		return ""
	}
	return *cgv.Status
}

// GetSetting returns the Setting value
func (cgv *ChannelGroupValues) GetSetting() *GroupSetting {
	return cgv.Setting
}

// GetMembers returns the Members value
func (cgv *ChannelGroupValues) GetMembers() GroupMembers {
	return cgv.Members
}

// ChannelGroupValues Setter Methods (support method chaining)
// SetName sets the Name value
func (cgv *ChannelGroupValues) SetName(value string) *ChannelGroupValues {
	cgv.Name = &value
	return cgv
}

// SetStatus sets the Status value
func (cgv *ChannelGroupValues) SetStatus(value string) *ChannelGroupValues {
	cgv.Status = &value
	return cgv
}

// SetSetting sets the Setting value
func (cgv *ChannelGroupValues) SetSetting(value *GroupSetting) *ChannelGroupValues {
	cgv.Setting = value
	return cgv
}

// SetMembers sets the Members value
func (cgv *ChannelGroupValues) SetMembers(value GroupMembers) *ChannelGroupValues {
	cgv.Members = value
	return cgv
}

// SetValues sets multiple ChannelGroupValues fields at once
func (cg *ChannelGroup) SetValues(values *ChannelGroupValues) *ChannelGroup {
	if values == nil {
		return cg
	}

	if cg.ChannelGroupValues == nil {
		cg.ChannelGroupValues = &ChannelGroupValues{}
	}

	// Set all fields from the provided values
	if values.Name != nil {
		cg.ChannelGroupValues.SetName(*values.Name)
	}
	if values.Status != nil {
		cg.ChannelGroupValues.SetStatus(*values.Status)
	}
	if values.Setting != nil {
		cg.ChannelGroupValues.SetSetting(values.Setting)
	}
	// Members is not a pointer, so we always set it
	cg.ChannelGroupValues.SetMembers(values.Members)

	return cg
}
