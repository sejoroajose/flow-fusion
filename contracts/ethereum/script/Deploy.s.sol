// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/FlowFusionEscrowFactory.sol";

contract DeployScript is Script {
    function run() external {
        vm.startBroadcast();
        
        address escrowFactory = 0x1111111111111111111111111111111111111111;
        address initialOwner = 0x9566141705e4c8D354468F0d221F7a040013E3BF;
        
        FlowFusionEscrowFactory factory = new FlowFusionEscrowFactory(
            IEscrowFactory(escrowFactory),
            initialOwner
        );
        
        console.log("FlowFusionEscrowFactory deployed at:", address(factory));
        
        vm.stopBroadcast();
    }
}