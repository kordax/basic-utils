# usize

Byte-sized data unit constants.

Use decimal units for SI sizes and binary units for power-of-two sizes:

```go
const databaseLoadChunkSize = 8 * usize.MiB
const payloadLimit = 10 * usize.MB
```

Short aliases are available for byte units (`KB`, `MB`, `KiB`, `MiB`) and bit units (`Mbit`, `Mibit`). Bit units are represented as byte counts.
