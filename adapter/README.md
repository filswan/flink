
# Chainlink Adapter - Filecoin Data DAO 

Data Dao signs the node for providing data feeds to the web3 blockchains like polygon, bsc, eth. On filecoin network data deal has an unique *deal id*. This data adapter provides deal ids for other blockchain systems to check if the data deal exists on filecoin network.


## Creating the adapter 

```bash
git clone https://github.com/filswan/filink
```

Enter into the newly-created directory

```bash
cd filink/adapter
```

You can remove the existing git history by running:

```bash
rm -rf .git
```

See [Install Locally](#install-locally) for a quickstart

## Input Params

- `deal`, or `dealId`: The dealId of filecoin transaction

## Output
Only need to take care of result part

```json
{
   "jobRunID":"2bb15c3f9cfc4336b95012872ff05092",
   "data":{
      "deal":{
         "deal_id":5210178,
         "deal_cid":"",
         "message_cid":"bafy2bzaceaotial6pogwzvm7woh5pf37sivrzm3fmp5teao365jl22z5q4pfc",
         "height":1697382,
         "piece_cid":"baga6ea4seaqjffbc2mmed2piulix5qfppyuhbqumnppme5ngj3q2ol4udijjqbq",
         "verified_deal":true,
         "storage_price_per_epoch":0,
         "signature":"",
         "signature_type":"",
         "created_at":1649227860,
         "piece_size":"1073741824",
         "start_height":1701360,
         "end_height":3234661,
         "client":"f1g463yb4ok3lq3tffkvvfmfyngcagpx4kg7c7rei",
         "client_collateral_format":"000000000000000000",
         "provider":"f067375",
         "provider_tag":"",
         "verified_provider":0,
         "provider_collateral_format":"000000000000000000",
         "status":0,
         "network_name":"filecoin_mainnet",
         "storage_price":0
      },
      "result":0
   },
   "result":0,
   "statusCode":200
}
```

## Install Locally

Install dependencies:

```bash
yarn
```
Default running port is 8080. By changing "port", you may set the new running port
Update hostname and port pointing to flink back-end in file adapter/config.json:

```bash
{
    "url" : "http://<host>:<port>/deal/"
    "port": <port>
}

```

### Test

Run the local tests:

```bash
yarn test
```
### Run Adapter in Production
We need pm2 to run and monitor node application 
1. Step 1: Check if pm2 has been installed. If not, continue step 2; else continue step 3
```bash
pm2 --version
```
2. Step 2: In order to install pm2, run below command
```bash
sudo yarn global add pm2 
```
3. Step 3: Start the adapter with below command
```bash
pm2 start app.js
```
4. Step 4: Verify if adapter is running
```bash
pm2 list
```

## Call the external adapter/API server
Flink is supporting both GET and POST api calls

POST /deal
```bash
curl -X POST -H "content-type:application/json" "http://localhost:<port>/deal" --data '{ "id": 0, "data": { "deal":"58160", "network":"filecoin_mainnet"} }'
```
Response format for POST
```bash
{
   "jobRunID":"2bb15c3f9cfc4336b95012872ff05092",
   "data":{
      "deal":{
         "deal_id":5210178,
         "deal_cid":"",
         "message_cid":"bafy2bzaceaotial6pogwzvm7woh5pf37sivrzm3fmp5teao365jl22z5q4pfc",
         "height":1697382,
         "piece_cid":"baga6ea4seaqjffbc2mmed2piulix5qfppyuhbqumnppme5ngj3q2ol4udijjqbq",
         "verified_deal":true,
         "storage_price_per_epoch":0,
         "signature":"",
         "signature_type":"",
         "created_at":1649227860,
         "piece_size":"1073741824",
         "start_height":1701360,
         "end_height":3234661,
         "client":"f1g463yb4ok3lq3tffkvvfmfyngcagpx4kg7c7rei",
         "client_collateral_format":"000000000000000000",
         "provider":"f067375",
         "provider_tag":"",
         "verified_provider":0,
         "provider_collateral_format":"000000000000000000",
         "status":0,
         "network_name":"filecoin_mainnet",
         "storage_price":0
      },
      "result":0
   },
   "result":0,
   "statusCode":200
}
```
GET /deal/{deal_id}?network=filecoin_mainnet
```bash
curl -X GET "http://localhost:<port>/deal/58160?network=filecoin_mainnet"
```
Response format for GET
```bash
{
   "jobRunID":"2bb15c3f9cfc4336b95012872ff05092",
   "data":{
      "deal":{
         "deal_id":5210178,
         "deal_cid":"",
         "message_cid":"bafy2bzaceaotial6pogwzvm7woh5pf37sivrzm3fmp5teao365jl22z5q4pfc",
         "height":1697382,
         "piece_cid":"baga6ea4seaqjffbc2mmed2piulix5qfppyuhbqumnppme5ngj3q2ol4udijjqbq",
         "verified_deal":true,
         "storage_price_per_epoch":0,
         "signature":"",
         "signature_type":"",
         "created_at":1649227860,
         "piece_size":"1073741824",
         "start_height":1701360,
         "end_height":3234661,
         "client":"f1g463yb4ok3lq3tffkvvfmfyngcagpx4kg7c7rei",
         "client_collateral_format":"000000000000000000",
         "provider":"f067375",
         "provider_tag":"",
         "verified_provider":0,
         "provider_collateral_format":"000000000000000000",
         "status":0,
         "network_name":"filecoin_mainnet",
         "storage_price":0
      },
      "result":0
   },
   "result":0,
   "statusCode":200
}
```
## Docker

If you wish to use Docker to run the adapter, you can build the image by running the following command:

```bash
docker build . -t external-adapter
```

Then run it with:

```bash
docker run -p 8080:8080 -it external-adapter:latest
```

## Serverless hosts

After [installing locally](#install-locally):

### Create the zip

```bash
zip -r external-adapter.zip .
```

### Install to AWS Lambda

- In Lambda Functions, create function
- On the Create function page:
  - Give the function a name
  - Use Node.js 12.x for the runtime
  - Choose an existing role or create a new one
  - Click Create Function
- Under Function code, select "Upload a .zip file" from the Code entry type drop-down
- Click Upload and select the `external-adapter.zip` file
- Handler:
    - index.handler for REST API Gateways
    - index.handlerv2 for HTTP API Gateways
- Add the environment variable (repeat for all environment variables):
  - Key: API_KEY
  - Value: Your_API_key
- Save

#### To Set Up an API Gateway (HTTP API)

If using a HTTP API Gateway, Lambda's built-in Test will fail, but you will be able to externally call the function successfully.

- Click Add Trigger
- Select API Gateway in Trigger configuration
- Under API, click Create an API
- Choose HTTP API
- Select the security for the API
- Click Add

#### To Set Up an API Gateway (REST API)

If using a REST API Gateway, you will need to disable the Lambda proxy integration for Lambda-based adapter to function.

- Click Add Trigger
- Select API Gateway in Trigger configuration
- Under API, click Create an API
- Choose REST API
- Select the security for the API
- Click Add
- Click the API Gateway trigger
- Click the name of the trigger (this is a link, a new window opens)
- Click Integration Request
- Uncheck Use Lamba Proxy integration
- Click OK on the two dialogs
- Return to your function
- Remove the API Gateway and Save
- Click Add Trigger and use the same API Gateway
- Select the deployment stage and security
- Click Add

### Install to GCP

- In Functions, create a new function, choose to ZIP upload
- Click Browse and select the `external-adapter.zip` file
- Select a Storage Bucket to keep the zip in
- Function to execute: gcpservice
- Click More, Add variable (repeat for all environment variables)
  - NAME: API_KEY
  - VALUE: Your_API_key
