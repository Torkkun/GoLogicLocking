import circuitgraph as cg
from logiclocking import locks, write_key

#c = cg.from_lib("c880")
c = cg.from_file("/home/ai-i2lab/projects/CircuitGraph test/workspace/addsysd.v")

c = cg.tx.syn(c, suppress_output=True)
num_keys = 2
cl, k = locks.xor_lock(c, num_keys)

cg.to_file(cl, "addsys_locked.v")
write_key(k, "addsys_locked_key.txt")
