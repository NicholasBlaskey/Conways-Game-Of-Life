# Conways-Game-Of-Life

![](example.gif)

Conways game of life visually implemented in Go using opengl. Different techniques are used to try to optimize framerate.

### Naive

Naive is rendering every square with a call to draw arrays setting a new uniform for position and color each time.

~500fps at 25x25

~80fps at 100x100

~3.7fps at 500x500

~0.8fps at 1000x1000

### One color

One color only renders one color's squares and the other color's squares are rendered because of the background being set to that each time.

~625fps at 25x25

~600fps at 100x100

~78fps at 500x500

~11fps at 100x1000

### Only updated tiles

Only updated tiles keeps track of the board before and the board now and only needs to render the tiles that actually get updated.

### Instancing 

Instancing takes advantage of using the Instancing feature in OpenGl to avoid CPU and GPU communication when drawing a lot of the same objects. 

~600fps at 25x25

~585fps at 100x100

~265fps at 500x500

~80fps at a 1000x1000





