package models

type MasterTransactionStatusStruct struct {
	Complete int
	Failed   int
}

var MasterTransactionStatusId = MasterTransactionStatusStruct{
	Complete: 1,
	Failed:   2,
}

type MasterTransactionStatusDbFieldStruct struct {
	Id   string
	Name string
}

var MasterTransactionStatusDbField = MasterTransactionStatusDbFieldStruct{
	Id:   "id",
	Name: "name",
}

type MasterTransactionStatus struct {
	Id   int    `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

func (MasterTransactionStatus) TableName() string {
	return "nexmedis_master_transaction_statuses"
}

type MasterTransactionStatusPrimaryId struct {
	Id int `gorm:"primaryKey"`
}

func (MasterTransactionStatusPrimaryId) TableName() string {
	return "nexmedis_master_transaction_statuses"
}
