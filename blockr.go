/*
 * Copyright (c) 2015, Shinya Yagyu
 * All rights reserved.
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 * 3. Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from this
 *    software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * See LICENSE file for the original license:
 *
 * This file also includes codes from https://github.com/soroushjp/go-bitcoin-multisig
 * copyrighted by Soroush Pour.
 */

package blockr

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bitgoin/address"
	"github.com/bitgoin/tx"
)

type unspent struct {
	Status string
	Data   struct {
		Address string
		Unspent []Unspent
	}
	Code    int
	Message string
}

type sendtx struct {
	Status  string
	Data    string
	Code    int
	Message string
}

//Unspent represents an available transaction.
type Unspent struct {
	Tx            string
	Amount        string
	N             int
	Confirmations int
	Script        string
}

var cacheUTXO = make(map[string][]Unspent)

//spent sets the cache that tx hash is already spent.
func spent(tra *tx.Tx) {
	for _, txin := range tra.TxIn {
		hh := tx.Reverse(txin.Hash)
		h := hex.EncodeToString(hh)
		for k, v := range cacheUTXO {
			for i, utxo := range v {
				if h == utxo.Tx {
					cacheUTXO[k] = append(v[0:i], v[i+1:]...)
					continue
				}
			}
		}
	}
}

//Service is a service using Blockr.io.
type Service struct {
	IsTestNet bool
}

//SendTX send a transaction using Blockr.io.
func (b *Service) SendTX(tra *tx.Tx) ([]byte, error) {
	btc := "btc"
	if b.IsTestNet {
		btc = "tbtc"
	}
	data, err := tra.Pack()
	if err != nil {
		return nil, err
	}
	log.Print(hex.EncodeToString(data))
	resp, err := http.PostForm("http://"+btc+".blockr.io/api/v1/tx/push",
		url.Values{"hex": {hex.EncodeToString(data)}})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Print(err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var u sendtx
	if err = json.Unmarshal(body, &u); err != nil {
		return nil, err
	}
	if u.Status != "success" {
		return nil, errors.New("blockr returns: " + u.Message)
	}
	spent(tra)
	return hex.DecodeString(u.Data)
}

//GetUTXO gets unspent transaction outputs by using Blockr.io.
func (b *Service) GetUTXO(addr string) ([]Unspent, error) {
	if cacheUTXO[addr] != nil {
		return cacheUTXO[addr], nil
	}
	btc := "btc"
	if b.IsTestNet {
		btc = "tbtc"
	}

	resp, err := http.Get("http://" + btc + ".blockr.io/api/v1/address/unspent/" + addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var u unspent
	if err = json.Unmarshal(body, &u); err != nil {
		return nil, err
	}
	if u.Status != "success" {
		return nil, errors.New("blockr returns: " + u.Message)
	}

	cacheUTXO[addr] = u.Data.Unspent
	return u.Data.Unspent, nil
}

//ToUTXO returns utxo in transaction package.
func ToUTXO(utxos []Unspent, privs *address.PrivateKey) (tx.UTXOs, error) {
	txs := make(tx.UTXOs, len(utxos))
	for i, utxo := range utxos {
		hash, err := hex.DecodeString(utxo.Tx)
		if err != nil {
			return nil, err
		}
		hash = tx.Reverse(hash)
		script, err := hex.DecodeString(utxo.Script)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(utxo.Amount, 64)
		if err != nil {
			return nil, err
		}
		txs[i] = &tx.UTXO{
			Key:     privs,
			Value:   uint64(amount * tx.Unit),
			TxHash:  hash,
			TxIndex: uint32(utxo.N),
			Script:  script,
		}
	}
	return txs, nil
}
