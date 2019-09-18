# V2 allows for (boolean) expressions
Resources:
- Haker news comments: https://news.ycombinator.com/item?id=13914218
- Javascritp TDP: https://crockford.com/javascript/tdop/tdop.html
- Intro to RDP: http://www.cs.binghamton.edu/~zdu/parsdemo/recintro.html
- Simpe RDP in Python: http://effbot.org/zone/simple-top-down-parsing.htm
- Parsing Expressions: http://www.craftinginterpreters.com/parsing-expressions.html
- Pratt Parsing: http://www.oilshell.org/blog/2016/11/01.html
- More Pratt Parsing: https://dev.to/jrop/pratt-parsing
- Pratt Parsers made easy: http://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
  - Golang pratt: https://github.com/richardjennings/prattparser
  - Port to go: https://github.com/moraes/bantam
- Ivy's parser and expression problem: https://www.youtube.com/watch?v=PXoG0WX0r_E
- Use Syntax Diagrams https://en.wikipedia.org/wiki/Syntax_diagram

## TODO 
[x] get field reading to work on some input object
[x] - get indexing operator to work
[x] - figure out if grammar requires strange precedence for dot operator. no
[x] - do we need an object type to implement dot operator for 
[x] - how to make that work for slices
[x] - how to present input, always add some 'hidden' $ token to the beginning? yes
[x] - implement a lexer for jqp v2
[x] - add an e2e function to query a value from `interface{}`
[x] - renamed object to map type, which is more what it is 
[x] - get basic query on syscall/js to work
[x] - simply code for syscall/js query
  [x] - allow tl js type to be evaluated
  [x] - allow index fetching from js array
  [x] - allow porting of native map/slice types
  [x] - re-analyse the types required to get ports to work
[x] - develop concat operator: ','
  [x] - lex it 
  [x] - make sure it is parsable as an expression
  [x] - make sure it evaluated, left to right
[x] - fix instable parsing on certain combinations of fields, indexes and calls
[x] - get calling a (js) object function to work
  [x] - add `call` operation/evaluation
    - [x] parse each argument separately
  [x] - make sure it can be used on js ports
[ ] - get array expansion to work 
[ ] - get pipelining of filters to work 
[ ] - get concat of filter output to work
[ ] - implement a decoder that reads js/interface values into tagged structs
 
# Clean up TODO
[ ] - Implement all other simple operators (mul, sub etc)
[ ] - Replace panics with error handling instead
[ ] - Add boolean type, use it to convert from JS
[ ] - Complete value.fromJS for undefined and symbol

# JQP 
Decode `syscall/js` values in structs using field tags written as jq-like queries.

# Usecases
- extracting data from dom events
- callDOM command return values
- readDOM command return value
- onDOM subscription messages

# Idea
Can we abstract current DOM event handling to a listenDOM(id, msg) "subscription"?

# TODO 
- [x] implement start of parser for the simple filter 
- [x] extend lexer with the position of tokens in the input
- [x] update lexer to replace `range` token with `colon` and `literal` tokens
- [x] finish parsing range operations
- [x] add pipe character for lexing, add pipeline type for chaining filters
- [x] make sure lexing allows for whitespaces between filters with pipes
- [ ] if after the filter a conditional (compare '>', '<', '<=', '>=' or conditional '==', '!=') is 
      found the filter turns into a gate. It will pass on if first or all(?) of the output 
      values return true for the conditional. The piping of filters allows for OR and AND constructions.
- [ ] add select builtin for conditional reading of values, but does require parsing a boolean expression