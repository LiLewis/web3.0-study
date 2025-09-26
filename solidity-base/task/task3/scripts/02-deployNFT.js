const { ethers } = require("hardhat");

// 部署 TestERC721
// npx hardhat run scripts/02-deployNFT.js --network sepolia
async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("🚀 部署人地址:", deployer.address);

  // 声明ERC721合约
  const NFT = await ethers.getContractFactory("TestERC721");
  const nft = await NFT.deploy(); // ✅ 无参数

  // 等待部署
  await nft.waitForDeployment();

  console.log("✅ TestERC721 部署成功，地址:", await nft.getAddress());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
