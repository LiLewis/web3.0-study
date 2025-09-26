// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

// NFT合约的NftAuction
// 核心拍卖合约，支持 ETH 和 ERC20 出价，可升级
contract NftAuction is Initializable, UUPSUpgradeable, OwnableUpgradeable, ReentrancyGuardUpgradeable {
    using SafeERC20 for IERC20Metadata;

    //类构造
    struct Auction {
        uint256 id;
        address seller;
        address nft;
        uint256 tokenId;
        address paymentToken;
        uint256 startTime;
        uint256 endTime;
        address highestBidder;
        uint256 highestBid;
        uint256 highestBidUsd;
        bool settled;
        uint256 reserve;
    }

    //变量映射
    uint256 public auctionCount;    
    mapping(uint256 => Auction) public auctions;
    mapping(address => uint256) public pendingReturns;
    mapping(address => address) public priceFeed;

    // ---- Events ---- 声明`event`事件
    event AuctionCreated(
        uint256 indexed id,
        address indexed nft,
        uint256 indexed tokenId,
        address seller
    );
    event BidPlaced(
        uint256 indexed id,
        address indexed bidder,
        uint256 amount,
        uint256 usdValue
    );
    event AuctionEnded(
        uint256 indexed id,
        address indexed winner,
        uint256 amount
    );
    event Withdrawn(
        address indexed who,
        uint256 amount
    );
    event PriceFeedSet(address token, address feed);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address admin) public initializer {
        __Ownable_init(admin);
        __UUPSUpgradeable_init(); //uups代理
        __ReentrancyGuard_init(); //防止重入攻击

        if (admin != address(0)) {
            transferOwnership(admin);
        }
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}

    function setPriceFeed(address token, address feed) external onlyOwner {
        priceFeed[token] = feed;
        emit PriceFeedSet(token, feed);
    }

    //创建合约
    function createAuction(
        uint256 durationSeconds,
        uint256 reserve,
        address nft,
        uint256 tokenId
    ) external returns (uint256) {
        require(durationSeconds > 0 && durationSeconds < 30 days, "bad duration");
        require(IERC721(nft).ownerOf(tokenId) == msg.sender, "not owner");

        auctionCount += 1;
        uint256 id = auctionCount;
        uint256 start = block.timestamp;
        uint256 end = start + durationSeconds;

        auctions[id] = Auction({
            id: id,
            seller: msg.sender,
            nft: nft,
            tokenId: tokenId,
            paymentToken: address(0),
            startTime: start,
            endTime: end,
            highestBidder: address(0),
            highestBid: 0,
            highestBidUsd: 0,
            settled: false,
            reserve: reserve
        });

        IERC721(nft).transferFrom(msg.sender, address(this), tokenId);

        emit AuctionCreated(id, nft, tokenId, msg.sender);
        return id;
    }

    function placeBid(uint256 auctionId) external payable nonReentrant {
        Auction storage a = auctions[auctionId];
        require(block.timestamp >= a.startTime && block.timestamp < a.endTime, "not active");
        require(a.paymentToken == address(0), "ETH not enabled");
        require(msg.value > 0, "zero");

        uint256 usd = _toUsdScaled(address(0), msg.value);
        require(usd > a.highestBidUsd && usd >= a.reserve, "bid too low");

        if (a.highestBidder != address(0)) {
            pendingReturns[a.highestBidder] += a.highestBid;
        }

        a.highestBid = msg.value;
        a.highestBidUsd = usd;
        a.highestBidder = msg.sender;

        emit BidPlaced(auctionId, msg.sender, msg.value, usd);
    }

    function bidERC20(uint256 auctionId, uint256 amount, address token) external nonReentrant {
        Auction storage a = auctions[auctionId];
        require(block.timestamp >= a.startTime && block.timestamp < a.endTime, "not active");
        require(a.paymentToken == token, "token mismatch");
        require(amount > 0, "zero");

        IERC20Metadata(token).safeTransferFrom(msg.sender, address(this), amount);

        uint256 usd = _toUsdScaled(token, amount);
        require(usd > a.highestBidUsd && usd >= a.reserve, "bid too low");

        if (a.highestBidder != address(0)) {
            pendingReturns[a.highestBidder] += a.highestBid;
        }

        a.highestBid = amount;
        a.highestBidUsd = usd;
        a.highestBidder = msg.sender;

        emit BidPlaced(auctionId, msg.sender, amount, usd);
    }

    //转账
    function withdraw() external nonReentrant {
        uint256 amount = pendingReturns[msg.sender];
        require(amount > 0, "no funds");
        pendingReturns[msg.sender] = 0;

        (bool ok, ) = msg.sender.call{value: amount}("");
        require(ok, "transfer failed");
        emit Withdrawn(msg.sender, amount);
    }

    //结束合约
    function endAuction(uint256 auctionId) external nonReentrant {
        Auction storage a = auctions[auctionId];
        require(block.timestamp >= a.endTime, "not ended");
        require(!a.settled, "already settled");
        a.settled = true;

        if (a.highestBidder == address(0)) {
            IERC721(a.nft).transferFrom(address(this), a.seller, a.tokenId);
            emit AuctionEnded(auctionId, address(0), 0);
            return;
        }

        IERC721(a.nft).transferFrom(address(this), a.highestBidder, a.tokenId);

        if (a.paymentToken == address(0)) {
            (bool ok, ) = a.seller.call{value: a.highestBid}("");
            require(ok, "transfer failed");
        } else {
            IERC20Metadata(a.paymentToken).safeTransfer(a.seller, a.highestBid);
        }

        emit AuctionEnded(auctionId, a.highestBidder, a.highestBid);
    }

    function _toUsdScaled(address token, uint256 amount) internal view returns (uint256) {
        address feed = priceFeed[token];
        require(feed != address(0), "no feed");
        (, int256 price, , , ) = AggregatorV3Interface(feed).latestRoundData();
        require(price > 0, "bad price");

        uint8 tokenDecimals = token == address(0) ? 18 : IERC20Metadata(token).decimals();
        uint8 feedDecimals = AggregatorV3Interface(feed).decimals();
        return (amount * uint256(price)) / (10 ** tokenDecimals);
    }

    receive() external payable {}
}
