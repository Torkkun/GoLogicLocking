module lut(in, out);

input [1:0] in;
output out;

always @(in)
case(in)
    2'b00 : out = 0;
    2'b01 : out = 1;
    2'b10 : out = 1;
    2'b11 : out = 0;
endcase

endmodule