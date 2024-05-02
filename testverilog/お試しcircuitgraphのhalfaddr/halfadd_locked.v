module half_adder (key_0, key_1, y, x, s, c);
  input key_0;
  input key_1;
  input y;
  input x;

  output s;
  output c;

  wire key_gate_0;
  wire key_gate_1;
  wire s;
  wire c;

  xnor g_0(key_gate_0, key_0, y);
  xnor g_1(key_gate_1, key_1, x);
  xor g_2(s, key_gate_0, key_gate_1);
  and g_3(c, key_gate_0, key_gate_1);
endmodule
