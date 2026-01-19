# Notabor Engine
a phsyics engine made in **golang** designed to make game dev easier by providing the following:
- ## Function based architecture
  there is a `Runnable` type function, which is a function that can return an error.
  
  a runnable can be appended to a Loop whether it's `FixedHzLoop` or `RenderLoop` by appending the function to the loop.
  
  the runnable then runs at a specific ticks per seconds (`Hz`) when the loop runs, no need for multiplying anything by delta time and no need for intense Object Oriented code.
- ## Loop based running
  as mentioned, a loop is `FixedHzLoop` or `RenderLoop`, it contains a list of `Runnable`, when it runs it and the Runnable returns an error, it deletes the Runnable from the list and continues, otherwise it runs without deleting
  this makes it so there is no need for any loops made by the player, and keeps the code clean and easy to read
- ## Native Multi-Window support
  a window can be created via a config, each window has one `RenderLoop` and zero or more `FixedHzLoop`, that way we can ensure separation between windows
  different window modes, `Windowed`, `BorderlessWindowed` and `FullScreen`
- ## Notagl
  It contains gl2 to render 2D objects and gl3 to render 3D objects. polygon is used to store po2 and draw shapes such as squares, rectangles, circles, and other polygons. mesh is used to store po3 and draw 3D shapes. render is responsible for rendering 2D or 3D objects using clean and easy-to-read code.
- ## Notamath
  contains mat3 to perform all 2D transformations using homogeneous coordinates. mat4 is similar to mat3, but the difference is that mat4 is specialized for 3D transformations. po2 is used to store 2D coordinates to draw polygons. po3 is similar to po2, but it is used to draw any 3D shape. vec2 and vec3 are used for collision tracking with the MTV (Minimum Translation Vector) algorithm.

# Documentation

## notacore
  - ### Objects:
    - `Settings` type: `struct`
    - `Engine` type: `struct`
    - `Runnable` type: `func() error`
    - `FixedHzLoop` type: `struct`
    - `RenderLoop` type: `struct`  
    - `WindowConfig` type: `struct`
    - `WindowType` type: `enum`
    - `windowRuntime2D` type: `struct`
    - `windowRunTime3D` type: `struct`
    - `Window` type: `interface`
    - `GlfwWindow2D` type: `struct`
    - `GlfwWindow3D` type: `struct`
    
  - ### Settings
    - `Vsync` type `bool`
  - ### Engine
    Handles the main engine operations like looping, initializing windows etc 
    - #### content
      - `Windows2D` type: `[]*GlfwWindow2D`
      - `Windows3D` type: `[]*GlfwWindow3D`
      - `Settings` type: `*Settings`
      - `WindowManager` type: `*GLFWWindowManager`
    - #### functions
      - `func (e *Engine) Run() error` runs all loops associated with all windows and keeps track of the lifecycle
      - `func (e *Engine) AllWindowsClosed() bool` checks if all windows are closed
      - `func (e *Engine) Shutdown()` terminates GLFW (stops the whole program)
      - `func (e *Engine) InitPlatform() error` locks the primary thread for rendering, uses the correct OS files and initializes window manager (MUST BE DONE BEFORE USING ANY ENGINE OPERATIONS)
      - `func (e *Engine) CreateWindow2D(cfg WindowConfig) (*GlfwWindow2D, error)` creates a new window which has a 2D renderer and GLBackend2D, and adds its pointer to `Windows2D`
      - `func (e *Engine) CreateWindow3D(cfg WindowConfig) (*GlfwWindow3D, error)` creates a new window which has a 3D renderer and GLBackend3D, and adds its pointer to `Windows3D`
  - ### Runnable
    signature: `type Runnable func() error` used by the engine to run logic
