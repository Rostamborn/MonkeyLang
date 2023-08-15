package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// basically an enum
const (
    _ int = iota
    LOWEST
    EQUALS      // ==
    LESSGREATER // > or <
    SUM         // +
    PRODUCT     // *
    PREFIX      // -X or !X
    CALL        // myFunction(X)
    INDEX       // array[index]
)

var precedences = map[token.TokenType]int {
    token.EQ: EQUALS,
    token.NOT_EQ: EQUALS,
    token.LT: LESSGREATER,
    token.GT: LESSGREATER,
    token.PLUS: SUM,
    token.MINUS: SUM,
    token.ASTERISK: PRODUCT,
    token.SLASH: PRODUCT,
}

type (
    prefixParseFunc func() ast.Expression
    infixParseFunc func(ast.Expression) ast.Expression
)

type Parser struct {
    lex *lexer.Lexer

    curToken token.Token
    peekToken token.Token
    errors []string

    prefixParseFuncs map[token.TokenType]prefixParseFunc
    infixParseFuncs map[token.TokenType]infixParseFunc
}

func NewParser(lex *lexer.Lexer) *Parser {
    p := &Parser{lex: lex, errors: make([]string, 0)}
    p.nextToken() // curToken is still nil
    p.nextToken() // after this second call, curToken is not nil anymore

    p.prefixParseFuncs = make(map[token.TokenType]prefixParseFunc)
    p.registerPrefix(token.IDENT, p.parseIdentifier)
    p.registerPrefix(token.INT, p.parseIntegerLiteral)
    p.registerPrefix(token.BANG, p.parsePrefixExpression)
    p.registerPrefix(token.MINUS, p.parsePrefixExpression)

    p.infixParseFuncs = make(map[token.TokenType]infixParseFunc)
    p.registerInfix(token.PLUS, p.parseInfixExpression)
    p.registerInfix(token.MINUS, p.parseInfixExpression)
    p.registerInfix(token.SLASH, p.parseInfixExpression)
    p.registerInfix(token.ASTERISK, p.parseInfixExpression)
    p.registerInfix(token.EQ, p.parseInfixExpression)
    p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)

    return p
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = make([]ast.Statement, 0)

    for !p.curTokenIs(token.EOF) {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.nextToken()
    }
    return program
}



func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.LET:
        return p.parseLetStatement()
    case token.RETURN:
        return p.parseReturnStatement()
    default:
        return p.parseExpressionStatement()
    }
}

// TODO: if you don't put semicolon, nothing happens!
// the computer just gets hot and loud
func (p *Parser) parseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token: p.curToken}

    if !p.expectPeek(token.IDENT) {
        return nil
    }

    stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

    if !p.expectPeek(token.ASSIGN) {
        return nil
    }

    for !p.curTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: p.curToken}

    // if !p.expectPeek(token.IDENT) {
    //     return nil
    // }
    //
    // stmt.ReturnValue = 
    p.nextToken()

    if !p.curTokenIs(token.SEMICOLON) {
        p.nextToken()
    }

    return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement{Token: p.curToken}

    stmt.Expression = p.parseExpression(LOWEST)

    if p.peekTokenIs(token.SEMICOLON) { // The smicolon is optional in expressions because it
        p.nextToken()                   // makes it easier to have stuff like "1 + 2" in the REPL
    }                                   // in that case we don't need to type "1 + 2;"

    return stmt
}

// Pratt Parsing in action
func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParseFuncs[p.curToken.Type]
    if prefix == nil {
        p.noPrefixFuncError(p.curToken.Type)
        return nil
    }
    leftExpression := prefix()

    for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
        infixFunc := p.infixParseFuncs[p.peekToken.Type]
        if infixFunc == nil {
            return leftExpression
        }

        p.nextToken()

        leftExpression = infixFunc(leftExpression)
    }

    return leftExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    literal := &ast.IntegerLiteral{Token: p.curToken}

    value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
    if err != nil {
        message := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
        p.errors = append(p.errors, message)
        return nil
    }

    literal.Value = value

    return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    expression := &ast.PrefixExpression{
        Token: p.curToken,
        Operator: p.curToken.Literal,
    }

    p.nextToken()

    expression.Right = p.parseExpression(PREFIX)

    return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression{
        Token: p.curToken,
        Operator: p.curToken.Literal,
        Left: left,
    }

    precedence := p.curPrecedence()
    p.nextToken()
    expression.Right = p.parseExpression(precedence)

    return expression
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.lex.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
    return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
    return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekTokenIs(t) {
        p.nextToken()
        return true
    } else {
        p.peekError(t)
        return false
    }
}

func (p *Parser) peekPrecedence() int {
    if p, ok := precedences[p.peekToken.Type]; ok {
        return p
    }
    return LOWEST
}

func (p *Parser) curPrecedence() int {
    if p, ok := precedences[p.curToken.Type]; ok {
        return p
    }
    return LOWEST
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
    message := fmt.Sprintf("expected next token to be {%s}, got {%s} instead", t, p.peekToken.Type)
    p.errors = append(p.errors, message)
}

func (p *Parser) noPrefixFuncError(t token.TokenType) {
    message := fmt.Sprintf("no prefix parse function for {%s} found", t)
    p.errors = append(p.errors, message)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFunc) {
    p.prefixParseFuncs[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFunc) {
    p.infixParseFuncs[tokenType] = fn
}

