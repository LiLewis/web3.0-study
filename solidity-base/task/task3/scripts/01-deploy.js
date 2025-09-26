const { ethers, upgrades } = require("hardhat");

//部署合约
//npx hardhat run scripts/01-deploy.js --network sepolia
async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("🚀 部署人地址:", deployer.address);

  // 部署NftAuction合约到ETH
  const NftAuction = await ethers.getContractFactory("NftAuction");

  // 部署代理合约，并且申明uups
  const nftAuction = await upgrades.deployProxy(NftAuction, [deployer.address], { kind: "uups" });

  // 等待代理合约部署完成
  await nftAuction.waitForDeployment();

  // 获取代理合约地址AUCTION_ADDRESS
  const proxyAddress = await nftAuction.getAddress();

  console.log("NftAuction deployed to (proxy):", proxyAddress);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});

