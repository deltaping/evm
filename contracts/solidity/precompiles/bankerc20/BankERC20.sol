// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.18;

import {IBankERC20, IBANKERC20_PRECOMPILE_ADDRESS} from "./IBankERC20.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/// @title BankERC20
/// @dev Abstract ERC-20 contract backed by the Cosmos bank module via the
/// BankERC20 precompile at 0x807. All balance state (mint, burn, transfer) is
/// delegated to the bank module so that Cosmos-side and EVM-side balances
/// are always in sync.
///
/// Constructor params serve as fallbacks: before the TokenPair is registered,
/// name()/symbol()/decimals() return the constructor values. After
/// registration, metadata is read from the bank module. The precompile's
/// metadata() returns found=false (never reverts) when no TokenPair exists,
/// so direct calls are safe without try/catch.
abstract contract BankERC20 is ERC20 {
    IBankERC20 constant bank = IBankERC20(IBANKERC20_PRECOMPILE_ADDRESS);

    uint8 private immutable _decimals;

    constructor(
        string memory name_,
        string memory symbol_,
        uint8 decimals_
    ) ERC20(name_, symbol_) {
        _decimals = decimals_;
    }

    /// @dev Returns the token name. Reads from bank metadata if the
    /// TokenPair is registered, otherwise falls back to the constructor value.
    function name() public view virtual override returns (string memory) {
        (bool found, string memory bankName,,) = bank.metadata(address(this));
        if (found) return bankName;
        return super.name();
    }

    /// @dev Returns the token symbol. Reads from bank metadata if the
    /// TokenPair is registered, otherwise falls back to the constructor value.
    function symbol() public view virtual override returns (string memory) {
        (bool found,, string memory bankSymbol,) = bank.metadata(address(this));
        if (found) return bankSymbol;
        return super.symbol();
    }

    /// @dev Returns the token decimals. Reads from bank metadata if the
    /// TokenPair is registered, otherwise falls back to the constructor value.
    function decimals() public view virtual override returns (uint8) {
        (bool found,,, uint8 bankDecimals) = bank.metadata(address(this));
        if (found) return bankDecimals;
        return _decimals;
    }

    /// @dev Returns the total supply from the bank module.
    function totalSupply() public view virtual override returns (uint256) {
        return bank.supplyOf(address(this));
    }

    /// @dev Returns the balance of an account from the bank module.
    function balanceOf(
        address account
    ) public view virtual override returns (uint256) {
        return bank.balanceOf(address(this), account);
    }

    /// @dev Overrides the internal _update hook to route mint, burn, and
    /// transfer operations through the BankERC20 precompile.
    function _update(
        address from,
        address to,
        uint256 value
    ) internal virtual override {
        if (from == address(0)) {
            require(bank.mint(to, value), "BankERC20: mint failed");
        } else if (to == address(0)) {
            require(bank.burn(from, value), "BankERC20: burn failed");
        } else {
            require(
                bank.transfer(from, to, value),
                "BankERC20: transfer failed"
            );
        }

        emit Transfer(from, to, value);
    }
}
