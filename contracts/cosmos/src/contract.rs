#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{
    to_json_binary, Addr, BankMsg, Binary, Coin, CosmosMsg, Deps, DepsMut, Env, MessageInfo, 
    Response, StdResult, Uint128, WasmMsg, ensure, Timestamp,
};
use cw2::set_contract_version;
use sha2::{Digest, Sha256};

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, EscrowResponse, EscrowsResponse};
use crate::state::{Config, Escrow, EscrowStatus, CONFIG, ESCROWS, ESCROW_COUNT};

// Contract name and version for migration
const CONTRACT_NAME: &str = "flow-fusion-cosmos";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    let config = Config {
        owner: info.sender.clone(),
        fee_collector: msg.fee_collector.unwrap_or(info.sender),
        fee_rate: msg.fee_rate.unwrap_or(0), // 0 = no fees for hackathon
    };

    CONFIG.save(deps.storage, &config)?;
    ESCROW_COUNT.save(deps.storage, &0u64)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("owner", config.owner)
        .add_attribute("fee_rate", config.fee_rate.to_string()))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::CreateEscrow {
            ethereum_tx_hash,
            hash_lock,
            time_lock,
            recipient,
            ethereum_sender,
        } => execute_create_escrow(
            deps, 
            env, 
            info, 
            ethereum_tx_hash,
            hash_lock, 
            time_lock, 
            recipient,
            ethereum_sender,
        ),
        ExecuteMsg::Withdraw { escrow_id, secret } => {
            execute_withdraw(deps, env, info, escrow_id, secret)
        }
        ExecuteMsg::Refund { escrow_id } => execute_refund(deps, env, info, escrow_id),
        ExecuteMsg::UpdateConfig { 
            fee_collector, 
            fee_rate 
        } => execute_update_config(deps, info, fee_collector, fee_rate),
    }
}

pub fn execute_create_escrow(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    ethereum_tx_hash: String,
    hash_lock: String,
    time_lock: u64,
    recipient: String,
    ethereum_sender: String,
) -> Result<Response, ContractError> {
    // Validate inputs
    ensure!(!ethereum_tx_hash.is_empty(), ContractError::InvalidInput("ethereum_tx_hash cannot be empty".to_string()));
    ensure!(!hash_lock.is_empty(), ContractError::InvalidInput("hash_lock cannot be empty".to_string()));
    ensure!(!recipient.is_empty(), ContractError::InvalidInput("recipient cannot be empty".to_string()));
    ensure!(!ethereum_sender.is_empty(), ContractError::InvalidInput("ethereum_sender cannot be empty".to_string()));
    ensure!(time_lock > 0, ContractError::InvalidInput("time_lock must be greater than 0".to_string()));

    // Validate recipient address
    let recipient_addr = deps.api.addr_validate(&recipient)?;

    // Ensure exactly one coin is sent
    ensure!(!info.funds.is_empty(), ContractError::NoFunds {});
    ensure!(info.funds.len() == 1, ContractError::MultipleDenoms {});
    
    let coin = &info.funds[0];
    ensure!(coin.amount > Uint128::zero(), ContractError::InsufficientFunds {});

    // Generate escrow ID
    let escrow_count = ESCROW_COUNT.load(deps.storage)?;
    let escrow_id = escrow_count + 1;
    ESCROW_COUNT.save(deps.storage, &escrow_id)?;

    // Calculate unlock time
    let unlock_time = env.block.time.plus_seconds(time_lock);

    // Create escrow
    let escrow = Escrow {
        id: escrow_id,
        ethereum_tx_hash: ethereum_tx_hash.clone(),
        sender: info.sender.clone(),
        recipient: recipient_addr.clone(),
        ethereum_sender: ethereum_sender.clone(),
        amount: coin.clone(),
        hash_lock: hash_lock.clone(),
        time_lock,
        unlock_time,
        status: EscrowStatus::Active,
        created_at: env.block.time,
    };

    ESCROWS.save(deps.storage, escrow_id, &escrow)?;

    Ok(Response::new()
        .add_attribute("method", "create_escrow")
        .add_attribute("escrow_id", escrow_id.to_string())
        .add_attribute("ethereum_tx_hash", ethereum_tx_hash)
        .add_attribute("sender", info.sender)
        .add_attribute("recipient", recipient_addr)
        .add_attribute("amount", coin.amount)
        .add_attribute("denom", &coin.denom)
        .add_attribute("hash_lock", hash_lock)
        .add_attribute("unlock_time", unlock_time.to_string()))
}

pub fn execute_withdraw(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    escrow_id: u64,
    secret: String,
) -> Result<Response, ContractError> {
    let mut escrow = ESCROWS.load(deps.storage, escrow_id)?;

    // Validate escrow status
    ensure!(escrow.status == EscrowStatus::Active, ContractError::EscrowNotActive {});

    // Validate secret against hash lock
    let secret_hash = hex::encode(Sha256::digest(secret.as_bytes()));
    ensure!(secret_hash == escrow.hash_lock, ContractError::InvalidSecret {});

    // Only recipient can withdraw
    ensure!(info.sender == escrow.recipient, ContractError::Unauthorized {});

    // Update escrow status
    escrow.status = EscrowStatus::Completed;
    ESCROWS.save(deps.storage, escrow_id, &escrow)?;

    // Transfer funds to recipient
    let transfer_msg = CosmosMsg::Bank(BankMsg::Send {
        to_address: escrow.recipient.to_string(),
        amount: vec![escrow.amount.clone()],
    });

    Ok(Response::new()
        .add_message(transfer_msg)
        .add_attribute("method", "withdraw")
        .add_attribute("escrow_id", escrow_id.to_string())
        .add_attribute("recipient", escrow.recipient)
        .add_attribute("amount", escrow.amount.amount)
        .add_attribute("secret_hash", secret_hash))
}

