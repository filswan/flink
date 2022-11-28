const createRequest = require('./index').createRequest
const configfile = require('./config.json')

const express = require('express')
const bodyParser = require('body-parser')
const app = express()
const portConfig = configfile.port

const port = portConfig || 8080

app.use(bodyParser.json())

app.post('/deal', (req, res) => {
  console.log('POST Data: ', req.body)
  createRequest('post', req.body, (status, result) => {
    console.log('Deal ID: ', result?.data?.deal?.deal_id)
    res.status(status).json(result)
  })
})

app.get('/deal/:deal_id', (req, res) => {
  const deal = req.params.deal_id
  const network = req.query.network
  console.log('GET Data: deal=', deal, ' and network=', network)
  const req_body = { id: 0, data: { deal: deal, network: network } }
  createRequest('get', req_body, (status, result) => {
    console.log('Deal ID: ', result?.data?.deal?.deal_id)
    res.status(status).json(result)
  })
})

app.listen(port, () => console.log(`Listening on port ${port}!`))
