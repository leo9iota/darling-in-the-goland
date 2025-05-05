--- src/gui/HUD.lua
-- @class HUD
-- Draws the heart counter and coin counter
local HUD = {}
local Player = require("src.entities.Player")

-- Love 2D load function 
function HUD:load()
    HUD.coinCounter = {} -- Coin counter table
    HUD.coinCounter.x = love.graphics.getWidth() - 200 -- Coin counter x position in HUD
    HUD.coinCounter.y = 50 -- Coin counter y position in HUD
    HUD.coinCounter.image = love.graphics.newImage("assets/world/coin.png") -- Load coin image
    HUD.coinCounter.width = HUD.coinCounter.image:getWidth() -- Get width of coin image
    HUD.coinCounter.height = HUD.coinCounter.image:getHeight() -- Get height of coin image
    HUD.coinCounter.scale = 3 -- Set scaling to 300%
    HUD.publicPixelFont = love.graphics.newFont("assets/fonts/public-pixel-font.ttf", 40) -- Coin counter font

    HUD.heartCounter = {} -- Heart counter table
    HUD.heartCounter.x = 0 -- Heart counter x position in HUD   
    HUD.heartCounter.y = 50 -- Heart counter y position in HUD  
    HUD.heartCounter.image = love.graphics.newImage("assets/world/heart.png") -- Load heart image
    HUD.heartCounter.width = HUD.heartCounter.image:getWidth() -- Get width of heart image
    HUD.heartCounter.height = HUD.heartCounter.image:getHeight() -- Get height of heart image
    HUD.heartCounter.scale = 3 -- Set scaling to 300%
    HUD.heartCounter.spaceBetween = HUD.heartCounter.width * HUD.heartCounter.scale + 30 -- Set the spacing between the hearts
end

-- Love 2D update function
function HUD:update(dt)

end

-- Love 2D draw function
function HUD:draw()
    HUD:drawCoinCounter() -- Draw coin counter to screen
    HUD:drawHeartCounter() -- Draw heart counter to screen
end

function HUD:drawHeartCounter()
    -- Iterate over the current player health table
    for i = 1, Player.health.current, 1 do
        local x = HUD.heartCounter.x + HUD.heartCounter.spaceBetween * i
        local shadowOffset = 2

        -- Draw shadow for heart counter with an offset of 3
        love.graphics.setColor(0, 0, 0, 0.5)
        love.graphics.draw(HUD.heartCounter.image, x + shadowOffset, HUD.heartCounter.y + shadowOffset, 0, HUD.heartCounter.scale,
            HUD.heartCounter.scale)
        love.graphics.setColor(1, 1, 1, 1)

        -- Draw heart image
        love.graphics.draw(HUD.heartCounter.image, x, HUD.heartCounter.y, 0, HUD.heartCounter.scale, HUD.heartCounter.scale)
    end
end

--[[
    This function is responsible for drawing the coin image and coin count to the screen.
]]
function HUD:drawCoinCounter()
    -- Set new font
    love.graphics.setFont(HUD.publicPixelFont)

    local shadowOffset = 3

    -- Draw shadow for coin counter with an offset of 3
    love.graphics.setColor(0, 0, 0, 0.5)
    love.graphics.draw(HUD.coinCounter.image, HUD.coinCounter.x + shadowOffset, HUD.coinCounter.y + 3, 0, HUD.coinCounter.scale,
        HUD.coinCounter.scale)
    love.graphics.setColor(1, 1, 1, 1)

    -- Draw coin image
    love.graphics.draw(HUD.coinCounter.image, HUD.coinCounter.x, HUD.coinCounter.y, 0, HUD.coinCounter.scale, HUD.coinCounter.scale)

    -- Coin counter position in the GUI
    local x = HUD.coinCounter.x + HUD.coinCounter.width * HUD.coinCounter.scale + 10
    local y = HUD.coinCounter.y + HUD.coinCounter.height / 2 * HUD.coinCounter.scale - HUD.publicPixelFont:getHeight() / 2

    -- Create a shadow around the coin
    love.graphics.setColor(0, 0, 0, 0.5)

    -- Offset shadow
    love.graphics.print(Player.coinCount, x + 3, y + 3)
    love.graphics.setColor(1, 1, 1, 1)
    love.graphics.print(Player.coinCount, x, y)
end

return HUD
