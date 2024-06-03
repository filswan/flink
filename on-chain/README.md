# FLink contract
Flink Contract is used by consumers get filecoin deal storage price from off-chain API.

By using chainlink any API, this Contract accesses a http get endpoint to get filecoin deal storage price

## How to get a deal info
1. Call requestDealInfo method with two parameters: network and deal id
2. Chainlink oracle read on-chain information of calling requestDealInfo
3. Chainlink oracle requests url to get storage price in deal information
4. Chainlink oracle puts the information on-chain by calling callback function fulfill
5. Other on-chain contract(e.g. SwanPayment contract) calls getPrice to get storage price
