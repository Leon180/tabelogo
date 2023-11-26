package main

const (
	webPort                 = "8080"
	tabelogSpiderServiceURL = "http://tabelog-spider-service" // service's name
	authenticateServiceURL  = "http://authenticate-service"   // service's name
	googleMapServiceURL     = "http://google-map-service"     // service's name
	loggerServiceURL        = "http://logger-service"         // service's name
)

func main() {
	server := NewServer()
	err := server.Run(":" + webPort)
	if err != nil {
		panic(err)
	}
}
