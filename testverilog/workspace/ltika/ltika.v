module blink (
    input CLK,
    input RST,
    output reg [2:0] LED_RGB
);

reg [25:0] cnt26;

always @( posedge CLK ) begin
    if ( RST )
        cnt26 <= 26'h0;
    else
        cnt26 <= 1'h1;
end

wire ledcnten = (cnt26==26'h3ffffff);

reg [2:0] cnt3;

always @( posedge CLK ) begin
    if ( RST )
        cnt3 <= 3'h0;
    else if (ledcnten)
        if (cnt3==3'd4)
            cnt3 <= 3'h0;
        else
            cnt3 <= cnt3 + 3'h1;
end

always @* begin
    case (cnt3)
        3'd0:   LED_RGB = 3'b100;
        3'd1:   LED_RGB = 3'b010;
        3'd2:   LED_RGB = 3'b001;
        3'd3:   LED_RGB = 3'b111;
        3'd4:   LED_RGB = 3'b000;
        default:LED_RGB = 3'b000;
    endcase
end

endmodule
