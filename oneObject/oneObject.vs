#version 410 core
layout (location = 0) in vec2 aPos;
layout (location = 1) in float fragColor;

//uniform vec2 aOffset;
//uniform float fragColor;

out float fColor;

void main()
{
    fColor = fragColor;
    gl_Position = vec4(aPos, 0.0, 1.0);
}