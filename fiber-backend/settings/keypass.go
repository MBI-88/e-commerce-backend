package settings


func GetKeypass() []byte{
	return []byte(Settings{}.Key)
}