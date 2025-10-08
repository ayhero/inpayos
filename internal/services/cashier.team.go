package services

import "sync"

type CasherTeamService struct {
}

var (
	casherTeamService     *CasherTeamService
	casherTeamServiceOnce sync.Once
)

func SetupCasherTeamService() {
	casherTeamServiceOnce.Do(func() {
		casherTeamService = &CasherTeamService{}
	})
}

// GetCasherTeamService 获取CasherTeam服务单例
func GetCasherTeamService() *CasherTeamService {
	if casherTeamService == nil {
		SetupCasherTeamService()
	}
	return casherTeamService
}
