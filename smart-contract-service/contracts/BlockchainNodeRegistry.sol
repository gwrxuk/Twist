// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title BlockchainNodeRegistry
 * @dev A smart contract for registering and managing blockchain nodes
 */
contract BlockchainNodeRegistry {
    // Enum for node status
    enum NodeStatus { Running, Stopped, Starting, Syncing, Error, Maintenance }
    
    // Enum for chain type
    enum ChainType { Ethereum, Polygon, Arbitrum, BSC, Custom }
    
    // Enum for cloud provider
    enum CloudProvider { AWS, GCP, Azure, DigitalOcean, OnPremise }
    
    // Struct to store node data
    struct Node {
        bytes32 id;
        string name;
        ChainType chainType;
        string endpointUrl;
        NodeStatus status;
        string version;
        uint256 currentBlock;
        uint256 highestBlock;
        uint256 registeredAt;
        uint256 updatedAt;
        string region;
        CloudProvider provider;
        address owner;
        bool isActive;
    }
    
    // Mapping of node ID to Node struct
    mapping(bytes32 => Node) public nodes;
    
    // Array of all node IDs
    bytes32[] public nodeIds;
    
    // Mapping of owner address to their node IDs
    mapping(address => bytes32[]) public ownerNodes;
    
    // Node count by chain type
    mapping(ChainType => uint256) public nodeCountByChain;
    
    // Events
    event NodeRegistered(bytes32 indexed id, string name, ChainType chainType, address indexed owner);
    event NodeUpdated(bytes32 indexed id, NodeStatus status, uint256 currentBlock, uint256 highestBlock);
    event NodeStatusChanged(bytes32 indexed id, NodeStatus oldStatus, NodeStatus newStatus);
    event NodeDeregistered(bytes32 indexed id, address indexed owner);
    
    // Modifiers
    modifier onlyNodeOwner(bytes32 _nodeId) {
        require(nodes[_nodeId].owner == msg.sender, "Not the node owner");
        _;
    }
    
    modifier nodeExists(bytes32 _nodeId) {
        require(nodes[_nodeId].id == _nodeId, "Node does not exist");
        _;
    }
    
    /**
     * @dev Register a new blockchain node
     * @param _name Name of the node
     * @param _chainType Type of blockchain
     * @param _endpointUrl URL endpoint of the node
     * @param _version Node software version
     * @param _region Geographic region
     * @param _provider Cloud provider
     * @return id The ID of the registered node
     */
    function registerNode(
        string memory _name,
        ChainType _chainType,
        string memory _endpointUrl,
        string memory _version,
        string memory _region,
        CloudProvider _provider
    ) external returns (bytes32) {
        // Generate a unique ID for the node
        bytes32 nodeId = keccak256(abi.encodePacked(msg.sender, _name, block.timestamp));
        
        // Ensure node doesn't already exist
        require(nodes[nodeId].id != nodeId, "Node already exists");
        
        // Create new node
        Node memory newNode = Node({
            id: nodeId,
            name: _name,
            chainType: _chainType,
            endpointUrl: _endpointUrl,
            status: NodeStatus.Starting,
            version: _version,
            currentBlock: 0,
            highestBlock: 0,
            registeredAt: block.timestamp,
            updatedAt: block.timestamp,
            region: _region,
            provider: _provider,
            owner: msg.sender,
            isActive: true
        });
        
        // Store node
        nodes[nodeId] = newNode;
        nodeIds.push(nodeId);
        ownerNodes[msg.sender].push(nodeId);
        nodeCountByChain[_chainType]++;
        
        // Emit event
        emit NodeRegistered(nodeId, _name, _chainType, msg.sender);
        
        return nodeId;
    }
    
    /**
     * @dev Update node status and block info
     * @param _nodeId ID of the node
     * @param _status New status
     * @param _currentBlock Current block number
     * @param _highestBlock Highest block number
     */
    function updateNodeStatus(
        bytes32 _nodeId,
        NodeStatus _status,
        uint256 _currentBlock,
        uint256 _highestBlock
    ) external onlyNodeOwner(_nodeId) nodeExists(_nodeId) {
        Node storage node = nodes[_nodeId];
        
        // Emit event if status changed
        if (node.status != _status) {
            emit NodeStatusChanged(_nodeId, node.status, _status);
            node.status = _status;
        }
        
        // Update block information
        node.currentBlock = _currentBlock;
        node.highestBlock = _highestBlock;
        node.updatedAt = block.timestamp;
        
        // Emit update event
        emit NodeUpdated(_nodeId, _status, _currentBlock, _highestBlock);
    }
    
    /**
     * @dev Deregister a node
     * @param _nodeId ID of the node to deregister
     */
    function deregisterNode(bytes32 _nodeId) external onlyNodeOwner(_nodeId) nodeExists(_nodeId) {
        Node storage node = nodes[_nodeId];
        
        // Mark node as inactive
        node.isActive = false;
        nodeCountByChain[node.chainType]--;
        
        // Emit event
        emit NodeDeregistered(_nodeId, msg.sender);
    }
    
    /**
     * @dev Get node count
     * @return Total number of registered nodes
     */
    function getNodeCount() external view returns (uint256) {
        return nodeIds.length;
    }
    
    /**
     * @dev Get node IDs owned by a specific owner
     * @param _owner Address of the owner
     * @return Array of node IDs
     */
    function getNodesByOwner(address _owner) external view returns (bytes32[] memory) {
        return ownerNodes[_owner];
    }
    
    /**
     * @dev Get node details
     * @param _nodeId ID of the node
     * @return Node details
     */
    function getNodeDetails(bytes32 _nodeId) external view nodeExists(_nodeId) returns (
        string memory name,
        ChainType chainType,
        string memory endpointUrl,
        NodeStatus status,
        string memory version,
        uint256 currentBlock,
        uint256 highestBlock,
        uint256 registeredAt,
        uint256 updatedAt,
        string memory region,
        CloudProvider provider,
        address owner,
        bool isActive
    ) {
        Node memory node = nodes[_nodeId];
        
        return (
            node.name,
            node.chainType,
            node.endpointUrl,
            node.status,
            node.version,
            node.currentBlock,
            node.highestBlock,
            node.registeredAt,
            node.updatedAt,
            node.region,
            node.provider,
            node.owner,
            node.isActive
        );
    }
    
    /**
     * @dev Calculate node sync percentage
     * @param _nodeId ID of the node
     * @return Sync percentage (0-100)
     */
    function getNodeSyncPercentage(bytes32 _nodeId) external view nodeExists(_nodeId) returns (uint256) {
        Node memory node = nodes[_nodeId];
        
        if (node.highestBlock == 0 || node.currentBlock > node.highestBlock) {
            return 100;
        }
        
        return (node.currentBlock * 100) / node.highestBlock;
    }
} 