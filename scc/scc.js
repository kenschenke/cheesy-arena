const WebSocket = require('ws');
const fs = require('fs');
const Gpio = require('onoff').Gpio;
var Socket = require('net').Socket;

const url = 'ws://10.0.100.5:8080/scc/websocket';
const alliancefn = '/boot/scc';

const RETRY_TIMEOUT = 3000;   // 30 secs
const DEBOUNCE_TIMEOUT = 10;  // 10 ms

const pinLED = 18;

// Scoring SCC
const pinEstop = 22;
const pinOff = 23;
const pinReset = 24;
const pinGreen = 25;
var eStop, btnOff, btnReset, btnGreen;

// Red or Blue SCC
const pinEstop1 = 23;
const pinEstop2 = 24;
const pinEstop3 = 25;
var eStop1, eStop2, eStop3;

var ws = null;
var alliance = '';

try {
    alliance = fs.readFileSync(alliancefn).toString().trim();
} catch (err) {
    console.error(err);
}

const led = new Gpio(pinLED, 'out');

if (alliance === "scoring") {
    eStop = new Gpio(pinEstop, 'in', 'both');
    btnOff = new Gpio(pinOff, 'in', 'rising', {debounceTimeout: DEBOUNCE_TIMEOUT});
    btnReset = new Gpio(pinReset, 'in', 'rising', {debounceTimeout: DEBOUNCE_TIMEOUT});
    btnGreen = new Gpio(pinGreen, 'in', 'rising', {debounceTimeout: DEBOUNCE_TIMEOUT});

    eStop.watch((err, value) => {
        if (err == null)
            SendButtonStatus(value!==0, false, false);
    });

    btnOff.watch((err, value) => {
        if (err === null)
            SendFieldLights('off');
    });
    btnReset.watch((err, value) => {
        if (err === null)
            SendFieldLights('purple');
    });
    btnGreen.watch((err, value) => {
        if (err === null)
            SendFieldLights('green');
    });
} else {
    eStop1 = new Gpio(pinEstop1, 'in', 'both');
    eStop2 = new Gpio(pinEstop2, 'in', 'both');
    eStop3 = new Gpio(pinEstop3, 'in', 'both');
    
    eStop1.watch((err, value) => {
        if (err === null)
            SendButtonStatus(value!==0, eStop2.readSync()!==0, eStop3.readSync()!==0);
    });
    eStop2.watch((err, value) => {
        if (err === null)
            SendButtonStatus(eStop1.readSync()!==0, value!==0, eStop3.readSync()!==0);
    });
    eStop3.watch((err, value) => {
        if (err === null)
            SendButtonStatus(eStop1.readSync()!==0, eStop2.readSync()!==0, value!==0);
    });
}

function SendButtonStatus(eStop1, eStop2, eStop3) {
    if (ws !== null && ws.isActive) {
        ws.send(JSON.stringify({
            type: 'sccupdate',
            data: { alliance, eStop1, eStop2, eStop3 }
        }))
    }
}

function SendFieldLights(color) {
    if (ws !== null && ws.isActive) {
        ws.send(JSON.stringify({
            type: 'setfieldlights',
            data: color
        }));
    }
}

function SetStatusLED(status) {
    led.writeSync(status);
    console.log('Status: %s', status!==0 ? 'ON' : 'OFF');
}

function sendColor(r, g, b) {
    var socket = new Socket();
    socket.setNoDelay();
    socket.connect(7890);

    var createOPCStream = require("opc");
    var stream = createOPCStream();
    stream.pipe(socket);

    var createStrand = require("opc/strand");
    var strand = createStrand(64);
    for (var i = 0; i < 64; i++) {
        strand.setPixel(i, r, g, b);
    }

    stream.writePixels(0, strand.buffer);
}

process.on('beforeExit', () => SetStatusLED(0));
process.on('SIGINT', () => {
    SetStatusLED(0);
    led.unexport();
    if (alliance === "scoring") {
        eStop.unexport();
        btnOff.unexport();
        btnReset.unexport();
        btnGreen.unexport();
    } else {
        eStop1.unexport();
        eStop2.unexport();
        eStop3.unexport();
    }
    process.exit();
});

function noop() {}

(function loop() {
    const interval = setInterval(() => {
        if (ws.isAlive == false) {
           ws.terminate();
        }
        ws.isAlive = false;
        ws.ping(noop);
    }, 10000);
    ws = new WebSocket(url);
    ws.isActive = false;
    ws.isAlive = false;
    ws.on('error', () => {
        ws.isActive = false;
        SetStatusLED(0);
        console.log('Failed to connect.  Waiting 3 seconds');
        clearInterval(interval);
        setTimeout(loop, RETRY_TIMEOUT);
    });
    ws.on('open', function open() {
        ws.isActive = true;
        ws.isAlive = true;
        SendButtonStatus(false, false, false);
        SetStatusLED(1);
    });
    ws.on('close', () => {
        if (ws.isActive) {
            console.log('Lost connection.  Attempting reconnect in 3 seconds.');
            SetStatusLED(0);
            clearInterval(interval);
            setTimeout(loop, RETRY_TIMEOUT);
        }
        ws.isActive = false;
    });
    ws.on('pong', () => {
        ws.isAlive = true;
        // console.log('PONG!');
    });
    
    if (alliance !== "scoring") {
        ws.on('message', function incoming(data) {
            try {
                message = JSON.parse(data);
                if (message.hasOwnProperty('type') && message.type == 'fieldLights') {
                    if (message.hasOwnProperty('data') && message.data.hasOwnProperty('Lights')) {
                        switch (message.data.Lights) {
                        case 'red':
                            sendColor(255, 0, 0);
                            break;
                        case 'green':
                            sendColor(0, 255, 0);
                            break;
                        case 'purple':
                            sendColor(100, 0, 100);
                            break;
                        case 'off':
                            sendColor(0, 0, 0);
                            break;
                        }
                    }
                    console.log('Set lights to %s', message.data.Lights);
                }
            } catch (e) {
                console.log('Bad message - ignoring');
            }
        });
    }
})();

