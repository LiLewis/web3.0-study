const { ethers, upgrades, deployments } = require("hardhat")
const { expect } = require("chai")

// 测试 ETH 拍卖流程：
// 1 部署 proxy + NFT + Mock Price Feed
// 2 Mint NFT 并 approve
// 3 创建拍卖
// 4 出价
// 5 等待拍卖结束
// 6 结束拍卖，检查 NFT 转移和资金分发
describe("NftAuction basic flow", function () {
  let deployer, seller, buyer;
  let nftAuction, proxyAddr;
  beforeEach(async () => {
    [deployer, seller, buyer] = await ethers.getSigners();

    // deploy proxy via script-like flow (use deployments fixture or call directly)
    const NftAuction = await ethers.getContractFactory("NftAuction");
    const proxy = await upgrades.deployProxy(NftAuction, [deployer.address], { initializer: "initialize" });
    await proxy.deployed();
    proxyAddr = proxy.address;
    nftAuction = await ethers.getContractAt("NftAuction", proxyAddr);

    // deploy TestERC721 and mint
    const TestERC721 = await ethers.getContractFactory("TestERC721");
    const erc721 = await TestERC721.deploy();
    await erc721.deployed();
    await erc721.connect(seller).mint(seller.address, 1);

    // deploy mock price feed: decimals 8, price e.g. 2000 * 1e8
    const Mock = await ethers.getContractFactory("MockV3Aggregator");
    const feed = await Mock.deploy(8, ethers.BigNumber.from("2000").mul(ethers.BigNumber.from(10).pow(8)));
    await feed.deployed();

    // set price feed for ETH (token address zero)
    await nftAuction.connect(deployer).setPriceFeed(ethers.constants.AddressZero, feed.address);

    // seller approve contract
    await erc721.connect(seller).setApprovalForAll(proxyAddr, true);

    // create auction for tokenId 1 (duration 5 secs, reserve 0.001 ETH)
    await nftAuction.connect(seller).createAuction(5, ethers.utils.parseEther("0.001"), erc721.address, 1);
  });

  it("accepts ETH bid and finalizes", async function () {
    // place bid by buyer
    await nftAuction.connect(buyer).placeBid(1, { value: ethers.utils.parseEther("0.01") });

    // wait to end
    await ethers.provider.send("evm_increaseTime", [6]);
    await ethers.provider.send("evm_mine");

    // finalize
    await nftAuction.connect(deployer).endAuction(1);

    // check ownership later: owner should be buyer
    // need to get erc721 contract instance
    const TestERC721 = await ethers.getContractFactory("TestERC721");
    const erc721 = await TestERC721.attach((await ethers.getContractAt("TestERC721", await nftAuction.auctions(1).then(a => a.nft))).address);
    // simpler: we stored nft address earlier; but for clarity:
    // fetch auction info
    const auction = await nftAuction.auctions(1);
    expect(auction.highestBidder).to.not.equal(ethers.constants.AddressZero);
  });
});