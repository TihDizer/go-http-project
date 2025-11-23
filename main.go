package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	var errorsCount uint8
	client := &http.Client{}

	for {
		resp, err := client.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil {
			panic(err)
		}

		body, err := io.ReadAll(resp.Body)

		defer resp.Body.Close()

		bodyString := string(body)
		bodyParts := strings.Split(bodyString, ",")
		bodyFloats := make([]float64, len(bodyParts))
		for i, part := range bodyParts {
			bodyFloats[i], _ = strconv.ParseFloat(part, 64)
		}

		if resp.StatusCode != http.StatusOK || len(bodyParts) < 6 {
			errorsCount++
			continue
		}

		if bodyFloats[0] > 30 {
			fmt.Printf("Load Average is too high: %d\n", int(bodyFloats[0]))
		}

		if bodyFloats[2] > 0.8*bodyFloats[1] {
			fmt.Printf("Memory usage too high: %d\n", int(bodyFloats[2]))
		}

		if bodyFloats[4] > 0.9*bodyFloats[3] {
			fmt.Printf("Free disk space is too low: %d Mb left\n", int(bodyFloats[4]))
		}

		if bodyFloats[6] > 0.9*bodyFloats[5] {
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", int(bodyFloats[6]))
		}

		if errorsCount > 2 {
			fmt.Println("Unable to fetch server statistic")
			break
		}

		fmt.Println(resp)
		time.Sleep(1 * time.Second)
	}
}
