const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("TwistToken", function () {
  let token;
  let owner;
  let user1;
  let user2;
  const initialSupply = ethers.utils.parseEther("10000000"); // 10 million tokens
  
  beforeEach(async function () {
    // Get signers
    [owner, user1, user2] = await ethers.getSigners();
    
    // Deploy token
    const TwistToken = await ethers.getContractFactory("TwistToken");
    token = await TwistToken.deploy(initialSupply, owner.address);
    await token.deployed();
  });
  
  describe("Deployment", function () {
    it("Should assign the total supply of tokens to the owner", async function () {
      const ownerBalance = await token.balanceOf(owner.address);
      expect(await token.totalSupply()).to.equal(ownerBalance);
      expect(ownerBalance).to.equal(initialSupply);
    });
    
    it("Should set the correct token name and symbol", async function () {
      expect(await token.name()).to.equal("Twist Token");
      expect(await token.symbol()).to.equal("TWIST");
    });
    
    it("Should set the max supply correctly", async function () {
      const maxSupply = await token.MAX_SUPPLY();
      expect(maxSupply).to.equal(ethers.utils.parseEther("100000000")); // 100 million
    });
    
    it("Should assign the correct roles", async function () {
      const adminRole = await token.DEFAULT_ADMIN_ROLE();
      const minterRole = await token.MINTER_ROLE();
      const pauserRole = await token.PAUSER_ROLE();
      
      expect(await token.hasRole(adminRole, owner.address)).to.be.true;
      expect(await token.hasRole(minterRole, owner.address)).to.be.true;
      expect(await token.hasRole(pauserRole, owner.address)).to.be.true;
    });
  });
  
  describe("Transfers", function () {
    it("Should transfer tokens between accounts", async function () {
      // Transfer 1000 tokens from owner to user1
      const transferAmount = ethers.utils.parseEther("1000");
      await token.transfer(user1.address, transferAmount);
      
      // Check balances
      const user1Balance = await token.balanceOf(user1.address);
      expect(user1Balance).to.equal(transferAmount);
      
      // Transfer 500 tokens from user1 to user2
      const secondTransferAmount = ethers.utils.parseEther("500");
      await token.connect(user1).transfer(user2.address, secondTransferAmount);
      
      // Check balances
      const user2Balance = await token.balanceOf(user2.address);
      expect(user2Balance).to.equal(secondTransferAmount);
      expect(await token.balanceOf(user1.address)).to.equal(transferAmount.sub(secondTransferAmount));
    });
    
    it("Should fail if sender doesn't have enough tokens", async function () {
      // Try to send more tokens than the user has
      const initialUser1Balance = await token.balanceOf(user1.address);
      expect(initialUser1Balance).to.equal(0);
      
      await expect(
        token.connect(user1).transfer(user2.address, ethers.utils.parseEther("1"))
      ).to.be.reverted;
    });
  });
  
  describe("Vesting", function () {
    const vestingAmount = ethers.utils.parseEther("10000"); // 10k tokens
    
    beforeEach(async function () {
      // Set up vesting for user1
      await token.addVesting(user1.address, vestingAmount);
    });
    
    it("Should add vesting correctly", async function () {
      const vestedAmount = await token.vestedAmount(user1.address);
      expect(vestedAmount).to.equal(vestingAmount);
      
      const totalVested = await token.totalVested();
      expect(totalVested).to.equal(vestingAmount);
    });
    
    it("Should calculate claimable amount correctly (time-based)", async function () {
      // At the start, a small amount should be claimable due to time elapsed
      const claimableAmount = await token.getClaimableAmount(user1.address);
      
      // Should be greater than 0 but less than the full amount
      // (This test is approximate due to block time variations)
      expect(claimableAmount).to.be.gt(0);
      expect(claimableAmount).to.be.lt(vestingAmount);
    });
    
    it("Should allow claiming vested tokens", async function () {
      // Claim vested tokens
      await token.connect(user1).claimVestedTokens();
      
      // Check balance
      const balance = await token.balanceOf(user1.address);
      expect(balance).to.be.gt(0);
      
      // Check claimed amount
      const claimed = await token.claimedAmount(user1.address);
      expect(claimed).to.equal(balance);
    });
    
    it("Should not allow claiming more than vested amount", async function () {
      // First claim
      await token.connect(user1).claimVestedTokens();
      
      // Time hasn't elapsed much, so a second claim should fail
      await expect(token.connect(user1).claimVestedTokens()).to.be.revertedWith("No tokens to claim");
    });
  });
  
  describe("Minting and Burning", function () {
    it("Should allow minting by minter role", async function () {
      const mintAmount = ethers.utils.parseEther("1000");
      const initialSupply = await token.totalSupply();
      
      await token.mint(user1.address, mintAmount);
      
      expect(await token.totalSupply()).to.equal(initialSupply.add(mintAmount));
      expect(await token.balanceOf(user1.address)).to.equal(mintAmount);
    });
    
    it("Should not allow minting beyond max supply", async function () {
      const maxSupply = await token.MAX_SUPPLY();
      const currentSupply = await token.totalSupply();
      const mintAmount = maxSupply.sub(currentSupply).add(1); // Exceed max supply by 1
      
      await expect(token.mint(user1.address, mintAmount)).to.be.revertedWith("Exceeds max supply");
    });
    
    it("Should allow burning tokens", async function () {
      // Transfer tokens to user1
      const transferAmount = ethers.utils.parseEther("1000");
      await token.transfer(user1.address, transferAmount);
      
      // Burn half of the tokens
      const burnAmount = ethers.utils.parseEther("500");
      await token.connect(user1).burn(burnAmount);
      
      // Check balance
      expect(await token.balanceOf(user1.address)).to.equal(transferAmount.sub(burnAmount));
      
      // Check total supply
      expect(await token.totalSupply()).to.equal(initialSupply.sub(burnAmount));
    });
  });
  
  describe("Pausing", function () {
    it("Should allow pausing and unpausing by pauser role", async function () {
      // Pause
      await token.pause();
      expect(await token.paused()).to.be.true;
      
      // Try transfer while paused
      await expect(
        token.transfer(user1.address, ethers.utils.parseEther("100"))
      ).to.be.reverted;
      
      // Unpause
      await token.unpause();
      expect(await token.paused()).to.be.false;
      
      // Transfer after unpausing should work
      await token.transfer(user1.address, ethers.utils.parseEther("100"));
      expect(await token.balanceOf(user1.address)).to.equal(ethers.utils.parseEther("100"));
    });
    
    it("Should not allow non-pausers to pause", async function () {
      await expect(token.connect(user1).pause()).to.be.reverted;
    });
  });
}); 