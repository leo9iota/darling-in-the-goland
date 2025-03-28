local Camera = {
    x = 0,
    y = 0,
    scale = 2 -- The value 2 translates to 200% zoom
}

function Camera:init()
    love.graphics.push()
    love.graphics.scale(self.scale, self.scale)
    love.graphics.translate(-self.x, self.y) -- Used to move the camera within the game world
end

function Camera:remove()
    love.graphics.pop()
end

function Camera:setPosition(x, y)
    self.x = x - love.graphics.getWidth() / 2 / self.scale -- Set window center to passed in value
    self.y = y

    local rightSide = self.x + love.graphics.getWidth() / 2

    if self.x < 0 then -- Prevent camera to go outside of left side
        self.x = 0
    elseif rightSide > MapWidth then
        self.x = MapWidth - love.graphics.getWidth() / 2
    end
end

return Camera
