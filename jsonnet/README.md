# Jsonnet

```bash
brew install jsonnet

# Basic examples
jsonnet landingpage.jsonnet

# From the commandline
jsonnet -e '{ x: 1 , y: self.x + 1 } { x: 10 }'

# Use of some jsonnet features
jsonnet example2.jsonnet

# JsonNet also has a stdlib:
# https://jsonnet.org/ref/stdlib.html
jsonnet std.jsonnet
```