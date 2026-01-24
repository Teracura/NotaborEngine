package notashader

const Vertex3D = `
#version 460 core
layout(location = 0) in vec3 aPos;
layout(location = 1) in vec4 aColor;

out vec4 vColor;

void main() {
    gl_Position = vec4(aPos, 1.0);
    vColor = aColor;
}`

const Fragment3D = `
#version 460 core

in vec4 vColor;
out vec4 FragColor;

void main()
{
    FragColor = vColor;
}
`
