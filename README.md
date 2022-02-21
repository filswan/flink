# FLink
FILink is a data provider transfers deals info to chainlink Oracle

## Why we need to have Fil-Chainlink proof
While a user is working on other blockchain like Ethereum, BSC,Polygon, they want to know the strage deal they send to filecoin network is online, however we do not have a native way to check if a storage deal is active on those blockchains.

## Solutions 
By build a filecoin data adapter to chainlink oracle. It can give other blockchain the capability of validation if a data is onchain

### Chainlink Adapter - DATA DAO 

Data Dao signs the node for providing data feeds to the web3 blockchains like polygoin, bsc, eth.
Every data deal has a deal ID, it is an unique id on filecoin network used for tracking deal info.
Filecoin data adatper provides deal ids for other blockchain system to check if the deal exists on filecoin network.

The adapter scans the data from Flecoin blockchain and post the deal info to polygon network, the deal info is by the following format:

```json
{
  "status": "success",
  "data": {
    "deal": {
      "deal_id": 4107,
      "deal_cid": "bafyreialcn7arcqxx5panzuodnbu6kn4wxipenopjiikoqlrc4w7ed2hfq",
      "message_cid": "bafy2bzaceaxzikfc6a5ax4vzh52oieercs6nhuowd424chhp7mfdcalvaihpu",
      "height": 102337,
      "piece_cid": "baga6ea4seaqooaumpt2kab7ltyyvftfgzyce5zzox3ou6n6qf34bkv42rj7rkly",
      "verified_deal": false,
      "storage_price":26207450000000,
      "storage_price_per_epoch": 50000000,
      "signature": "BIE8UTJHBRq6qkBPOtz/IhydTl//+oZP2gMMiGIPck4ZkUPQ/41QLwmTNuNqnUB62j85njIlnRxWLMXv8HJgtQE=",
      "signature_type": "secp256k1",
      "created_at": 1616875110,
      "piece_size_format": "512.00 bytes",
      "start_height": 108199,
      "end_height": 632348,
      "client": "t1e3mn5c6v7d3wjxiuycozkfbnqvzwysjmimltm5a",
      "client_collateral_format": "000000000000000000",
      "provider": "t05661",
      "provider_tag": "",
      "verified_provider": 0,
      "provider_collateral_format": "000000000000000000",
      "status": 0,
      "network_name": "filecoin_calibration"
    }
  }
}
```

## Sample Use Case
### Polygon NFTs USDC payment for Filecoin storage
![FilLInk](https://user-images.githubusercontent.com/8363795/143550092-bc10f493-b6c5-48e0-ac46-5bbd49a11731.png)


### Deal Matching
scheduler to update status to trigger DAO signature for unlock event
* get deal_id by proposal_cid
* get deal_id from Chainlink filecoin adapter
    * if matches trigger DAO signature
        * match client_address
        * match deal_cid (proposal_cid)
    * else waiting for next check cycle

## Sponsors
This project is sponspor by Filecoin Foundation & Chainlink:

[RFP : chainlink and filecoin data bounties](https://github.com/filecoin-project/devgrants/pull/290
)

<img src="filecoin.png" width="200"> &nbsp&nbsp&nbsp  <img src="chainlink.png" width="200">

