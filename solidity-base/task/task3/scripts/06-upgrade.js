const { ethers, upgrades } = require("hardhat");
require("dotenv").config();

const AUCTION_PROXY = process.env.AUCTION_ADDRESS;

//升级合约
// npx hardhat run scripts/06-upgrade.js --network sepolia
async function main() {
  console.log("🔄 开始升级合约...");

  const AuctionV2 = await ethers.getContractFactory("NftAuctionV2");
  const upgraded = await upgrades.upgradeProxy(AUCTION_PROXY, AuctionV2);

  await upgraded.waitForDeployment();
  console.log("✅ 升级成功，地址 (proxy 不变):", await upgraded.getAddress());
  
  // console.log("🆕 V2 版本:", await upgraded.version());

  // 调用新方法验证升级成功
  const msg = await upgraded.testHello();
  console.log("🆕 V2 返回:", msg);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
