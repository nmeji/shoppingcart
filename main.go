package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	sim "github.com/nmeji/exam/amaysim"
)

type Command struct {
	Usage       string
	Description string
	Run         func(string)
}

var (
	cart     sim.Cart
	commands map[string]*Command
	catalog  *sim.Products
)

func init() {
	commands = make(map[string]*Command)
	commands["help"] = &Command{
		Usage:       "help [command]",
		Description: "Show help",
		Run:         ShowHelp,
	}
	commands["catalog"] = &Command{
		Usage:       "catalog",
		Description: "Show product catalog",
	}
	commands["add"] = &Command{
		Usage:       "add item_id [quantity]",
		Description: "Adds item(s) to cart",
		Run:         AddItem,
	}
	commands["promo"] = &Command{
		Usage:       "promo code",
		Description: "Applies promo code",
		Run:         UsePromo,
	}
	commands["cart"] = &Command{
		Usage:       "cart",
		Description: "Show cart summary",
		Run:         ViewCart,
	}
	commands["remove"] = &Command{
		Usage:       "remove item_id [quantity]",
		Description: "Removes item(s) from cart",
		Run:         RemoveItem,
	}
	commands["exit"] = &Command{
		Usage:       "exit",
		Description: "Terminates the application",
	}
}

func ShowHelp(s string) {
	if cmd, ok := commands[s]; !ok {
		if s != "" {
			fmt.Printf("Command '%s' is not found\n\n", s)
		}
		fmt.Println("List of all commands:")
		fmt.Printf("%s\t\t\t\t%s\n", commands["catalog"].Usage, commands["catalog"].Description)
		fmt.Printf("%s\t\t%s\n", commands["add"].Usage, commands["add"].Description)
		fmt.Printf("%s\t%s\n", commands["remove"].Usage, commands["remove"].Description)
		fmt.Printf("%s\t\t\t\t%s\n", commands["cart"].Usage, commands["cart"].Description)
		fmt.Printf("%s\t\t\t%s\n", commands["promo"].Usage, commands["promo"].Description)
		fmt.Printf("%s\t\t\t%s\n", commands["help"].Usage, commands["help"].Description)
		fmt.Printf("%s\t\t\t\t%s\n", commands["exit"].Usage, commands["exit"].Description)
		fmt.Println()
	} else {
		fmt.Printf("Usage:\t%s\n", cmd.Usage)
		fmt.Printf("    %s\n", cmd.Description)
		fmt.Println()
	}
}

func ReadIntArgs(s string) (int, int) {
	args := strings.Split(s, " ")
	itemID, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		panic(err)
	}
	quantity := 1
	if len(args) > 1 {
		qty, err := strconv.ParseInt(args[1], 10, 32)
		quantity = int(qty)
		if err != nil {
			panic(err)
		}
	}
	return int(itemID), quantity
}

func AddItem(s string) {
	itemID, quantity := ReadIntArgs(s)
	product := catalog.Catalog[itemID-1]
	cart.Add(sim.CartItem{Quantity: quantity, SubTotal: 0.0, Product: product})
}

func UsePromo(s string) {
	cart.AddPromoCode(s)
}

func RemoveItem(s string) {
	itemID, quantity := ReadIntArgs(s)
	product := catalog.Catalog[itemID-1]
	cart.Remove(sim.CartItem{Quantity: quantity, SubTotal: 0.0, Product: product})
}

func ViewCart(s string) {
	added := cart.ItemsAdded()
	if len(added) == 0 {
		fmt.Println("You haven't added any item to your cart.")
		return
	}
	totalPrice := cart.Total()
	totalItems := cart.Items()
	fmt.Println("You've added,")
	for _, item := range added {
		fmt.Printf("%d x %s\n", item.Quantity, item.Name)
	}
	fmt.Printf("\nTotal Price is $%.02f\n\n", totalPrice)
	fmt.Println("Your total item(s):")
	for _, item := range totalItems {
		fmt.Printf("%d x %s\n", item.Quantity, item.Name)
	}
	fmt.Println()
}

func main() {
	f, err := os.Open("product-catalog.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	catalog, err = sim.InitProducts(bufio.NewReader(f))
	if err != nil {
		panic(err)
	}

	shoppingCart := sim.ShoppingCart{
		PricingRules: sim.AllRules,
		Added:        make(map[string]*sim.CartItem),
		PromoCode:    make(map[string]int),
	}
	cart = shoppingCart
	cont := true

	commands["catalog"].Run = func(a string) {
		catalog.PrintCatalog()
	}
	commands["exit"].Run = func(a string) {
		cont = false
	}
	cliReader := bufio.NewReader(os.Stdin)
	for cont {
		fmt.Print("> ")
		line, _ := cliReader.ReadString('\n')
		line = strings.TrimSpace(line[0 : len(line)-1])
		key := line
		args := ""
		if sep := strings.IndexRune(line, ' '); sep > -1 {
			key = line[0:sep]
			args = line[sep+1:]
		}
		if cmd, ok := commands[key]; ok {
			cmd.Run(args)
		} else {
			fmt.Printf("Unknown command '%s', use 'help' to show available commands\n", key)
		}
	}
}
