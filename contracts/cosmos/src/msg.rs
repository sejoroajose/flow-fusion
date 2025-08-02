use cosmwasm_schema::{cw_serde, QueryResponses};
use crate::state::{Config, Escrow};

#[cw_serde]
pub struct InstantiateMsg {
    pub fee_collector: Option<String>,
    pub fee_rate: Option<u64>, // Basis points (10000 = 100%)
}

#[cw_serde]
pub enum ExecuteMsg {
    CreateEscrow {
        ethereum_tx_hash: String,
        hash_lock: String,
        time_lock: u64, // Seconds
        recipient: String,
        ethereum_sender: String,
    },
    Withdraw {
        escrow_id: u64,
        secret: String,
    },
    Refund {
        escrow_id: u64,
    },
    UpdateConfig {
        fee_collector: Option<String>,
        fee_rate: Option<u64>,
    },
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(EscrowResponse)]
    GetEscrow { escrow_id: u64 },
    #[returns(EscrowsResponse)]
    GetEscrows {
        start_after: Option<u64>,
        limit: Option<u32>,
    },
    #[returns(Config)]
    GetConfig {},
}

// Response types
#[cw_serde]
pub struct EscrowResponse {
    pub escrow: Escrow,
}

#[cw_serde]
pub struct EscrowsResponse {
    pub escrows: Vec<Escrow>,
}