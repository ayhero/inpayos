package services

import (
	"errors"
	"inpayos/internal/log"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
)

// ChangeMerchantPassword 修改商户密码
func ChangeMerchantPassword(email, newPassword string) error {
	// 根据邮箱获取商户信息
	merchant := models.GetMerchantByEmail(email)
	if merchant == nil {
		//log.Get().Errorf("ChangeMerchantPassword: Merchant not found for email: %s", email)
		return errors.New("商户不存在")
	}
	merchant.Salt = utils.GenerateSalt()
	// 设置新密码和盐
	merchant.SetPassword(newPassword)

	// 加密密码
	merchant.Encrypt()

	// 更新密码和盐
	if err := models.GetDB().Model(&models.Merchant{}).Where("email = ?", email).Updates(map[string]interface{}{
		"password": merchant.GetPassword(),
		"salt":     merchant.Salt,
	}).Error; err != nil {
		//log.Get().Errorf("ChangeMerchantPassword: Update password error: %v", err)
		return errors.New("密码更新失败")
	}

	//log.Get().Infof("ChangeMerchantPassword: Password changed successfully for email: %s", email)
	return nil
}

// ResetMerchantPassword 重置商户密码
func ResetMerchantPassword(email string) (string, error) {
	// 检查商户是否存在
	merchant := models.GetMerchantByEmail(email)
	if merchant == nil {
		return "", errors.New("商户不存在")
	}

	// 生成新密码
	newPassword, err := utils.GenerateRandomPassword(12)
	if err != nil {
		log.Get().Error("Failed to generate password: ", err)
		return "", errors.New("生成新密码失败")
	}

	// 更新商户密码
	merchant.SetPassword(newPassword)
	merchant.Encrypt()

	if err := models.WriteDB.Model(merchant).Update("password", merchant.Password).Error; err != nil {
		log.Get().Error("Failed to update password: ", err)
		return "", errors.New("密码更新失败")
	}

	return newPassword, nil
}

// SendNewPasswordEmail 发送新密码邮件
func SendNewPasswordEmail(email, newPassword string) error {
	msgService := GetMessageService()
	msg := &Message{
		Type: protocol.MsgTypeNewPassword,
		To:   email,
		Params: map[string]any{
			"to":       email,
			"password": newPassword,
		},
	}

	return msgService.SendEmailMessage(msg)
}
