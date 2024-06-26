/* Generated by Yosys 0.38+141 (git sha1 078b876f5, clang++ 14.0.0-1ubuntu1.1 -fPIC -Os) */

module hlfaddr(A, B, C, S);
  input A;
  wire A;
  input B;
  wire B;
  output C;
  wire C;
  output S;
  wire S;
  wire s1;
  wire s2;
  wire s3;
  assign s3 = ~s2;
  assign s1 = A | B;
  assign s2 = A & B;
  assign S = s1 & s3;
  assign C = s2;
endmodule
