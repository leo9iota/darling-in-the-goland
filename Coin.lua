Coin = {}
Coin.__index = Coin
ActiveCoins = {}

--[[
    This function acts like a constructor (similar to Java). We have to do it this way
    because Lua does not support OOP natively. It has two parameters, x and y, which define
    the position of the coin.
]]
function Coin.new(x, y)
    local instance = setmetatable({}, Coin)
    instance.x = x
    instance.y = y
    instance.image = love.graphics.newImage("assets/coin.png")
    instance.width = instance.image:getWidth()
    instance.height = instance.image:getHeight()
    instance.scaleX = 1
    instance.randomTimeOffset = math.random(0, 100)
    instance.isCoinRemovable = false

    instance.physics = {}
    instance.physics.body = love.physics.newBody(World, instance.x, instance.y, "static")
    instance.physics.shape = love.physics.newRectangleShape(instance.width, instance.height)
    instance.physics.fixture = love.physics.newFixture(instance.physics.body, instance.physics.shape)
    instance.physics.fixture:setSensor(true)
    table.insert(ActiveCoins, instance)
end

--[[
    This function is responsible for removing a coin from the `ActiveCoins` table is the
    player touched the coin. The physical body of the coin is stored inside of the World
    and doesn't get removed even though we remove the coin instance. To get rid of it, we
    also need to utilize the LÖVE 2D function `body:destroy`.    
]]
function Coin:removeCoin()
    for index, instance in ipairs(ActiveCoins) do
        -- Check if the current instance equals to itself
        if instance == self then
            Player:collectCoin()
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
        self:removeCoin()
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
    the table and update each instance (coin) with the `update(dt)` function provided
    by LÖVE 2D.
]]
function Coin.updateAllCoins(dt)
    for index, instance in ipairs(ActiveCoins) do
        instance:update(dt)
    end
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
function Coin.drawAllCoins()
    for index, instance in ipairs(ActiveCoins) do
        instance:draw()
    end
end

--[[
    This function is responsible for handling the collision between the player and the coin.
    We loop through the `ActiveCoins` table and check if the player has collided with a
    coin.
]]
function Coin.beginContact(fixtureA, fixtureB, collision)
    for index, instance in ipairs(ActiveCoins) do
        if fixtureA == instance.physics.fixture or fixtureB == instance.physics.fixture then
            if fixtureA == Player.physics.fixture or fixtureB == Player.physics.fixture then
                -- This is the variable that checks if a object should be removed
                instance.isCoinRemovable = true
                return true
            end
        end
    end
end
