// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

/**
 * @title TwistToken
 * @dev ERC20 token for the Twist blockchain infrastructure platform
 * This token provides governance and utility within the Twist ecosystem
 */
contract TwistToken is ERC20, ERC20Burnable, Pausable, AccessControl {
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    
    // Maximum supply of 100 million tokens
    uint256 public constant MAX_SUPPLY = 100_000_000 * 10**18;
    
    // Emission parameters for vesting
    uint256 public vestingStartTime;
    uint256 public vestingEndTime;
    uint256 public totalVested;
    
    // Mapping for vesting
    mapping(address => uint256) public vestedAmount;
    mapping(address => uint256) public claimedAmount;
    
    // Events
    event TokensVested(address indexed beneficiary, uint256 amount);
    event TokensClaimed(address indexed beneficiary, uint256 amount);
    
    /**
     * @dev Initialize the token contract
     * @param _initialSupply Initial supply to mint
     * @param _admin Admin address
     */
    constructor(uint256 _initialSupply, address _admin) ERC20("Twist Token", "TWIST") {
        require(_initialSupply <= MAX_SUPPLY, "Initial supply exceeds maximum");
        
        _grantRole(DEFAULT_ADMIN_ROLE, _admin);
        _grantRole(PAUSER_ROLE, _admin);
        _grantRole(MINTER_ROLE, _admin);
        
        // Mint initial supply
        _mint(_admin, _initialSupply);
        
        // Set vesting period for 2 years
        vestingStartTime = block.timestamp;
        vestingEndTime = block.timestamp + 730 days; // ~2 years
    }
    
    /**
     * @dev Add vesting schedule for a beneficiary
     * @param _beneficiary Address of beneficiary
     * @param _amount Amount to vest
     */
    function addVesting(address _beneficiary, uint256 _amount) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(_beneficiary != address(0), "Invalid address");
        require(_amount > 0, "Amount must be greater than 0");
        require(totalSupply() + _amount <= MAX_SUPPLY, "Exceeds max supply");
        
        vestedAmount[_beneficiary] += _amount;
        totalVested += _amount;
        
        emit TokensVested(_beneficiary, _amount);
    }
    
    /**
     * @dev Claim vested tokens
     */
    function claimVestedTokens() external whenNotPaused {
        address beneficiary = msg.sender;
        uint256 vestedTotal = vestedAmount[beneficiary];
        uint256 claimed = claimedAmount[beneficiary];
        
        require(vestedTotal > 0, "No vested tokens");
        require(claimed < vestedTotal, "All tokens claimed");
        
        uint256 claimable;
        if (block.timestamp >= vestingEndTime) {
            // Vesting period complete - claim all remaining tokens
            claimable = vestedTotal - claimed;
        } else {
            // Calculate vested amount based on time elapsed
            uint256 elapsed = block.timestamp - vestingStartTime;
            uint256 totalPeriod = vestingEndTime - vestingStartTime;
            uint256 vestedSoFar = (vestedTotal * elapsed) / totalPeriod;
            
            claimable = vestedSoFar - claimed;
        }
        
        require(claimable > 0, "No tokens to claim");
        
        claimedAmount[beneficiary] += claimable;
        _mint(beneficiary, claimable);
        
        emit TokensClaimed(beneficiary, claimable);
    }
    
    /**
     * @dev Mint new tokens (only by MINTER_ROLE)
     * @param to Recipient address
     * @param amount Amount to mint
     */
    function mint(address to, uint256 amount) 
        external 
        onlyRole(MINTER_ROLE) 
    {
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        _mint(to, amount);
    }
    
    /**
     * @dev Pause token transfers
     */
    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
    }
    
    /**
     * @dev Unpause token transfers
     */
    function unpause() external onlyRole(PAUSER_ROLE) {
        _unpause();
    }
    
    /**
     * @dev Override to add pause functionality
     */
    function _beforeTokenTransfer(address from, address to, uint256 amount)
        internal
        whenNotPaused
        override
    {
        super._beforeTokenTransfer(from, to, amount);
    }
    
    /**
     * @dev Get claimable amount for a beneficiary
     * @param _beneficiary Beneficiary address
     * @return Amount claimable
     */
    function getClaimableAmount(address _beneficiary) external view returns (uint256) {
        uint256 vestedTotal = vestedAmount[_beneficiary];
        uint256 claimed = claimedAmount[_beneficiary];
        
        if (vestedTotal == 0 || claimed >= vestedTotal) {
            return 0;
        }
        
        if (block.timestamp >= vestingEndTime) {
            return vestedTotal - claimed;
        }
        
        uint256 elapsed = block.timestamp - vestingStartTime;
        uint256 totalPeriod = vestingEndTime - vestingStartTime;
        uint256 vestedSoFar = (vestedTotal * elapsed) / totalPeriod;
        
        return vestedSoFar > claimed ? vestedSoFar - claimed : 0;
    }
} 