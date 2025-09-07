// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

// 1.确保你的 MetaMask 已经切换到 Sepolia 测试网。
// 2.连接小狐狸，部署合约。
// 3.增发代币（Mint）
// 4.调用 mint 函数，to = 你的钱包地址, amount = 1000000000000000000 （=1 MTK，因为有 18 位小数）
// 5.查询余额
// 6.调用 balanceOf(你的钱包地址)，会返回 1000000000000000000。
// 7.转账调用 transfer(目标地址, 1000000000000000000)，把 1 MTK 转给别人。
// 8.授权 & 代扣先调用 approve(某个地址, 500000000000000000) 授权别人代扣 0.5 MTK。
// 9.再用被授权的地址调用 transferFrom(你的地址, 其他地址, 500000000000000000)。
// 10.复制合约的部署地址。打开 MetaMask → 导入代币 → 自定义代币。粘贴合约地址，就能看到你创建的代币。
contract MyERC20 {
    // 代币名字
    string public name = "LULUB";
    // 代币符号
    string public symbol = "MTK";
    // 小数位，通常是18
    uint8 public decimals = 18;
    // 代币总量
    uint256 public totalSupply;
    // 合约部署者
    address public owner;

    // 账户余额
    mapping(address => uint256) private balances;

    // 授权额度： owner => (spender => amount)
    mapping(address => mapping(address => uint256)) private allowances;

    // 事件
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    constructor() {
        // 部署者为初始所有者
        owner = msg.sender;

        // 初始 1000 个代币
        uint256 initialSupply = 1000 * (10 ** uint256(decimals));
        totalSupply = initialSupply;
        balances[owner] = initialSupply;
        // 记录 mint 事件
        emit Transfer(address(0), owner, initialSupply);
    }

    // 查询余额
    function balanceOf(address account) public view returns (uint256) {
        return balances[account];
    }

    // 转账
    function transfer(address to, uint256 amount) public returns (bool) {
        require(balances[msg.sender] >= amount, "Not enough balance");
        balances[msg.sender] -= amount;
        balances[to] += amount;
        emit Transfer(msg.sender, to, amount);
        return true;
    }

    // 授权
    function approve(address spender, uint256 amount) public returns (bool) {
        allowances[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    // 查询授权额度
    function allowance(address _owner, address spender) public view returns (uint256) {
        return allowances[_owner][spender];
    }

    // 代扣转账
    function transferFrom(address from, address to, uint256 amount) public returns (bool) {
        require(balances[from] >= amount, "Not enough balance");
        require(allowances[from][msg.sender] >= amount, "Allowance exceeded");

        balances[from] -= amount;
        balances[to] += amount;
        allowances[from][msg.sender] -= amount;

        emit Transfer(from, to, amount);
        return true;
    }

    // 增发代币（只有所有者可以调用）
    function mint(address to, uint256 amount) public {
        require(msg.sender == owner, "Only owner can mint");
        totalSupply += amount;
        balances[to] += amount;
        emit Transfer(address(0), to, amount);
    }

}