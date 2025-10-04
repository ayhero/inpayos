package protocol

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// channel_account
type TargetSummary struct {
	Target              string
	TimeIndexSummaryLib map[string]*TimeIndexSummary
}

func (t *TargetSummary) GetIndex(time_idx, data_idx string) decimal.Decimal {
	dt, ok := t.TimeIndexSummaryLib[time_idx]
	if ok {
		return dt.GetIndex(data_idx)
	}
	return decimal.Zero
}

// last_10_mm
type TimeIndexSummary struct {
	Name string
	Data map[string]*DataIndexSummary
}

func (t *TimeIndexSummary) GetIndex(idx string) decimal.Decimal {
	dt, ok := t.Data[idx]
	if !ok {
		return decimal.Zero
	}
	return dt.GetTz(DefaultTimeZone.String())
}

// succ_count
type DataIndexSummary struct {
	Name string
	Data MapData
}

//tz:time:data
//tz:last_10_mm:succ_rate/succ_count

func (t *DataIndexSummary) GetTz(time_zone string) decimal.Decimal {
	return *t.Data.GetDecimal(time_zone)
}

func NewSummary(target string) *TargetSummary {
	return &TargetSummary{
		Target: target,
	}
}

type TimeIndex struct {
	Range  string `json:"range"`
	Target int64  `json:"target"`
	Unit   string `json:"unit"`
}

func (t *TimeIndex) GenTimeRange(params MapData) {
	query_start_unix := int64(0)
	query_end_unix := int64(0)
	defer func() {
		params.Set(QUERY_IDX_TIME_INDEX, fmt.Sprintf("%v_%v_%v", t.Range, t.Target, t.Unit))
		params.Set(QUEYRY_IDX_QUERY_START_UNIX, query_start_unix)
		params.Set(QUERY_IDX_QUERY_END_UNIX, query_end_unix)
	}()
	current_unix := params.GetInt64(QUERY_IDX_REFRESH_AT)
	switch t.Range {
	case IDX_LAST:
		zero_unix := params.GetInt64(QUERY_IDX_TODAY_ZERO_UNIX)
		today_age_unix := params.GetInt64(QUERY_IDX_TODAY_AGE_UNIX)
		step := int64(0)
		switch t.Unit {
		case IDX_TIME_SECOND:
			step = t.Target
		case IDX_TIME_MINUTE:
			step = t.Target * 60
		case IDX_TIME_HOUR:
			step = t.Target * 60 * 60
		case IDX_TIME_DAY:
			step = t.Target * 24 * 60 * 60
		}
		current_idx := today_age_unix / step
		query_start_unix = zero_unix + (current_idx-1)*step
		query_end_unix = zero_unix + (current_idx)*step
	case IDX_RECENT:
		query_end_unix = current_unix
		switch t.Unit {
		case IDX_TIME_SECOND:
			query_start_unix = current_unix - t.Target
		case IDX_TIME_MINUTE:
			query_start_unix = current_unix - t.Target*60
		case IDX_TIME_HOUR:
			query_start_unix = current_unix - t.Target*60*60
		}
	}
}
