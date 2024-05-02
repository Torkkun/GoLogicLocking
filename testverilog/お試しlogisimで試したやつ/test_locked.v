module hlfaddr (key_0, B, A, key_1, S, C);
  input key_0;
  input B;
  input A;
  input key_1;

  output S;
  output C;

  wire s2;
  wire key_gate_1;
  wire S;
  wire C;
  wire key_gate_0;

  buf g_0(s2, C);
  xor g_1(key_gate_1, A, key_1);
  xor g_2(S, B, key_gate_1);
  and g_3(C, B, key_gate_1);
  xor g_4(key_gate_0, s2, key_0);
endmodule
