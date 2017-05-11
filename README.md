# Godzilla: Go running JavaScript

## Overview

Godzilla is a JavaScript to Go source code transpiler and runtime that is intended to be a near drop-in replacement for ES2015.
It compiles ES2015 source code to Go source code which is then compiled to native code.
The compiled Go source code is a series of calls to the Godzilla runtime, a Go library serving a similar purpose to Node.js.

## Compiling

```
make
```

## Running

```
echo "console.log('hello, Godzilla')" | bin/godzilla run
```

## Related Arts

* [grumpy](https://github.com/google/grumpy)
