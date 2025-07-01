grammar selectlang;

// -------------------------------------------------------------------
// Parser Rules
// -------------------------------------------------------------------

// A complete SELECT statement over a stream.
select_stmt
    : SELECT columns FROM stream where_clause? EOF   #SelectStatement
    ;

// Currently, only "*" is supported for columns.
columns
    : STAR                                            #AllColumns
    ;

// The stream name (e.g. Notification).
stream
    : IDENTIFIER                                      #StreamName
    ;

// The optional WHERE clause.
where_clause
    : WHERE expression                                #WhereClause
    ;

// Boolean expressions with OR (lowest precedence).
expression
    : expression OR and_expr                          #OrExpression
    | and_expr                                        #AndToExpression
    ;

// AND expressions (higher precedence).
and_expr
    : and_expr AND primary_expr                       #AndExpression
    | primary_expr                                    #PrimaryExpression
    ;

// Primary expressions: either a parenthesized expression or a predicate.
primary_expr
    : LPAREN expression RPAREN                        #ParenExpression
    | predicate                                       #PredicateExpression
    ;

// A predicate can be a comparison, regex match, an IN predicate, or a HAS predicate.
predicate
    : field comp_operator value                       #ComparisonPredicate
    | field regex_operator regex_value                #RegexPredicate
    | field IN value_list                             #InPredicate
    | field HAS value                                 #HasPredicate
    ;

// A comma‐separated list of literal values.
value_list
    : value (COMMA value)*                            #ValueList
    ;

// The field rule distinguishes between notification object fields and log fields.
field
    : obj_field                                       #ObjField
    | log_field                                       #LogField
    ;

// The notification object "obj" (allowed properties: ID, Name, Operation).
obj_field
    : OBJ DOT (ID_FIELD | NAME_FIELD | OP_FIELD)      #ObjFieldAccess
    ;

// The log entry "log" (allowed properties: Operation, Path, Name, Value).
log_field
    : LOG DOT (OP_FIELD | PATH_FIELD | NAME_FIELD | VAL_FIELD)  #LogFieldAccess
    ;

// A literal value: either a number or a string.
value
    : NUMBER                                          #NumberValue
    | STRING                                          #StringValue
    ;

// Comparison operators.
comp_operator
    : EQ                                              #EqualsOp
    | NE                                              #NotEqualsOp
    | GT                                              #GreaterThanOp
    | LT                                              #LessThanOp
    | GE                                              #GreaterOrEqualOp
    | LE                                              #LessOrEqualOp
    ;

// The regex operator.
regex_operator
    : REGEX_OP                                        #RegexOp
    ;

// A regex value is provided as a string literal.
regex_value
    : STRING                                          #RegexValue
    ;

// -------------------------------------------------------------------
// Lexer Rules
// -------------------------------------------------------------------

SELECT: 'SELECT';
FROM: 'FROM';
WHERE: 'WHERE';
AND: 'AND';
OR: 'OR';
IN: 'IN';
HAS: 'HAS';
STAR: '*';
DOT: '.';

EQ: '==';
NE: '!=';
GT: '>';
LT: '<';
GE: '>=';
LE: '<=';

// The regex operator is represented by '~='.
REGEX_OP: '~=';

// Parentheses and comma.
LPAREN: '(';
RPAREN: ')';
COMMA: ',';

// ---
// Fixed tokens for object fields. These must be defined before IDENTIFIER.
// "obj" is the notification object.
OBJ: 'obj';
// "log" is a single value change entry.
LOG: 'log';

// Allowed properties for the notification object.
ID_FIELD: 'ID';
NAME_FIELD: 'Name';
OP_FIELD: 'Operation';

// Allowed properties for the log entry.
PATH_FIELD: 'Path';
VAL_FIELD: 'Value';

// An identifier for stream names (e.g. "Notification"). This comes after fixed tokens.
IDENTIFIER: [a-zA-Z_][a-zA-Z0-9_]*;

// A number literal.
NUMBER: [0-9]+ ('.' [0-9]+)?;

// A string literal enclosed in single quotes (supports escapes).
STRING: '\'' ( ~['\\] | '\\' . )* '\'';

WS: [ \t\r\n]+ -> skip;