--- src/gui/GUI.lua
-- @class GUI
-- Draws the heart counter and coin counter
local GUI = {}
local Player = require("src.Player")

-- Love 2D load function 
function GUI:load()
    self.coinCounter = {} -- Coin counter table
    self.coinCounter.x = love.graphics.getWidth() - 200 -- Coin counter x position in GUI
    self.coinCounter.y = 50 -- Coin counter y position in GUI
    self.coinCounter.image = love.graphics.newImage("assets/world/coin.png") -- Load coin image
    self.coinCounter.width = self.coinCounter.image:getWidth() -- Get width of coin image
    self.coinCounter.height = self.coinCounter.image:getHeight() -- Get height of coin image
    self.coinCounter.scale = 3 -- Set scaling to 300%
    self.publicPixelFont = love.graphics.newFont("assets/fonts/public-pixel-font.ttf", 40) -- Coin counter font

    self.heartCounter = {} -- Heart counter table
    self.heartCounter.x = 0 -- Heart counter x position in GUI
    self.heartCounter.y = 50 -- Heart counter y position in GUI
    self.heartCounter.image = love.graphics.newImage("assets/world/heart.png") -- Load heart image
    self.heartCounter.width = self.heartCounter.image:getWidth() -- Get width of heart image
    self.heartCounter.height = self.heartCounter.image:getHeight() -- Get height of heart image
    self.heartCounter.scale = 3 -- Set scaling to 300%
    self.heartCounter.spaceBetween = self.heartCounter.width * self.heartCounter.scale + 30 -- Set the spacing between the hearts

end

-- Love 2D update function
function GUI:update(dt)

end

-- Love 2D draw function
function GUI:draw()
    self:drawCoinCounter() -- Draw coin counter to screen
    self:drawHeartCounter() -- Draw heart counter to screen
end

function GUI:drawHeartCounter()
    -- Iterate over the current player health table
    for i = 1, Player.health.current, 1 do
        local x = self.heartCounter.x + self.heartCounter.spaceBetween * i
        local shadowOffset = 2

        -- Draw shadow for heart counter with an offset of 3
        love.graphics.setColor(0, 0, 0, 0.5)
        love.graphics.draw(self.heartCounter.image, x + shadowOffset, self.heartCounter.y + shadowOffset, 0, self.heartCounter.scale,
            self.heartCounter.scale)
        love.graphics.setColor(1, 1, 1, 1)

        -- Draw heart image
        love.graphics.draw(self.heartCounter.image, x, self.heartCounter.y, 0, self.heartCounter.scale, self.heartCounter.scale)
    end
end

--[[
    This function is responsible for drawing the coin image and coin count to the screen.
]]
function GUI:drawCoinCounter()
    -- Set new font
    love.graphics.setFont(self.publicPixelFont)

    local shadowOffset = 3

    -- Draw shadow for coin counter with an offset of 3
    love.graphics.setColor(0, 0, 0, 0.5)
    love.graphics.draw(self.coinCounter.image, self.coinCounter.x + shadowOffset, self.coinCounter.y + 3, 0, self.coinCounter.scale,
        self.coinCounter.scale)
    love.graphics.setColor(1, 1, 1, 1)

    -- Draw coin image
    love.graphics.draw(self.coinCounter.image, self.coinCounter.x, self.coinCounter.y, 0, self.coinCounter.scale, self.coinCounter.scale)

    -- Coin counter position in the GUI
    local x = self.coinCounter.x + self.coinCounter.width * self.coinCounter.scale + 10
    local y = self.coinCounter.y + self.coinCounter.height / 2 * self.coinCounter.scale - self.publicPixelFont:getHeight() / 2

    -- Create a shadow around the coin
    love.graphics.setColor(0, 0, 0, 0.5)

    -- Offset shadow
    love.graphics.print(Player.coinCount, x + 3, y + 3)
    love.graphics.setColor(1, 1, 1, 1)
    love.graphics.print(Player.coinCount, x, y)
end

return GUI
