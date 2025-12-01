package dao

// 模拟数据库
var database = map[string]string{}

func AddUser(username string, password string) {
	database[username] = password
}
func FindUser(username string, password string) bool {
	if password01, ok := database[username]; ok {
		if password01 == password {
			return true
		} else {
			return false
		}
	}
	return false
}
func SelectPasswordFromUsername(username string) string {
	return database[username]
}
