package main

import (
	"context"
	"fmt"
	"github.com/Concordium/concordium-go-sdk"
)

func main() {
	ctx := context.Background()

	cli, err := concordium.NewClient(ctx, "35.184.87.228:10003", "rpcadmin")
	if err != nil {
		panic(err)
	}

	if false {
		r, err := cli.GetConsensusStatus(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
		fmt.Printf("%v\n", r.BestBlockHeight)
	}

	if false {
		r, err := cli.GetBlockInfo(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetAncestors(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25", 10)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetBranches(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetBlocksAtHeight(ctx, 88794)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetAccountList(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetInstances(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetAccountInfo(ctx,
			"4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25",
			"3djqZmm3jFEfMHXj4RtuTYLfr7VJ5ZwmVGmNot8sbadxFrA5eW",
		)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetInstanceInfo(ctx,
			"4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25",
			&concordium.ContractAddress{Index: 0, SubIndex: 0},
		)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if true {
		r, err := cli.InvokeContract(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25", &concordium.ContractContext{
			Invoker:   concordium.WrapAccountAddress("3djqZmm3jFEfMHXj4RtuTYLfr7VJ5ZwmVGmNot8sbadxFrA5eW"),
			Contract:  &concordium.ContractAddress{Index: 0, SubIndex: 0},
			Amount:    concordium.NewAmountZero(),
			Method:    "a.func",
			Parameter: "",
			Energy:    10000000,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetRewardStatus(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetBirkParameters(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetModuleList(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetModuleSource(ctx,
			"4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25",
			"85a8a9242518e07617763de99e5c6bdf39d82fa534a8838929a2167655002813",
		)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetIdentityProviders(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetAnonymityRevokers(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}

	if false {
		r, err := cli.GetCryptographicParameters(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", r)
	}

	if false {
		r, err := cli.GetAccountNonFinalizedTransactions(ctx, "3djqZmm3jFEfMHXj4RtuTYLfr7VJ5ZwmVGmNot8sbadxFrA5eW")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", r)
	}

	if false {
		r, err := cli.GetBlockSummary(ctx, "4eccaca49abab6df9d24ac8f0da973d4b2dbe6180810842b15cd1cc2078d0b25")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", r)
	}

	if false {
		r, err := cli.GetNextAccountNonce(ctx, "3djqZmm3jFEfMHXj4RtuTYLfr7VJ5ZwmVGmNot8sbadxFrA5eW")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", r)
	}
}
