# uos

Operating-system helpers, mostly focused on environment variables.

Includes required env accessors, optional env accessors, typed parsers, default values, duration/time/URL/bool parsing, and CPU count helpers.

## Example

```go
port := uos.RequireEnvAs("PORT", uos.MapStringToInt)
debug := uos.GetEnvOptBool("DEBUG").OrElse(false)
```
