package transfer

const createTransfer = `INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
)`

const createEntry = `INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)`

const addAccountBalance = `UPDATE accounts
SET balance = balance + $1
WHERE id = $2`
