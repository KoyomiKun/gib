package model

import (
	"Blog/util"
	"Blog/util/errmsg"

	"github.com/jinzhu/gorm"
)

type Category struct {
	ID   uint   `json:"id,omitempty"`
	Name string `gorm:"type:varchar(20);not null" json:"name,omitempty"`
}

// 查询类别是否存在
func ExistsCategory(c *Category) errmsg.Code {
	var cate Category
	db.Select("id").Where("name = ?", c.Name).First(&cate)
	if cate.ID > 0 {
		return errmsg.ERROR_USERNAME_USED
	}
	return errmsg.SUCCESS
}

// 新增类别
func InsertCategory(c *Category) errmsg.Code {
	if err := db.Create(&c).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 查询类别表
func GetCategorysList(pageSize int, pageNum int) ([]Category, errmsg.Code) {
	var categorys []Category
	// limit: 指定需要多少条； offset：指定从哪一条开始
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&categorys).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errmsg.ERROR
	}
	return categorys, errmsg.SUCCESS
}

// 编辑类别信息
func EditCategory(c *Category) errmsg.Code {
	dataMap := util.Struct2Map(*c)
	logger.Debug("Edit category:", dataMap)
	err := db.Model(c).Updates(dataMap).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除类别
func DeleteCategory(id int) errmsg.Code {
	err = db.Delete(&Category{}, id).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
