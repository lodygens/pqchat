# PQChat

PQChat is a quantum threat resistant peer-to-peer chat application.

This is a case study permetting not only to play with [libp2p](https://github.com/libp2p) but also to test [Open Quantum Safe library](https://openquantumsafe.org/).

This is open source.

This should not be used in production.


---

# License

Apache 2.0

---

# Requirements:



* Fight [Harvest Now Decrypt Later](https://en.wikipedia.org/wiki/Harvest_now,_decrypt_later) threat with a quantum-resistant transport with [ML-KEM](https://csrc.nist.gov/pubs/fips/203/final)
* Ensure quantum-resistant identity with [ML-DSA](https://csrc.nist.gov/pubs/fips/204/final)
* User may have an optional pseudo
* Messages may be broadcasted to all, or unicasted to one or more users
* Connect using P2P mechanisms


Each user has an identity keypair:

* Algorithm: e.g. `ML-DSA-65` (NIST level 1)
* Public key: `id_pub`
* Private key: `id_priv`
* Pseudo: user-chosen string

User can:

* generate a keypair before running the application; the application will load them
* let the application generate the key pair

We define a user ID as:

```text
user_id = SHA256( id_pub || pseudo )
```

Every message is signed:

```text
signature = ML-DSA.Sign(id_priv, message_bytes)
```

and verified by others with `id_pub`.

---

## Session keys (ML-KEM + AES-GCM)

Each pair of users (A,B) has a shared symmetric key:

1. A and B establish a libp2p stream (Noise-encrypted transport)

2. They run ML-KEM-768 handshake *inside* the stream:

   * A (server role) → generate ML-KEM keypair `(pkA, skA)`
   * A → send `pkA` to B
   * B → `Encap(pkA)` → `(ct, ss)`
   * B → send `ct` to A
   * A → `Decap(skA, ct)` → `ss`

3. Both have the same shared secret `ss`

4. Derive AES key:

```text
aes_key = HKDF( ss, "PQCHAT-SESSION", 32 bytes )
```

5. Use AES-GCM (128 or 256 bits) for message confidentiality:

   * `ciphertext = AES_GCM_Encrypt(aes_key, nonce, plaintext, aad)`
   * `plaintext = AES_GCM_Decrypt(...)`

> This gives a PQC-resistant channel between each pair of users.

For broadcast, simplest PoC:

* each message is individually encrypted for each recipient with their per-peer AES key
* then sent directly over that connection
* that’s O(N) per broadcast, but fine for a demo.

---

## Message formats


### Control messages: presence / identity

Example: `HELLO`

```json
{
  "type": "HELLO",
  "pseudo": "oleg",
  "user_id": "hex(SHA256(id_pub||pseudo))",
  "ml_dsa_pub": "base64(...)", 
  "libp2p_peer_id": "12D3KooW...",
  "sig": "base64( ML-DSA.Sign(id_priv, canonical_json_without_sig) )"
}
```

Peers store `(user_id → {pseudo, ml_dsa_pub, peer_id})` once verified.

---

### Chat messages

Before encryption/signature, “logical” message:

```json
{
  "type": "CHAT",
  "from": "user_id",
  "to": ["user_id_1", "user_id_2"],  // empty or ["*"] = broadcast
  "body": "Hello world",
  "timestamp": 1732970000
}
```

Then:

1. Serialize to canonical JSON → `m`
2. Compute signature:

```text
sig = ML-DSA.Sign(id_priv, m)
```

3. Build signed envelope:

```json
{
  "msg": { ...as above... },
  "sig": "base64(sig)",
  "pub": "base64(id_pub)" // or omit if cached from HELLO
}
```

4. Encrypt envelope with AES-GCM session key → ciphertext bytes
5. Prepend framing (len, etc.) and send over libp2p stream.

Receiver:

* decrypt AES-GCM
* parse JSON
* verify `ML-DSA.Verify(pub, msg, sig)`
* display if ok.

---

## Runtime & CLI UX

One binary, two roles:

* Node: always a P2P peer (libp2p host)
* Pseudos / keypairs:

```bash
pqchat \
  -pseudo "oleg" \
  -ml-dsa-priv ./keys/oleg-ml-dsa-priv.bin \
  -ml-dsa-pub  ./keys/oleg-ml-dsa-pub.bin \
  -listen "/ip4/0.0.0.0/tcp/0" \
  -peer "/ip4/…/tcp/4001/p2p/12D3KooW…"  # bootstrap/relay or other peers
```

If keys don’t exist, the program can:

* generate ML-DSA keypair
* save to files
* print a warning (“new identity created”).

Chat UX in terminal:

* `hello everyone` → default: broadcast
* `@alice hi` → send only to user `alice`
* `@alice @bob secret` → send to both

Map pseudos → user_ids → peers.

---

# System Requirements

You need a machine running:

* macOS (Intel or ARM)
* Ubuntu/Linux
* Raspberry Pi OS (ARM64) → supported after installing liboqs manually
* Go 1.21+

---

# Dependencies

PQChat depends on:

| Component                   | Purpose                                       |
| --------------------------- | --------------------------------------------- |
| OpenSSL 3.x             | PQC provider support (liboqs or oqs-provider) |
| liboqs (C library)      | PQC primitives (ML-KEM, ML-DSA, BIKE, Frodo…) |
| liboqs-go (Go wrappers) | Go bindings calling liboqs                    |
| pkg-config              | Required for liboqs-go compilation            |
| libp2p                  | P2P networking                                |
| cgo                     | To call liboqs from Go                        |
| Make                    | For building binaries                         |

If you run:

```
go run examples/kem/kem.go
```

and obtain shared secrets → everything is OK.

---

# Installation

## Install liboqs and oqs provider

One can use 
- [install_oqs_openssl_mac.sh](https://github.com/lodygens/cafe/tree/main/src/poc/pq-scan/install_oqs_openssl_mac.sh)
- [install_oqs_openssl_debian.sh](https://github.com/lodygens/cafe/tree/main/src/poc/pq-scan/install_oqs_openssl_debian.sh)


## Install liboqs-go (Go bindings + C layer)

Please follow [liboqs-go README](https://github.com/open-quantum-safe/liboqs-go.git)

Execute the following to check the installation:

```bash
go run liboqs-go/examples/kem/kem.go
```

Expected output:

```
liboqs version: 0.15.0-rc1
Enabled KEMs: [ML-KEM-512 ML-KEM-768 …]
Shared secrets coincide? true
```

If this works → PQChat will compile.

---

# Building PQChat

From the project root:

```
make all
```

This produces:

```
bin/relayer
bin/pqchat
```

---

# Running the relay

On machine A (can be behind NAT):

```bash
./bin/relayer
```

You will see:

```
Relay PeerID: 12D3K...
  /ip4/192.168.1.X/tcp/4001
  /ip4/xxx.xxx.xxx.xxx/tcp/4001
```

Copy this address.

---

# Running PQChat

On machine B:

```bash
./bin/pqchat --relay /ip4/.../p2p/12D3K...
```

You will see:

```
✔ Connected to relay
✔ PQC handshake (ML-KEM-768) complete
Your PQ identity (ML-DSA-65): <hex>
```

Then you can:

* send messages
* broadcast
* direct-message specific peers
* inspect PQ handshake logs

---

# Security Model

* PQ identity = ML-DSA public key
* PQ handshake = ML-KEM-768 (Kyber)
* Session key = AES-GCM using ML-KEM shared secret
* No classical crypto fallback
* No Noise → fully PQC end-to-end
* Each peer only processes tasks for its own identity (future work)


---

# Run

You can now:

* run a decentralized relay
* chat securely across NAT
* use ML-KEM for P2P key exchange
* sign PQ identities with ML-DSA
* experiment with PQC + P2P networking


# Local test without any relay

libp2p n’a *pas besoin* d’un relay si les deux peers sont *sur la même machine* et *sont en TCP/IP direct*.

1. Launch first peer

In a first terminal :

```bash
./bin/pqchat 
⚠️ No relay configured, running in direct TCP mode.
Local PeerID: 12D3KooWMBjeJ5aewPjZTseS6rugXB2PJRDCfGKJR4xUmFmiM6m9
Listening on:
    /ip4/127.0.0.1/tcp/56142
    /ip4/192.168.64.1/tcp/56142
```

Recreate its multiaddr with the port and the PeerID.

```
/ip4/192.168.64.1/tcp/56142/p2p/12D3KooWMBjeJ5aewPjZTseS6rugXB2PJRDCfGKJR4xUmFmiM6m9
```

---

1. Launch a second peer and connect to the first one

In another terminal (use the multiaddr of the first node):

```bash
./bin/pqchat -connect /ip4/192.168.64.1/tcp/56142/p2p/12D3KooWMBjeJ5aewPjZTseS6rugXB2PJRDCfGKJR4xUmFmiM6m9
⚠️ No relay configured, running in direct TCP mode.
Local PeerID: 12D3KooWLUSKEb5WTijoqTVYUJXjDsP8SSBk8qPGw7vpqLPaSkPA
Listening on:
    /ip4/127.0.0.1/tcp/56268
    /ip4/192.168.64.1/tcp/56268
Connecting to peer: 12D3KooWMBjeJ5aewPjZTseS6rugXB2PJRDCfGKJR4xUmFmiM6m9
Running PQC client handshake…
PQC session established (client side)
```

# Local test with a local relay

1. Launch the relay

In a first terminal  :

```bash
./bin/relayer

relay listening on:
  /ip4/127.0.0.1/tcp/4001/p2p/12D3KooWrelay...
```

**→ Copy the multiaddr of the relay**.


2. Launch a first peer 

In a second terminal:

```bash
./bin/pqchat -pseudo alice -relay /ip4/127.0.0.1/tcp/4001/p2p/12D3KooWrelay
```

Alice publie ses adresses.

---

3. Launch a second peer

In a third terminal:

```bash
./bin/pqchat -pseudo bob -relay /ip4/127.0.0.1/tcp/4001/p2p/12D3KooWrelay \
  -connect /ip4/127.0.0.1/tcp/XXXXX/p2p/<PeerIDAlice>
```

