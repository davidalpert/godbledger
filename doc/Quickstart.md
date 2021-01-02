On Ubuntu 18.04

**Download the latest release and run the server**
```
$ wget https://github.com/darcys22/godbledger/releases/download/v0.5.0/godbledger-linux-x64-v0.5.0.tar.gz
$ tar -xvzf godbledger-linux-x64-v0.5.0.tar.gz 
$ cd ~
$ ln -s godbledger-linux-x64-v0.5.0 godbledger

$ ~/godbledger/godbledger init
$ ~/godbledger/godbledger
```
This will start the server running a sqlite database backend. `godbledger init` will create a config file in `~/.ledger` and the sqlite file will be in `~/.ledger/ledgerdata`.

**Send transactions to the server**
This needs to be run in a separate terminal with the GoDBLedger server running. `ledger_cli` will communicate with the running `godbledger`

There are several ways to do so:

*Using JSON*
```
$ ~/godbledger/ledger_cli jsonjournal '{"Payee":"Darcy Financial","Date":"2019-06-30T00:00:00Z","AccountChanges":[{"Name":"Asset:Cash","Description":"Cash is better","Currency":"USD","Balance":"100"},{"Name":"Revenue:Sales","Description":"Income is good yo","Currency":"USD","Balance":"-100"}]}'
```

*Using the Wizard*
```
$ ~/godbledger/ledger_cli journal
Journal Entry Wizard
--------------------
Enter the date (yyyy-mm-dd): 2019-06-30
Enter the Journal Descripion: Get Money Get Paid!

Line item #1
Enter the line Descripion: Income is good yo 
Enter the Account: Revenue:Sales
Enter the Amount: -1000
Would you like to enter more line items? (n to stop): 

Line item #2
Enter the line Descripion: Cash is better
Enter the Account: Asset:Cash
Enter the Amount: 1000
Would you like to enter more line items? (n to stop): n


&{Get Money Get Paid!
 2019-06-30 00:00:00 +0000 UTC [{Revenue:Sales
 Income is good yo
 -1000/1} {Asset:Cash
 Cash is better
 1000/1}]}
```
This will send an example transaction to the server. If all goes well you should have a transaction in your database now for 1000 income and increasing the cash account also by 1000

**To view this transaction run**
```
$ ~/godbledger/reporter trialbalance

     ACCOUNT    | BALANCE AT 20 FEBRUARY 2020
----------------+------------------------------
  Asset:Cash    |                        1000
                |
  Revenue:Sales |                       -1000
                |
```

## Installation Tutorial on Youtube
https://youtu.be/D8vDRxGn5v8
