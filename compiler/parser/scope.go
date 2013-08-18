package parser

type Scope struct {
	def    map[string]*Symbol
	parent *Scope
	p      *Parser
}

func (s *Scope) define(n *Symbol) *Symbol {
	t, ok := s.def[n.val.(string)]
	if ok {
		if t.res {
			s.p.error(t, "already reserved")
		} else {
			s.p.error(t, "already defined")
		}
	}
	s.def[n.val.(string)] = n
	n.res = false
	n.lbp = 0
	n.nudfn = itself
	n.ledfn = nil
	n.stdfn = nil
	return n
}

// The find method is used to find the definition of a name. It starts with the
// current scope and seeks, if necessary, back through the chain of parent scopes
// and ultimately to the Symbol table. It returns Symbol_table["(name)"] if it
// cannot find a definition.
func (s *Scope) find(id string) *Symbol {
	for scp := s; scp != nil; scp = scp.parent {
		if o, ok := scp.def[id]; ok {
			return o
		}
	}
	if o, ok := s.p.tbl[id]; ok {
		return o
	}
	return s.p.tbl[_SYM_NAME]
}

func (s *Scope) reserve(n *Symbol) {
	if n.ar != arName || n.res {
		return
	}
	val, ok := n.val.(string)
	if !ok {
		s.p.error(n, "expected a string value")
	}
	if t, ok := s.def[val]; ok {
		if t.res {
			return
		}
		if t.ar == arName {
			s.p.error(n, "already defined")
		}
	}
	s.def[val] = n
	n.res = true
}
