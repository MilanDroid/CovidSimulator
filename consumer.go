package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type PacientInfo struct {
	Day    int
	Minute int
}

type StationBussy struct {
	Station       int
	AttentionTime int
}

type DayInfo struct {
	Day                       int
	NotAttendedPacients       int
	AttendedPacients          int
	AttendedPacientsMorning   int
	AttendedPacientsAfternoon int
	AttentionMedian           int
	AttentionMedianMorning    int
	AttentionMedianAfternoon  int
	WaitingTimeMedian         int
}

func getParamethers() (int, int) {
	stations := 0
	resources := 0

	fmt.Print("Stations [1-15]: ")
	fmt.Scanln(&stations)

	fmt.Print("Resources: ")
	fmt.Scanln(&resources)

	return stations, resources
}

func getResourcesByTime(stations int, resources int) (int, int) {
	difference := resources - stations
	if difference < 0 {
		return resources, 0
	} else if difference > stations {
		return stations, stations
	} else {
		return stations, difference
	}
}

func minuteToHour(minute int) string {
	workday := " AM"
	if minute >= 720 {
		workday = " PM"
	}
	return strconv.Itoa(minute/60) + ":" + strconv.Itoa(minute%60) + workday
}

func divide(a int, b int) int {
	if b > 0 {
		return a / b
	} else {
		return 0
	}
}

