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

const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    go.run(result.instance); 
});

document.getElementById("rom").addEventListener("change", async (e) => {
    const buf = await e.target.files[0].arrayBuffer();
    loadRom(new Uint8Array(buf));
});
