# Monkey Lang

### Overview
This is a language mostly based on the fantastic book 
[Writing An Interpreter In Go](https://interpreterbook.com/). 
I do this for educational purposes so I can build my own stuff 
from ground up later. 
I've implemented the language using [Go](https://go.dev/) and [Odin](http://odin-lang.org/)(still in development) programming languages. 

## Installation

simply clone the repo and: 
```sh
cd MonkeyLang/Go
go run .
```
by doing this, you will run the REPL(Interactive mode) 


Warning: The above code works for Go. If you want to try out 
the Odin version: 
```sh
cd MonkeyLang/Odin
odin run .
```


## Features
the language is mostly gonna be based on the book, but shall 
also include more features like support for else-if expressions. 
The Odin version must have a dedicated GC(garbage collector).


## Structure

#### Lexer
Given an input string, we first tokenize it with the help of 
the Lexer. At this step we don't care about the syntax or the semantics 
of our language. We only define what is considered an identifier, a keyword,
an operator and etc.


#### Parser
After tokenization, the Parser will basically check for the correctness of 
the syntax. It reads through the tokens and builds an [Abstract Syntax Tree](https://en.wikipedia.org/wiki/Abstract_syntax_tree). 
The tree is an abstract summarization of our grammar for the given input. it is 
abstract because it does not contain some inessential parts of the syntax like 
punctuation(braces, brackets and etc.) 
It must be noted that [Pratt Parsing](https://en.wikipedia.org/wiki/Operator-precedence_parser#Pratt_parsing) is the author's choice 
for the parsing algorithm.


#### Eval
We give meaning to whatever we accept as valid input. This means that we need 
to implement some sort of object system to keep track of things. In monkey everything is 
an object. We even have first class functions(we can pass them around and assign them to variable 
. we even have closures). the Eval() function evaluates the input recursively. Environments are created 
as a way to keep track of identifiers and the values associated with them(kind of like scopes). 
Functions have their own environments that are the extension of the environment they were defined 
in. That is precisely the reason why closures are possible.


#### GC
This is a feature exclusive to the Odin version as it is a manual memory management language. 
The Go version takes care of deallocation thanks to the Go's GC. 
We basically have to take care of the objects created in memory through monkey scripts and make sure 
we avoid memory leaks and avoid crashes and undefined behaviour. 


## Current State for Go version
| Feature | State |
| ------- | -------|
| Lexer | working state |
| Parser | working state |
| Eval/Semantics | working state |


## Current State for Odin version
| Feature | State |
| ------- | -------|
| Lexer | working state |
| Parser | not developed yet |
| Eval/Semantics | not developed yet |
| Garbage Collection | not developed yet |
