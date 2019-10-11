package model

// Member 用户模型
type Member struct {
	QQ         string `gorm:"unique;not null;PRIMARY_KEY"`
	NickName   string
	Sponsor    bool `gorm:"default:0"`
	Role       int  `gorm:"default:2"`
	Money      float32
	GameServer int
	GameName   string
}

// GetUser get user
func GetUser(qq string) (Member, error) {
	var member Member
	result := DB.Where("QQ = ?", qq).Find(&member)
	return member, result.Error
}

// UpdateUser update ones privillage
func UpdateUser(qq string, member Member) error {
	return DB.Where("QQ = ?", qq).Find(new(Member)).Updates(member).Error
}
