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
        value interface{}
    }{
        {"!5;", "!", 5},
        {"-15;", "-", 15},
        {"!true;", "!", true},
        {"!false;", "!", false},
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

        if !testLiteralExpression(t, exp.Right, tt.value) {
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

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
    ident, ok := exp.(*ast.Identifier)
    if !ok {
        t.Errorf("exp not *ast.Identifier. got=%T", exp)
        return false
    }

    if ident.Value != value {
        t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
        return false
    }

    if ident.TokenLiteral() != value {
        t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
        return false
    }

    return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
    boolean, ok := exp.(*ast.Boolean)
    if !ok {
        t.Errorf("exp not *ast.Boolean. got=%T", exp)
        return false
    }

    if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
        t.Errorf("boolean.TokenLiteral not %s. got=%s", fmt.Sprintf("%t", value), boolean.TokenLiteral())
        return false
    }

    if boolean.Value != value {
        t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
        return false
    }

    return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
    switch v := expected.(type) {
    case int:
        return testIntegerLiteral(t, exp, int64(v))
    case int64:
        return testIntegerLiteral(t, exp, v)
    case string:
        return testIdentifier(t, exp, v)
    case bool:
        return testBooleanLiteral(t, exp, v)
    }

    t.Errorf("type of exp not handled. got=%T", exp)
    return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
    opExpression, ok := exp.(*ast.InfixExpression)
    if !ok {
        t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
        return false
    }

    if !testLiteralExpression(t, opExpression.Left, left) {
        return false
    }

    if opExpression.Operator != operator {
        t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExpression.Operator)
        return false
    }

    if !testLiteralExpression(t, opExpression.Right, right) {
        return false
    }

    return true
}

func TestParsingInfixExpressions(t *testing.T) {
    infixTests := []struct {
        input string
        leftValue interface{}
        operator string
        rightValue interface{}
    }{
        {"8 + 8;", 8, "+", 8},
        {"8 - 8;", 8, "-", 8},
        {"8 * 8;", 8, "*", 8},
        {"8 / 8;", 8, "/", 8},
        {"8 > 8;", 8, ">", 8},
        {"8 < 8;", 8, "<", 8},
        {"8 == 8;", 8, "==", 8},
        {"8 != 8;", 8, "!=", 8},
        {"true == true", true, "==", true},
        {"true != false", true, "!=", false},
        {"false == false", false, "==", false},
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

        if !testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue) {
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
        // {"(a + b) * c", ""},
        {"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
        {"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
        {"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
        {"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
        {"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
        {"true", "true"},
        {"false", "false"},
        {"3 > 5 == false", "((3 > 5) == false)"},
        {"3 < 5 == true", "((3 < 5) == true)"},
        {"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
        {"(5 + 5) * 2", "((5 + 5) * 2)"},
        {"2 / (5 + 5)", "(2 / (5 + 5))"},
        {"-(5 + 5)", "(-(5 + 5))"},
        {"!(true == true)", "(!(true == true))"},
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

func TestBooleanExpression(t *testing.T) {
    input := "false; true;"

    lex := lexer.NewLexer(input)
    p := NewParser(lex)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 2 {
        t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    boolean, ok := stmt.Expression.(*ast.Boolean)
    if !ok {
        t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
    }

    if boolean.TokenLiteral() != "false" {
        t.Errorf("boolean.TokenLiteral not %s. got=%s", "false", boolean.TokenLiteral())
    }

    if boolean.Value != false {
        t.Errorf("boolean.Value not %v. got=%v", false, boolean.Value)
    }

    stmt, ok = program.Statements[1].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[1] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    boolean, ok = stmt.Expression.(*ast.Boolean)
    if !ok {
        t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
    }

    if boolean.TokenLiteral() != "true" {
        t.Errorf("boolean.TokenLiteral not %s. got=%s", "true", boolean.TokenLiteral())
    }

    if boolean.Value != true {
        t.Errorf("boolean.Value not %v. got=%v", false, boolean.Value)
    }
}

func TestIfExpression(t *testing.T) {
    input := "if (x < y) { x }"

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

    exp, ok := stmt.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
    }

    if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
        return
    }

    if len(exp.Consequence.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
    }

    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
    }

    if !testIdentifier(t, consequence.Expression, "x") {
        return
    }

    if exp.Alternative != nil {
        t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
    }
}

func TestIfElseExpression(t *testing.T) {
    input := "if (x < y) { x } else { y }"

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

    exp, ok := stmt.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("exp not *ast.IfExpression. got=%T", stmt.Expression)
    }

    if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
        return
    }

    if len(exp.Consequence.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
    }

    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
    }

    if !testIdentifier(t, consequence.Expression, "x") {
        return
    }

    alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
    }

    if !testIdentifier(t, alternative.Expression, "y") {
        return
    } 
}
