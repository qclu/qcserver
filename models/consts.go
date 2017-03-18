package models

const TIME_FMT = "2006-01-02 15:04:05"

//db tables
const (
	DB_T_ADMINISTRATOR = "qc_administrator"
	DB_T_HOSPITAL      = "qc_hospital"
	DB_T_DEPARTMENT    = "qc_department"
	DB_T_METHODOLOGY   = "qc_methodology"
	DB_T_DEVMODEL      = "qc_dev_model"
	DB_T_REGMODEL      = "qc_reagent_model"
	DB_T_REGREL        = "qc_reagent_rel"
	DB_T_QCPRODUCT     = "qc_qc_product"
	DB_T_QCLOTNUM      = "qc_qc_lotnum"
	DB_T_REGPRODUCE    = "qc_reagent_produce"
	DB_T_DEVREL        = "qc_dev_rel"
	DB_T_HWVERSION     = "qc_hw_version"
	DB_T_SWVERSION     = "qc_sw_version"
	DB_T_DEVLOG        = "qc_dev_log"
	DB_T_LOGTYPE       = "qc_log_type"
	DB_T_DEVSTAT       = "qc_dev_stat"

	DB_V_DEVREL = "qc_v_dev_rel"
)

//error definition
const (
	ERR_OBJ_NOT_EXIST = "Object not exist"
)
