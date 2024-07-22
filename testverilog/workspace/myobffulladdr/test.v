module test (key_1, A, key_0, B, CIN, COUT, Q);
//no parameter
  input key_1;
  input A;
  input key_0;
  input B;
  input CIN;
  output COUT;
  output Q;
//no register decl
//no event decl
  wire key_gate_0;
  wire key_gate_1;
  wire _01_;
  wire _02_;
  wire _03_;
  wire _04_;
  wire _00_;
  assign COUT = _02_ | _03_;
  assign key_gate_0 = CIN ^ _04_;
  assign _04_ = B ^ A;
  assign _00_ = B & A;
  assign key_gate_1 = CIN & B;
  assign _03_ = _01_ | _00_;
  assign _02_ = A & CIN;
  assign Q = ~(key_gate_0 ^ key_0);
  assign _01_ = key_gate_1 ^ key_1;

endmodule