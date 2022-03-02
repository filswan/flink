
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
   "jobRunID":0,
   "data":{
      "code":200,
      "message":"",
      "data":{
         "DealId":58160,
         "dealCid":"bafyreifo2pp5d4se44xu32p5ikm3qjzmfv7ihbmdsilz7j5wii7h7ne3gm",
         "messageCid":"bafy2bzacebyuyuxr23e2n3njhtltcuc7sc73cuumzd2fww4mt4ivtzg2zn6um",
         "height":455035,
         "pieceCid":"baga6ea4seaql3pcfitmlane3nbrlcitb4ffzdkkswy4e2tn4tf67muicdcueiki",
         "verifiedDeal":false,
         "storagePricePerEpoch":"976562 AttoFIL",
         "signature":"rWJeBFkmGZTAnIitvM6NiRpn8vqwlRAjr4PpMHmfL6Kb86qeXU99DtHWmjW8WyARAFn3mTUtB4+rlibfEUFlts4cAESxfHPiuOciVj0r0d8Y3te0axEZETGsJeLQPPkY",
         "signatureType":"bls",
         "createdAt":"2021-11-24 07:57:30",
         "pieceSizeFormat":"2.00 MiB",
         "satrtHeight":466508,
         "endHeight":1980929,
         "client":"t3u7pumush376xbytsgs5wabkhtadjzfydxxda2vzyasg7cimkcphswrq66j4dubbhwpnojqd3jie6ermpwvvq",
         "clientCollateralFormat":"0 AttoFIL",
         "provider":"t024557",
         "providerTag":"",
         "providerIsVerified":0,
         "providerCollateralFormat":"0 AttoFIL",
         "status":0
      },
      "result":{
         "dealCid":"bafyreifo2pp5d4se44xu32p5ikm3qjzmfv7ihbmdsilz7j5wii7h7ne3gm",
         "price":"976562 AttoFIL"
      }
   },
   "result":{
      "dealCid":"bafyreifo2pp5d4se44xu32p5ikm3qjzmfv7ihbmdsilz7j5wii7h7ne3gm",
      "price":"976562 AttoFIL"
   },
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

### Run

```bash
yarn start
```

## Call the external adapter/API server
Flink is supporting both GET and POST api calls

POST /deal
```bash
curl -X POST -H "content-type:application/json" "http://localhost:<port>/deal" --data '{ "id": 0, "data": { "deal":"58160", "network":"filecoin_mainnet"} }'
```
GET /deal/{deal_id}?network=filecoin_mainnet
```bash
curl -X GET "http://localhost:<port>/deal/58160?network=filecoin_mainnet"
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
