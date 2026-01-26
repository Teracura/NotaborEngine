# Notabor Engine
a phyics engine made in **golang** designed to make game dev easier by providing the following:
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
      - `Windows` type `[]Window`
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
- ### Objects
  - `Mat3`type`struct`
  - `Mat4`type`struct`
  - `Po2`type`struct`
  - `Po3`type`struct`
  - `Transform2D`type`struct`
  - `Transform3D `type `struct`
  - `AxisMask` type `enum`
  - `Vec2`type`struct`
  - `Vec3`type`struct`
- ### Mat3
  - #### content
    - `Mat3` type `[9]float32` 
  - #### functions
    - `func Mat3Identity() Mat3` Returns an identity 3×3 matrix.
    - `func Mat3Translation(t Vec2) Mat3` Translates a shape by vector t
    - `func Mat3Scale(s Vec2) Mat3` Scales a shape by vector s.
    - `func Mat3Rotation(rad float32) Mat3` Rotates a shape by rad radians.
    - `func Mat3Shear(kx, ky float32) Mat3` Shears a shape by kx and ky.
    - `func Mat3TRS(pos Vec2, rot float32, scale Vec2) Mat3`Translates by pos, rotates by rot, and scales by scale.
    - `func (m Mat3) Mul(b Mat3) Mat3`Multiplies matrix m by matrix b  (note: multiplication order may affect the result).
    - `func (m Mat3) TransformPo2(p Po2) Po2` Applies a linear transformation to a point.
    - `func (m Mat3) TransformVec2(v Vec2) Vec2` Applies a linear transformation to a vector.
    - `func (m Mat3) Transpose() Mat3` Returns the transposed matrix.
    - `func (m Mat3) Det() float32` Returns the matrix determinant.
    -  `func (m Mat3) InverseAffine() Mat3` Returns the inverse affine transformation.
    -  `func (m Mat3) String() string` Returns a string representation of the matrix.
-  ### Mat4
  - #### content
    - `Mat4` type `[16]float32` 
  - #### functions
    - `func Mat4Identity() Mat4` Returns an identity 4×4 matrix.
    - `func Mat4Translation(t Vec3) Mat4` Translates a shape by vector t.
    - `func Mat4Scale(s Vec3) Mat4` Scales a shape by vector s.
    - `func Mat4RotationAxisAngle(axis Vec3, angle float32) Mat4` Rotates a shape by rad radians.
    - `func (m Mat4) Mul(b Mat4) Mat4` Multiplies matrix m by matrix b  (note: multiplication order may affect the result).
    - `func (m Mat4) SmartMul(b Mat4) Mat4` Optimized matrix multiplication using identity, translation, and scale shortcuts
    - `func (m Mat4) TransformPo3(p Po3) Po3` Applies a linear transformation to a point.
    - `func (m Mat4) TransformVec3(v Vec3) Vec3`  Applies a linear transformation to a vector.
    - `func Mat4TRS(pos Vec3, axis Vec3, angle float32, scale Vec3) Mat4` Translates by pos, rotates by rot, and scales by scale.
    - `func Mat4Perspective(fovY, aspect, near, far float32) Mat4` Creates perspective projection matrix for 3D depth rendering
    - `func Mat4LookAt(eye Vec3, center Vec3, up Vec3) Mat4` Creates a view matrix from camera position and orientation
    - `func Mat4Ortho(left, right, bottom, top, near, far float32) Mat4` Creates an orthographic projection matrix for 3D rendering
    - `func (m Mat4) InverseAffine() Mat4` Returns the inverse affine transformation.
    - `func (m Mat4) NormalMatrix() Mat3` Computes inverse and transpose for normal vectors
- ### Po2
  - #### content
     - `Po2` type `X,Y float32`
  - #### functions
     - `func (p Po2) Add(v Vec2) Po2` Adds a vector to a 2D point
     - `func (p Po2) Sub(q Po2) Vec2` Produces a vector by subtracting a 2D point from a 2D point
     - `func (p Po2) DistanceSquared(q Po2) float32` Returns squared distance between two 2D points
     - `func (p Po2) Distance(q Po2) float32` Returns distance between two 2D points
     - `func (p Po2) Equals(q Po2, eps float32) bool` Checks approximate 2D point equality within epsilon tolerance
     - `func (p Po2) String() string` Returns formatted string representation of a 2D point
     - `func Orient(a, b, c Po2) float32` Computes 2D orientation / signed triangle area
