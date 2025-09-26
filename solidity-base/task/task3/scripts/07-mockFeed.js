const { ethers } = require("hardhat");
require("dotenv").config();

// npx hardhat run scripts/07-mockFeed.js --network localhost
// 或者 sepolia 也可以，不过 sepolia 上链要花费 gas

async function main() {
  console.log("🚀 部署 MockV3Aggregator...");

  // 初始化：8 位小数，价格 = 2000 * 10^8 (即 $2000)
  const decimals = 8;
  const initialPrice = ethers.parseUnits("2000", decimals); 

  const MockV3Aggregator = await ethers.getContractFactory("MockV3Aggregator");
  const mock = await MockV3Aggregator.deploy(decimals, initialPrice);
  await mock.waitForDeployment();

  console.log("✅ MockV3Aggregator 部署成功:", await mock.getAddress());

  // 读取当前价格
  let roundData = await mock.latestRoundData();
  console.log("📊 初始价格:", roundData[1].toString());

  // 更新价格
  const tx = await mock.updateAnswer(ethers.parseUnits("2500", decimals)); // $2500
  await tx.wait();

  // 再读一次
  roundData = await mock.latestRoundData();
  console.log("📊 更新后的价格:", roundData[1].toString());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
