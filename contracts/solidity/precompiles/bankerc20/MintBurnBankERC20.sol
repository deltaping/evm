// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.18;

import {BankERC20} from "./BankERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title MintBurnBankERC20
/// @dev Concrete ERC-20 token backed by the Cosmos bank module with
/// owner-controlled minting. Anyone can burn their own tokens.
///
/// Deployment flow:
///   1. Deploy this contract with name, symbol, decimals, and initial owner.
///   2. Register a TokenPair via governance (registerERC20WithDenom) binding
///      this contract's address to the bank denom. The proposal queries
///      name()/symbol()/decimals() to populate bank metadata.
///   3. The owner can then call mint() to create new tokens.
contract MintBurnBankERC20 is Ownable, BankERC20 {
    constructor(
        string memory name_,
        string memory symbol_,
        uint8 decimals_,
        address initialOwner
    ) BankERC20(name_, symbol_, decimals_) Ownable(initialOwner) {}

    /// @dev Mints new tokens to the specified address. Only the owner can call this.
    function mint(address to, uint256 amount) public virtual onlyOwner {
        _mint(to, amount);
    }

    /// @dev Burns tokens from the caller's account.
    function burn(uint256 value) public virtual {
        _burn(_msgSender(), value);
    }

    /// @dev Burns tokens from the specified account, deducting from the caller's allowance.
    function burnFrom(address account, uint256 value) public virtual {
        _spendAllowance(account, _msgSender(), value);
        _burn(account, value);
    }
}
