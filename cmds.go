package main

/*
AMB supported commands
*/

var (
	supportedRpcsBlockchain = map[string]bool{
		"getbestblockhash":      true,
		"getblock":              true,
		"getblockchaininfo":     true,
		"getblockcount":         true,
		"getblockhash":          true,
		"getblockfilter":        true,
		"getblockheader":        true,
		"getblockstats":         true,
		"getchaintips":          true,
		"getdifficulty":         true,
		"getmempoolancestors":   true,
		"getmempooldescendants": true,
		"getmempoolentry":       true,
		"getmempoolinfo":        true,
		"getrawmempool":         true,
		"gettxout":              true,
		"gettxoutproof":         true,
	}

	supportedRpcsRawTransactions = map[string]bool{

		"createrawtransaction": true,
		"decoderawtransaction": true,
		"decodescript":         true,
		"getrawtransaction":    true,
		"sendrawtransaction":   true,
	}

	supportedRpcsUtil = map[string]bool{
		"createmultisig":   true,
		"estimatesmartfee": true,
		"validateaddress":  true,
		"verifymessage":    true,
	}
)
