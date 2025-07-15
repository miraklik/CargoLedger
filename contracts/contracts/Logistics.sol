// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

//TODO: Добавить Chainlink Oracle

contract Logistics is AccessControl, Pausable{
    enum CargoStatus {
        Created,
        InTransit, 
        Delivered,
        Cancelled
    }

    struct Cargo {
        uint256 id;
        address sender;
        address carrier;
        address receiver;
        string descriptionIpfsHash;
        CargoStatus status;
        uint256 timestamp;
    }
    
    error InvalidIpfsHash();
    error CargoNotFound();
    error OnlyCarrierCanUpdate();
    error AlreadyInTransit();
    error InvalidAddress();
    error NotSender();

    uint256 public nextCargoId;
    mapping(uint256 => Cargo) public cargos;

    event CargoCreated(
        uint256 indexed id,
        address indexed sender,
        address indexed receiver,
        string descriptionIpfsHash
    );

    event CargoUpdated(
        uint256 indexed id,
        CargoStatus newStatus
    );

    modifier validStatusTransition(CargoStatus current, CargoStatus newStatus) {
        require(
            (current == CargoStatus.Created && (newStatus == CargoStatus.InTransit || newStatus == CargoStatus.Cancelled)) ||
            (current == CargoStatus.InTransit && (newStatus == CargoStatus.Delivered || newStatus == CargoStatus.Cancelled)),
            "Invalid status transition"
        );
        _;
    }

    bytes32 public constant SENDER_ROLE = keccak256(bytes("SENDER_ROLE"));
    bytes32 public constant RECEIVER_ROLE = keccak256(bytes("RECEIVER_ROLE"));
    bytes32 public constant CARRIER_ROLE = keccak256(bytes("CARRIER_ROLE"));

    constructor(address carrier, address receiver) {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(SENDER_ROLE, msg.sender);
        _grantRole(CARRIER_ROLE, carrier);
        _grantRole(RECEIVER_ROLE, receiver);
    }

    /**
     * @dev Validates IPFS hash format
     * @param ipfsHash The IPFS hash to validate
     * @return bool True if valid format
     */
     function _isValidIpfsHash(string memory ipfsHash) private pure returns (bool) {
        bytes memory hashBytes = bytes(ipfsHash);
    
        if (hashBytes.length < 32 || hashBytes.length > 64) {
            return false;
        }
        
        return true;
    }

    /**
     * @dev a function for creating cargo
     * @param receiver cargo recipient
     * @param carrier Cargo carrier
     * @param descriptionIpfsHash description hash in ipfs
     */
    function createCargo(
        address receiver, 
        address carrier, 
        string memory descriptionIpfsHash
        ) external whenNotPaused onlyRole(SENDER_ROLE){
        require(receiver != address(0) && carrier != address(0), InvalidAddress());
        if (!_isValidIpfsHash(descriptionIpfsHash)) {
            revert InvalidIpfsHash();
        }

        cargos[nextCargoId] = Cargo({
            id: nextCargoId,
            sender: msg.sender,
            carrier: carrier,
            receiver: receiver,
            descriptionIpfsHash: descriptionIpfsHash,
            status: CargoStatus.Created,
            timestamp: block.timestamp
        });

        emit CargoCreated(nextCargoId, msg.sender, receiver, descriptionIpfsHash);
        nextCargoId++;
    }

    /**
     * @dev Updates cargo status
     * @param id Cargo ID
     * @param newStatus New cargo status
     */
    function updateCargo(
        uint256 id, 
        CargoStatus newStatus
        ) external whenNotPaused onlyRole(CARRIER_ROLE) validStatusTransition(cargos[id].status, newStatus){
        Cargo storage s = cargos[id];
        require(s.carrier != address(0), CargoNotFound());
        require(msg.sender == s.carrier, OnlyCarrierCanUpdate());

        s.status = newStatus;
        s.timestamp = block.timestamp;
        emit CargoUpdated(id, newStatus);
    }

    /**
     * @dev Cancel cargo
     * @param id Cargo ID 
     */
    function cancelCargo(uint256 id) external onlyRole(SENDER_ROLE) {
        Cargo storage s = cargos[id];
        require(s.sender == msg.sender, NotSender());
        require(s.status == CargoStatus.Created, AlreadyInTransit());
        s.status = CargoStatus.Cancelled;
        s.timestamp = block.timestamp;
        emit CargoUpdated(id, CargoStatus.Cancelled);
    }

    /**
     * @dev Confirm cargo delivery
     * @param id Cargo ID
     */
    function confirmDelivery(uint256 id) external onlyRole(RECEIVER_ROLE) {
        Cargo storage s = cargos[id];
        require(s.receiver == msg.sender, InvalidAddress());
        require(s.status == CargoStatus.InTransit, AlreadyInTransit());
        s.status = CargoStatus.Delivered;
        s.timestamp = block.timestamp;
        emit CargoUpdated(id, CargoStatus.Delivered);
    }

    /**
     * @dev Get specific cargo
     * @param id Cargo ID
     * @return Cargo
     */
    function getCargo(uint256 id) public view onlyRole(RECEIVER_ROLE) returns (Cargo memory) {
        require(cargos[id].carrier != address(0), CargoNotFound());
        return cargos[id];
    }

    /**
     * @dev Get cargo status
     * @param id Cargo ID
     * @return CargoStatus
     */
    function getCargoStatus(uint256 id) external view returns (CargoStatus) {
        require(cargos[id].carrier != address(0), CargoNotFound());
        return cargos[id].status;
    }

    /**
     * @dev Get cargo participants
     * @param id Cargo ID
     * @return sender 
     * @return carrier 
     * @return receiver 
     */
    function getParticipants(uint256 id) external view returns (address sender, address carrier, address receiver) {
        require(cargos[id].carrier != address(0), CargoNotFound());
        Cargo memory cargo = cargos[id];
        return (cargo.sender, cargo.carrier, cargo.receiver);
    }

    /**
     * @dev Get IPFS hash for specific cargo
     * @param id Cargo ID
     * @return string IPFS hash
     */
    function getCargoIpfsHash(uint256 id) public view returns (string memory) {
        Cargo memory cargo = cargos[id];
        if (cargo.carrier == address(0)) {
            revert CargoNotFound();
        }
        return cargo.descriptionIpfsHash;
    }

    /**
     * @dev Get all cargos for a specific address (as sender)
     * @param user Address to check
     * @return cargoIds Array of cargo IDs
     */
    function getCargosBySender(address user) external view returns (uint256[] memory) {
        uint256[] memory tempCargos = new uint256[](nextCargoId);
        uint256 count = 0;
        
        for (uint256 i = 0; i < nextCargoId; i++) {
            if (cargos[i].sender == user) {
                tempCargos[count] = i;
                count++;
            }
        }
        
        uint256[] memory result = new uint256[](count);
        for (uint256 i = 0; i < count; i++) {
            result[i] = tempCargos[i];
        }
        
        return result;
    }

    /**
     * @dev Get all cargos for a specific address (as carrier)
     * @param user Address to check
     * @return cargoIds Array of cargo IDs
     */
    function getCargosByCarrier(address user) external view returns (uint256[] memory) {
        uint256[] memory tempCargos = new uint256[](nextCargoId);
        uint256 count = 0;
        
        for (uint256 i = 0; i < nextCargoId; i++) {
            if (cargos[i].carrier == user) {
                tempCargos[count] = i;
                count++;
            }
        }
        
        uint256[] memory result = new uint256[](count);
        for (uint256 i = 0; i < count; i++) {
            result[i] = tempCargos[i];
        }
        
        return result;
    }

    /**
     * @dev Get all cargos for a specific address (as receiver)
     * @param user Address to check
     * @return cargoIds Array of cargo IDs
     */
    function getCargosByReceiver(address user) external view returns (uint256[] memory) {
        uint256[] memory tempCargos = new uint256[](nextCargoId);
        uint256 count = 0;
        
        for (uint256 i = 0; i < nextCargoId; i++) {
            if (cargos[i].receiver == user) {
                tempCargos[count] = i;
                count++;
            }
        }
        
        uint256[] memory result = new uint256[](count);
        for (uint256 i = 0; i < count; i++) {
            result[i] = tempCargos[i];
        }
        
        return result;
    }

    /**
     * @dev Pauses the contract
     */
    function pause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    /**
     * @dev Unpause the contract
     */
    function unpause() public whenPaused onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }
}