package main

const (
	webPort         = "80"
	radisPort       = "6379"
	maxCollectLinks = 3
)

func main() {
	server := NewServer()
	err := server.Run(":" + webPort)
	if err != nil {
		panic(err)
	}
}
