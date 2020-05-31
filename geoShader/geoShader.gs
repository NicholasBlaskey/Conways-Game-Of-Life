#version 410 core
layout (points) in;
layout (triangle_strip, max_vertices = 4) out;

in VS_OUT {
    float color;
} gs_in[];

out float fColor;

uniform float xOffset;
uniform float yOffset;

void build_triangle(vec4 position)
{    
    fColor = gs_in[0].color; // gs_in[0] since there's only one input vertex
    gl_Position = position + vec4(-xOffset, -yOffset, 0.0, 0.0);    
    EmitVertex();   
    gl_Position = position + vec4(xOffset, -yOffset, 0.0, 0.0);    
    EmitVertex();
    gl_Position = position + vec4(-xOffset, yOffset, 0.0, 0.0);    
    EmitVertex();
    gl_Position = position + vec4(xOffset, yOffset, 0.0, 0.0);    
    EmitVertex();
    
    EndPrimitive();
}

void main() {    
    build_triangle(gl_in[0].gl_Position);
}
