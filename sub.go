package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func subNewBlock() {
	ctx := context.Background()
	esClient, err := elastic.NewClient(elastic.SetURL(config.ElasticURL), elastic.SetSniff(config.ElasticSniff))
	if err != nil {
		log.Fatalln(err.Error())
	}
	// esClient.DeleteIndex("eth_sub_address").Do(ctx)
	csv2es(ctx, esClient)

	nodeClient, err := ethclient.Dial(config.EthRPC)
	if err != nil {
		log.Fatalln(err.Error())
	}

	blockCh := make(chan *types.Header)
	sub, err := nodeClient.SubscribeNewHead(ctx, blockCh)
	if err != nil {
		log.Error(err.Error())
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalln(err.Error())
		case head := <-blockCh:
			block, err := nodeClient.BlockByNumber(ctx, head.Number)
			if err != nil {
				log.Fatalln(err.Error())
			}
			txs := block.Transactions()
			for _, tx := range txs {
				fmt.Println("tx", tx.Hash().Hex())
				fmt.Println("to", tx.To().Hex())
			}
		}
	}
}