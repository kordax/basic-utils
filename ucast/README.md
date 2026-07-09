# ucast

Conversion helpers for basic Go types.

Includes string conversions, numeric conversions, boolean conversions, and generic helpers for working with basic pointer/value types.

## Example

```go
port, err := ucast.String[int]("8080")
enabled := ucast.Type(true) // "true"
```
