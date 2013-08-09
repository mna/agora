%{
package yacc

import (
	"fmt"
)
%}

%union {
	lit string
}

%token <lit> string_lit identifier

%%

Module : ImportStmt

ImportStmt : "import" ImportSpec ";" |
					 	 "import" "(" ImportSpec ";" ")"

ImportSpec : identifier ImportPath |
						 ImportPath

ImportPath : string_lit { fmt.Println($1) }
