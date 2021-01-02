The following CREATE TABLE commands show the SQL code from GoDBLedger that creates the general ledger structure.

### Transactions Table

Transactions are transfers of commodities between accounts and are the main table to record everything in the general ledger. The actual heavy work is starts here with various other tables each linking to the transaction. A transaction doesn't contain the recorded amount, that should be recorded in the split table. 

```
CREATE TABLE IF NOT EXISTS transactions (
		transaction_id VARCHAR(255) NOT NULL,
		postdate DATETIME NOT NULL,
		brief VARCHAR(255),
		poster_user_id VARCHAR(255),
		PRIMARY KEY(transaction_id),
                FOREIGN KEY (poster_user_id) REFERENCES users (user_id) ON DELETE RESTRICT
	);
```
| Description    | Type   | Comment                                                                                                                                      |
|----------------|--------|----------------------------------------------------------------------------------------------------------------------------------------------|
| transaction_id | String | Identifier of the transaction, GoDBLedger uses github.com/rs/xid to generate the identifiers for all transactions                            |
| postdate       | Date   | Date that the transaction was received by GoDBLedger                                                                                         |
| brief          | String | Brief description of the transaction, for larger descriptions the description should be saved to transaction_body table                      |
| poster_user_id | String | Foreign Key referencing the User who posted the transaction, currently GoDBLedger only has a single user "MainUser" which is used by default |


### Splits Table

Splits tie an amount of a commodity to the account in a transaction. Must be linked to a transaction using the `transaction_id` key. The amount is positive for a Debit and negative for a Credit. Its recorded as an integer with the intention to protect against floating point rounding errors, the user should reference the currency table to identify how many decimal points the currency has (Dollars will be recorded as their amounts in cents and the currency reference will have "2" as the number of decimals). 
```
CREATE TABLE IF NOT EXISTS splits (
		split_id VARCHAR(255) NOT NULL,
		split_date DATETIME,
		description VARCHAR(255),
		currency VARCHAR(255),
		amount BIGINT,
		transaction_id VARCHAR(255),
		FOREIGN KEY(transaction_id) REFERENCES transactions(transaction_id) ON DELETE CASCADE,
		PRIMARY KEY(split_id)
	);
```
| Description    | Type    | Comment                                                                                                                                                                                                          |
|----------------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| split_id       | String  | Unique Identifier, uses github.com/rs/xid                                                                                                                                                                        |
| split_date     | Date    | transaction date, this is referenced by the GL and will be used by reports such as the profit and loss etc                                                                                                       |
| description    | String  | Optional brief comment on the split, should be used to differentiate the line items in a transaction if they are different. Most of the overall detail however should be saved into the transaction description  |
| currency       | String  | References the Currency Table, Also should be a human readable key ie "USD"                                                                                                                                      |
| amount         | Integer | Positive for Debits, Negative for Credit. Detailing the smallest unit of a currency (Cents) which is then divided by the "Decimals" specified in the currency table.                                             |
| transaction_id | String  | Foreign Key referencing the transactions table. All splits must reference a transaction                                                                                                                          |

### Accounts Table

Accounts are the containing entity for transactions
```
CREATE TABLE IF NOT EXISTS accounts (
		account_id VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		PRIMARY KEY(account_id)
	);
```
| Description | Type   | Comment                                 |
|-------------|--------|-----------------------------------------|
| account_id  | String | User defined identifier for the account |
| name        | String | Display name of the account             |


### Split Accounts Table
Linking the splits to the accounts is the split_accounts table. This enables a split to have multiple accounts although the current implementation only utilises a single account per split.
```
CREATE TABLE IF NOT EXISTS split_accounts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		split_id VARCHAR(255),
		account_id VARCHAR(255),
		FOREIGN KEY(split_id) REFERENCES splits(split_id) ON DELETE CASCADE,
		FOREIGN KEY(account_id) REFERENCES accounts(account_id) ON DELETE CASCADE
	);
```
| Description | Type    | Comment                                                     |
|-------------|---------|-------------------------------------------------------------|
| id          | Integer | Auto generated identifier for the many to many relationship |
| split_id    | String  | Foreign Key linking to the Splits table                     |
| account_id  | String  | Foreign Key linking to the Accounts table                   |

### Currencies Table

Currencies describe the type of asset referenced in the splits. Could include any type of commodity (Gold, Shells, Rum). GoDBLedger has built in a few common currencies for ease of use.
```
CREATE TABLE IF NOT EXISTS currencies (
		name VARCHAR(255) NOT NULL,
		decimals INT NOT NULL,
		PRIMARY KEY(name)
	);
```
| Description | Type    | Comment                                                                                          |
|-------------|---------|--------------------------------------------------------------------------------------------------|
| name        | String  | Unique identifier for the currency, user defined but recommend to follow ISO 4217 Currency Codes |
| decimals    | Integer | Number of decimals that the currency uses (USD = 2 decimals, BTC = 8 decimals)                   |

### Tags Table

Tags are the implementation detail to allow for transactions and accounts to be grouped in a user defined manner. They allow for more complicated analysis and reporting. An example is the "Void" tag that has been implemented by default to allow for transactions to be hidden but not deleted. Ideally accounts would be tagged according to their category in the financials (Asset, Liability etc) which would allow for automatic compiling of financial statements. Additionally transactions would be tagged according to the front end implementation so they can be easily extracted (ie the protocol version).

```
	CREATE TABLE IF NOT EXISTS tags (
		tag_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
		tag_name VARCHAR(100) NOT NULL UNIQUE
	);
```
| Description | Type    | Comment                                                                                                                                         |
|-------------|---------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| tag_id      | Integer | Auto generated identifier for the Tag                                                                                                           |
| tag_name    | String  | User defined Tag, ideally used to filter transactions and accounts. For example GoDBLedger automatically excludes transactions tagged as "Void" |

And the corresponding Many to Many Relationship tables for both transactions and accounts
```	
CREATE TABLE IF NOT EXISTS account_tag (
    account_id VARCHAR(255) NOT NULL,
    tag_id INTEGER NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts (account_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (tag_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    PRIMARY KEY (account_id, tag_id)
	);
CREATE TABLE IF NOT EXISTS transaction_tag (
    transaction_id VARCHAR(255) NOT NULL,
    tag_id INTEGER NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions (transaction_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (tag_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    PRIMARY KEY (transaction_id, tag_id)
	);
```

There is a Transactions Body table to allow for larger amounts of data to be recorded on each transaction. Idea is that a transaction should be able to record everything relevant in the description which might include full contracts which are hundreds of pages long.

```
CREATE TABLE IF NOT EXISTS transactions_body (
		transaction_id VARCHAR(255) NOT NULL,
		body TEXT,
		FOREIGN KEY(transaction_id) REFERENCES transactions(transaction_id) ON DELETE CASCADE
	);
```

Finally there is an entities table which is not utilised. Thought being that everything related to the business the database is referring to is stored here, and also would allow for multiple entities to record their data in the same database. 
```
CREATE TABLE IF NOT EXISTS entities (
		entity_id VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		tag VARCHAR(255),
		type VARCHAR(255),
		description VARCHAR(255),
		PRIMARY KEY(entity_id)
	);`
```

'''Note:''' SQLite3 doesn't support date or timestamp data types, so in SQLite3 databases those fields are represented as <tt>CHAR[8]</tt> and <tt>CHAR[14]</tt> respectively.