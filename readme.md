# Redis Server in Go

A Redis server implementation in Go with support for Strings, Lists, Sets, and Hashes.

## Features

- RESP (Redis Serialization Protocol) compatible
- Thread-safe operations
- 25 Redis commands across 4 data types
- Works with any Redis client (redis-cli, client libraries)

## Quick Start

### 1. Start the Server

```bash
go run .
```

The server will start on `localhost:6379`.

### 2. Connect with redis-cli

```bash
redis-cli
```

Or use any Redis client library in your preferred language.

---

## Supported Commands (25 Total)

### Connection Commands (2)

#### PING
Check if the server is alive.

```bash
127.0.0.1:6379> PING
PONG

127.0.0.1:6379> PING "hello"
"hello"
```

- **Returns**: `PONG` or the message you send
- **Complexity**: O(1)
- **Use case**: Health checks, latency measurement

#### ECHO
Echo back the given string.

```bash
127.0.0.1:6379> ECHO "Hello, Redis!"
"Hello, Redis!"
```

- **Returns**: The string you provided
- **Complexity**: O(1)
- **Use case**: Testing, debugging

---

### String Commands (6)

Strings are simple key-value pairs.

#### SET
Set a string value.

```bash
127.0.0.1:6379> SET name "John Doe"
OK

127.0.0.1:6379> SET counter "0"
OK
```

- **Syntax**: `SET key value`
- **Returns**: `OK`
- **Complexity**: O(1)
- **Note**: Overwrites existing value

#### GET
Get a string value.

```bash
127.0.0.1:6379> GET name
"John Doe"

127.0.0.1:6379> GET nonexistent
(nil)
```

- **Syntax**: `GET key`
- **Returns**: Value or `nil` if key doesn't exist
- **Complexity**: O(1)

#### INCR
Increment an integer value by 1.

```bash
127.0.0.1:6379> SET counter "10"
OK

127.0.0.1:6379> INCR counter
(integer) 11

127.0.0.1:6379> INCR counter
(integer) 12

127.0.0.1:6379> INCR newcounter
(integer) 1
```

- **Syntax**: `INCR key`
- **Returns**: New value after increment
- **Complexity**: O(1)
- **Note**: Creates key with value 0 if it doesn't exist, then increments to 1
- **Error**: Returns error if value is not an integer

#### DECR
Decrement an integer value by 1.

```bash
127.0.0.1:6379> SET counter "10"
OK

127.0.0.1:6379> DECR counter
(integer) 9

127.0.0.1:6379> DECR counter
(integer) 8

127.0.0.1:6379> DECR newcounter
(integer) -1
```

- **Syntax**: `DECR key`
- **Returns**: New value after decrement
- **Complexity**: O(1)
- **Note**: Creates key with value 0 if it doesn't exist, then decrements to -1

#### EXISTS
Check if a key exists.

```bash
127.0.0.1:6379> SET name "John"
OK

127.0.0.1:6379> EXISTS name
(integer) 1

127.0.0.1:6379> EXISTS nonexistent
(integer) 0
```

- **Syntax**: `EXISTS key`
- **Returns**: `1` if exists, `0` if not
- **Complexity**: O(1)
- **Note**: Works for all data types (strings, lists, sets, hashes)

#### DEL
Delete a key.

```bash
127.0.0.1:6379> SET name "John"
OK

127.0.0.1:6379> DEL name
(integer) 1

127.0.0.1:6379> DEL name
(integer) 0
```

- **Syntax**: `DEL key`
- **Returns**: `1` if deleted, `0` if key didn't exist
- **Complexity**: O(1)
- **Note**: Works for all data types

---

### List Commands (6)

Lists are ordered collections of strings. You can push/pop from both ends.

```
mylist: ["first", "second", "third"]
         ↑                        ↑
       LEFT                     RIGHT
```

#### LPUSH
Push values to the left (head) of a list.

