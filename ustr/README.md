# ustr

String helpers.

Includes defaulting, concatenation, blank checks, pointer helpers, trimming, split-and-trim, and whitespace normalization.

## Example

```go
name := ustr.DefaultIfBlank(input.Name, "anonymous")
tags := ustr.SplitAndTrim("go, utils, cache", ",")
```
