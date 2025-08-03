package config

func GetRedisAddr() string {
	// return os.Getenv("REDIS_ADDR")
	return "localhost:6379"
}

func GetMysqlDSN() string {
	// return os.Getenv("MYSQL_DSN")
	return "root:yang@tcp(127.0.0.1:3306)/game?charset=utf8mb4&parseTime=True&loc=Local"
}

func GetJwtSecret() []byte {
	// return []byte(os.Getenv("JWT_SECRET"))
	return []byte("secret")
}
