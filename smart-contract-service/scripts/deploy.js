// Deployment script for Twist smart contracts
const { ethers } = require("hardhat");

async function main() {
  console.log("Starting deployment...");
  
  // Get the deployer account
  const [deployer] = await ethers.getSigners();
  console.log(`Deploying contracts with the account: ${deployer.address}`);
  
  const deployerBalance = await deployer.getBalance();
  console.log(`Account balance: ${ethers.utils.formatEther(deployerBalance)} ETH`);
  
  // Deploy the BlockchainNodeRegistry contract
  console.log("Deploying BlockchainNodeRegistry...");
  const BlockchainNodeRegistry = await ethers.getContractFactory("BlockchainNodeRegistry");
  const nodeRegistry = await BlockchainNodeRegistry.deploy();
  await nodeRegistry.deployed();
  console.log(`BlockchainNodeRegistry deployed to: ${nodeRegistry.address}`);
  
  // Deploy the TwistToken contract with 10 million initial supply
  console.log("Deploying TwistToken...");
  const initialSupply = ethers.utils.parseEther("10000000"); // 10 million tokens with 18 decimals
  const TwistToken = await ethers.getContractFactory("TwistToken");
  const token = await TwistToken.deploy(initialSupply, deployer.address);
  await token.deployed();
  console.log(`TwistToken deployed to: ${token.address}`);
  
  // Log deployment details for verification
  console.log("\nDeployment complete!");
  console.log("Contract Addresses:");
  console.log("====================");
  console.log(`BlockchainNodeRegistry: ${nodeRegistry.address}`);
  console.log(`TwistToken: ${token.address}`);
  
  // Verify contracts on Etherscan if not on a local network
  if (network.name !== "hardhat" && network.name !== "localhost") {
    console.log("\nWaiting for block confirmations...");
    // Wait for some confirmations
    await nodeRegistry.deployTransaction.wait(5);
    await token.deployTransaction.wait(5);
    
    console.log("\nVerifying contracts on Etherscan...");
    
    try {
      await hre.run("verify:verify", {
        address: nodeRegistry.address,
        constructorArguments: [],
      });
      
      await hre.run("verify:verify", {
        address: token.address,
        constructorArguments: [initialSupply, deployer.address],
      });
      
      console.log("Verification complete!");
    } catch (error) {
      console.error("Error verifying contracts:", error);
    }
  }
}

// Execute the deployment
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  }); 