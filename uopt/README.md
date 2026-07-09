# uopt

Generic optional value implementation.

Includes optional constructors, JSON/SQL support, default helpers, conditional constructors, and functional helpers such as `Map`, `FlatMap`, and `Filter`.

## Example

```go
name := uopt.Map(userOpt, func(user User) string {
	return user.Name
})
value := name.OrElse("anonymous")
```
