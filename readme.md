Commands to implement:

# Commands

* **PING**

  * What it does: checks connection liveness between client and server.
  * Typical reply: `PONG`.
  * Use case: health check, latency measurement.
  * Complexity: O(1).
  * Notes: `PING <message>` echoes the message back.

* **ECHO**

  * What it does: returns the given string.
  * Example: `ECHO "hello"` → `"hello"`.
  * Use case: test round-trip, measure latency, simple debugging.
  * Complexity: O(1).

* **SET**

  * What it does: assign a value to a key (creates or overwrites).
  * Syntax examples: `SET key value` ; `SET key value EX seconds` ; `SET key value PX milliseconds` ; `SET key value NX` (only set if key doesn't exist) ; `SET key value XX` (only set if key exists).
  * Value types: stores a string (binary-safe). Other Redis types are separate (list, set, hash).
  * Return: `"OK"` on success.
  * Atomicity: single command is atomic.
  * Complexity: O(1) for common cases (amortized).
  * Notes: expiration options (EX/PX) set TTL; NX/XX are useful for locks & CAS patterns.

* **GET**

  * What it does: retrieves the string value of a key.
  * Example: `GET key` → `value` or `nil` if key missing.
  * Error: type error if the key holds a non-string data type.
  * Complexity: O(1).
  * Notes: GET does not change TTL unless using special commands like `GETEX`.

* **EXISTS**

  * What it does: checks whether a key (or keys) exist.
  * Behavior:

    * Older Redis versions returned `1` (exists) or `0` (does not).
    * Modern Redis accepts multiple keys and returns the count of keys that exist (integer ≥ 0).
  * Example: `EXISTS a b c` → `2` if two of them exist.
  * Complexity: O(N) when called with N keys (checks each key); O(1) for a single key.
  * Notes: useful for conditional flows; beware of race conditions in concurrent clients.

* **DEL**

  * What it does: deletes one or more keys.
  * Return: number of keys actually removed (integer).
  * Example: `DEL key1 key2` → `1` if only one existed and was removed.
  * Complexity: O(N) with N = number of keys to delete (plus cost proportional to the size of data removed in some implementations).
  * Notes: deleting many large keys can be costly; `UNLINK` is available in newer Redis to unlink asynchronously.

* **INCR**

  * What it does: increments the integer value stored at key by 1.
  * Behavior:

    * If key does not exist → key is set to `0` then incremented (result `1`).
    * If key holds a string that represents an integer → increments and stores new value.
    * If value is not an integer-like string → returns an error.
  * Return: the new value (integer).
  * Atomicity: atomic — safe under concurrency for counters.
  * Complexity: O(1).
  * Notes: use `INCRBY` to add arbitrary integer; use `INCRBYFLOAT` for floating increments.

* **DECR**

  * What it does: decrements the integer value stored at key by 1.
  * Behavior: mirrors INCR (creates key with 0 then `-1` if missing).
  * Return: the new value (integer).
  * Atomicity: atomic.
  * Complexity: O(1).
  * Notes: use `DECRBY` for custom step.

* **LPUSH**

  * What it does: insert one or more values at the head (left) of a list stored at key.
  * Behavior:

    * If key does not exist → a new list is created.
    * If key exists but is not a list → type error.
  * Return: new length of the list (integer).
  * Order: `LPUSH key a b c` results in list `c, b, a, ...` (last argument becomes head).
  * Complexity: O(1) for each pushed element (amortized).
  * Use cases: implement stacks (LIFO), queues (with `RPOP`), or producer-consumer patterns.

* **RPUSH**

  * What it does: insert one or more values at the tail (right) of a list.
  * Behavior: same creation/type rules as LPUSH.
  * Return: new list length.
  * Order: `RPUSH key a b c` results in list `..., a, b, c` (pushes to tail).
  * Complexity: O(1) per element (amortized).
  * Use cases: implement FIFO queues (with `LPOP`) and append operations.

* **SAVE**

  * What it does: synchronously save the dataset to disk (create RDB snapshot).
  * Behavior: blocks the server during the save operation until snapshot finishes.
  * Return: `"OK"` on success.
  * Complexity: depends on dataset size — can be expensive and causes pause.
  * Notes:

    * `BGSAVE` is the non-blocking alternative (forks a child to write snapshot).
    * Modern deployments rely on `AOF` or periodic snapshots; `SAVE` is rarely used in production because it blocks clients.

# Extra practical notes

* **Type errors:** many commands will return a type error if the key exists but holds a different Redis data type.
* **Atomicity:** single Redis commands are atomic; use transactions (`MULTI`/`EXEC`) or Lua scripts for multi-step atomic operations.
* **Performance:** most basic commands are O(1); some commands scale with input size (check docs if doing bulk ops).
* **Persistence:** `SET`, `INCR`, `LPUSH` and others change in-memory state; whether those changes survive a restart depends on your persistence settings (RDB snapshots, AOF).
* **Concurrency:** `INCR`/`DECR` are safe for counters under concurrency because they are atomic.

If you want, I can now:

* give short concrete examples showing request → response for each command, or
* show typical error cases and how to handle them (one command at a time). Which would you like?


---


Run: go run main.go


Connect using: telnet localhost 6379