- ### Po3
  - #### content
    - `Po3` type `X,Y,Z float32`
  - #### functions
     - `func (p Po3) Add(v Vec3) Po3` Adds a vector to a 3D point
     - `func (p Po3) SubPo(q Po3) Vec3` Produces a vector by subtracting a 3D point from a 3D point
     - `func (p Po3) SubVec(q Vec3) Vec3` Subtracts vector from point, returns displacement vector
     - `func (p Po3) DistanceSquared(q Po3) float32` Returns squared distance between two 3D points
     - `func (p Po3) Distance(q Po3) float32` Returns distance between two 3D points
     - `func (p Po3) Equals(q Po3, eps float32) bool` Checks approximate 3D point equality within epsilon tolerance
     - `func (p Po3) String() string` Returns formatted string representation of a 3D point
- ### Transform2D
  - ##### content
     - `Position` type `Vec2`
     - `Rotation` type `float32`
     - `Scale` type `Vec2`
     - `Dirty` type `bool`
  - #### functions
     - `func NewTransform2D() Transform2D` Initializes Transform2D with unit scale and identity matrix
     - `func (t *Transform2D) SetPosition(p Vec2)` Sets position and marks transform as dirty
     - `func (t *Transform2D) SetRotation(r float32)` Sets rotation and marks transform as dirty
     - `func (t *Transform2D) SetScale(s Vec2)` Sets scale and marks transform as dirty
     - `func (t *Transform2D) Matrix() Mat3` Returns cached transform matrix, recomputing only when marked dirty
     - `func (t *Transform2D) TransformPoint(p Po2) Po2` Marks transform as dirty, forcing matrix recomputation
     - `func (t *Transform2D) TransformVector(v Vec2) Vec2` Transforms vector using current 2D transformation matrix
     - `func (t *Transform2D) TranslateBy(delta Vec2)` Adds delta to position, marking transform dirty
     - `func (t *Transform2D) RotateBy(delta float32)` Adds delta to rotation, marks transform dirty
     - `func (t *Transform2D) ScaleBy(factor Vec2)` Multiplies scale by factor, marking transform dirty
     - `func (t *Transform2D) Snapshot()` Stores current transform values for later interpolation
     - `func (t *Transform2D) InterpolatedMatrix(alpha float32) Mat3` Builds interpolated transform matrix between previous and current states
- ### Transform3D
  - #### content
     - `Position` type `Vec3`
     - `RotationAxis` type `Vec3`
     - `Rotation` type `float32`
     - `Scale ` type `Vec3`
     - `Dirty` type `bool`
     - `AxisMask` type `uint8` (Bitmask enum representing selectable X, Y, Z axis combinations)  
  - #### functions
     - `func NewTransform3D() Transform3D` Initializes Transform3D with unit scale and identity matrix
     - `func (t *Transform3D) SetPosition(p Vec3)` Sets position and marks transform as dirty
     - `func (t *Transform3D) SetRotationAxis(r Vec3)` Sets normalized rotation axis, marks transform dirty
     - `func (t *Transform3D) SetScale(s Vec3)` Sets scale and marks transform as dirty
     - `func (t *Transform3D) Matrix() Mat4` Returns cached transform matrix, recomputing only when marked dirty
     - `func (t *Transform3D) TransformPo3(p Po3) Po3` Marks transform as dirty, forcing matrix recomputation
     - `func (t *Transform3D) TransformVec3(v Vec3) Vec3` Transforms vector using current 3D transformation matrix
     - `func (t *Transform3D) TranslateBy(delta Vec3)`  Adds delta to position, marking transform dirty
     - `func (t *Transform3D) RotateBy(delta float32)` Adds delta to rotation, marks transform dirty
     - `func (t *Transform3D) ScaleBy(factor Vec3)` Multiplies scale by factor, marking transform dirty
     - `func (t *Transform3D) Snapshot()`Stores current transform values for later interpolation
     - `func (t *Transform3D) InterpolatedMatrix(alpha float32) Mat4`  Builds interpolated transform matrix between previous and current states







    
    	             