```bash
127.0.0.1:6379> LPUSH tasks "task3"
(integer) 1

127.0.0.1:6379> LPUSH tasks "task2" "task1"
(integer) 3

127.0.0.1:6379> LRANGE tasks 0 -1
1) "task1"
2) "task2"
3) "task3"
```

- **Syntax**: `LPUSH key value [value ...]`
- **Returns**: Length of list after push
- **Complexity**: O(1) per element
- **Note**: Last value becomes the head. `LPUSH key a b c` results in `[c, b, a]`

#### RPUSH
Push values to the right (tail) of a list.

```bash
127.0.0.1:6379> RPUSH tasks "task1" "task2" "task3"
(integer) 3

127.0.0.1:6379> LRANGE tasks 0 -1
1) "task1"
2) "task2"
3) "task3"
```

- **Syntax**: `RPUSH key value [value ...]`
- **Returns**: Length of list after push
- **Complexity**: O(1) per element
- **Use case**: Building queues (FIFO with LPOP)

#### LPOP
Pop a value from the left (head) of a list.

```bash
127.0.0.1:6379> RPUSH tasks "task1" "task2" "task3"
(integer) 3

127.0.0.1:6379> LPOP tasks
"task1"

127.0.0.1:6379> LPOP tasks
"task2"

127.0.0.1:6379> LPOP emptylist
(nil)
```

- **Syntax**: `LPOP key`
- **Returns**: The popped value or `nil` if list is empty
- **Complexity**: O(1)
- **Use case**: Queue processing (with RPUSH)

#### RPOP
Pop a value from the right (tail) of a list.

```bash
127.0.0.1:6379> RPUSH tasks "task1" "task2" "task3"
(integer) 3

127.0.0.1:6379> RPOP tasks
"task3"

127.0.0.1:6379> RPOP tasks
"task2"
```

- **Syntax**: `RPOP key`
- **Returns**: The popped value or `nil` if list is empty
- **Complexity**: O(1)
- **Use case**: Stack operations (LIFO with RPUSH)

#### LRANGE
Get a range of elements from a list.

```bash
127.0.0.1:6379> RPUSH mylist "a" "b" "c" "d" "e"
(integer) 5

127.0.0.1:6379> LRANGE mylist 0 2
1) "a"
2) "b"
3) "c"

127.0.0.1:6379> LRANGE mylist 0 -1
1) "a"
2) "b"
3) "c"
4) "d"
5) "e"

127.0.0.1:6379> LRANGE mylist -3 -1
1) "c"
2) "d"
3) "e"
```

- **Syntax**: `LRANGE key start stop`
- **Returns**: Array of elements
- **Complexity**: O(S+N) where S is start offset and N is number of elements
- **Note**: Indices are 0-based. Negative indices count from the end (-1 is last element)

#### LLEN
Get the length of a list.

```bash
127.0.0.1:6379> RPUSH mylist "a" "b" "c"
(integer) 3

127.0.0.1:6379> LLEN mylist
(integer) 3

127.0.0.1:6379> LLEN nonexistent
(integer) 0
```

- **Syntax**: `LLEN key`
- **Returns**: Length of list or 0 if key doesn't exist
- **Complexity**: O(1)

---

### Set Commands (5)

Sets are unordered collections of unique strings. No duplicates allowed!

```
myset: {"apple", "banana", "orange"}
```

#### SADD
Add members to a set.

```bash
127.0.0.1:6379> SADD tags "redis" "golang" "tutorial"
(integer) 3

127.0.0.1:6379> SADD tags "redis"
(integer) 0

127.0.0.1:6379> SADD tags "database" "nosql"
(integer) 2
```

- **Syntax**: `SADD key member [member ...]`
- **Returns**: Number of members actually added (excludes duplicates)
- **Complexity**: O(1) per member
- **Note**: Automatically ignores duplicates

#### SMEMBERS
Get all members of a set.

