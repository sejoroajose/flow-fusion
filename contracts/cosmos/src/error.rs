use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},

    #[error("Invalid input: {0}")]
    InvalidInput(String),

    #[error("No funds sent")]
    NoFunds {},

    #[error("Multiple denoms not allowed")]
    MultipleDenoms {},

    #[error("Insufficient funds")]
    InsufficientFunds {},

    #[error("Escrow not active")]
    EscrowNotActive {},

    #[error("Invalid secret")]
    InvalidSecret {},

    #[error("Time lock not expired")]
    TimeLockNotExpired {},

    #[error("Escrow not found")]
    EscrowNotFound {},

    #[error("Custom Error val: {val:?}")]
    CustomError { val: String },
}