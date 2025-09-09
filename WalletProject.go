package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// структуры

type Wallet struct {
	Name string
	Quantity map[string]float64
}

type Asset struct {
	Ticker string
	Price float64
}

type Transaction struct {
	From string
	To string
	Ticker string
	Quantity float64
}

type TotalAsset struct {
	Ticker string
	TotalPrice float64
}

var totalAssets []TotalAsset // слайс структур для хранения и сортировки стоимости активов в кошельке

// методы

func (w Wallet) ShowWalletInfo() {
	fmt.Printf("Кошелек: ")
	fmt.Printf("Имя: %s", w.Name)
	fmt.Printf("\nБаланс:\n")
	for key, value := range w.Quantity {
		fmt.Printf("Тикер актива: %s, Количество: %.2f\n", key, value)
	}
	fmt.Printf("-------------------------------")
}

func (asset Asset) PrintAssetInfo() string {
	result := fmt.Sprintf("Актив - %s, Цена - %.2f", asset.Ticker, asset.Price)
	return result
}

func (w Wallet) CalcBalance(prices []Asset) float64 {
	total := 0.0
	for _, asset := range prices {
		quantity := w.Quantity[asset.Ticker]
		total += quantity * asset.Price
	}
	return total
}

func (w *Wallet) Deposit(ticker string, quantity float64) {
	w.Quantity[ticker] += quantity
}

func (w *Wallet) Withdraw(ticker string, quantity float64) {
	defer func() {
		error := recover()
		if error != nil {
			fmt.Println("Недостаточно Средств!")
		}
	}()
	if w.Quantity[ticker] < quantity {
		panic("Error")
	}
	w.Quantity[ticker] -= quantity
}

func (w Wallet) PrintAssets() {
	for key, value := range(w.Quantity) {
		fmt.Printf("Актив: %s, Кол-во: %.2f\n", key, value)
	}
}

func (t Transaction) PrintTransaction() string {
	result := fmt.Sprintf("%s -> %s: Актив: %s, Кол-во: %.2f", t.From, t.To, t.Ticker, t.Quantity)
	return result
}

func GetUserTransaction(name string, txs []Transaction) []Transaction {
	var result []Transaction

	for _, t := range txs {
		if t.From == name || t.To == name {
			result = append(result, t)
		}
	}
	return result
}

func Exchange(quantity, rate float64) float64 {
	defer func() {
		error := recover()
		if error != nil {
			fmt.Println("Неверный курс!")
		}
	}()
	if rate <= 0 {
		panic("Error")
	}
	return quantity * rate
}

func GetPrice(ticker string, assets []Asset) float64 {
	for _, asset := range assets {
		if asset.Ticker == ticker {
			return asset.Price
		}
	}
	return 0.0
}

func (w *Wallet) Buy(tickerToBuy string, quantity float64, tickerToPay string, assets []Asset) error {
	tickerToBuyPrice := GetPrice(tickerToBuy, assets)
	tickerToPayPrice := GetPrice(tickerToPay, assets)

	if tickerToBuyPrice == 0 || tickerToPayPrice == 0 {
		return fmt.Errorf("Ошибка! Проверте корректность тикера")
	}

	// переводим стоимость и наш баланс в рубли
	costInRubles := tickerToBuyPrice * quantity
	payInRubles := tickerToPayPrice * w.Quantity[tickerToPay]

	// переводим стоимость в валюту за которую покупаем 
	CostToPayCurrency := (costInRubles * w.Quantity[tickerToPay]) / payInRubles
	// =  tickerToBuyPrice * quantity * w.Quantity[tickerToPay] / tickerToPayPrice * w.Quantity[tickerToPay] = 
	// = tickerToBuyPrice * quantity / tickerToPayPrice - та же самая формула что и в примере с гпт :)

	// сравниваем оба в валюте покупки
	if w.Quantity[tickerToPay] < CostToPayCurrency {
		return fmt.Errorf("Недостаточно средств: Необходимо %.2f %s, имееется %.2f %s", CostToPayCurrency, tickerToPay, w.Quantity[tickerToPay], tickerToPay)
	}

	// ту валюту которую купили, добавляем к балансу кошелька, 
	// а ту которую потратили - вычитаем
	w.Quantity[tickerToBuy] += quantity
	w.Quantity[tickerToPay] -= CostToPayCurrency

	fmt.Printf("Куплено %.2f %s за %.2f %s\n", quantity, tickerToBuy, CostToPayCurrency, tickerToPay)
	return nil
}

