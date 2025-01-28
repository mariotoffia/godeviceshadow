grammar selectlang;

filter
    : expression EOF
    ;

expression
    : primaryExpr
    | loggerExpr
    | LPAREN expression RPAREN
    | expression AND expression
    | expression OR expression
    ;

primaryExpr
    : idExpr
    | nameExpr
    | operationExpr    
    ;

idExpr
    : ID regexOrString
    ;

nameExpr
    : NAME regexOrString
    ;

operationExpr
    : OPERATION operations
    ;

operations
    : ('report' | 'desired' | 'delete' | 'all') (',' ('report' | 'desired' | 'delete'))*
    ;

loggerExpr
    : loggerOp (':' regex)? mapVarExpr? loggerConstraints? 
    ;

mapVarExpr
    : EQ STRING
    ;

loggerOp
    : ('add' | 'remove' | 'update' | 'acknowledge' | 'no-change' | 'all') 
      (',' ('add' | 'remove' | 'update' | 'acknowledge' | 'no-change' | 'all'))*
    ;

loggerConstraints
    : WHERE LPAREN valueComparison RPAREN
    ;

valueComparison
    : valueCondition ( OR valueCondition )*
    ;

valueCondition
    : valueFactor ( AND valueFactor )*
    ;

valueFactor
    :'value' compareOp constantOrRegex
    | LPAREN valueComparison RPAREN
    ;

compareOp
    : EQ | NE | GT | LT | GE | LE | BEFORE | AFTER
    ;

constantOrRegex
    : NUMBER                               # NumericLiteral
    | STRING                               # StringLiteral
    | TIME                                 # TimeLiteral
    | regex                                # RegexLiteral
    ;

regexOrString
    : REGEX
    | STRING
    ;

regex
    : REGEX
    ;

WHERE:
    'WHERE'
    ;

LPAREN
    : '('
    ;

RPAREN
    : ')'
    ;

EQ
    : '=='
    ;

NE
    : '!='
    ;

GT
    : '>'
    ;

LT
    : '<'
    ;

GE
    : '>='
    ;

LE
    : '<='
    ;

BEFORE
    : 'before'
    ;

AFTER
    : 'after'
    ;

ID
    : 'id:'
    ;

AND:
    'AND'
    ;

OR:
    'OR'
    ;

NOT:
    'NOT'
    ;

NAME
    : 'name:'
    ;

OPERATION
    : 'operation:'
    ;

NUMBER
    : ('-')? [0-9]+ ('.' [0-9]+)?
    ;

STRING
    : '\'' .*? '\''
    ;

TIME
    : [0-9]{4} '-' [0-9]{2} '-' [0-9]{2} 'T' [0-9]{2} ':' [0-9]{2} ':' [0-9]{2} 'Z'
    ;

REGEX
    : '/' .*? '/'
    ;

WS
    : [ \t\r\n]+ -> skip
    ;