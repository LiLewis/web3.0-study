const { ethers } = require("hardhat");
require("dotenv").config();

const NFT_ADDRESS = process.env.NFT_ADDRESS;
const ERC20_ADDRESS = process.env.ERC20_ADDRESS;
const AUCTION_ADDRESS = process.env.AUCTION_ADDRESS;

const ADDR1 = process.env.ADDR1;
const ADDR2 = process.env.ADDR2;

// 创建 & 出价流程
// npx hardhat run scripts/04-auctionFlow.js --network sepolia
async function main() {
  const [owner] = await ethers.getSigners();
  console.log("👑 NFT Owner:", owner.address);

  const addr1 = await ethers.getSigner(ADDR1);
  const addr2 = await ethers.getSigner(ADDR2);
  console.log("👤 Addr1:", addr1.address);
  console.log("👤 Addr2:", addr2.address);

  const nft = await ethers.getContractAt("TestERC721", NFT_ADDRESS);
  const erc20 = await ethers.getContractAt("TestERC20", ERC20_ADDRESS);
  const auction = await ethers.getContractAt("NftAuction", AUCTION_ADDRESS);

  // === ETH 拍卖流程 ===
  console.log("\n=== ETH 拍卖流程 ===");
  const tokenIdETH = 1029;
  await (await nft.mint(owner.address, tokenIdETH)).wait();
  console.log("✅ Mint 成功: tokenId =", tokenIdETH);

  await (await nft.approve(AUCTION_ADDRESS, tokenIdETH)).wait();
  console.log("✅ NFT 已批准给拍卖合约:", AUCTION_ADDRESS);

  const txETH = await auction.createAuction(60, ethers.parseEther("0.1"), NFT_ADDRESS, tokenIdETH);
  const receiptETH = await txETH.wait();
  // const auctionIdETH = receiptETH.events?.find(e => e.event === "AuctionCreated")?.args?.auctionId;
  // console.log("📢 创建 ETH 拍卖成功, AuctionId =", auctionIdETH.toString());

  // 确保 mint 和 approve 都用 NFT owner
  console.log("✅ NFT Owner:", owner.address);
  await (await nft.connect(owner).mint(owner.address, tokenIdETH, { gasLimit: 500_000 })).wait();
  console.log("✅ Mint 成功: tokenId =", tokenIdETH);

  await (await nft.connect(owner).approve(AUCTION_ADDRESS, tokenIdETH, { gasLimit: 100_000 })).wait();
  console.log("✅ NFT 已批准给拍卖合约:", AUCTION_ADDRESS);

  // 创建拍卖也用 owner
  const auctionIdETH = await auction.connect(owner).createAuction(
      60, 
      ethers.parseEther("0.1").toString(), 
      NFT_ADDRESS, 
      tokenIdETH
  );
  console.log("📢 ETH 拍卖创建成功, AuctionId =", auctionIdETH.toString());

  await (await auction.connect(addr1).placeBid(auctionIdETH, { value: ethers.parseEther("0.2") })).wait();
  console.log("💰 Addr1 出价 0.2 ETH");
  await (await auction.connect(addr2).placeBid(auctionIdETH, { value: ethers.parseEther("0.3") })).wait();
  console.log("💰 Addr2 出价 0.3 ETH");

  console.log("⏳ 等待拍卖时间结束后执行 endAuction(", auctionIdETH.toString(), ") ...");

  // === ERC20 拍卖流程 ===
  console.log("\n=== ERC20 拍卖流程 ===");
  const tokenIdERC20 = 1022;
  await (await nft.mint(owner.address, tokenIdERC20)).wait();
  console.log("✅ Mint 成功: tokenId =", tokenIdERC20);

  await (await nft.approve(AUCTION_ADDRESS, tokenIdERC20)).wait();
  console.log("✅ NFT 已批准给拍卖合约:", AUCTION_ADDRESS);

  // 给 Addr1 和 Addr2 分发 ERC20
  await (await erc20.transfer(addr1.address, ethers.parseEther("100"))).wait();
  await (await erc20.transfer(addr2.address, ethers.parseEther("100"))).wait();
  console.log("✅ ERC20 已分发给 Addr1/Addr2");

  const txERC20 = await auction.createAuction(60, ethers.parseEther("10"), NFT_ADDRESS, tokenIdERC20);
  const receiptERC20 = await txERC20.wait();
  // const auctionIdERC20 = receiptERC20.events?.find(e => e.event === "AuctionCreated")?.args?.auctionId;
  // console.log("📢 创建 ERC20 拍卖成功, AuctionId =", auctionIdERC20.toString());
  const auctionIdERC20 = await auction.createAuction(
    60,
    ethers.parseEther("10").toString(), 
    NFT_ADDRESS,
    tokenIdERC20
  );
  console.log("📢 创建 ERC20 拍卖成功, AuctionId =", auctionIdERC20.toString());


  await (await erc20.connect(addr1).approve(AUCTION_ADDRESS, ethers.parseEther("20"))).wait();
  await (await auction.connect(addr1).placeBid(auctionIdERC20)).wait();
  console.log("💰 Addr1 出价 20 ERC20");

  await (await erc20.connect(addr2).approve(AUCTION_ADDRESS, ethers.parseEther("30"))).wait();
  await (await auction.connect(addr2).placeBid(auctionIdERC20)).wait();
  console.log("💰 Addr2 出价 30 ERC20");

  console.log("⏳ 等待拍卖时间结束后执行 endAuction(", auctionIdERC20.toString(), ") ...");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
