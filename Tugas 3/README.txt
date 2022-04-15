PAY TO PUBLIC KEY HASH

Pay-to-PubKey-Hash (P2PKH) contract is used to send bitcoins to a bitcoin address.
It is the most common contract on the Bitcoin network.
Such contracts are unlocked by the public key and a signature created by the corresponding private key.

ScriptPubKey= OP_DUP OP_HASH160 <Public KeyHash> OP_EQUAL OP_CHECKSIG
ScriptSig= <Signature> <Public Key>
OP_DUP OP_HASH160 <pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG is Constants are added to the stack.
OP_HASH160 <pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG Top stack item is duplicated.
<pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG top stack item is hashed.
OP_EQUALVERIFY OP_CHECKSIG constant added.
