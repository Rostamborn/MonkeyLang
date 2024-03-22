# Monkey Lang

### Overview
This is a language mostly based on the fantastic book 
[Writing An Interpreter In Go](https://interpreterbook.com/). 
I do this for educational purposes so I can build my own stuff 
from ground up later. 

## Installation

simply clone the repo and: 
```sh
cd MonkeyLang
go run .
```
by doing this, you will run the REPL(Interactive mode) 
## Features
the language is mostly gonna be based on the book, but shall 
also include more features like support for else-if expressions. 


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


## Current State for Go version
| Feature | State |
| ------- | -------|
| Lexer | working state |
| Parser | working state |
| Eval/Semantics | working state |
