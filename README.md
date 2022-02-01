# Playground for golang-eth methods

## Test minting NFTs `./nft`

### Requirements 
[Infura](https://infura.io/) - Infura used as JSON RPC to interact with the ethereum blockchain and network (essentially a service provider)

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
