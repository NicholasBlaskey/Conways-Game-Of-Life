#version 410 core
layout (location = 0) in vec2 aPos;
layout (location = 1) in vec2 aOffset;
layout (location = 2) in float fragColor;

out float fColor;

void main()
{
    fColor = fragColor;
    gl_Position = vec4(aPos + aOffset, 0.0, 1.0);
}