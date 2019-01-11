# Jiffy API

> Interact with contracts in a jiffy

Jiffy backend will store contract addresses, ABI's ,recorded labels for registered contracts and meta transactions so that users can easily manage.

#### Services provided

- [ ] AUTH using ethereum signatures
- [x] Registering new contract and address with unique URL
- [x] Registering labels per user address and contract
- [x] Assigning jiffy subdomains to contract
- [x] Transactions and calls sorted per contract address
- [x] Transaction/Call per label id
- [ ] Adding proper handling for ABI
- [ ] Adding API versioning

## Installation Instructions

```bash
$ dep ensure -v
```

## Run server

```bash
$ go run main.go
```
