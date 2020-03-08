package scheduler

const insertNewEventSQL = `INSERT INTO events ("when", message) VALUES (?, ?)`

const selectEventToPrecessSQL = `
SELECT id, message
FROM events
WHERE "when" <= ?
  AND not done
ORDER BY "when" ASC
LIMIT 1
FOR UPDATE SKIP LOCKED
`

const markEventAsDoneSQL = `UPDATE events SET done = TRUE WHERE id = ?`
