import circuitgraph as cg
from logiclocking import locks, write_key

c = cg.from_file("./addveri/muladd.v")

c = cg.tx.syn(c, suppress_output=True)
num_keys = 2
cl, k = locks.xor_lock(c, num_keys)
