--- src/Camera.lua
-- @class Coin
-- Coin game objects that the player picks up
local Coin = {}
Coin.__index = Coin
local ActiveCoins = {}
local Player = require("src.Player")

--[[
    This function acts like a constructor (similar to Java). We have to do it this way
    because Lua does not support OOP natively. It has two parameters, x and y, which define
    the position of the coin.
]]
function Coin.new(x, y)
    local coin = setmetatable({}, Coin)
    coin.x = x
    coin.y = y
    coin.image = love.graphics.newImage("assets/world/coin.png")
    coin.width = coin.image:getWidth()
    coin.height = coin.image:getHeight()
    coin.scaleX = 1
    coin.randomTimeOffset = math.random(0, 100)
    coin.isCoinRemovable = false

    coin.physics = {}
    coin.physics.body = love.physics.newBody(World, coin.x, coin.y, "static")
    coin.physics.shape = love.physics.newRectangleShape(coin.width, coin.height)
    coin.physics.fixture = love.physics.newFixture(coin.physics.body, coin.physics.shape)
    coin.physics.fixture:setSensor(true)
    table.insert(ActiveCoins, coin)
end

--[[
    This function is responsible for removing a coin from the `ActiveCoins` table is the
    player touched the coin. The physical body of the coin is stored inside of the World
    and doesn't get removed even though we remove the coin coin. To get rid of it, we
    also need to utilize the Love2D function `body:destroy`.    
]]
function Coin:remove()
    for index, coin in ipairs(ActiveCoins) do
        -- Check if the current coin equals to itself
        if coin == self then
            Player:incrementCoinCount()
            print("Coin Count: ", Player.coinCount)
            self.physics.body:destroy()
            -- If thats the case remove from table
            table.remove(ActiveCoins, index)
        end
    end
end

function Coin:update(dt)
    self:spinAnim(dt)
    self:checkCoinRemoval()
end

function Coin.removeAll()
    for i, v in ipairs(ActiveCoins) do -- Loop trough active coin table and destroy them
        v.physics.body:destroy()
    end

    ActiveCoins = {} -- Destroying only reference to coin, which will trigger Lua's garbage collector
end

--[[
    This function removes the coin.
]]
function Coin:checkCoinRemoval()
    if self.isCoinRemovable then
        --[[
            --- FIX ---
            Use colon operator instead of dot operator. The colon operator implicitly
            passses it's own table as first argument. We need pass the table as an
            argument, to access all it's methods and variables.
        ]]
        self:remove()
    end
end

--[[
    This function is responsible for making the coin appear to be spinning with the help
    of scaling.
]]
function Coin:spinAnim(dt)
    self.scaleX = math.sin(love.timer.getTime() * 4 + self.randomTimeOffset)
end

--[[
    This function is responsible for updating all the coins on the map. We loop through
    the table and update each coin (coin) with the `update(dt)` function provided
    by LÖVE 2D.
]]
function Coin.updateAll(dt)
    for index, coin in ipairs(ActiveCoins) do coin:update(dt) end
end

--[[
    This function is responsible for drawing the coins. We also set the origin point to the
    center, due to LÖVE 2D and Box2D (Tiled) implementing origin points differently.
    - LÖVE 2D: top-left
    - Box2D: center
]]
function Coin:draw()
    love.graphics.draw(self.image, self.x, self.y, 0, self.scaleX, 1, self.width / 2, self.height / 2)
end

--[[
    This function is responsible for drawing all coins to the screen which are stored inside
    the `ActiveCoins` table.
]]
function Coin.drawAll()
    for index, coin in ipairs(ActiveCoins) do coin:draw() end
end

--[[
    This function is responsible for handling the collision between the player and the coin.
    We loop through the `ActiveCoins` table and check if the player has collided with a
    coin.
]]
function Coin.beginContact(fixtureA, fixtureB, collision)
    for index, coin in ipairs(ActiveCoins) do
        if fixtureA == coin.physics.fixture or fixtureB == coin.physics.fixture then
            if fixtureA == Player.physics.fixture or fixtureB == Player.physics.fixture then
                -- This is the variable that checks if a object should be removed
                coin.isCoinRemovable = true
                return true
            end
        end
    end
end

return Coin
