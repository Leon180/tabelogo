package main

const (
	webPort                 = "8080"
	tabelogSpiderServiceURL = "http://tabelog-spider-service"
)

func main() {
	server := NewServer()
	err := server.Run(":" + webPort)
	if err != nil {
		panic(err)
	}
}
