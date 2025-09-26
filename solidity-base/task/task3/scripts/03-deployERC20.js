const { ethers } = require("hardhat");

// 部署 TestERC20
// npx hardhat run scripts/03-deployERC20.js --network sepolia
async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("🚀 部署人地址:", deployer.address);

  // 声明ERC20合约
  const ERC20 = await ethers.getContractFactory("TestERC20");
  const erc20 = await ERC20.deploy(); // ✅ 无参数

  // 等待部署
  await erc20.waitForDeployment();

  console.log("✅ TestERC20 部署成功，地址:", await erc20.getAddress());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
