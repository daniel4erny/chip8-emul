let audioCtx = null;
let osc = null;

function startBeep() {
    audioCtx ??= new AudioContext();
    osc = audioCtx.createOscillator();
    osc.type = "square";
    osc.frequency.value = 440;
    osc.connect(audioCtx.destination);
    osc.start();
}

function stopBeep() {
    if (osc) {
        osc.stop();
        osc = null;
    }
}

/** @type {HTMLCanvasElement} */
const canvas = document.getElementById("ch8Display")
const canvasCtx = canvas.getContext("2d");
const imageData = canvasCtx.createImageData(64, 32);
const data = imageData.data;

function drawScreen(pixels) {
    for (let i = 0; i < 2048; i++) {
        const v = pixels[i];          
        data[i*4] = v;
        data[i*4 + 1] = v;
        data[i*4 + 2] = v;
        data[i*4 + 3] = 255;
    }
    canvasCtx.putImageData(imageData, 0, 0)
}


const keyMap = {
    Digit1: 0x1, Digit2: 0x2, Digit3: 0x3, Digit4: 0xC,
    KeyQ:   0x4, KeyW:   0x5, KeyE:   0x6, KeyR:   0xD,
    KeyA:   0x7, KeyS:   0x8, KeyD:   0x9, KeyF:   0xE,
    KeyZ:   0xA, KeyX:   0x0, KeyC:   0xB, KeyV:   0xF,
};

const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance); 
});

document.getElementById("rom").addEventListener("change", async (e) => {
    const buf = await e.target.files[0].arrayBuffer();
    loadRom(new Uint8Array(buf));
});

document.addEventListener("keydown", (e) => {
    if (e.code in keyMap) {
        setKey(keyMap[e.code], true);
        e.preventDefault();
    }
});

document.addEventListener("keyup", (e) => {
    if (e.code in keyMap) {
        setKey(keyMap[e.code], false);
        e.preventDefault();
    }
});
