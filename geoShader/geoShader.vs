#version 410 core
layout (location = 0) in vec2 aPos;
layout (location = 1) in float aColor;

out VS_OUT {
    float color;
} vs_out;

void main()
{
    vs_out.color = aColor;
    gl_Position = vec4(aPos.x, aPos.y, 0.0, 1.0); 
}