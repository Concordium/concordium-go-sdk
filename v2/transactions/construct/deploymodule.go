package construct

// TODO: fix it.
// DeployModule deploys the given Wasm module. The module is given
// as a binary source, and no processing is done to the module.
//func DeployModule(numSigs uint32, sender v2.AccountAddress, nonce v2.Nonce, expiry v2.TransactionTime,
//	module v2.WasmModule) *v2.PreAccountTransaction {
//	moduleSize := module.Source.Size()
//	payload := &v2.AccountTransactionPayload{Payload: &v2.DeployModulePayload{Module: module}}
//	energy := &v2.AddEnergy{
//		NumSigs: numSigs,
//		Energy:  costs.DeployModuleCost(moduleSize),
//	}
//	return makeTransaction(sender, nonce, expiry, energy, payload)
//}
