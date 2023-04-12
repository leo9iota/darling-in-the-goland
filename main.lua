-- Import STI library to import maps from Tiled
local STI = require("sti")

function love.load()
    Map = STI("maps/map-1.lua")
end

function love.update(dt)
    
end

function love.draw()
    love.graphics.clear(62/255, 151/255, 186/255, 255/255)
    Map:draw(0, 0, 2, 2)

    --[[
        The "push()" adds the scaling to the stack and the "pop()" function removes the scaling
        from the stack, which means everything outside of these functions will not be affected
        by the scaling
    ]]
    love.graphics.push()
    love.graphics.scale(2, 2)
    love.graphics.pop()

end