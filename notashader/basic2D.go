package notashader

const Vertex2D = `#version 460 core
layout(location = 0) in vec2 aPos;

void main() {
    gl_Position = vec4(aPos, 0.0, 1.0);
}`

const Vertex2DTransform = `#version 460 core
layout(location = 0) in vec2 aPos;
uniform mat3 uTransform;

void main() {
    vec3 pos = uTransform * vec3(aPos, 1.0);
    gl_Position = vec4(pos.xy, 0.0, 1.0);
}`

const Fragment2D = `#version 460 core
out vec4 FragColor;

void main() {
    FragColor = vec4(1.0, 0.5, 0.2, 1.0);
}`
