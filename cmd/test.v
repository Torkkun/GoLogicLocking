module test (key_0, key_1, A, B, CIN, COUT, Q);
//no parameter
  input key_0;
  input key_1;
  input A;
  input B;
  input CIN;
  output COUT;
  output Q;
//no register decl
//no event decl
  wire _01_;
  wire _02_;
  wire _03_;
  wire _04_;
  wire _00_;
  wire llw_0;
  wire llw_1;
  assign COUT = _02_ | _03_;
  assign llw_1 = B ^ A;
  assign Q = CIN ^ _04_;
  assign _01_ = CIN & B;
  assign llw_0 = A & CIN;
  assign _03_ = _01_ | _00_;
  assign _00_ = B & A;
  assign _02_ = llw_0 ^ key_0;
  assign _04_ = ~(llw_1 ^ key_1);

endmodule