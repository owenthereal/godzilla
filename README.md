# Godzilla: Go running JavaScript

## Overview

Godzilla is a [ES2015](http://babeljs.io/learn-es2015) to Go source code transpiler and runtime that is intended to be a near drop-in replacement for [Node.js](https://nodejs.org).
It compiles ES2015 source code to Go source code which is then compiled to native code.
The compiled Go source code is a series of calls to the Godzilla runtime, a Go library serving a similar purpose to Node.js.

Godzilla parses ES2015 source code with the awesome [babylon](https://github.com/babel/babylon).
That means at the moment Node.js is required for compilation.
As Godzilla becomes mature, `babylon` will be compiled to Go source code using Godzilla itself so that the Node.js dependency can be dropped.

**Note that Godzilla is at a very early stage and only very few language features are implemented**

[![asciicast](https://asciinema.org/a/120175.png)](https://asciinema.org/a/120175)

## Compiling

Make sure Go and Node.js are installed properly, then run:

```
make
```

## Running

```
echo "console.log('Hello, Godzilla')" | bin/godzilla run
```

## Related Arts

* [grumpy](https://github.com/google/grumpy)
