package notashader

const Vertex2D = `#version 460 core
layout(location = 0) in vec2 aPos;
layout(location = 1) in vec4 aColor;

out vec4 vColor;

void main() {
    gl_Position = vec4(aPos, 0.0, 1.0);
    vColor = aColor;
}`

const Fragment2D = `#version 460 core
in vec4 vColor;
out vec4 FragColor;

void main() {
    FragColor = vColor;
}`

const Circle2DVertex = `
    #version 460 core

    layout(location = 0) in vec2 aPos;
    layout(location = 1) in vec4 aColor;

    out vec2 vLocalPos;
    out vec4 vColor;

    void main() {
        gl_Position = vec4(aPos, 0.0, 1.0);

        vLocalPos = aPos; 
        vColor = aColor;
    }
    `

const Circle2DFragment = `
    #version 460 core

    in vec2 vLocalPos;
    in vec4 vColor;

    out vec4 FragColor;

    void main() {
        float dist = length(vLocalPos);

        float radius = 0.5;
        float edge = 0.01; 
        float alpha = 1.0 - smoothstep(radius - edge, radius, dist);

        if (alpha <= 0.0)
            discard;

        FragColor = vec4(vColor.rgb, vColor.a * alpha);
    }
    `
