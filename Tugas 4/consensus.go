// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)
// ChainHeaderReader adalah sekumpulan kecil metode yang diperlukan untuk mengakses blockchain lokal selama verifikasi header
type ChainHeaderReader interface {
	// Mengambil kembali konfigurasi chain blockchain
	Config() *params.ChainConfig

  // CurrentHeader digunakan untuk mengambil current header dari local chain
	CurrentHeader() *types.Header

  // GetHeader untuk mengambil block header dari database berdasarkan hash dan nomor.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber mengambil block number dari database berdasarkan nomor.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash mengambil header blok dari database berdasarkan hash-nya.
	GetHeaderByHash(hash common.Hash) *types.Header

	// GetTd untuk total difficulty dari database berdasarkan hash dan nomor.
	GetTd(hash common.Hash, number uint64) *big.Int
}
// Secara keseluruhan Interface ChainHeaderReader Untuk melakukan verifikasi header diperlukan pengembilan header dari rantai lokal,
// dari database dengan nomor, dari database dengan hash dan mengambil total kesulitan dari database dengan hash dan nomor.


// ChainReader adalah sekumpulan kecil metode yang diperlukan untuk mengakses blockchain lokal selama verifikasi header
type ChainReader interface {
	ChainHeaderReader

	// GetBlock mengambil blok dari database berdasarkan hash dan nomor.
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine adalah mesin konsensus agnostik algoritma.
type Engine interface {
  // Penulis mengambil alamat Ethereum dari akun yang mencetak blok yang diberikan,
  // yang mungkin berbeda dari header coinbase jika mesin konsensus didasarkan pada tanda tangan.
	Author(header *types.Header) (common.Address, error)

  // VerifyHeader memeriksa apakah header sesuai dengan aturan konsensus.
  // Memverifikasi segel dapat dilakukan secara opsional atau secara eksplisit melalui metode VerifySeal.
	VerifyHeader(chain ChainHeaderReader, header *types.Header, seal bool) error

	// VerifyHeaders mirip dengan VerifyHeader, tetapi VerifyHeaders memverifikasi sekumpulan header secara bersamaan.
        // VerifyHeaders mengembalikan quit channel untuk membatalkan operasi dan 
	// a results channel untuk mengambil async verifications
	VerifyHeaders(chain ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles memverifikasi bahwa block's uncles yang diberikan sesuai dengan aturan konsensus
	VerifyUncles(chain ChainReader, block *types.Block) error

	// Menginisialisasi consensus fields of a block header sesuai dengan aturan mesin
	// Perubahan dijalankan sebaris
	Prepare(chain ChainHeaderReader, header *types.Header) error

	// Finalize menjalankan modifikasi status pasca transaksi.
	// Note: Block header dan state database harus diperbarui untuk mencermikan aturan konsensus yang terjadi pada finalisasi.
	Finalize(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header)

	// FinalizeAndAssemble menjalankan modifikasi status pasca transaksi dan merakit final block.
	// Note: Block header dan state database harus diperbarui untuk mencermikan aturan konsensus yang terjadi pada finalisasi.
	FinalizeAndAssemble(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

	// Seal generate permintaan penyegelan baru untuk input block yang di berikan
  // dan mendorong hasilnya ke channel yang diberikan.
	
	// Note, metode ini mengembalikan segera dan akan mengirimkan hasil async.
  // Lebih dari satu hasil juga dapat dikembalikan tergantung pada algoritma konsensus
	Seal(chain ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

	// SealHash mengembalikan hash dari sebuah blok sebelum disegel.
	SealHash(header *types.Header) common.Hash

	// CalcDifficulty adalah algoritma penyesuaian difficulty.
  // Ini akan mengembalikan difficulty yang harus dimiliki oleh blok baru
	CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int

	// APIs mengembalikan API RPC yang disediakan mesin konsensus ini
	APIs(chain ChainHeaderReader) []rpc.API

	// Close untuk mengakhiri
	Close() error
}

// PoW adalah mesin konsensus berdasarkan bukti kerja.
type PoW interface {
	Engine

	// Hashrate mengembalikan the current mining hashrate dari PoW consensus engine.
	Hashrate() float64
}
