# unumber

Numeric abstraction helpers.

Includes flexible number parsing, denomination helpers, and support for regular and big numeric values.

## Example

```go
amount, err := unumber.FromString("123.45", false)
cents, err := unumber.AsDenom[int64](123.45, 2)
```
