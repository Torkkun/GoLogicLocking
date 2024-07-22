import circuitgraph as cg
from logiclocking import attacks

c = cg.from_file("./test.v")
#key = {"key_0": True, "key_1": False}
assumption = {"key_gate_0", "key_gate_1"}
#attacks.mitter_attack(c, key, )
c2 = c.copy()
m = cg.tx.miter(c, c2)
#r = cg.sat.solve(m, assumptions={"sat":True})
#cg.sat.solve(c, assumption={})
print(r)