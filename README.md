# Playground for golang-eth methods

## Test minting NFTs `./nft`

### Requirements 
[Infura](https://infura.io/) - infura is used as our JSON rpc to interact with the Ethereum BlockChain.

- Grab the ABI (spec for a contract) via `fetch_abi.py` e.g

`python fetch_abi.py <contract_address> -o abi.json`

### Example Config
```json
{
  "infura_api_url": "infura_api_url",
  "value": 30000000000000000,
  "wallet_private_key": "wallet_private",
  "contract_address": "contract_address",
  "gas_limit": 5000000, (in wei)
  "abi": "abi.json"
}
```
