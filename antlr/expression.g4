
// 扩展语法时，运算符优先级参考：https://baike.baidu.com/item/%E8%BF%90%E7%AE%97%E7%AC%A6%E4%BC%98%E5%85%88%E7%BA%A7/4752611

grammar expression;


node
    : expression
    ;

expression
    : binary
    ;

binary
    : level2_binary (First_level_op level2_binary)*
    ;

// 思考：为什么 参数也必须是 signedAtom
//      因为递归推导可以保证可以推导出任意预期的表达式
//      使用是优先级更低的 非终结符 会破坏文法对优先级的定义、比如 level2_binary，那么 1 * 2+3 中 2+3 是level2_binary，但将其作为一个整体完全有问题
//      使用优先级更高的 非终结符 则会使得文法不完整，比如 1 * 2 * 3 则可能永远推导不出（这个有待论证，因为递归的推导、估计也是可以的）
//
//在代码中使用 signedAtomBinaryxx 方法进行解析，因为该二元表达式的 参数是 signedAtom
level2_binary
    : signedAtom (Second_level_op signedAtom)*
    ;

// note 如果加入其他优先级运算符，则将当前最高优先级、比如 level2_binary 的 signedAtom 替换为 level_n_binary
//      level_n_binary 定义如下
//      level_n_binary
//              : signedAtom (n_level_op signedAtom)*
//              ;

signedAtom
    : unary
    | atom
    ;

unary
    :UnaryOp atom
    ;

atom
    : Variable
    | String
    | Number
    | func
    | sub_node
    ;

func
    : Func_name LParen node (Comma node)*  RParen
    ;

sub_node
    : LParen expression RParen
    ;

UnaryOp
    : '+'
    | '-'
    ;

First_level_op
    : '+'
    | '-'
    ;

Second_level_op
    : '*'
    | '/'
    | '%'
    ;

Func_name
    : (Letter | '_')+ Letter*
    ;

Letter
     : 'a' .. 'z' | 'A' .. 'Z'
     ;

LParen
    : '('
    ;

RParen
    : ')'
    ;

Comma
    : ','
    ;

Variable
    : Letter
    ;

// todo 貌似应该以 "" 起始
String
    : Letter
    ;

// 范围选择符使用两个点（..）来表示范围的开始和结束
Number
    : ('0'..'9')+
    ;