// scripts/08-local-auctionFlow.js
const { ethers } = require("hardhat");

async function main() {
  console.log("🏁 Starting local auction flow on Hardhat node...");

  // 获取前 3 个钱包
  const [deployer, seller, bidder] = await ethers.getSigners();
  console.log("Deployer:", deployer.address);
  console.log("Sample seller:", seller.address);
  console.log("Sample bidder:", bidder.address);

  // 1️⃣ 部署 MockV3Aggregator
  console.log("🚀 Deploying MockV3Aggregator...");
  const MockV3Aggregator = await ethers.getContractFactory("MockV3Aggregator", deployer);
  const mockPriceFeed = await MockV3Aggregator.deploy(8, 200_000_000_000); // 初始价格 2000 USD * 1e8
  await mockPriceFeed.waitForDeployment();
  console.log("✅ MockV3Aggregator deployed:", await mockPriceFeed.getAddress());

  // 2️⃣ 部署 TestERC20
  console.log("🚀 Deploying TestERC20...");
  const TestERC20 = await ethers.getContractFactory("TestERC20", deployer);
  const erc20 = await TestERC20.deploy();
  await erc20.waitForDeployment();
  console.log("✅ TestERC20 deployed:", await erc20.getAddress());

  // 给 bidder 转一些 ERC20
  const bidAmount = ethers.parseUnits("1000", 18);
  await erc20.transfer(bidder.address, bidAmount);
  console.log(`💰 Transferred ${bidAmount} ERC20 to bidder`);

  // 3️⃣ 部署 TestERC721
  console.log("🚀 Deploying TestERC721...");
  const TestERC721 = await ethers.getContractFactory("TestERC721", seller);
  const erc721 = await TestERC721.deploy();
  await erc721.waitForDeployment();
  console.log("✅ TestERC721 deployed:", await erc721.getAddress());

  // 4️⃣ 部署 NftAuction
  console.log("🚀 Deploying NftAuction...");
  const NftAuction = await ethers.getContractFactory("NftAuction", deployer);
  const auction = await NftAuction.deploy(await mockPriceFeed.getAddress());
  await auction.waitForDeployment();
  console.log("✅ NftAuction deployed:", await auction.getAddress());

  // 5️⃣ Mint NFT 给 seller
  console.log("🎨 Minting NFT to seller...");
  const tokenId = 1;
  await erc721.connect(seller).mint(tokenId);
  console.log("✅ NFT minted, tokenId =", tokenId);

  // 6️⃣ Seller 批准 NFT 给拍卖合约
  await erc721.connect(seller).setApprovalForAll(await auction.getAddress(), true);
  console.log("✅ NFT approved to auction contract");

  // 7️⃣ 创建拍卖
  const reserve = ethers.parseEther("0.01");
  const duration = 60 * 60; // 1 小时
  const tx = await auction.connect(seller).createAuction(duration, reserve, await erc721.getAddress(), tokenId);
  const receipt = await tx.wait();
  const auctionId = receipt.logs[0].args.id;
  console.log("🏷️ Auction created, id =", auctionId);

  // 8️⃣ Bidder 出价
  const bidValue = ethers.parseEther("0.02");
  await auction.connect(bidder).createBid(auctionId, { value: bidValue });
  console.log(`💸 Bidder placed bid: ${bidValue} ETH`);

  // 9️⃣ 结算拍卖
  // 快进时间模拟拍卖结束
  await ethers.provider.send("evm_increaseTime", [duration + 1]);
  await ethers.provider.send("evm_mine");

  await auction.connect(deployer).settleAuction(auctionId);
  console.log("🏁 Auction settled");

  // 查询 NFT 新 owner
  const newOwner = await erc721.ownerOf(tokenId);
  console.log("🎉 NFT new owner:", newOwner);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
