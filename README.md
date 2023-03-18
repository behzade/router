# Simple Router
This is a simple router that is compatible with the standard library http.Handler interface.

## Features
** Remove extra slashes: multiple slashes or trailing slashes have no effect on the routing
** Handle route parameters: add named paramters to the path with {var} syntax. Path and query params are captured and saved in the context as a value.
** Builtin Options handler: get the Allow header with an OPTIONS request or on MethodNotAllowed responses.
** Can use custom not found and method not allowed handlers.
** Add global middlewares to the router with a score to be sorted by.

## Limitations
** Only supports lowercase a-z and 0-9 and "-" as allowed characters in path.

## Example Usage
```go


```

## Issues
** Feel free to open an issue with a bug/feature request.
