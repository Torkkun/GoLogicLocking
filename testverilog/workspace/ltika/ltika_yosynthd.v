/* Generated by Yosys 0.40+50 (git sha1 0f9ee20ea, g++ 11.4.0-1ubuntu1~22.04 -fPIC -Os) */

module blink(CLK, RST, LED_RGB);
  input CLK;
  wire CLK;
  output [2:0] LED_RGB;
  wire [2:0] LED_RGB;
  input RST;
  wire RST;
  wire [25:0] cnt26;
  wire [2:0] cnt3;
  assign LED_RGB = 3'h4;
  assign cnt26[25:1] = 25'h0000000;
  assign cnt3 = 3'h0;
endmodule
