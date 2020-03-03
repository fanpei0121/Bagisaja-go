package model

type GoodsCat struct {
	Id       int    `gorm:"primary_key" json:"id"`
	ParentId int    `json:"parent_id"`
	Name     string `json:"name"`
	TypeId   int    `json:"type_id"`
	Sort     int    `json:"sort"`
	ImageId  string `json:"image_id"`
	Status   int    `json:"status"`
	Utime    int    `json:"utime"`
}

type CatData struct {
	GoodsCat
	Url string `json:"url"`
}

// 获取分类
func (t *GoodsCat) GetCats(parentId int) []CatData {
	var result []CatData
	recordNotFound := DB.Table("jshop_goods_cat a").
		Select("a.*,b.url").
		Joins("left join jshop_images b ON a.image_id=b.id").
		Order("a.sort asc").
		Where("a.parent_id = ? AND a.status = 1", parentId).
		Scan(&result).RecordNotFound()
	if recordNotFound {
		return nil
	}
	return result
}
