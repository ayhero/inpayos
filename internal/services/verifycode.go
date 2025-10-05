package services

import (
	"fmt"
	"inpayos/internal/config"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"math/rand"
	"time"
)

const (
	// Redis key 模板
	verifyCodeSendCountKey = "verify_code:count:%s:%s:%s" // 验证码发送次数key：verify_code:count:type:email:20060102
	verifyCodeLastSendKey  = "verify_code:last:%s:%s"     // 最后发送时间key：verify_code:last:type:email
	verifyCodeKey          = "verify_code:code:%s:%s"     // 验证码key：verify_code:code:type:email
)

// VerifyCodeService 验证码服务
type VerifyCodeService struct {
	config     *config.VerifyCodeConfig
	MsgService *MessageService
}

var (
	verifyCodeService *VerifyCodeService
)

func SetupVerifyCodeService() {
	cfg := config.Get().VerifyCode
	if cfg == nil {
		panic("email verify code config is nil")
	}
	verifyCodeService = &VerifyCodeService{
		config:     cfg,
		MsgService: GetMessageService(),
	}
}

// GetVerifyCodeService 获取验证码服务实例
func GetVerifyCodeService() *VerifyCodeService {
	if verifyCodeService == nil {
		SetupVerifyCodeService()
	}
	return verifyCodeService
}

// generateCode 生成验证码
func (s *VerifyCodeService) generateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := ""
	for range s.config.Length {
		code += fmt.Sprintf("%d", r.Intn(10))
	}
	return code
}

// canSendCode 检查是否可以发送验证码
func (s *VerifyCodeService) canSendCode(typ, email string) error {
	// 检查发送间隔
	lastSendKey := fmt.Sprintf(verifyCodeLastSendKey, typ, email)
	lastTime, err := models.GetInt(lastSendKey)
	if err == nil && lastTime > 0 {
		if time.Now().Unix()-lastTime < int64(s.config.SendInterval) {
			return fmt.Errorf("retry after %v seconds", s.config.SendInterval-int(time.Now().Unix()-lastTime))
		}
	}

	// 检查当天发送次数
	today := time.Now().Format(time.DateOnly)
	countKey := fmt.Sprintf(verifyCodeSendCountKey, typ, email, today)
	sendCount, err := models.GetInt(countKey)
	if err == nil && sendCount >= int64(s.config.MaxSendTimes) {
		return fmt.Errorf("out of send limit, max %d times per day", s.config.MaxSendTimes)
	}

	return nil
}

// SendEmailVerifyCode 发送验证码
func SendEmailVerifyCode(typ, email string) error {
	s := GetVerifyCodeService()
	// 检查发送限制
	if err := s.canSendCode(typ, email); err != nil {
		return err
	}

	// 生成验证码
	code := s.generateCode()

	log.Get().Infof("Sending verify code: %s to %s", code, email)
	// 发送邮件
	// 构造邮件内容
	msg := &Message{
		Type: protocol.MsgTypeVerifyCode,
		To:   email,
		Params: map[string]any{
			"to":   email,
			"code": code,
		},
	}
	if err := s.MsgService.SendEmailMessage(msg); err != nil {
		return err
	}

	// 更新发送记录
	today := time.Now().Format(time.DateOnly)
	countKey := fmt.Sprintf(verifyCodeSendCountKey, typ, email, today)
	lastSendKey := fmt.Sprintf(verifyCodeLastSendKey, typ, email)
	codeKey := fmt.Sprintf(verifyCodeKey, typ, email)

	// 更新发送次数
	sendCount, _ := models.GetInt(countKey)
	models.SetInt(countKey, sendCount+1, time.Hour*24)

	// 更新最后发送时间
	models.SetInt(lastSendKey, time.Now().Unix(), time.Duration(s.config.SendInterval)*time.Second)

	// 存储验证码
	models.SetCache(codeKey, code, time.Duration(s.config.Expiration)*time.Minute)

	return nil
}

// VerifyEmailCode 验证验证码
func VerifyEmailCode(typ, email, code string) (result bool) {
	result = false
	if code == "" {
		return
	}
	codeKey := fmt.Sprintf(verifyCodeKey, typ, email)
	storedCode, err := models.GetCache(codeKey)
	if err != nil {
		return
	}
	defer func() {
		if result {
			models.Delete(codeKey)
		}
	}()
	// 验证成功后删除验证码
	result = storedCode == code
	return
}
