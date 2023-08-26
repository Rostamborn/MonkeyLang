package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type Node interface {
    TokenLiteral() string
    String() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

func (p *Program) String() string {
    var out bytes.Buffer

    for _, str := range p.Statements {
        out.WriteString(str.String())
    }

    return out.String()
}

type LetStatement struct {
    Token token.Token
    Name *Identifier
    Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}

func (ls *LetStatement) statementNode() {

}

func (ls *LetStatement) String() string {
    var out  bytes.Buffer

    out.WriteString(ls.TokenLiteral() + " ")
    out.WriteString(ls.Name.String())
    out.WriteString(" = ")

    if ls.Value != nil {
        out.WriteString(ls.Value.String())
    }

    out.WriteString(";")

    return out.String()
}



type ReturnStatement struct {
    Token token.Token
    ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string {
    return rs.Token.Literal
}

func (rs *ReturnStatement) statementNode() {
    
}

func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.TokenLiteral() + " ")

    if rs.ReturnValue != nil {
        out.WriteString(rs.ReturnValue.String())
    }

    out.WriteString(";")

    return out.String()
}

type ExpressionStatement struct {
    Token token.Token
    Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
    return es.Token.Literal
}

func (es *ExpressionStatement) statementNode() {

}

func (es *ExpressionStatement) String() string {
    if es.Expression != nil {
        return es.Expression.String()
    }

    return ""
}

type BlockStatement struct {
    Token token.Token
    Statements []Statement
}

func (bs *BlockStatement) TokenLiteral() string {
    return bs.Token.Literal
}

func (bs *BlockStatement) statementNode() {
    
}

func (bs *BlockStatement) String() string {
    var out bytes.Buffer

    for _, stmt := range bs.Statements {
        out.WriteString(stmt.String())
    }

    return out.String()
}

type Identifier struct {
    Token token.Token
    Value string
}

func (i *Identifier) TokenLiteral() string {
    return i.Token.Literal
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) String() string {
    return i.Value
}

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

func (il *IntegerLiteral) TokenLiteral() string {
    return il.Token.Literal
}

func (il *IntegerLiteral) expressionNode() {

}

func (il *IntegerLiteral) String() string {
    return il.Token.Literal
}

type StringLiteral struct {
    Token token.Token
    Value string
}

func (sl *StringLiteral) TokenLiteral() string {
    return sl.Token.Literal
}

func (sl *StringLiteral) expressionNode() {

}

func (sl *StringLiteral) String() string {
    var out bytes.Buffer

    out.WriteString("\"")
    out.WriteString(sl.Value)
    out.WriteString("\"")

    return out.String()
}

type PrefixExpression struct {
    Token token.Token
    Operator string // -, !
    Right Expression
}

func (pe *PrefixExpression) TokenLiteral() string {
    return pe.Token.Literal
}

func (pe *PrefixExpression) expressionNode() {

}

func (pe *PrefixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(pe.Operator)
    out.WriteString(pe.Right.String())
    out.WriteString(")")

    return out.String()
}

type InfixExpression struct {
    Token token.Token
    Left Expression
    Operator string
    Right Expression
}

func (ie *InfixExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *InfixExpression) expressionNode() {

}

func (ie *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString(" " + ie.Operator + " ")
    out.WriteString(ie.Right.String())
    out.WriteString(")")

    return out.String()
}

type Boolean struct {
    Token token.Token
    Value bool
}

func (b *Boolean) TokenLiteral() string {
    return b.Token.Literal
}

func (b *Boolean) expressionNode() {

}

func (b *Boolean) String() string {
    return b.Token.Literal
}

type IfExpression struct {
    Token token.Token
    Condition Expression
    Consequence *BlockStatement
    Alternative []*IfExpression
    Default *BlockStatement
}

func (ie *IfExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *IfExpression) expressionNode() {

}

func (ie *IfExpression) String() string {
    var out bytes.Buffer

    out.WriteString("if")
    out.WriteString(ie.Condition.String())
    out.WriteString(" ")
    out.WriteString(ie.Consequence.String())


    for _, alt := range ie.Alternative {
        // if alt != nil {
        out.WriteString("else if ")
        out.WriteString(alt.String())
        // }
    }
    
    if ie.Default != nil {
        out.WriteString("else ")
        out.WriteString(ie.Default.String())
    }

    return out.String()
}

type FunctionLiteral struct {
    Token token.Token
    Parameters []*Identifier
    Body *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string {
    return fl.Token.Literal
}

func (fl *FunctionLiteral) expressionNode() {

}

func (fl *FunctionLiteral) String() string {
    var out bytes.Buffer

    params := []string{}

    for _, p := range fl.Parameters {
        params = append(params, p.String())
    }

    out.WriteString("fn")
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") ")
    out.WriteString(fl.Body.String())

    return out.String()
}

type CallExpression struct {
    Token token.Token
    Function Expression
    Arguments []Expression
}

func(ce *CallExpression) TokenLiteral() string {
    return ce.Token.Literal
}

func (ce *CallExpression) expressionNode() {

}

func (ce *CallExpression) String() string {
    var out bytes.Buffer

    args := []string{}
    for _, a := range ce.Arguments {
        args = append(args, a.String())
    }

    out.WriteString(ce.Function.String())
    out.WriteString("(")
    out.WriteString(strings.Join(args, ", "))
    out.WriteString(")")

    return out.String()
}

type ArrayLiteral struct {
    Token token.Token
    Elements []Expression
}

func (al *ArrayLiteral) TokenLiteral() string {
    return al.Token.Literal
}

func (al *ArrayLiteral) expressionNode() {

}

func (al *ArrayLiteral) String() string {
    var out bytes.Buffer

    elements := []string{}
    for _, e := range al.Elements {
        elements = append(elements, e.String())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")

    return out.String()
}

type IndexExpression struct {
    Token token.Token
    Left Expression
    Index Expression
}

func (ie *IndexExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *IndexExpression) expressionNode() {

}

func (ie *IndexExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString("[")
    out.WriteString(ie.Index.String())
    out.WriteString("])")

    return out.String()
}
