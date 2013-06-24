local c = {}
function c.d(arg)
  io.write(arg)
end

local a = {}
a.b = c

a.b.d("toto")