func ShowTransactionsByUser(trx []Transaction) {
	trxUser := make(map[string][]Transaction)
	for _, t := range trx {
		trxUser[t.From] = append(trxUser[t.From], t)
		trxUser[t.To] = append(trxUser[t.To], t)
	}

	for user, trx := range trxUser {
		fmt.Printf("Пользователь %s:\n", user)
		for _, t := range trx {
			fmt.Println(" -", t.PrintTransaction())
		}
		fmt.Println()
	}
}

func (w Wallet) CalcTotalValue(prices []Asset) {
	fmt.Printf("Общая стоимость активов: %.2f", w.CalcBalance(prices))
}

func runMenu(trx []Transaction, w *Wallet, assets []Asset) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("0. Информация о кошельке")
		fmt.Println("1. Обмен валют")
		fmt.Println("2. Пополнение актива в кошельке (в единицах актива)")
		fmt.Println("3. Снять актив из кошелька (в единицах актива)")
		fmt.Println("4. Показать транзакции по пользователям")
		fmt.Println("5. Купить актив")
		fmt.Println("6. Показать все активы в кошельке")
		fmt.Println("7. Посчитать общую стоимость активов")
		fmt.Println("exit. Выйти")
		fmt.Println()

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "0":
			w.ShowWalletInfo()
		case "1":
			fmt.Println("Результат обмена:", Exchange(12_900, 0.12))
		case "2":
			fmt.Println("Введите тикер актива: ")
			scanner.Scan()
			ticker := strings.TrimSpace(scanner.Text())
			fmt.Println("Введите количество: ")
			var quantity float64
			fmt.Scanln(&quantity)
			w.Deposit(ticker, quantity)
		case "3":
			fmt.Println("Введите тикер актива: ")
			scanner.Scan()
			ticker := strings.TrimSpace(scanner.Text())
			fmt.Println("Введите количество: ")
			var quantity float64
			fmt.Scanln(&quantity)
			w.Withdraw(ticker, quantity)
		case "4":
			ShowTransactionsByUser(trx)
		case "5":
			fmt.Println("Введите тикер актива, который хотите купить: ")
			scanner.Scan()
			tickerToBuy := strings.TrimSpace(scanner.Text())
			fmt.Println("Введите количество: ")
			var quantity float64
			fmt.Scanln(&quantity)
			fmt.Println("Введите тикер актива, за который будет произведена покупка: ")
			scanner.Scan()
			tickerToPay := strings.TrimSpace(scanner.Text())

			err := w.Buy(tickerToBuy, quantity, tickerToPay, assets)
			if err != nil {
				fmt.Println("Ошибка:", err)
			}
		case "6":
			w.PrintAssets()
		case "7":
			w.CalcTotalValue(assets)
		case "exit":
			fmt.Println("Всего доброго!")
			return 
		default:
			fmt.Println("Неверная команда!")
		}

	}
}

func main() {
	wallet1 := Wallet{
		Name: "Anton", 
		Quantity: map[string]float64{
			"BTC": 2.1,
			"ETH": 5,
			"USDT": 3500.21,
		},
	}
	assets := []Asset{
		{Ticker: "BTC", Price: 8_589_330.55,},
		{Ticker: "ETH", Price: 201_971.84,},
		{Ticker: "USDT", Price: 78.84,},
	}

	trx := []Transaction{
		{From: "Anton", To: "Elena", Ticker: "Rub", Quantity: 12000},
		{From: "Max", To: "Anton", Ticker: "USDT", Quantity: 500.95},
		{From: "Bob", To: "Elena", Ticker: "Rub", Quantity: 230_123},
		{From: "Max", To: "Elena", Ticker: "BTC", Quantity: 1.5},
	}

	runMenu(trx, &wallet1, assets)
}