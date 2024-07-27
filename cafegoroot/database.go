package main

type Product struct {
	Id          int
	Name        string
	Price       int
	Description string
}

func getProducts() []Product {
	return []Product{
		{Id: 1, Name: "Americano", Price: 100, Description: "Espresso, diluted for a lighter experience"},
		{Id: 2, Name: "Cappuccino", Price: 110, Description: "Espresso with steamed milk"},
		{Id: 3, Name: "Espresso", Price: 90, Description: "A strong shot of coffee"},
	}
}
