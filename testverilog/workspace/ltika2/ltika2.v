module blink(
    input CLK,
    output [3:0] OUT
);

parameter CNT_1SEC = 27'd124999999;  // 125MHz clk for 1sec
reg [26:0] cnt = 27'd0;
reg onoff = 1'd0;

always @(posedge CLK) begin
    if (cnt == CNT_1SEC) begin
        cnt <= 27'd0;
        onoff <= ~onoff;
    end
    else begin
        cnt <= cnt + 27'd1;
    end
end

assign OUT = {onoff, onoff, onoff, onoff};

endmodule
