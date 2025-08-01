use cosmwasm_schema::cw_serde;
use cosmwasm_std::{Addr, Coin, Timestamp};
use cw_storage_plus::{Item, Map};

#[cw_serde]
pub struct Config {
    pub owner: Addr,
    pub fee_collector: Addr,
    pub fee_rate: u64, // Basis points (10000 = 100%)
}

#[cw_serde]
pub struct Escrow {
    pub id: u64,
    pub ethereum_tx_hash: String,
    pub sender: Addr,
    pub recipient: Addr,
    pub ethereum_sender: String, // Ethereum address as string
    pub amount: Coin,
    pub hash_lock: String,       // SHA256 hash of the secret
    pub time_lock: u64,          // Duration in seconds
    pub unlock_time: Timestamp,  // When the escrow can be refunded
    pub status: EscrowStatus,
    pub created_at: Timestamp,
}

#[cw_serde]
pub enum EscrowStatus {
    Active,
    Completed,
    Refunded,
}

// Storage keys
pub const CONFIG: Item<Config> = Item::new("config");
pub const ESCROWS: Map<u64, Escrow> = Map::new("escrows");
pub const ESCROW_COUNT: Item<u64> = Item::new("escrow_count");