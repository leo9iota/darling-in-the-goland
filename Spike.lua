local Spike = {}
Spike.__index = Spike
local ActiveSpikes = {}
local Player = require("Player")

function Spike.new(x, y)
    local spike = setmetatable({}, Spike)
    spike.x = x
    spike.y = y
    spike.image = love.graphics.newImage("assets/spikes.png")
    spike.width = spike.image:getWidth()
    spike.height = spike.image:getHeight()

    spike.damage = 1 -- The amount of damage the player takes when colliding with a spike

    spike.physics = {}
    spike.physics.body = love.physics.newBody(World, spike.x, spike.y, "static")
    spike.physics.shape = love.physics.newRectangleShape(spike.width, spike.height)
    spike.physics.fixture = love.physics.newFixture(spike.physics.body, spike.physics.shape)
    spike.physics.fixture:setSensor(true)
    table.insert(ActiveSpikes, spike)
end

function Spike:update(dt)

end

function Spike.updateAll(dt)
    for index, spike in ipairs(ActiveSpikes) do spike:update(dt) end
end

function Spike:draw()
    love.graphics.draw(self.image, self.x, self.y, 0, self.scaleX, 1, self.width / 2, self.height / 2)
end

function Spike.drawAll()
    for index, spike in ipairs(ActiveSpikes) do spike:draw() end
end

function Spike.removeAll()
    for i, v in ipairs(ActiveSpikes) do v.physics.body:destroy() end

    ActiveSpikes = {}
end

function Spike.beginContact(fixtureA, fixtureB, collision)
    for index, spike in ipairs(ActiveSpikes) do
        if fixtureA == spike.physics.fixture or fixtureB == spike.physics.fixture then
            if fixtureA == Player.physics.fixture or fixtureB == Player.physics.fixture then
                Player:takeDamage(spike.damage)
                return true
            end
        end
    end
end

return Spike
