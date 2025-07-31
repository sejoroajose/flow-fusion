use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::state::{Config, Escrow};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct InstantiateMsg {
    pub fee_collector: Option<String>,
    pub fee_rate: Option<u64>, // Basis points (10000 = 100%)
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
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

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
    GetEscrow {
        escrow_id: u64,
    },
    GetEscrows {
        start_after: Option<u64>,
        limit: Option<u32>,
    },
    GetConfig {},
}

// Response types
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct EscrowResponse {
    pub escrow: Escrow,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct EscrowsResponse {
    pub escrows: Vec<Escrow>,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct ConfigResponse {
    pub config: Config,
}