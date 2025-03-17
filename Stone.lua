local Stone = {}
Stone.__index = Stone
local ActiveStone = {}
local Player = require("Player")

function Stone.new(x, y)
    local stone = setmetatable({}, Stone)
    stone.x = x
    stone.y = y
    stone.image = love.graphics.newImage("assets/stone.png")
    stone.width = stone.image:getWidth()
    stone.height = stone.image:getHeight()

    stone.damage = 1 -- The amount of damage the player takes when colliding with a spike

    stone.physics = {}
    stone.physics.body = love.physics.newBody(World, spike.x, spike.y, "static")
    stone.physics.shape = love.physics.newRectangleShape(spike.width, spike.height)
    stone.physics.fixture = love.physics.newFixture(spike.physics.body, spike.physics.shape)
    stone.physics.fixture:setSensor(true)
    table.insert(ActiveStone, stone)
end

function Stone:update(dt)

end

function Stone.updateAllSpikes(dt)
    for index, spike in ipairs(ActiveSpikes) do spike:update(dt) end
end

function Stone:draw()
    love.graphics.draw(Stone.image, Stone.x, Stone.y, 0, Stone.scaleX, 1, Stone.width / 2, Stone.height / 2)
end

function Stone.drawAllSpikes()
    for index, spike in ipairs(ActiveSpikes) do spike:draw() end
end

function Stone.beginContact(fixtureA, fixtureB, collision)
    for index, spike in ipairs(ActiveSpikes) do
        if fixtureA == spike.physics.fixture or fixtureB == spike.physics.fixture then
            if fixtureA == Player.physics.fixture or fixtureB == Player.physics.fixture then
                Player:takeDamage(spike.damage)
                return true
            end
        end
    end
end

return Stone
