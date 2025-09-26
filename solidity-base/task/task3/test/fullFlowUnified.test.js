const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("Full NFT Auction Flow (ETH/ERC20) + Upgrade", function () {
  let nftContract, erc20Contract, auctionContract;
  let owner, addr1, addr2;

  before(async function () {
    [owner, addr1, addr2] = await ethers.getSigners();

    // 部署 NFT 合约
    const NFT = await ethers.getContractFactory("TestERC721");
    nftContract = await NFT.deploy();
    await nftContract.waitForDeployment();

    // 部署 ERC20 合约
    const ERC20 = await ethers.getContractFactory("TestERC20");
    erc20Contract = await ERC20.deploy();
    await erc20Contract.waitForDeployment();

    // 部署拍卖合约（UUPS Proxy）
    const Auction = await ethers.getContractFactory("NftAuction");
    auctionContract = await upgrades.deployProxy(
      Auction,
      [owner.address], // 初始化 admin
      { initializer: "initialize" }
    );
    await auctionContract.waitForDeployment();
  });

  it("should handle ETH auction flow", async function () {
    // mint 一个 NFT 给 owner
    const mintTx = await nftContract.mint(owner.address, 1);
    await mintTx.wait();

    // 授权拍卖合约转移 NFT
    await nftContract.approve(await auctionContract.getAddress(), 1);

    // 创建 ETH 拍卖
    const createTx = await auctionContract.createAuction(
      await nftContract.getAddress(),
      1,
      ethers.ZeroAddress, // ZeroAddress = ETH 支付
      ethers.parseEther("0.1"),
      60
    );
    await createTx.wait();

    // addr1 出价 0.2 ETH
    await auctionContract.connect(addr1).placeBid(0, {
      value: ethers.parseEther("0.2"),
    });

    // addr2 出价 0.3 ETH
    await auctionContract.connect(addr2).placeBid(0, {
      value: ethers.parseEther("0.3"),
    });

    // 快进时间
    await ethers.provider.send("evm_increaseTime", [70]);
    await ethers.provider.send("evm_mine");

    // 结束拍卖
    await auctionContract.endAuction(0);

    // 检查 NFT 归属
    expect(await nftContract.ownerOf(1)).to.equal(addr2.address);
  });

  it("should handle ERC20 auction flow", async function () {
    // mint NFT 给 owner
    await nftContract.mint(owner.address, 2);
    await nftContract.approve(await auctionContract.getAddress(), 2);

    // mint ERC20 给 addr1 和 addr2
    await erc20Contract.mint(addr1.address, ethers.parseEther("100"));
    await erc20Contract.mint(addr2.address, ethers.parseEther("100"));

    // 创建 ERC20 拍卖
    await auctionContract.createAuction(
      await nftContract.getAddress(),
      2,
      await erc20Contract.getAddress(), // ERC20 支付
      ethers.parseEther("10"),
      60
    );

    // addr1 出价
    await erc20Contract.connect(addr1).approve(await auctionContract.getAddress(), ethers.parseEther("20"));
    await auctionContract.connect(addr1).placeBid(1, { value: 0 });

    // addr2 出价更高
    await erc20Contract.connect(addr2).approve(await auctionContract.getAddress(), ethers.parseEther("30"));
    await auctionContract.connect(addr2).placeBid(1, { value: 0 });

    // 快进时间
    await ethers.provider.send("evm_increaseTime", [70]);
    await ethers.provider.send("evm_mine");

    // 结束拍卖
    await auctionContract.endAuction(1);

    // 检查 NFT 归属
    expect(await nftContract.ownerOf(2)).to.equal(addr2.address);
  });

  it("should upgrade the contract", async function () {
    const AuctionV2 = await ethers.getContractFactory("NftAuctionV2");
    const upgraded = await upgrades.upgradeProxy(
      await auctionContract.getAddress(),
      AuctionV2
    );

    await upgraded.waitForDeployment();

    // 调用 V2 新增函数
    expect(await upgraded.version()).to.equal("V2");
  });
});
