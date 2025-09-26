const { ethers } = require("hardhat");

const AUCTION_ADDRESS = process.env.AUCTION_ADDRESS;

//手动输入拍卖 ID 来结束，并打印赢家信息
//结束拍卖 #0 or #1  npx hardhat run scripts/05-endAuction.js --network sepolia 0
async function main() {
  const [caller] = await ethers.getSigners();
  console.log("📢 调用人:", caller.address);

  const auction = await ethers.getContractAt("NftAuction", AUCTION_ADDRESS);

  // 👉 从命令行传入拍卖 ID
  const auctionId = process.argv[2];
  if (!auctionId) {
    console.error("❌ 请传入拍卖 ID: npx hardhat run scripts/endAuction.js --network sepolia <auctionId>");
    return;
  }

  console.log(`⏳ 正在结束拍卖 #${auctionId} ...`);
  const tx = await auction.endAuction(auctionId);
  await tx.wait();

  console.log(`✅ 拍卖 #${auctionId} 已结束`);

  // 获取拍卖信息
  const auctionInfo = await auction.auctions(auctionId);
  const winner = auctionInfo.highestBidder;
  const amount = auctionInfo.highestBid;

  console.log("🏆 胜利者:", winner);
  console.log("💰 出价金额:", ethers.formatEther(amount));
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
