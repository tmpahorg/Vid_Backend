package models

type PassRecord struct {
	Uid           int    `gorm:"primary_key"`
	EncryptedPass string `gorm:"type:char(128);not null"`
}

// Rename Table
func (PassRecord) TableName() string {
	// Default: pass_record
	return "tbl_passrecord"
}

// @override
func (p *PassRecord) CheckValid() bool {
	return p.Uid != 0 && p.EncryptedPass != ""
}