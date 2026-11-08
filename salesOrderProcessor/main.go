package main

import (
	"fmt"
)

var store = map[string]float32{
	"t-shirt": 19.99,
	"pants":   59.99,
	"coat":    199.99,
	"shoes":   159.99,
}

const saleActive = true
const discountRate = 0.1

func addItemToCart(itemCode string) float32 {
	price, found := store[itemCode]

	if !found {
		fmt.Println("Item not found " + itemCode)
		return 0.0
	}

	fmt.Printf("Item '%s' added to cart. Price: $%.2f\n", itemCode, price)
	return price
}

func main() {

	order := []string{
		"t-shirt", "pants", "coat", "shoes", "hat",
	}

	var subtotal float32

	for _, item := range order {
		subtotal += addItemToCart(item)
	}

	var discountAmount float32

	if saleActive {
		discountAmount = subtotal * discountRate
		fmt.Printf("Order subtotal: $%.2f\n", subtotal)
		fmt.Printf("Discount applied: $%.2f\n", discountAmount)
		fmt.Printf("Order total: $%.2f\n", subtotal-discountAmount)
	} else {
		fmt.Printf("Order total: $%.2f", subtotal)
	}
}
