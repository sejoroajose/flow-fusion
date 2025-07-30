// SPDX-License-Identifier: MIT
pragma solidity 0.8.30;

import {Ownable} from "./@openzeppelin/contracts/access/Ownable.sol";
import {IERC20} from "./@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "./@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "./@openzeppelin/contracts/security/ReentrancyGuard.sol";

import {IEscrowFactory} from "../lib/cross-chain-swap/contracts/interfaces/IEscrowFactory.sol";
import {IBaseEscrow} from "../lib/cross-chain-swap/contracts/interfaces/IBaseEscrow.sol";
import {TimelocksLib, Timelocks} from "../lib/cross-chain-swap/contracts/libraries/TimelocksLib.sol";

/**
 * @title FlowFusionEscrowFactory
 * @dev Enhanced escrow factory for Cosmos bridge with TWAP support
 * @custom:security-contact security@flow-fusion.com
 */
contract FlowFusionEscrowFactory is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;
    using TimelocksLib for Timelocks;

    // Events
    event TWAPOrderCreated(
        bytes32 indexed orderId,
        address indexed maker,
        address indexed token,
        uint256 totalAmount,
        uint256 timeWindow,
        uint256 intervalCount,
        string cosmosRecipient
    );

    event TWAPIntervalExecuted(
        bytes32 indexed orderId,
        uint256 indexed intervalIndex,
        uint256 amount,
        bytes32 secretHash
    );

    event CosmosEscrowCreated(
        bytes32 indexed orderId,
        IBaseEscrow.Immutables immutables,
        string cosmosRecipient
    );

    // Structs
    struct TWAPConfig {
        uint256 totalAmount;      // Total amount to swap
        uint256 timeWindow;       // Total time window in seconds
        uint256 intervalCount;    // Number of execution intervals
        uint256 maxSlippage;      // Maximum slippage per interval (basis points)
        uint256 startTime;        // When execution should start
    }

    struct TWAPOrder {
        address maker;
        address token;
        TWAPConfig config;
        string cosmosRecipient;   // Cosmos address to receive tokens
        uint256 executedIntervals;
        uint256 totalExecuted;
        bool cancelled;
        mapping(uint256 => bytes32) intervalSecrets;
    }

    // State variables
    IEscrowFactory public immutable escrowFactory;
    mapping(bytes32 => TWAPOrder) public twapOrders;
    mapping(address => bool) public authorizedResolvers;

    // Constants
    uint256 public constant MAX_INTERVALS = 100;
    uint256 public constant MIN_INTERVAL_DURATION = 60; // 1 minute
    uint256 public constant MAX_SLIPPAGE = 1000; // 10%

    error InvalidTWAPConfig();
    error UnauthorizedResolver();
    error OrderAlreadyExists();
    error OrderNotFound();
    error IntervalAlreadyExecuted();
    error OrderCancelled();
    error InvalidInterval();

    constructor(
        IEscrowFactory _escrowFactory,
        address initialOwner
    ) Ownable(initialOwner) {
        escrowFactory = _escrowFactory;
    }

    /**
     * @notice Create a new TWAP order for cross-chain swap to Cosmos
     */
    function createTWAPOrder(
        bytes32 orderId,
        address token,
        TWAPConfig calldata config,
        string calldata cosmosRecipient
    ) external payable {
        require(token != address(0), "Invalid token address");
        require(bytes(cosmosRecipient).length >= 10 && bytes(cosmosRecipient).length <= 64, "Invalid Cosmos recipient");
        require(config.totalAmount > 0, "Invalid amount");
        require(config.intervalCount > 0 && config.intervalCount <= MAX_INTERVALS, "Invalid interval count");
        require(config.timeWindow / config.intervalCount >= MIN_INTERVAL_DURATION, "Invalid interval duration");
        require(config.maxSlippage <= MAX_SLIPPAGE, "Invalid slippage");
        require(config.startTime >= block.timestamp, "Invalid start time");

        // Validate TWAP configuration
        if (
            config.totalAmount == 0 ||
            config.intervalCount == 0 ||
            config.intervalCount > MAX_INTERVALS ||
            config.timeWindow / config.intervalCount < MIN_INTERVAL_DURATION ||
            config.maxSlippage > MAX_SLIPPAGE ||
            bytes(cosmosRecipient).length == 0
        ) {
            revert InvalidTWAPConfig();
        }

        // Ensure order doesn't exist
        if (twapOrders[orderId].maker != address(0)) {
            revert OrderAlreadyExists();
        }

        // Transfer tokens from maker
        IERC20(token).safeTransferFrom(msg.sender, address(this), config.totalAmount);

        // Create TWAP order
        TWAPOrder storage order = twapOrders[orderId];
        order.maker = msg.sender;
        order.token = token;
        order.config = config;
        order.cosmosRecipient = cosmosRecipient;

        emit TWAPOrderCreated(
            orderId,
            msg.sender,
            token,
            config.totalAmount,
            config.timeWindow,
            config.intervalCount,
            cosmosRecipient
        );
    }

    /**
     * @notice Execute a TWAP interval (called by authorized resolver)
     */
    function executeTWAPInterval(
        bytes32 orderId,
        uint256 intervalIndex,
        bytes32 secretHash,
        IBaseEscrow.Immutables calldata immutables
    ) external payable nonReentrant {
        require(secretHash != bytes32(0), "Invalid secret hash");
        require(immutables.amount > 0, "Invalid escrow amount");
        require(immutables.maker != address(0), "Invalid maker address");

        // Check authorization
        if (!authorizedResolvers[msg.sender]) {
            revert UnauthorizedResolver();
        }

        TWAPOrder storage order = twapOrders[orderId];
        
        // Validate order
        if (order.maker == address(0)) {
            revert OrderNotFound();
        }
        if (order.cancelled) {
            revert OrderCancelled();
        }
        if (intervalIndex >= order.config.intervalCount) {
            revert InvalidInterval();
        }
        if (order.intervalSecrets[intervalIndex] != bytes32(0)) {
            revert IntervalAlreadyExecuted();
        }

        // Calculate interval amount
        uint256 intervalAmount = order.config.totalAmount / order.config.intervalCount;
        
        // Handle remainder in last interval
        if (intervalIndex == order.config.intervalCount - 1) {
            intervalAmount = order.config.totalAmount - order.totalExecuted;
        }

        // Create escrow for this interval
        IBaseEscrow.Immutables memory intervalImmutables = immutables;
        intervalImmutables.amount = intervalAmount;
        intervalImmutables.hashLock = secretHash;
        intervalImmutables.maker = order.maker;

        // Transfer tokens to escrow
        IERC20(order.token).safeTransfer(address(escrowFactory), intervalAmount);
        
        // Create escrow via factory
        escrowFactory.createSrcEscrow{value: msg.value}(intervalImmutables, block.timestamp);

        // Update order state
        order.intervalSecrets[intervalIndex] = secretHash;
        order.executedIntervals++;
        order.totalExecuted += intervalAmount;

        emit TWAPIntervalExecuted(orderId, intervalIndex, intervalAmount, secretHash);
        emit CosmosEscrowCreated(orderId, intervalImmutables, order.cosmosRecipient);
    }

    event TWAPOrderCancelled(bytes32 indexed orderId, address indexed maker, uint256 refundedAmount);
    /**
     * @notice Cancel a TWAP order (only maker can cancel)
     */
    function cancelTWAPOrder(bytes32 orderId) external nonReentrant {
        TWAPOrder storage order = twapOrders[orderId];
        
        require(order.maker == msg.sender, "Unauthorized");
        require(!order.cancelled, "Order already cancelled");

        order.cancelled = true;
        uint256 remainingAmount = order.config.totalAmount - order.totalExecuted;
        if (remainingAmount > 0) {
            IERC20(order.token).safeTransfer(order.maker, remainingAmount);
            emit TWAPOrderCancelled(orderId, order.maker, remainingAmount);
        }
    }

    /**
     * @notice Get TWAP order details
     */
    function getTWAPOrder(bytes32 orderId) external view returns (
        address maker,
        address token,
        TWAPConfig memory config,
        string memory cosmosRecipient,
        uint256 executedIntervals,
        uint256 totalExecuted,
        bool cancelled
    ) {
        TWAPOrder storage order = twapOrders[orderId];
        return (
            order.maker,
            order.token,
            order.config,
            order.cosmosRecipient,
            order.executedIntervals,
            order.totalExecuted,
            order.cancelled
        );
    }

    /**
     * @notice Check if interval is executed
     */
    function isIntervalExecuted(bytes32 orderId, uint256 intervalIndex) external view returns (bool) {
        return twapOrders[orderId].intervalSecrets[intervalIndex] != bytes32(0);
    }

    /**
     * @notice Add authorized resolver
     */
    function addResolver(address resolver) external onlyOwner {
        authorizedResolvers[resolver] = true;
    }

    /**
     * @notice Remove authorized resolver
     */
    function removeResolver(address resolver) external onlyOwner {
        authorizedResolvers[resolver] = false;
    }

    /**
     * @notice Emergency withdrawal (owner only)
     */
    function emergencyWithdraw(address token, uint256 amount) external onlyOwner {
        IERC20(token).safeTransfer(owner(), amount);
    }
}