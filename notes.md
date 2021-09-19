# Creating a simple Pratt Parser

## Concpepts

- Prefix-Parsing-Functions
    A prefix-parsing-function is a function that is called to parse an expression from the start.
    e.g when parsing `1 + 2 + 3;`, there needs to be a prefix parsing function to parse the 1, 2 and 3 ie (1 + 2 + 3)
    (2 + 3) (3)
- Infix-Parsing-Function is a function that takes in an already parsed expression and adds it to another expression creating a new expression eg. There  needs to be an infix parsing function that takes in `1` and parses `(2 + 3)` to create `(1 + (2 + 3))` 

- You need prefix-parsing-functions to parse
    - identifiers eg `name`
    - integers eg `20`
    - bang ie `!`
    - booleans ie `true and false`

- You need infix-parsing-functions for infix operators such as `+,-,*`

## Precedence and Associativity

This is how the pratt parser handles operator precedence and associativity. Given that every operator has a given binding-power-value, in the inner loop for the `parseExpression(precedence int)` function, the binding-power of the current operator  is compared to the binding power of the next operator and if the next operator has a greater binding power than the current operator then the next operator is parsed first, thus resultig in it being nested inside another epression thus making it deeper in the ast