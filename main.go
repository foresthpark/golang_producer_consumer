package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const numberOfPizzas = 10

var pizzasMade, pizzasFailed, totalPizzas int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error

}

type PizzaOrder struct {
	pizzaNumber int
	message string
	success bool
}

func (p *Producer) Close() error {
	ch := make(chan error)

	p.quit <- ch

	return <- ch

}

func pizzeria(pizzaMaker *Producer) {
	var i int 
	// keep track of which pizza we're making

	// run forever until it receives a quit signal
	// try to make pizzas
	for {
		// try to make a pizza
		currentPizza := makePizza(i)

		// decision structure
		if currentPizza != nil {
			i = currentPizza.pizzaNumber

			select {
				// we tried to make a pizza(we sent something to the data channel)
				case pizzaMaker.data <- *currentPizza:
				
				// we want to quit, so we send a quit signal to the quit channel (a  chan error)
				case quitChan := <- pizzaMaker.quit:
					close(pizzaMaker.data)
					close(quitChan)
					return
			}
		}



		
	}
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= numberOfPizzas {
		delay := rand.Intn(5) + 1

		fmt.Printf("Received order number #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd > 5 {
			pizzasFailed++
		} else {
			pizzasMade++
			
		}

		totalPizzas++

		fmt.Printf("Making pizza #%d, it will take %d seconds to complete.\n", pizzaNumber, delay)

		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza number #%d ***", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza number #%d ***", pizzaNumber)

		} else {
			success = true
			msg = fmt.Sprintf("*** Pizza number #%d is ready! ***", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message: msg,
			success: success,
		}

		return &p

	} 
	
	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func main() {
	 // seed the random number generator
	 rand.Seed(time.Now().UnixNano())
	 
	 // Print message that the program has started
	 color.Cyan("The Pizzeria is open for business!")
	 color.Cyan("----------------------------------")

	 // create a producer
	 pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	 }

	 // run the producer in the background as a go routine
	 go pizzeria(pizzaJob)

	 // create and run consumer
	 for i := range pizzaJob.data {
		if i.pizzaNumber <= numberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is ready!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("Order #%d failed!", i.pizzaNumber)
			}
		} else {
			color.Cyan("Done making pizzas...")
			if err := pizzaJob.Close(); err != nil {
				color.Red("Error closing the pizzeria: %s", err) 
			}
			
		}
	 }

	 // print out the ending message
	 color.Cyan("All order DONE")
	 color.Yellow("We made %d pizzas, but failed to make %d pizzas, with a total of %d attempts", pizzasMade, pizzasFailed, totalPizzas)
	 color.Cyan("----------------------------------")
}