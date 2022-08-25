package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("%s: %s", "Error. Connection failed.", err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Error. Can't create the channel", err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare("Vacunation", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Error creating queue", err)
	}

	days := getParamethers()

	// Analizamos en el intervalo de dias propuesto
	for day := 1; day <= days; day++ {
		pacients := 0

		fmt.Printf("Day %d \n", day)
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)

		// for minute := 270; minute <= 1200; minute++ {
		for minute := 390; minute <= 1200; minute++ {
			// Para calcular aleatoriamente la llegada de un paciente
			random := r.Intn(100)
			arrived := false

			// if minute < 450 { // Intervalo de tiempo 4:30AM a 7:30AM .31 por minuto
			if minute <= 450 { // Intervalo de tiempo 6:30AM a 7:30AM .31 por minuto
				if random < 31 {
					arrived = true
				}
			} else if minute <= 630 { // Intervalo de tiempo 7:31 AM a 10:30 AM .46 por minuto
				if random < 46 {
					arrived = true
				}
			} else if minute <= 720 { // Intervalo de tiempo 10.31 AM a 12:00 PM .55 por minuto
				if random < 55 {
					arrived = true
				}
			} else if minute <= 810 { // Intervalo de tiempo 12:00 PM a 1:30 PM .0 por minuto
			} else if minute <= 1110 { // Intervalo de tiempo 1:31 PM a 6:30 PM .73 por minuto
				if random < 73 {
					arrived = true
				}
			} else if minute <= 1200 { // Intervalo de tiempo 6:30 PM a 8:00 PM .88 por minuto
				if random < 88 {
					arrived = true
				}
			}

			//Envio de paciente a la cola
			if arrived {
				pacients++
				fmt.Println(fmt.Sprint("Pushed to queue. Day: ", day, ". Hour: ", minuteToHour(minute), ". Minute: ", minute))
				err := channel.Publish("", queue.Name, false, false,
					amqp.Publishing{
						Headers:     nil,
						ContentType: "text/plain",
						Expiration:  "36000",
						Body:        []byte(fmt.Sprint("{\"day\":", day, ",\"minute\":", minute, "}")),
					})
				if err != nil {
					log.Fatalf("%s: %s", "Error sending message", err)
					pacients--
				}
			}
		}
		log.Printf("Finished day with " + fmt.Sprint(pacients) + " pacients.")
	}

	log.Printf("Finished")
}

func getParamethers() int {
	days := 0
	fmt.Print("Days: ")
	fmt.Scanln(&days)
	return days
}

func minuteToHour(minute int) string {
	workday := " AM"
	if minute >= 720 {
		workday = " PM"
	}
	return strconv.Itoa(minute/60) + ":" + strconv.Itoa(minute%60) + workday
}
