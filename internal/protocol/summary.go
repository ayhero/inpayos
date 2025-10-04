package protocol

import (
	"time"
)

type TimeZone string

func (t TimeZone) String() string {
	return string(t)
}

func (t TimeZone) ToLocation() *time.Location {
	loc, err := time.LoadLocation(t.String())
	if err == nil {
		return loc
	}
	return nil
}
func GetTimeLocation(time_zone string) *time.Location {
	loc, ok := TimeZoneLib[time_zone]
	if ok {
		return loc
	}
	tz := TimeZone(time_zone)
	loc = tz.ToLocation()
	if loc == nil {
		TimeZoneLib[time_zone] = loc
	}
	return loc
}

const (
	Asia_ShangHai TimeZone = "Asia/Shanghai"
)

var (
	DefaultTimeZone TimeZone = Asia_ShangHai
	TimeZoneLib              = map[string]*time.Location{}
)

func init() {
	TimeZoneLib[Asia_ShangHai.String()] = Asia_ShangHai.ToLocation()
}

const (
	IDX_RECENT      = "recent"
	IDX_LAST        = "last"
	IDX_TIME_MINUTE = "mm"
	IDX_TIME_SECOND = "ss"
	IDX_TIME_HOUR   = "hh"
	IDX_TIME_DAY    = "d"
	IDX_TIME_MONTH  = "m"
	IDX_TIME_YEAR   = "y"
	IDX_AVG         = "avg"
	IDX_MAX         = "max"
	IDX_MIN         = "min"

	QUERY_IDX_DATA_INDEX        = "data_index"
	QUERY_IDX_TIME_INDEX        = "time_index"
	QUERY_IDX_REFRESH_AT        = "refresh_at"
	QUERY_IDX_TIME_ZONE         = "time_zone"
	QUERY_IDX_TODAY_ZERO_UNIX   = "today_zero_unix"
	QUERY_IDX_TODAY_AGE_UNIX    = "today_age_unix"
	QUEYRY_IDX_QUERY_START_UNIX = "query_start_unix"
	QUERY_IDX_QUERY_END_UNIX    = "query_end_unix"

	DAY_UNIX = 86400

	STAT_IDX_TOTAL_COUNT      = "total_count"
	STAT_IDX_TOTAL_AMOUNT     = "total_amount"
	STAT_IDX_TOTAL_USD_AMOUNT = "total_usd_amount"

	STAT_IDX_SUCC_COUNT      = "succ_count"
	STAT_IDX_SUCC_AMOUNT     = "succ_amount"
	STAT_IDX_SUCC_RATE       = "succ_rate"
	STAT_IDX_SUCC_USD_AMOUNT = "succ_usd_amount"

	STAT_IDX_FAIL_COUNT      = "fail_count"
	STAT_IDX_FAIL_AMOUNT     = "fail_amount"
	STAT_IDX_FAIL_RATE       = "fail_rate"
	STAT_IDX_FAIL_USD_AMOUNT = "fail_usd_amount"

	STAT_IDX_PND_COUNT      = "pnd_count"
	STAT_IDX_PND_AMOUNT     = "pnd_amount"
	STAT_IDX_PND_RATE       = "pnd_rate"
	STAT_IDX_PND_USD_AMOUNT = "pnd_usd_amount"
)

var (
	DefaultSummaryIndexList = []*TimeIndex{
		{
			Range:  IDX_LAST,
			Target: 1,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_LAST,
			Target: 10,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_LAST,
			Target: 15,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_LAST,
			Target: 30,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_LAST,
			Target: 45,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_LAST,
			Target: 1,
			Unit:   IDX_TIME_HOUR,
		},
		{
			Range:  IDX_LAST,
			Target: 90,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_LAST,
			Target: 2,
			Unit:   IDX_TIME_HOUR,
		}, {
			Range:  IDX_RECENT,
			Target: 1,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_RECENT,
			Target: 10,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_RECENT,
			Target: 15,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_RECENT,
			Target: 30,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_RECENT,
			Target: 45,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_RECENT,
			Target: 1,
			Unit:   IDX_TIME_HOUR,
		},
		{
			Range:  IDX_RECENT,
			Target: 90,
			Unit:   IDX_TIME_MINUTE,
		},
		{
			Range:  IDX_RECENT,
			Target: 2,
			Unit:   IDX_TIME_HOUR,
		},
	}
)
