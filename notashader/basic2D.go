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
        
        // We need to know where we are relative to the center of the quad.
        // Since Fixate() centers your vertices around (0,0), 
        // aPos is already "local enough" for a simple distance check.
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
        // Calculate distance from the local origin (0,0)
        float dist = length(vLocalPos);

        // Since we don't have scale here, we'll assume a radius of 0.5 
        // (matching your CreateRectangle setup)
        float radius = 0.5;
        float edge = 0.01; 
        float alpha = 1.0 - smoothstep(radius - edge, radius, dist);

        if (alpha <= 0.0)
            discard;

        FragColor = vec4(vColor.rgb, vColor.a * alpha);
    }
    `
