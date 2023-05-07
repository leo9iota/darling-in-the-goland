GUI = {}

function GUI:load()
    self.coinCounter = {}
    self.coinCounter.x = 50
    self.coinCounter.y = 50
    self.coinCounter.image = love.graphics.newImage("assets/coin.png")
    self.coinCounter.width = self.coinCounter.image:getWidth()
    self.coinCounter.height = self.coinCounter.image:getHeight()
    self.coinCounter.scale = 3 -- Equivalent to 300%
    self.publicPixelFont = love.graphics.newFont("assets/fonts/public-pixel-font.ttf", 40)
end

function GUI:update(dt)

end

function GUI:draw()
    self:displayCoinCounter() -- Draw coin counter to screen
end

--[[
    This function is responsible for drawing the coin image and coin count to the screen.
]]
function GUI:displayCoinCounter()
    --[[
        This creates a shadow around the coin image
    ]]
    love.graphics.setColor(0, 0, 0, 0.5)
    -- Draw coin image
    love.graphics.draw(self.coinCounter.image, self.coinCounter.x + 3, self.coinCounter.y + 3, 0, self.coinCounter.scale, self.coinCounter.scale)
    love.graphics.setColor(1, 1, 1, 1)
    love.graphics.draw(self.coinCounter.image, self.coinCounter.x, self.coinCounter.y, 0, self.coinCounter.scale, self.coinCounter.scale)
    
    -- Set new font
    love.graphics.setFont(self.publicPixelFont)

    -- Coin counter alignment
    local x = self.coinCounter.x + self.coinCounter.width * self.coinCounter.scale
    local y = self.coinCounter.y + self.coinCounter.height / 2 * self.coinCounter.scale - self.publicPixelFont:getHeight() / 2

    --[[
        This creates a shadow around the coin value
    ]]
    love.graphics.setColor(0, 0, 0, 0.5)
    -- Print coin count
    love.graphics.print(":" .. Player.coinCount, x + 3, y + 3)
    love.graphics.setColor(1, 1, 1, 1)
    love.graphics.print(":" .. Player.coinCount, x, y)
end