```bash
127.0.0.1:6379> SADD tags "redis" "golang" "tutorial"
(integer) 3

127.0.0.1:6379> SMEMBERS tags
1) "redis"
2) "golang"
3) "tutorial"

127.0.0.1:6379> SMEMBERS nonexistent
(empty array)
```

- **Syntax**: `SMEMBERS key`
- **Returns**: Array of all members
- **Complexity**: O(N) where N is set size
- **Note**: Order is not guaranteed

#### SISMEMBER
Check if a member exists in a set.

```bash
127.0.0.1:6379> SADD tags "redis" "golang"
(integer) 2

127.0.0.1:6379> SISMEMBER tags "redis"
(integer) 1

127.0.0.1:6379> SISMEMBER tags "python"
(integer) 0
```

- **Syntax**: `SISMEMBER key member`
- **Returns**: `1` if member exists, `0` if not
- **Complexity**: O(1)
- **Use case**: Fast membership checks

#### SREM
Remove members from a set.

```bash
127.0.0.1:6379> SADD tags "redis" "golang" "tutorial"
(integer) 3

127.0.0.1:6379> SREM tags "tutorial"
(integer) 1

127.0.0.1:6379> SREM tags "python"
(integer) 0

127.0.0.1:6379> SMEMBERS tags
1) "redis"
2) "golang"
```

- **Syntax**: `SREM key member [member ...]`
- **Returns**: Number of members actually removed
- **Complexity**: O(1) per member

#### SCARD
Get the cardinality (size) of a set.

```bash
127.0.0.1:6379> SADD tags "redis" "golang" "tutorial"
(integer) 3

127.0.0.1:6379> SCARD tags
(integer) 3

127.0.0.1:6379> SCARD nonexistent
(integer) 0
```

- **Syntax**: `SCARD key`
- **Returns**: Number of members in the set
- **Complexity**: O(1)

---

### Hash Commands (6)

Hashes are maps of field-value pairs. Perfect for representing objects!

```
user:1000 → {
  "name": "John",
  "email": "john@example.com",
  "age": "30"
}
```

#### HSET
Set a field in a hash.

```bash
127.0.0.1:6379> HSET user:1 name "John Doe"
(integer) 1

127.0.0.1:6379> HSET user:1 email "john@example.com"
(integer) 1

127.0.0.1:6379> HSET user:1 name "Jane Doe"
(integer) 0
```

- **Syntax**: `HSET key field value`
- **Returns**: `1` if new field, `0` if field was updated
- **Complexity**: O(1)

#### HGET
Get a field from a hash.

```bash
127.0.0.1:6379> HSET user:1 name "John Doe"
(integer) 1

127.0.0.1:6379> HGET user:1 name
"John Doe"

127.0.0.1:6379> HGET user:1 age
(nil)
```

- **Syntax**: `HGET key field`
- **Returns**: Value of field or `nil` if field doesn't exist
- **Complexity**: O(1)

#### HGETALL
Get all fields and values from a hash.

```bash
127.0.0.1:6379> HSET user:1 name "John Doe"
(integer) 1

127.0.0.1:6379> HSET user:1 email "john@example.com"
(integer) 1

127.0.0.1:6379> HSET user:1 age "30"
(integer) 1

127.0.0.1:6379> HGETALL user:1
1) "name"
2) "John Doe"
3) "email"
4) "john@example.com"
5) "age"
6) "30"
```

- **Syntax**: `HGETALL key`
- **Returns**: Flat array of [field1, value1, field2, value2, ...]
- **Complexity**: O(N) where N is hash size
- **Note**: Returns empty array if key doesn't exist

#### HDEL
Delete fields from a hash.

```bash
127.0.0.1:6379> HSET user:1 name "John" email "john@example.com" age "30"
(integer) 3

127.0.0.1:6379> HDEL user:1 age
(integer) 1

127.0.0.1:6379> HDEL user:1 phone
(integer) 0

127.0.0.1:6379> HGETALL user:1
1) "name"
2) "John"
3) "email"
4) "john@example.com"
```

