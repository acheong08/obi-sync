const crypto = require("crypto");

function B(e) {
    return e.slice(e.byteOffset, e.byteOffset + e.byteLength);
}



async function _(e) {
    const t = await j(e);
    return G(t);
}

async function j(e) {
    const digest = await crypto.subtle.digest("SHA-256", new Uint8Array(e));
    return new Uint8Array(digest);
}


function G(e) {
    const t = new Uint8Array(e);
    const n = [];
    for (let i = 0; i < t.length; i++) {
        const r = t[i];
        n.push((r >>> 4).toString(16));
        n.push((15 & r).toString(16));
    }
    return n.join("");
}

const K = "aes-256-gcm";

async function Y(e) {
    return crypto.subtle.importKey("raw", e, K, false, ["encrypt", "decrypt"]);
}

const $ = 32768;
const Q = {
    N: $,
    r: 8,
    p: 1,
    maxmem: 67108864
};

async function J(e, t) {
    const normalizedE = e.normalize("NFKC");
    const normalizedT = t.normalize("NFKC");
    const n = await new Promise((resolve, reject) => {
        crypto.scrypt(Buffer.from(normalizedE, 'utf-8'), Buffer.from(normalizedT, 'utf-8'), 32, Q, (err, key) => {
            if (err) {
                reject(err);
            } else {
                resolve(key);
            }
        });
    });
    return B(n);
}

async function makeKeyHash(e, t) {
    const n = await J(e, t);
    const hash = await j(n);
    return G(hash);
}

makeKeyHash("ZsSjgKx4yaeBNCFipS)T", "jePEuEPhNsr8zguY3%98").then((hash) => {
    console.log(hash);
});