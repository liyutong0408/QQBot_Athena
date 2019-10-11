package model

func migration() {
	DB.AutoMigrate(&Member{})
}
