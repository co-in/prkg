## Pseudorandom Key Generation

### Motivation:
I noticed that the library (https://github.com/tyler-smith/go-bip39) which I used for BIP39 had previously disappeared.
I did not want to use forks, so everything started with mnemonic.
Ready-made solutions are very focused on a specific type of curve (like sec256k1 or ed25519), or on compliance with
certain standards that are not relevant to me. Therefore, the next step was to abandon HD in favor of Pseudorandom Key
Generation.
Path derivation no longer conforms to BIP-33, BIP-44, SLIP-10, so a custom format was used for **Path**, visually
distinct to avoid collisions and deception, but key generation is not tied to it

### Pros:
- BIP39 compatible, but with the ability to fully customize your dictionary (though I don't see the point in that right
  now)
- Hardened derivation as a single standard for **any** curves and cases, instead of the vulnerable scenario of **leaking
  a parent extended public key**

### Cons:
- Does not meet existing standards
- Child keys cannot be created autonomously from extended

### Plans:

- Add split mnemonics to cards
- Add comments, tests, usage examples