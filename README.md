<div align="center">
<p style="font-size:24px; font-weight: bold;">Oasis Signature Tools</p>
<p>
      <img alt="Tessellated Logo" src="media/tessellated-logo.png" />
<small>By <a href="https://tessellated.io" target="_blank"> Tessellated // tessellated.io</a></small>
</p>
</div>

---


This repository contains tooling and documentation for manually generating and signing airgapped transactions on [Oasis Protocol](https://oasisprotocol.org/).

## Usage 

This program formulates bytes to sign from the outputs of Oasis network's tools.

Usage:

```shell
go run main.go <base64 raw value>
```

## Transaction Signing Mini Manual

You need a few things to sign a transaction. 

### Creating the bytes to sign

First, get the Oasis CLI tools:
```shell
git clone https://github.com/oasisprotocol/cli
cd cli && make && mv ./oasis ~/go/bin/oasis-cli
```

Grab your nonce from [OasisScan](https://oasisscan.com)
```shell
export NONCE=31
```

Use Oasis tools to generate a dummy account. This account prevents the CLI from complaining you don't have a private key. 
```shell
oasis-cli wallet create insecure_dummy
```

Use the Oasis tool to build a transaction from the dummy account. We do this so that we can get the raw bytes of the transaction.
The CLI will estimate gas and fees automatically, but we have to override the nonce since the CLI is signign from a different account:
> Note: You can mess around with the CLI's node endpoints by running `./oasis networks list` and related commands

Here's an example for casting a governace vote, but this will work with any transaction at the consensus layer:
```shell
oasis-cli network governance cast-vote --account insecure_dummy --format json --nonce $NONCE --output-file tmp-dummy-tx.json 4 yes
```

Extract the raw bytes and cleanup the dummy transaction file
```shell
export RAW_BYTES=$(cat tmp-dummy-tx.json | jq -r .untrusted_raw_value)
rm tmp-dummy-tx.json
```

The program in this repository will generate bytes to sign for the transaction.
> Note: This program only works with consensus layer transactions. If you need to make a different transaction, consider introspecting [this code](https://github.com/oasisprotocol/oasis-core/blob/master/go/common/crypto/signature/signer.go#L362) in a debugger.

From this repository's root:
```shell
go run main.go $RAW_BYTES
```

### Sign Tx Bytes

Sign the bytes output from above in whatever way you want to. [This tool](https://cyphr.me/ed25519_tool/ed.html) can help you check that your signatures is valid. 

### Assemble the Final Transaction

The final transaction JSON file is in the following format. 

```json
{
  "untrusted_raw_value": "",
  "signature": {
    "public_key": "",
    "signature": ""
  }
}
```

Set: 
- `untrusted_raw_value`: To be the raw value outputted by the Oasis CLI (the value we saved in `$RAW_BYTES`)
- `public_key`: To be your public key
- `signature`: To be the base64 encoded signature you created

Here's an example:
```json
{
  "untrusted_raw_value": "pGNmZWWiY2dhcxkE2WZhbW91bnRAZGJvZHmiYmlkBGR2b3RlAWVub25jZRgeZm1ldGhvZHNnb3Zlcm5hbmNlLkNhc3RWb3Rl",
  "signature": {
    "public_key": "RS+saioy8ukcohygmbUH0IKGXhLjaE6BWrkdOBBqpDc=",
    "signature": "s7Y5qOojP4U+ku+ZKyRpB0m071vFsKbaDs1XeVeLLh5PDzWCyculhaK7HJIop2KQr7kgD4w/Ef0ll30ZSvFyDA=="
  }
}
```

Save this to a file, like `tx.json`.

### Broadcast

Lastly, broadcast with the Oasis CLI and cleanup.
```shell
oasis-cli tx submit tx.json
rm tx.json
```
