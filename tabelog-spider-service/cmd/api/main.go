package main

const (
	webPort         = "80"
	maxCollectLinks = 5
)

func main() {
	server := NewServer()
	err := server.Run(":" + webPort)
	if err != nil {
		panic(err)
	}
}
