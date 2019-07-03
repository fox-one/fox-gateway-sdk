package main

import (
	"context"

	gateway "github.com/fox-one/fox-gateway-sdk"
	log "github.com/sirupsen/logrus"
)

const (
	api = "https://openapi.fox.one"

	// 商户 key secret
	merchantKey    = ""
	merchantSecret = ""

	// 交易所业务
	exchangeService = "exchange"
)

func main() {
	ctx := context.Background()

	merchantSrv := gateway.NewMerchantClient(api).WithSession(merchantKey, merchantSecret)

	// create member
	output, err := merchantSrv.CreateMember(ctx, true)
	if err != nil {
		log.Panic(err)
	}

	member := output.Member

	// update member profile
	if _, err := merchantSrv.Member(member.ID).UpdateProfile(ctx, "name", "http://avatar.com"); err != nil {
		log.Panicf("update profile failed: %v", err)
	}

	// set pin first time
	if err := merchantSrv.Member(member.ID).UpdatePin(ctx, "", "123456"); err != nil {
		log.Panicf("set pin failed: %v", err)
	}

	// update pin
	if err := merchantSrv.Member(member.ID).UpdatePin(ctx, "123456", "654321"); err != nil {
		log.Panicf("update pin failed: %v", err)
	}

	// wallet
	{
		srv := merchantSrv.Member(member.ID).Service(exchangeService)
		srvWithPin := merchantSrv.Member(member.ID).ServiceWithPin(exchangeService, "123456")

		// read assets
		if assets, err := srv.ReadAssets(ctx, 0); err == nil {
			log.Infof("read %d assets", len(assets))
		} else {
			log.Panicf("read assets failed: %v", err)
		}

		// read asset detail
		BTC := "c6d0c728-2624-429b-8e0d-d9d19b6592fa"
		if asset, err := srv.ReadAsset(ctx, BTC); err == nil {
			log.Infof("btc balance: %s", asset.Balance)
		} else {
			log.Panicf("read asset failed: %v", err)
		}

		op := gateway.WalletAssetOperation{
			AssetID: BTC,
			Amount:  "100",
			Memo:    "test",
		}

		// transfer
		transfer := &gateway.WalletTransferOperation{
			OpponentID:           "", // mixin id of target wallet
			WalletAssetOperation: op,
		}

		if snapshot, err := srvWithPin.Transfer(ctx, transfer); err != nil {
			log.Errorf("transfer failed: %v", err)
		} else {
			log.Infof("snapshot id: %s", snapshot.SnapshotID)
		}

		// withdraw
		withdraw := &gateway.WalletWithdrawOperation{
			PublicKey:            "", // withdraw address
			WalletAssetOperation: op,
		}

		if snapshot, err := srvWithPin.Withdraw(ctx, withdraw); err != nil {
			log.Errorf("withdraw failed: %v", err)
		} else {
			log.Infof("snapshot id: %s", snapshot.SnapshotID)
		}
	}
}
