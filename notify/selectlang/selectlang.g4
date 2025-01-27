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
    | NOT expression
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
    : loggerOp (':' regex)? mapVarExpr? ('AND' LPAREN valueComparison RPAREN)? 
    ;

mapVarExpr
    : EQ STRING
    ;

loggerOp
    : ('add' | 'remove' | 'update' | 'acknowledge' | 'no-change' | 'all') 
      (',' ('add' | 'remove' | 'update' | 'acknowledge' | 'no-change' | 'all'))*
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
    : EQ | NE | '>' | '<' | '>=' | '<=' | 'before' | 'after' | 'regexp'
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