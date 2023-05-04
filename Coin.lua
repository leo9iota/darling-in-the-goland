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
    table.insert(ActiveCoins, instance) -- Insert all instances into `ActiveCoins` table
end

function Coin:update(dt)
    
end

--[[
    This function is responsible for making the coin appear to be spinning with the help
    of scaling.
]]
function Coin:spinAnim()
    
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
    love.graphics.draw(self.image, self.x, self.y, 0, 1, 1, self.width / 2, self.height / 2)
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