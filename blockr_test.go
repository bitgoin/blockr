/*
 * Copyright (c) 2016, Shinya Yagyu
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
 */

package blockr

import (
	"log"
	"testing"

	"github.com/bitgoin/address"
	"github.com/bitgoin/tx"
)

func TestSend(t *testing.T) {
	wif := "928Qr9J5oAC6AYieWJ3fG3dZDjuC7BFVUqgu4GsvRVpoXiTaJJf"
	txKey, err := address.FromWIF(wif, address.BitcoinTest)
	if err != nil {
		t.Errorf(err.Error())
	}
	adr := txKey.PublicKey.Address()
	log.Println("address for tx=", adr)
	if adr != "n3Bp1hbgtmwDtjQTpa6BnPPCA8fTymsiZy" {
		t.Errorf("invalid address")
	}
	blk := Service{IsTestNet: true}
	txs, err := blk.GetUTXO(adr)
	if err != nil {
		t.Error(err)
	}
	log.Println("UTXO:")
	for _, tx := range txs {
		log.Println("hash", tx.Tx)
		log.Println("amount", tx.Amount)
		log.Println("index", tx.N)
		log.Println("script", tx.Script)
	}
	utxo, err := ToUTXO(txs, txKey)
	if err != nil {
		t.Error(err)
	}
	send := []*tx.Send{
		&tx.Send{
			Addr:   "n2eMqTT929pb1RDNuqEnxdaLau1rxy3efi",
			Amount: 0.01 * tx.Unit,
		},
		&tx.Send{
			Addr:   adr,
			Amount: 0,
		},
	}

	tx, err := tx.NewP2PK(0.0001*tx.Unit, utxo, 0, send...)
	if err != nil {
		t.Error(err)
	}
	_, err = blk.SendTX(tx)
	if err != nil {
		t.Error(err)
	}
}