pub fn execute_refund(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    escrow_id: u64,
) -> Result<Response, ContractError> {
    let mut escrow = ESCROWS.load(deps.storage, escrow_id)?;

    // Validate escrow status
    ensure!(escrow.status == EscrowStatus::Active, ContractError::EscrowNotActive {});

    // Check if time lock has expired
    ensure!(env.block.time >= escrow.unlock_time, ContractError::TimeLockNotExpired {});

    // Only original sender can refund
    ensure!(info.sender == escrow.sender, ContractError::Unauthorized {});

    // Update escrow status
    escrow.status = EscrowStatus::Refunded;
    ESCROWS.save(deps.storage, escrow_id, &escrow)?;

    // Transfer funds back to sender
    let refund_msg = CosmosMsg::Bank(BankMsg::Send {
        to_address: escrow.sender.to_string(),
        amount: vec![escrow.amount.clone()],
    });

    Ok(Response::new()
        .add_message(refund_msg)
        .add_attribute("method", "refund")
        .add_attribute("escrow_id", escrow_id.to_string())
        .add_attribute("sender", escrow.sender)
        .add_attribute("amount", escrow.amount.amount))
}

pub fn execute_update_config(
    deps: DepsMut,
    info: MessageInfo,
    fee_collector: Option<String>,
    fee_rate: Option<u64>,
) -> Result<Response, ContractError> {
    let mut config = CONFIG.load(deps.storage)?;

    // Only owner can update config
    ensure!(info.sender == config.owner, ContractError::Unauthorized {});

    if let Some(collector) = fee_collector {
        config.fee_collector = deps.api.addr_validate(&collector)?;
    }

    if let Some(rate) = fee_rate {
        ensure!(rate <= 10000, ContractError::InvalidInput("fee_rate cannot exceed 10000 (100%)".to_string()));
        config.fee_rate = rate;
    }

    CONFIG.save(deps.storage, &config)?;

    Ok(Response::new()
        .add_attribute("method", "update_config")
        .add_attribute("fee_collector", config.fee_collector)
        .add_attribute("fee_rate", config.fee_rate.to_string()))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetEscrow { escrow_id } => to_json_binary(&query_escrow(deps, escrow_id)?),
        QueryMsg::GetEscrows { 
            start_after, 
            limit 
        } => to_json_binary(&query_escrows(deps, start_after, limit)?),
        QueryMsg::GetConfig {} => to_json_binary(&CONFIG.load(deps.storage)?),
    }
}

fn query_escrow(deps: Deps, escrow_id: u64) -> StdResult<EscrowResponse> {
    let escrow = ESCROWS.load(deps.storage, escrow_id)?;
    Ok(EscrowResponse { escrow })
}

fn query_escrows(
    deps: Deps,
    start_after: Option<u64>,
    limit: Option<u32>,
) -> StdResult<EscrowsResponse> {
    let limit = limit.unwrap_or(10).min(100) as usize;
    let start = start_after.map(|id| id + 1).unwrap_or(1);

    let escrows: StdResult<Vec<_>> = ESCROWS
        .range(deps.storage, Some(start.into()), None, cosmwasm_std::Order::Ascending)
        .take(limit)
        .map(|item| item.map(|(_, escrow)| escrow))
        .collect();

    Ok(EscrowsResponse {
        escrows: escrows?,
    })
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{coins, from_json};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            fee_collector: None,
            fee_rate: None,
        };
        let info = mock_info("creator", &coins(1000, "earth"));

        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // Query config
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetConfig {}).unwrap();
        let config: Config = from_json(&res).unwrap();
        assert_eq!("creator", config.owner);
    }

    #[test]
    fn create_and_complete_escrow() {
        let mut deps = mock_dependencies();
        
        // Initialize
        let msg = InstantiateMsg {
            fee_collector: None,
            fee_rate: None,
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Create escrow
        let info = mock_info("sender", &coins(100, "uatom"));
        let msg = ExecuteMsg::CreateEscrow {
            ethereum_tx_hash: "0x123".to_string(),
            hash_lock: "hash123".to_string(),
            time_lock: 3600,
            recipient: "recipient".to_string(),
            ethereum_sender: "0xeth123".to_string(),
        };

        let res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!("create_escrow", res.attributes[0].value);

        // Query escrow
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetEscrow { escrow_id: 1 }).unwrap();
        let escrow_response: EscrowResponse = from_json(&res).unwrap();
        assert_eq!(1, escrow_response.escrow.id);
        assert_eq!("0x123", escrow_response.escrow.ethereum_tx_hash);
    }
}