- **Syntax**: `HDEL key field [field ...]`
- **Returns**: Number of fields actually deleted
- **Complexity**: O(1) per field

#### HEXISTS
Check if a field exists in a hash.

```bash
127.0.0.1:6379> HSET user:1 name "John"
(integer) 1

127.0.0.1:6379> HEXISTS user:1 name
(integer) 1

127.0.0.1:6379> HEXISTS user:1 age
(integer) 0
```

- **Syntax**: `HEXISTS key field`
- **Returns**: `1` if field exists, `0` if not
- **Complexity**: O(1)

#### HLEN
Get the number of fields in a hash.

```bash
127.0.0.1:6379> HSET user:1 name "John" email "john@example.com"
(integer) 2

127.0.0.1:6379> HLEN user:1
(integer) 2

127.0.0.1:6379> HLEN nonexistent
(integer) 0
```

- **Syntax**: `HLEN key`
- **Returns**: Number of fields in the hash
- **Complexity**: O(1)

---

## Some More Examples

### Example 1: Task Queue

```bash
RPUSH queue:tasks "process-payment-123"
RPUSH queue:tasks "send-email-456"
RPUSH queue:tasks "generate-report-789"

LPOP queue:tasks

LLEN queue:tasks
```

### Example 2: User Profile

```bash
HSET user:1000 name "Alice"
HSET user:1000 email "alice@example.com"
HSET user:1000 created_at "2025-01-01"

HGET user:1000 email

HGETALL user:1000
```

### Example 3: Tagging System

```bash
SADD post:1:tags "redis" "golang" "tutorial"
SADD post:2:tags "redis" "python"

SISMEMBER post:1:tags "golang"

SMEMBERS post:1:tags
```

### Example 4: Page View Counter

```bash
INCR pageviews:home
INCR pageviews:home
INCR pageviews:home

GET pageviews:home
```

### Example 5: Shopping Cart

```bash
RPUSH cart:user123 "product:laptop"
RPUSH cart:user123 "product:mouse"
RPUSH cart:user123 "product:keyboard"

LRANGE cart:user123 0 -1

HSET product:laptop name "MacBook Pro"
HSET product:laptop price "2499"

HGETALL product:laptop
```

---

## Architecture

### Data Storage

```go
type store struct {
    strings map[string]string              
    lists   map[string][]string            
    sets    map[string]map[string]struct{} 
    hashes  map[string]map[string]string   
    mu      sync.RWMutex
}
```

- **Strings**: Simple key-value map
- **Lists**: Go slices for ordered collections
- **Sets**: `map[string]struct{}` for O(1) lookups with zero memory overhead
- **Hashes**: Nested maps for structured data
- **Thread-safe**: All operations protected by RWMutex

### RESP Protocol

All communication uses the Redis Serialization Protocol (RESP):

- **Simple Strings**: `+OK\r\n`
- **Errors**: `-ERR message\r\n`
- **Integers**: `:42\r\n`
- **Bulk Strings**: `$5\r\nhello\r\n`
- **Arrays**: `*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n`

---

## Testing

Run the test suite:

```bash
go test -v
```

Run specific tests:

```bash
go test -v -run TestLists
go test -v -run TestSets
go test -v -run TestHashes
```

---

## Implementation Notes

### Memory Management
- Empty data structures are automatically deleted to save memory
- Keys are removed when their last element/field is deleted

### Thread Safety
- Read operations use `RLock()` for concurrent reads
- Write operations use `Lock()` for exclusive access
- All operations are atomic

### Differences from Real Redis
- No persistence (in-memory only)
- No TTL/expiration
- No pub/sub
- No transactions (MULTI/EXEC)
- No Lua scripting
- No sorted sets
- No set operations (SUNION, SINTER, SDIFF)

---

## License

MIT
