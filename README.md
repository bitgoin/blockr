[![Build Status](https://travis-ci.org/bitgoin/blockr.svg?branch=master)](https://travis-ci.org/bitgoin/blockr)
[![GoDoc](https://godoc.org/github.com/bitgoin/blockr?status.svg)](https://godoc.org/github.com/bitgoin/blockr)
[![GitHub license](https://img.shields.io/badge/license-BSD-blue.svg)](https://raw.githubusercontent.com/bitgoin/blockr/master/LICENSE)


# Blockr 

## Overview

This is a library to send transactions and gether unspent transaction outputs(UTXO) by [blockr](https://blockr.io/) web API.


## Requirements

This requires

* git
* go 1.3+


## Installation

     $ go get github.com/bitgoin/blockr


## Example
(This example omits error handlings for simplicity.)

## Key Handling

```go

import "github.com/bitgoin/blockr"

func main(){
	//prepare private key
 	txKey, err := address.FromWIF("some wif", address.BitcoinTest)

    //make service struct
	blk := Service{IsTestNet: true}

    //get utxos.
	txs, err := blk.GetUTXO(adr)

    //convert to utxo which can be used in the tx package.
	utxo, err := ToUTXO(txs, txKey)

    //prepare send info.
	send := []*tx.Send{
		&tx.Send{
			Addr:   "n3Bp1hbgtmwDtjQTpa6BnPPCA8fTymsiZy",
			Amount: 0.05 * tx.Unit,
		},
		&tx.Send{
			Addr:   "n2eMqTT929pb1RDNuqEnxdaLau1rxy3efi",
			Amount: 0.01 * tx.Unit,
		},
		&tx.Send{
			Addr:   adr,
			Amount: 0,
		},
	}

    //create tx.
 	tx, err := tx.NewP2PK(0.0001*tx.Unit, utxo, 0, send...)

	//send tx.
	txhash, err = blk.SendTX(tx)
}
```

## Normal Payment

```go
import gocoin

func main(){
	key, _ := gocoin.GenerateKey(true)

	//get unspent transactions
	service := gocoin.NewBlockrService(true)
	txs, _ := service.GetUTXO(adr,nil)
	
	//Normal Payment
	gocoin.Pay([]*Key{txKey}, []*gocoin.Amounts{&{gocoin.Amounts{"n2eMqTT929pb1RDNuqEnxdaLau1rxy3efi", 0.01*gocoin.BTC}}, service)
}
```

## M of N Multisig

```go
import gocoin

func main(){
	key, _ := gocoin.GenerateKey(true)
	service := gocoin.NewBlockrService(true)

	//2 of 3 multisig
	key1, _ := gocoin.GenerateKey(true)
	key2, _ := gocoin.GenerateKey(true)
	key3, _ := gocoin.GenerateKey(true)
	rs, _:= gocoin.NewRedeemScript(2, []*PublicKey{key1.Pub, key2.Pub, key3.Pub})
	//make a fund
	rs.Pay([]*Key{txKey}, 0.05*gocoin.BTC, service)

    //get a raw transaction for signing.
	rawtx, tx, _:= rs.CreateRawTransactionHashed([]*gocoin.Amounts{&{gocoin.Amounts{"n3Bp1hbgtmwDtjQTpa6BnPPCA8fTymsiZy", 0.05*gocoin.BTC}}, service)

	//spend the fund
	sign1, _:= key2.Priv.Sign(rawtx)
	sign2, _:= key3.Priv.Sign(rawtx)
	rs.Spend(tx, [][]byte{nil, sign1, sign2}, service)
}
```


## Micropayment Channel

```go
import gocoin

func main(){
	service := gocoin.NewBlockrService(true)

	key1, _ := gocoin.GenerateKey(true) //payer
	key2, _ := gocoin.GenerateKey(true) //payee

	payer, _:= gocoin.NewMicropayer(key1, key2.Pub, service)
	payee, _:= gocoin.NewMicropayee(key2, key1.Pub, service)

	txHash, _:= payer.CreateBond([]*Key{key1}, 0.05*BTC)

	locktime := time.Now().Add(time.Hour)
	sign, _:= payee.SignToRefund(txHash, 0.05*gocoin.BTC-gocoin.Fee, uint32(locktime.Unix()))
	payer.SendBond(uint32(locktime.Unix()), sign) //return an error if payee's sig is invalid

	signIP, _:= payer.SignToIncrementedPayment(0.001 * gocoin.BTC)
	payee.IncrementPayment(0.001*gocoin.BTC, signIP) //return an error if payer's sig is invalid
	//more payments

	payee.SendLastPayment()
	//or
	//	payer.SendRefund() after locktime

}
```

# Contribution
Improvements to the codebase and pull requests are encouraged.


