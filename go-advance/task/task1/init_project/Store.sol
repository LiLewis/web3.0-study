pragma solidity ^0.8.26;

//1.nvm use v22.19.0
//2.npm install -g solc
//3.$env:PATH += ";D:\blockchain\nodejs\node_global\"
//4.solcjs --version
//5.solcjs --bin Store.sol
//6.go install github.com/ethereum/go-ethereum/cmd/abigen@latest
//7.abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=store.go
//8.go get github.com/ethereum/go-ethereum, 还爆红快速修复命令：go mod tidy -e

contract Store {
  event ItemSet(bytes32 key, bytes32 value);

  string public version;
  mapping (bytes32 => bytes32) public items;

  constructor(string memory _version) {
    version = _version;
  }

  function setItem(bytes32 key, bytes32 value) external {
    items[key] = value;
    emit ItemSet(key, value);
  }
}