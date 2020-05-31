# Conways-Game-Of-Life

![](example.gif)

Conways game of life visually implemented in Go using opengl. Different techniques are used to try to optimize framerate.

| Method               | fps at 25x25 | fps at 100x100 | fps at 500x500 | fps at 1000x1000 |
| -------------------- | -----: | -----: | -----: | -----: |
| Naive                |  500            | 80  | 3.7 | 0.8 |
| One Color            |  625          | 600 | 78 | 11 |
| Only Updated Tiles*  | 400-430       | 250-350 | 8-60 | 1-13|
| Instancing           | 600           | 585 | 265 | 80 |
| One Color Instancing | 725           | 700 | 300 | 110 |
| Chunks               | 650 | 325 | 21 | 6 |
| One Object           | 725 | 600 | 175 | 56 |
| Geometry Shader      | 675 | 625 | 225 | 100 |

(*Two records are shown for start framerate and the framerate after 30 seconds) 

### Naive

Naive is rendering every square with a call to draw arrays setting a new uniform for position and color each time.


### One Color

One color only renders one color's squares and the other color's squares are rendered because of the background being set to that each time.


### Only Updated Tiles

Only updated tiles keeps track of the board before and the board now and only needs to render the tiles that actually get updated. To do this I needed to utilize a framebuffer because just using the screen was running into issues. As a result of the framebuffer it was very diffucult to get resizing working properly so I didn't.

There are two fps readings measured. One is from the start and the second is 30 seconds later. This is relevant to rendering only updated tiles because at the start there is a 50-50 chance for tiles to be alive or dead. This causes a lot of updating to happen until it settles more. This causes a lot more time to be used at the start.

### Instancing 

Instancing takes advantage of using the Instancing feature in OpenGl to avoid CPU and GPU communication when drawing a lot of the same objects. 

### One Color Instancing

This method is the same as instancing but we are only drawing the white squares. Instead of updating a color VBO each frame we are updating a position VBO each frame.

### Chunks

This method will take a chunk size (we are going to assume that numX divides this size for simplicity). We will then generate 2^chunkSize VAOs that represent each of the possible color combinations of the chunksize. We will then lookup the chunksize.

### One object

This method creates a single object of vertices and only the color attrib is updated each iteration. Only one rendering call is needed but a lot more vertices are stored.

### Geometry Shader

This method uses the geometry shader. Instead of needing to pass 6 vertices in we can pass in a single point for each triangle.

### Shoutout

https://learnopengl.com/

Concepts and code ideas taken from the LearnOpenGL tutorial.
