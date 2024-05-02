import circuitgraph as cg
from logiclocking import locks, write_key

a = cg.logic.half_adder()
c = cg.tx.syn(a, suppress_output=True)

num_keys = 2
cl, k = locks.xor_lock(c, num_keys)

cg.to_file(cl, "halfadd_locked.v")
write_key(k, "halfadder_locked_key.txt")
