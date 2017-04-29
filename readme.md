# Hydrator

Very simple demo to show how to use the template engine of Golang. It's called hydrator
as it could form the basis of a templated/JSON data hydrator.

# Purpose

1. Takes in JSON input (try `source.json`)
2. Go through the JSON and apply the transformations dictated in external template (try `transform.tmpl`)
3. Print formatted JSON to standard out

Hacked together in about 30 mins so don't judge me - and if you find it useful,
open a PR and tidy it up!

# Getting Started

`go get -u github.com/17twenty/hydrator`

```
hydrator -s <source JSON file> -t <transformation template> [-w <output.json>]
```

# Todo

It would be really cool to use the `plugins` capability to dynamically add new
template functions. I'll leave that to someone else tho.