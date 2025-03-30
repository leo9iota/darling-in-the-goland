-- In your main.lua, add this line near your other requires
local SlabGUI = require("src.SlabGUI")

-- In your love.load function, add:
function love.load()
    -- Your existing code
    -- ...

    SlabGUI:load()
end

-- In your love.update function:
function love.update(dt)
    -- Your existing code
    -- ...

    SlabGUI:update(dt)
end

-- In your love.draw function:
function love.draw()
    -- Your existing code
    -- ...

    SlabGUI:draw()
end

-- Add these callbacks for Slab to work properly
function love.keypressed(key, scancode, isrepeat)
    Player:jump(key)

    if key == "escape" then
        love.event.quit()
    end

    SlabGUI:keypressed(key, scancode, isrepeat)
end

function love.keyreleased(key)
    SlabGUI:keyreleased(key)
end

function love.textinput(text)
    SlabGUI:textinput(text)
end

function love.mousepressed(x, y, button)
    SlabGUI:mousepressed(x, y, button)
end

function love.mousereleased(x, y, button)
    SlabGUI:mousereleased(x, y, button)
end

function love.mousemoved(x, y, dx, dy)
    SlabGUI:mousemoved(x, y, dx, dy)
end
