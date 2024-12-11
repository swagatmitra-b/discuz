package main

import "log"

func main() {

	dbDriver, err := connectDB()

	if err != nil {
		log.Fatal(err)
	}

	defer dbDriver.db.Close()

	server := createAPIServer(":9000", dbDriver)
	server.launch()
}