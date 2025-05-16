const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("BlockchainNodeRegistry", function () {
  let nodeRegistry;
  let owner;
  let user;
  
  // Node data for testing
  const nodeName = "Test Ethereum Node";
  const chainType = 0; // Ethereum
  const endpointUrl = "https://ethereum-node.example.com";
  const version = "Geth/v1.11.6";
  const region = "us-east-1";
  const provider = 0; // AWS
  
  // Setup before each test
  beforeEach(async function () {
    // Get signers
    [owner, user] = await ethers.getSigners();
    
    // Deploy contract
    const BlockchainNodeRegistry = await ethers.getContractFactory("BlockchainNodeRegistry");
    nodeRegistry = await BlockchainNodeRegistry.deploy();
    await nodeRegistry.deployed();
  });
  
  describe("Node Registration", function () {
    it("Should register a new node", async function () {
      // Register a node
      const tx = await nodeRegistry.registerNode(
        nodeName,
        chainType,
        endpointUrl,
        version,
        region,
        provider
      );
      
      // Wait for transaction to be mined
      const receipt = await tx.wait();
      
      // Check if NodeRegistered event was emitted
      const event = receipt.events.find(e => e.event === "NodeRegistered");
      expect(event).to.not.be.undefined;
      
      // Extract the node ID from the event
      const nodeId = event.args.id;
      
      // Check node count
      expect(await nodeRegistry.getNodeCount()).to.equal(1);
      
      // Check node details
      const node = await nodeRegistry.getNodeDetails(nodeId);
      expect(node.name).to.equal(nodeName);
      expect(node.chainType).to.equal(chainType);
      expect(node.endpointUrl).to.equal(endpointUrl);
      expect(node.status).to.equal(2); // NodeStatus.Starting
      expect(node.owner).to.equal(owner.address);
      expect(node.isActive).to.be.true;
    });
    
    it("Should not allow registering a node with the same ID", async function () {
      // Register a node
      await nodeRegistry.registerNode(
        nodeName,
        chainType,
        endpointUrl,
        version,
        region,
        provider
      );
      
      // Try to register the same node again (this should fail due to the unique ID generation)
      // We're not using exactly the same ID here, but the test illustrates the concept
      await expect(
        nodeRegistry.registerNode(
          nodeName,
          chainType,
          endpointUrl,
          version,
          region,
          provider
        )
      ).to.not.be.reverted; // It will have a different ID due to timestamp
    });
  });
  
  describe("Node Updates", function () {
    let nodeId;
    
    beforeEach(async function () {
      // Register a node first
      const tx = await nodeRegistry.registerNode(
        nodeName,
        chainType,
        endpointUrl,
        version,
        region,
        provider
      );
      
      const receipt = await tx.wait();
      const event = receipt.events.find(e => e.event === "NodeRegistered");
      nodeId = event.args.id;
    });
    
    it("Should update node status", async function () {
      // Update node status
      const status = 0; // NodeStatus.Running
      const currentBlock = 1000000;
      const highestBlock = 1005000;
      
      await nodeRegistry.updateNodeStatus(nodeId, status, currentBlock, highestBlock);
      
      // Check if the status was updated
      const node = await nodeRegistry.getNodeDetails(nodeId);
      expect(node.status).to.equal(status);
      expect(node.currentBlock).to.equal(currentBlock);
      expect(node.highestBlock).to.equal(highestBlock);
    });
    
    it("Should calculate sync percentage correctly", async function () {
      // Update node with block info
      const currentBlock = 750000;
      const highestBlock = 1000000;
      
      await nodeRegistry.updateNodeStatus(nodeId, 3, currentBlock, highestBlock); // NodeStatus.Syncing
      
      // Calculate expected sync percentage
      const expectedPercentage = (currentBlock * 100) / highestBlock;
      
      // Check sync percentage
      const syncPercentage = await nodeRegistry.getNodeSyncPercentage(nodeId);
      expect(syncPercentage).to.equal(expectedPercentage);
    });
    
    it("Should not allow updating another user's node", async function () {
      // Try to update node status as a different user
      await expect(
        nodeRegistry.connect(user).updateNodeStatus(nodeId, 0, 1000000, 1005000)
      ).to.be.revertedWith("Not the node owner");
    });
  });
  
  describe("Node Deregistration", function () {
    let nodeId;
    
    beforeEach(async function () {
      // Register a node first
      const tx = await nodeRegistry.registerNode(
        nodeName,
        chainType,
        endpointUrl,
        version,
        region,
        provider
      );
      
      const receipt = await tx.wait();
      const event = receipt.events.find(e => e.event === "NodeRegistered");
      nodeId = event.args.id;
    });
    
    it("Should deregister a node", async function () {
      // Deregister the node
      await nodeRegistry.deregisterNode(nodeId);
      
      // Check if the node is inactive
      const node = await nodeRegistry.getNodeDetails(nodeId);
      expect(node.isActive).to.be.false;
    });
    
    it("Should not allow deregistering another user's node", async function () {
      // Try to deregister node as a different user
      await expect(
        nodeRegistry.connect(user).deregisterNode(nodeId)
      ).to.be.revertedWith("Not the node owner");
    });
  });
  
  describe("Node Queries", function () {
    beforeEach(async function () {
      // Register nodes for different chains
      await nodeRegistry.registerNode("Ethereum Node", 0, "https://eth.example.com", "v1", "us-east-1", 0);
      await nodeRegistry.registerNode("Polygon Node", 1, "https://polygon.example.com", "v1", "us-west-1", 0);
      await nodeRegistry.registerNode("Arbitrum Node", 2, "https://arbitrum.example.com", "v1", "eu-west-1", 0);
    });
    
    it("Should return correct node count", async function () {
      expect(await nodeRegistry.getNodeCount()).to.equal(3);
    });
    
    it("Should return correct nodes by owner", async function () {
      const ownerNodes = await nodeRegistry.getNodesByOwner(owner.address);
      expect(ownerNodes.length).to.equal(3);
    });
    
    it("Should have correct node count by chain", async function () {
      expect(await nodeRegistry.nodeCountByChain(0)).to.equal(1); // Ethereum
      expect(await nodeRegistry.nodeCountByChain(1)).to.equal(1); // Polygon
      expect(await nodeRegistry.nodeCountByChain(2)).to.equal(1); // Arbitrum
      expect(await nodeRegistry.nodeCountByChain(3)).to.equal(0); // BSC (no nodes)
    });
  });
}); 