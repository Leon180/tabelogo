package main

const (
	webPort         = "80"
	maxCollectLinks = 4
)

func main() {
	server := NewServer()
	err := server.Run(":" + webPort)
	if err != nil {
		panic(err)
	}
}
