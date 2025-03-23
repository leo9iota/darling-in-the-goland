local Stone = {
    image = love.graphics.newImage("assets/stone.png") -- NOTE: Load image here and not within the "new()" method to avoid unnecessary memory usage
}
Stone.__index = Stone

Stone.width = Stone.image:getWidth()
Stone.height = Stone.image:getHeight()

local ActiveStones = {}

function Stone.new(x, y)
    local stone = setmetatable({}, Stone)
    stone.x = x
    stone.y = y
    stone.rotation = 0

    stone.physics = {}
    stone.physics.body = love.physics.newBody(World, stone.x, stone.y, "dynamic")
    stone.physics.shape = love.physics.newRectangleShape(stone.width, stone.height)
    stone.physics.fixture = love.physics.newFixture(stone.physics.body, stone.physics.shape)
    stone.physics.body:setMass(25)
    table.insert(ActiveStones, stone)
end

function Stone:update(dt)
    self:syncPhysics()
end

function Stone:syncPhysics()
    self.x, self.y = self.physics.body:getPosition()
    self.rotation = self.physics.body:getAngle()
end

function Stone:draw()
    love.graphics.draw(self.image, self.x, self.y, self.rotation, self.scaleX, 1, self.width / 2, self.height / 2)
end

function Stone.updateAll(dt)
    for i, instance in ipairs(ActiveStones) do instance:update(dt) end
end

function Stone.drawAll()
    for i, instance in ipairs(ActiveStones) do instance:draw() end
end

function Stone.removeAll()
    for i, v in ipairs(ActiveStones) do v.physics.body:destroy() end

    ActiveStones = {}
end

return Stone