- ### FixedHzLoop
  Handles repeated execution of tasks such as logic updates.
  - #### content
    - `Loop` type `interface`
    - `Runnables` type: `[]Runnable`
    - `Hz` type `float32`
  - #### functions
    - `func (l *FixedHzLoop) Start()` Uses concurrency and multithreading  to execute runnables without blocking the main thread and handles removal of runnables that return errors
    - `func (l *FixedHzLoop) Stop()` Stops the loop and cleans up resources
    - `func (l *FixedHzLoop) Remove(i int)` Removes a runnable from the loop by index.
    - `func (l *FixedHzLoop) Alpha(now time.Time) float32` Returns an interpolation factor between the last fixed logic tick and the next one,
    used for smooth rendering.
- ### RenderLoop
  - #### content
    - `Runnables` type: `[]Runnable`
    - `MaxHz` type: `float32`
    - `LastTime` type: `time.Time`
  - #### functions
    - `func (r *RenderLoop) Render()` Runs all Runnables once per call in main thread
- ### WindowConfig
  - #### content
    - `X` type: `int` (X coordinate of the origin of the screen)
    - `Y` type: `int` (Y coordinate of the origin of the screen)
    - `W` type: `int` (screen width)
    - `H` type: `int` (screen height)
    - `Title` type: `string`
    - `Resizable` type: `bool`
    - `Type` type: `WindowType`
    - `LogicLoops` type: `[]*FixedHzLoop`
    - `RenderLoop` type `*RenderLoop`
- ### WindowType
  - #### values
    - `Windowed`
    - `Fullscreen`
    - `Borderless`
 - ### windowRuntime2D
   - #### content
     - `Renderer` type: `*notagl.Renderer2D`
- ### windowRuntime3D
   - #### content
     - `Renderer` type: `*notagl.Renderer3D`
- ### Window
  - #### functions
  - `GetConfig() *WindowConfig`
	- `GetRuntime() *WindowBaseRuntime`
	- `MakeContextCurrent()`
	- `SwapBuffers()`
	- `ShouldClose() bool`
	- `RunRenderer()`
- ### GlfwWindow2D
  implements `Window`
  - #### content
    - `ID` type: `int`
    - `Handle` type: `*glfw.Window`
    - `Config` type:  `WindowConfig`
    - `RunTime` type: `windowRunTime2D`
  - #### functions
    - `func (w *GlfwWindow2D) MakeContextCurrent()` makes the window marked as `current` for glfw
    - `func (w *GlfwWindow2D) SwapBuffers()` swaps the buffers of the window
    - `func (w *GlfwWindow2D) ShouldClose() bool` checks if window should close
    - `func (w *GlfwWindow2D) Close()` sets `ShouldClose` to true
    - `func (w *GlfwWindow2D) Size() (int, int)` returns width and height of the window
    - `func (w *GlfwWindow2D) Position() (int, int)` returns the x and y coordinates of the window origin
    - `func (w *GlfwWindow2D) SetFullscreen(fullscreen bool) error` sets window to fullscreen mode
    - `func (w *GlfwWindow2D) SetBorderless(borderless bool) error` sets window to borderless mode
    - `func (w *GlfwWindow2D) SetWindowed(x, y, width, height int) error` sets window to windowed mode
- ### GlfwWindow3D
  the same as GlfwWindow2D, except `RunTime` is of type: `windowRunTime3D`

## notagl
  - ### Objects
    - `DrawOrder2D` type: `struct`
    - `DrawOrder3D` type: `struct`
    - `Renderer2D` type: `struct`
    - `Renderer3D` type: `struct`
    - `Polygon` type: `struct`
    - `func Triangulate2D(polygon []notamath.Po2) []notamath.Po2` type: `function`. Takes a list of points of a polygon, and links all vertices via triangles for GPU rendering
  - ### DrawOrder2D
    - #### content
      - `Vertices` type: `[]notamath.Po2` vertices of the polygon drawn
  - ### DrawOrder3D
    - #### content
      - `Vertices` type: `[]notamath.Po3` vertices of the mesh drawn
  - ### Renderer2D
    - #### content
      - `Orders` type: `[]DrawOrder2D`
    - #### functions
      - `func (r *Renderer2D) Submit(p Polygon, alpha float32)` creates a new order and appends to the renderer orders it with the alpha value of the polygon (note: set alpha to 1 if static)
  - ### Renderer3D
    (WORK IN PROGRESS)
## notamath
