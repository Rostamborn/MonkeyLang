package parser

import (
    "fmt"
    "testing"
    "monkey/ast"
    "monkey/lexer"
)

func TestLetStatements(t *testing.T) {
    input := `
    let x = 5;
    let y = 10;
    let foobar = 838383;
    `

    lex := lexer.NewLexer(input)
    p := NewParser(lex)

    program := p.ParseProgram()
    checkParserErrors(t, p)
    // if program == nil {
    //     t.Fatalf("ParseProgram() returned nil")
    // }
    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
    }

    tests := []struct {
        expectedIdentifier string
    }{
        {"x"},
        {"y"},
        {"foobar"},
    }

    for i, tt := range tests {
        stmt := program.Statements[i]
        if !testLetStatement(t, stmt, tt.expectedIdentifier) {
            return
        }
    }
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
    if s.TokenLiteral() != "let" {
        t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
        return false
    }

    letStmt, ok := s.(*ast.LetStatement)
    if !ok {
        t.Errorf("s not *ast.LetStatement. got=%T", s)
        return false
    }

    if letStmt.Name.Value != name {
        t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
        return false
    }

    if letStmt.Name.TokenLiteral() != name {
        t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
        return false
    }

    return true
}

func TestReturnStatements(t *testing.T) {
    input := `
    return 5;
    return 32;
    return 329846;
    return ;
    `
    lex := lexer.NewLexer(input)
    p := NewParser(lex)

    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 4 {
        t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
    }

    for _, stmt := range program.Statements {
        rtStatement, ok := stmt.(*ast.ReturnStatement) 
        if !ok {
            t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
            continue
        }
        if rtStatement.TokenLiteral() != "return" {
            t.Errorf("rtStatement.TokenLiteral not 'return', got %q", rtStatement.TokenLiteral())
        }
    }
}

func checkParserErrors(t *testing.T, p *Parser) {
    errors := p.Errors()

    if len(errors) == 0 {
        return
    }

    t.Errorf("parser has %d errors", len(errors))
    for _, msg := range errors {
        t.Errorf("parser error: %q", msg)
    }
    t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
    input := "mate;"

    lex := lexer.NewLexer(input)
    p := NewParser(lex)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    ident, ok := stmt.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
    }

    if ident.TokenLiteral() != "mate" {
        t.Errorf("ident.TokenLiteral not %s. got=%s", "mate", ident.TokenLiteral())
    }

    if ident.Value != "mate" {
        t.Errorf("ident.Value not %s. got=%s", "mate", ident.Value)
    }
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "17;"

    lex := lexer.NewLexer(input)
    p := NewParser(lex)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    literal, ok := stmt.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
    }

    if literal.Value != 17 {
        t.Errorf("literal.Value not %d. got=%d", 17, literal.Value)
    }

    if literal.TokenLiteral() != "17" {
        t.Errorf("literal.TokenLiteral not %s. got=%s", "17", literal.TokenLiteral())
    }
}

func TestParsingPrefixExpression(t *testing.T) {
    prefixTests := []struct {
        input string
        operator string
        value int64
    }{
        {"!5;", "!", 5},
        {"-15;", "-", 15},
    }

    for _, tt := range prefixTests {
        lex := lexer.NewLexer(tt.input)
        p := NewParser(lex)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
        }

        exp, ok := stmt.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("exp not *ast.PrefixExpression. got=%T", stmt.Expression)
        }

        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not %s. got=%s", tt.operator, exp.Operator)
        }

        if !testIntegerLiteral(t, exp.Right, tt.value) {
            return
        }
    }
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
    integer, ok := il.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
        return false
    }

    if integer.Value != value {
        t.Errorf("integ.Value not %d. got=%d", value, integer.Value)
        return false
    }

    if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
        t.Errorf("integ.TokenLiteral not %d. got=%s", value, integer.TokenLiteral())
        return false
    }

    return true
}

func TestParsingInfixExpressions(t *testing.T) {
    infixTests := []struct {
        input string
        leftValue int64
        operator string
        rightValue int64
    }{
        {"8 + 8;", 8, "+", 8},
        {"8 - 8;", 8, "-", 8},
        {"8 * 8;", 8, "*", 8},
        {"8 / 8;", 8, "/", 8},
        {"8 > 8;", 8, ">", 8},
        {"8 < 8;", 8, "<", 8},
        {"8 == 8;", 8, "==", 8},
        {"8 != 8;", 8, "!=", 8},
    }

    for _, tt := range infixTests {
        lex := lexer.NewLexer(tt.input)
        p := NewParser(lex)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
        }

        exp, ok := stmt.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("exp not *ast.InfixExpression. got=%T", stmt.Expression)
        }

        if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
            return
        }

        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not %s. got=%s", tt.operator, exp.Operator)
        }

        if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
            return
        }
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {"-a * b", "((-a) * b)"},
        {"!-a", "(!(-a))"},
        {"a + b + c", "((a + b) + c)"},
        {"a + b - c", "((a + b) - c)"},
        {"a * b * c", "((a * b) * c)"},
        {"a * b / c", "((a * b) / c)"},
        {"a + b / c", "(a + (b / c))"},
        {"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
        {"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
        {"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
        {"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
        {"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
        // {"true", "true"},
        // {"false", "false"},
        // {"3 > 5 == false", "((3 > 5) == false)"},
        // {"3 < 5 == true", "((3 < 5) == true)"},
        // {"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
        // {"(5 + 5) * 2", "((5 + 5) * 2)"},
        // {"2 / (5 + 5)", "(2 / (5 + 5))"},
        // {"-(5 + 5)", "(-(5 + 5))"},
        // {"!(true == true)", "(!(true == true))"},
    }

    for _, tt := range tests {
        lex := lexer.NewLexer(tt.input)
        p := NewParser(lex)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        actual := program.String()
        if actual != tt.expected {
            t.Errorf("expected=%q, got=%q", tt.expected, actual)
        }
    }
}
