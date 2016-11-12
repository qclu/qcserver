package models

const TIME_FMT = "2006-01-02 15:04:05"

//db tables
const (
	DB_T_ADMINISTRATOR = "qc_administrator"
	DB_T_HOSPITAL      = "qc_hospital"
	DB_T_DEPARTMENT    = "qc_department"
	DB_T_METHODOLOGY   = "qc_methodology"
	DB_T_DEVMODEL      = "qc_dev_model"
	DB_T_REGMODEL      = "qc_reg_model"
	DB_T_REGREL        = "qc_reg_rel"
	DB_T_QCPRODUCT     = "qc_qc_product"
	DB_T_QCLOTNUM      = "qc_qc_lotnum"
	DB_T_REGPRODUCE    = "qc_reg_produce"
	DB_T_DEVREL        = "qc_dev_rel"
	DB_T_HWVERSION     = "qc_hw_version"
	DB_T_SWVERSION     = "qc_sw_version"
)

//error definition
const (
	ERR_OBJ_NOT_EXIST = "Object not exist"
)
