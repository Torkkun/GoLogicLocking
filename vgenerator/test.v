module test_module (A, B, C, D);
//no parameter
  input A;
  input B;
  input C;
  output D;
  input [4:0] E;
//no register decl
//no event decl
  wire _1_;
  wire [2:0] _2_;
  assign out1 = in1 & in2;
  assign out1 = ~(in1 & in2);
  assign out1 = in1 | in2;
  assign out1 = ~(in1 | in2);
  assign out1 = in1 ^ in2;
  assign out1 = ~(in1 ^ in2);
  assign out1 = in1;
  assign out1 = ~in1;

endmodule