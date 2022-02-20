package model

import (
	"blogservice/global"
	"blogservice/pkg/setting"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}

	db.SingularTable(true)
	db.Callback().Create().Replace("gorm:uptate_time_stamp", updateTimeStampForCreateCallBack)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallBack)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil

}

func updateTimeStampForCreateCallBack(scope *gorm.Scope) {
	if !scope.HasError() {
		_ = scope.SetColumn("CreatedOn", time.Now().Unix())
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}
func updateTimeStampForUpdateCallBack(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_colume"); !ok {
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string

		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeleteOnFiled := scope.FieldByName("DeleteOn")
		isDelFiled, hasIsDelFiled := scope.FieldByName("IsDel")

		if !scope.Search.Unscoped && hasDeleteOnFiled && hasIsDelFiled {
			now := time.Now().Unix()
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelFiled.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExit(scope.CombinedConditionSql()),
				addExtraSpaceIfExit(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExit(scope.CombinedConditionSql()),
				addExtraSpaceIfExit(extraOption),
			)).Exec()
		}
	}
}
func addExtraSpaceIfExit(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
