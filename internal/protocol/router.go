package protocol

const (
	RouterStrategyAll  = "all"
	RouterStrategyOnce = "once"
)

type RouterInfo struct {
	Mid             string            `json:"mid"`
	ChannelGroup    string            `json:"channel_group"`
	ChannelAccounts []string          `json:"channel_accounts"`
	Strategy        string            `json:"strategy"`
	ChannelCodeLib  map[string]string `json:"channel_code_lib"`
}
