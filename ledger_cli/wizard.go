package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

var commandWizardJournal = cli.Command{
	Name:      "journal",
	Usage:     "creates and submits a single transaction",
	ArgsUsage: "[]",
	Description: `
`,
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {

		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Journal Entry Wizard")
		fmt.Println("--------------------")

		fmt.Print("Enter the date (yyyy-mm-dd): ")
		datetext, _ := reader.ReadString('\n')
		date, err := time.Parse("2006-01-02", strings.TrimSpace(datetext))
		if err != nil {
			panic(err)
		}

		fmt.Print("Enter the Journal Descripion: ")
		desc, _ := reader.ReadString('\n')
		fmt.Println("")

		count := 0

		var transactionLines []Account

		for {
			count += 1
			fmt.Printf("Line item #%d\n", count)
			fmt.Print("Enter the line Descripion: ")
			lineDesc, _ := reader.ReadString('\n')

			fmt.Print("Enter the Account: ")
			lineAccount, _ := reader.ReadString('\n')

			fmt.Print("Enter the Amount: ")
			var i int64
			_, err := fmt.Scanf("%d", &i)
			if err != nil {
				panic(err)
			}
			lineAmount := big.NewRat(i, 1)

			transactionLines = append(transactionLines, Account{
				Name:        lineAccount,
				Description: lineDesc,
				Balance:     lineAmount,
			})

			fmt.Print("Would you like to enter more line items? (n to stop): ")
			exit, _ := reader.ReadString('\n')
			fmt.Println("")
			if strings.ToLower(strings.TrimSpace(exit)) == "n" {
				fmt.Println("")
				break
			}
		}

		req := &Transaction{
			Date:           date,
			Payee:          desc,
			AccountChanges: transactionLines,
			Signature:      "stuff",
		}

		fmt.Printf("%v\n", req)

		err = Send(req)
		if err != nil {
			log.Fatalf("could not send: %v", err)
		}

		return nil
	},
}