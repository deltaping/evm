// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.18;

/// @dev The IBankERC20 contract's address.
address constant IBANKERC20_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000807;

/// @dev The IBankERC20 contract's instance.
IBankERC20 constant IBANKERC20_CONTRACT = IBankERC20(IBANKERC20_PRECOMPILE_ADDRESS);

/**
 * @author Deltaping Team
 * @title BankERC20 Precompile Interface
 * @dev Precompile for minting, burning, and transferring bank-backed ERC-20
 * tokens. Only callable by contracts (not EOAs). The caller contract's address
 * is used to resolve the target denom via a registered TokenPair.
 */
interface IBankERC20 {
    /// @dev mint mints bank coins for the denom bound to the caller contract's
    /// TokenPair and sends them to the given address.
    /// Only callable by contracts (not EOAs). Requires a registered TokenPair.
    /// @param to the recipient address.
    /// @param amount the amount of coins to mint.
    /// @return success true if the operation succeeded.
    function mint(address to, uint256 amount) external returns (bool success);

    /// @dev burn burns bank coins for the denom bound to the caller contract's
    /// TokenPair, deducting them from the given address.
    /// Only callable by contracts (not EOAs). Requires a registered TokenPair.
    /// @param from the address whose coins will be burned.
    /// @param amount the amount of coins to burn.
    /// @return success true if the operation succeeded.
    function burn(address from, uint256 amount) external returns (bool success);

    /// @dev transfer transfers bank coins for the denom bound to the caller
    /// contract's TokenPair from one address to another.
    /// Only callable by contracts (not EOAs). Requires a registered TokenPair.
    /// @param from the sender address.
    /// @param to the recipient address.
    /// @param amount the amount of coins to transfer.
    /// @return success true if the operation succeeded.
    function transfer(
        address from,
        address to,
        uint256 amount
    ) external returns (bool success);

    /// @dev balanceOf returns the bank module balance of a specific token
    /// for a given account.
    /// @param token the ERC20 contract address (used to resolve the denom via TokenPair).
    /// @param account the address to query.
    /// @return balance the balance amount.
    function balanceOf(
        address token,
        address account
    ) external view returns (uint256 balance);

    /// @dev supplyOf returns the total supply of a specific token from the
    /// bank module.
    /// @param token the ERC20 contract address (used to resolve the denom via TokenPair).
    /// @return totalSupply the total supply.
    function supplyOf(
        address token
    ) external view returns (uint256 totalSupply);

    /// @dev metadata returns the bank module metadata for a specific token.
    /// Returns found=false (with zero-value fields) when no TokenPair is
    /// registered, so callers can distinguish "not registered" from
    /// "registered with empty/zero values".
    /// @param token the ERC20 contract address (used to resolve the denom via TokenPair).
    /// @return found true if a registered TokenPair with metadata exists.
    /// @return name the token name.
    /// @return symbol the token symbol.
    /// @return decimals the token decimals.
    function metadata(
        address token
    ) external view returns (bool found, string memory name, string memory symbol, uint8 decimals);
}
