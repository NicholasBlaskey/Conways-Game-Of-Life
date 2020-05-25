#version 410 core
layout (location = 0) in vec2 aPos;

uniform vec2 aOffset;
uniform float fragColor;

out float fColor;

void main()
{
    fColor = fragColor;
    gl_Position = vec4(aPos + aOffset, 0.0, 1.0);
}