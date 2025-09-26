require("@nomicfoundation/hardhat-toolbox");
require('hardhat-deploy');
require('@openzeppelin/hardhat-upgrades');
require("dotenv").config(); // 用 .env 存 RPC 和私钥

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    version: "0.8.28",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  namedAccounts: {
    deployer: {
      default: 0, // 本地默认第一个账户
    },
    user1: {
      default: 1,
    },
    user2: {
      default: 2,
    },
  },
  networks: {
    hardhat: {
      chainId: 31337,
    },
    sepolia: {
      url: process.env.SEPOLIA_RPC || "",
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
      chainId: 11155111,
    },
  },
  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY || "",
  },
};
