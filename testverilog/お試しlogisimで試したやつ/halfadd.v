module hlfaddr(
    input A,
    input B,
    output C,
    output S
);

wire s1, s2, s3;

assign s1 = A | B;
assign s2 = A & B;
assign s3 = ~s2;
assign S = s1 & s3;
assign C = s2;
endmodule
