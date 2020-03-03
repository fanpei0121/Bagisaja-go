package model

type PromotionResult struct {
	Id          int    `gorm:"primary_key" json:"id"`
	PromotionId int    `json:"promotion_id"`
	Code        string `json:"code"`
	Params      string `json:"params"`
}

// 根据促销id获取结果列表
func (t *PromotionResult) GetListByPromotionId(promotionId int) []PromotionResult {
	var results []PromotionResult
	if DB.Where("promotion_id = ?", promotionId).Find(&results).RecordNotFound() {
		return nil
	}
	return results
}