func main() {
	const Filename = "response.json"
	const MorningTurnEndTime = 720
	const LunchEndTime = 810
	const AfternoonTurnEndTime = 1200

	maxAttentionTime := 10
	minAttentionTime := 5

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("%s: %s", "Error. Connection failed", err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Error. Can't connect to the channel", err)
	}
	defer channel.Close()

	channel.QueueDeclare("Vacunation", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Error creating queue", err)
	}

	stations, resources := getParamethers()
	afternoonResources, morningResources := getResourcesByTime(stations, resources)

	stationsQueue := make([]StationBussy, 0, stations)
	less := func(i, j int) bool {
		return stationsQueue[i].AttentionTime < stationsQueue[j].AttentionTime
	}

	isMorning := false
	isAfternoon := false
	attending := true
	currentTime := 0
	waitingTimeGlobal := 0
	currentDay := DayInfo{
		Day:                       0,
		NotAttendedPacients:       0,
		AttendedPacients:          0,
		AttendedPacientsMorning:   0,
		AttendedPacientsAfternoon: 0,
		AttentionMedian:           0,
		AttentionMedianMorning:    0,
		AttentionMedianAfternoon:  0,
		WaitingTimeMedian:         0,
	}
	daysReport := make([]DayInfo, 0, 365)

	for station := 1; station <= morningResources; station++ {
		stationInfo := StationBussy{Station: station, AttentionTime: 0}
		stationsQueue = append(stationsQueue, stationInfo)
	}

	messages, err := channel.Consume("Vacunation", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Error creating consumer channel", err)
	}

	go func() {
		for message := range messages {
			var info PacientInfo
			json.Unmarshal([]byte(string(message.Body)), &info)

			if currentDay.Day == 0 {
				currentDay.Day = info.Day
				isMorning = true
				isAfternoon = false
				fmt.Println("Day: " + fmt.Sprint(currentDay.Day))
			}
			if currentDay.Day != info.Day {
				currentDay.AttendedPacientsAfternoon = currentDay.AttendedPacients - currentDay.AttendedPacientsMorning
				daysReport = append(daysReport, currentDay)

				currentDay.Day = info.Day
				currentDay.NotAttendedPacients = 0
				currentDay.AttendedPacients = 0
				currentDay.AttendedPacientsMorning = 0
				currentDay.AttendedPacientsAfternoon = 0
				currentDay.AttentionMedian = 0
				currentDay.AttentionMedianMorning = 0
				currentDay.AttentionMedianAfternoon = 0
				currentDay.WaitingTimeMedian = 0

				currentTime = info.Minute
				waitingTimeGlobal = 0
				isMorning = true
				isAfternoon = false
				attending = true

				// Estableciendo estaciones de la maniana en base a los recursos disponibles
				stationsQueue = stationsQueue[:0]
				for station := 1; station <= morningResources; station++ {
					stationInfo := StationBussy{Station: station, AttentionTime: 0}
					stationsQueue = append(stationsQueue, stationInfo)
				}
				fmt.Println("Day: " + fmt.Sprint(currentDay.Day))
			}

			if currentTime == 0 || currentTime < info.Minute {
				currentTime = info.Minute
			}

			// Cambio de turno a la tarde, terminar de atender pacientes en estaciones
			if currentTime >= MorningTurnEndTime && isMorning && attending {
				for _, station := range stationsQueue {
					if station.AttentionTime != 0 {
						station.AttentionTime = 0
					}
				}

				// Estableciendo estaciones de la tarde en base a los recursos disponibles
				stationsQueue = stationsQueue[:0]
				for station := 1; station <= afternoonResources; station++ {
					stationInfo := StationBussy{Station: station, AttentionTime: 0}
					stationsQueue = append(stationsQueue, stationInfo)
				}

				currentDay.AttendedPacientsMorning = currentDay.AttendedPacients
				currentTime = LunchEndTime
				isMorning = false
				isAfternoon = true
				fmt.Println("Morning turn finished. Attended: " + fmt.Sprint(currentDay.AttendedPacientsMorning))
			}
			// Finaliza el dia, termina de atender pacientes en las estaciones
			if currentTime >= AfternoonTurnEndTime && isAfternoon && attending {
				for _, station := range stationsQueue {
					if station.AttentionTime != 0 {
						station.AttentionTime = 0
					}
				}
				currentDay.AttendedPacientsAfternoon = currentDay.AttendedPacients - currentDay.AttendedPacientsMorning
				currentDay.AttentionMedian = divide(currentDay.AttendedPacients, (morningResources + afternoonResources))
				currentDay.AttentionMedianMorning = divide(currentDay.AttendedPacientsMorning, morningResources)
				currentDay.AttentionMedianAfternoon = divide(currentDay.AttendedPacientsAfternoon, afternoonResources)
				currentDay.WaitingTimeMedian = divide(waitingTimeGlobal, currentDay.AttendedPacients)
				attending = false
				fmt.Println("Afternoon turn finished. Attended: " + fmt.Sprint(currentDay.AttendedPacientsAfternoon))
			}

			if attending && len(stationsQueue) > 0 {
				source := rand.NewSource(time.Now().UnixNano())
				r := rand.New(source)
				attentionTime := r.Intn(maxAttentionTime-minAttentionTime) + minAttentionTime

				sort.Slice(stationsQueue, less)
				spentTimeOnCycle := stationsQueue[0].AttentionTime
				currentTime = currentTime + spentTimeOnCycle

				fmt.Println("Stations status start: " + fmt.Sprint(stationsQueue))
				fmt.Println("SpentTimeOnCycle: " + fmt.Sprint(spentTimeOnCycle))
				for station := range stationsQueue {
					if stationsQueue[station].AttentionTime != 0 {
						stationsQueue[station].AttentionTime = stationsQueue[station].AttentionTime - spentTimeOnCycle
					}
				}

				stationsQueue[0].AttentionTime = attentionTime
				fmt.Println("Stations status end: " + fmt.Sprint(stationsQueue))
				currentDay.AttendedPacients++

				waitingTime := (currentTime + attentionTime) - info.Minute
				waitingTimeGlobal = waitingTimeGlobal + waitingTime
				fmt.Println("Attended. Day: " + fmt.Sprint(info.Day) + ", went at: " + fmt.Sprint(info.Minute) + " - " + minuteToHour(info.Minute) + ". Attention time: " + fmt.Sprint(attentionTime) + ". End time: " + minuteToHour(currentTime+attentionTime) + ". Waiting time: " + fmt.Sprint(waitingTime) + " minutes. Station: " + fmt.Sprint(stationsQueue[0].Station) + ". Current time: " + fmt.Sprint(currentTime) + ". Total attended: " + fmt.Sprint(currentDay.AttendedPacients))
			} else {
				currentDay.NotAttendedPacients++
				fmt.Println("Not attended. Day: " + fmt.Sprint(info.Day) + ", went at: " + fmt.Sprint(info.Minute) + " - " + minuteToHour(info.Minute) + ". Total not attended: " + fmt.Sprint(currentDay.NotAttendedPacients))
			}

			message.Ack(true)
		}

		currentDay.AttendedPacientsAfternoon = currentDay.AttendedPacients - currentDay.AttendedPacientsMorning
		daysReport = append(daysReport, currentDay)
		fmt.Println("Stations: " + fmt.Sprint(stations) + ". Resources: " + fmt.Sprint(resources))
		fmt.Println("Used resources: " + fmt.Sprint(afternoonResources+morningResources) + ". Unused resources: " + fmt.Sprint(resources-(afternoonResources+morningResources)))
		fmt.Println("Results matrix: " + fmt.Sprint(daysReport))
		//MarshalIndent
		empJSON, err := json.MarshalIndent(daysReport, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("Report: %s\n", string(empJSON))

		// f, err := os.Create("response-" + fmt.Sprint(stations) + "-" + fmt.Sprint(resources) + ".json")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer f.Close()

		// _, err2 := f.WriteString(string(empJSON))
		// if err2 != nil {
		// 	log.Fatal(err2)
		// }
		// fmt.Printf("Output: response-" + fmt.Sprint(stations) + "-" + fmt.Sprint(resources) + ".json")
	}()

	// Si ponemos un scan, evitaremos que el programa se cierre y recibiremos los nuevos mensajes
	fmt.Scanln()
}
