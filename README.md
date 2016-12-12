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


# Contribution
Improvements to the codebase and pull requests are encouraged.


