# Notabor Engine
an phsyics engine made in **goland** designed to make game dev easier by providing the following:
- ## Function based architecture
  there is a `Runnable` type function, which is a function that can return an error.
  
  a runnable can be appended to a Loop whether it's `FixedHzLoop` or `RenderLoop` by appending the function to the loop.
  
  the runnable then runs at a specific ticks per seconds `Hz` when the loop runs, no need for multiplying anything by delta time and no need for intense Object Oriented code.
- ## Loop based running
  as mentioned, a loop is `FixedHzLoop` or `RenderLoop`, it contains a list of `Runnable`, when it runs it and the Runnable returns an error, it deletes the Runnable from the list and continues, otherwise it runs without deleting
  this makes it so there is no need for any for loops made by the player, and keeps the code clean and easy to read
- ## Native Multi-Window support
  a window can be created via a config, each window has one `RenderLoop` and zero or more `FixedHzLoop`, that way we can ensure separation between windows
  different window modes, `Windowed`, `BorderlessWindowed` and `FullScreen`
