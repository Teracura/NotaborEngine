package notashader

const TexturedVertex2D = `#version 460 core
layout(location = 0) in vec2 aPos;
layout(location = 1) in vec4 aColor;
layout(location = 2) in vec2 aUV;

out vec4 vColor;
out vec2 vUV;

void main() {
    gl_Position = vec4(aPos, 0.0, 1.0);
    vColor = aColor;
    vUV = aUV;
}`

const TexturedFragment2D = `#version 460 core
in vec4 vColor;
in vec2 vUV;
out vec4 FragColor;

uniform sampler2D uTexture;
uniform bool uUseTexture = true;

void main() {
    vec4 texColor = vec4(1.0, 1.0, 1.0, 1.0);
    if (uUseTexture) {
        texColor = texture(uTexture, vUV);
    }
    FragColor = vColor * texColor;
}`
