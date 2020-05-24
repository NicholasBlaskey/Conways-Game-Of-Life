# Conways-Game-Of-Life

![](example.gif)

Conways game of life visually implemented in Go using opengl. Different techniques are used to try to optimize framerate.

### Instancing 

Instancing takes advantage of using the Instancing feature in OpenGl to avoid CPU and GPU communication when drawing a lot of the same objects. 

~80fps at a 1000 by 1000 board


