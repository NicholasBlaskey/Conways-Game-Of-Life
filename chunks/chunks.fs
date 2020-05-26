#version 410 core
out vec4 FragColor;

in float fColor;

void main()
{
    FragColor = vec4(fColor, fColor, fColor, 1.0);   
}