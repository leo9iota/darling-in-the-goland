Spike = {}
Spike.__index = Spike
ActiveSpikes = {}

function Spike.new(x, y)
    local instance = setmetatable({}, Spike)
    instance.x = x
    instance.y = y
    instance.image = love.graphics.newImage("assets/spikes.png")
    instance.width = instance.image:getWidth()
    instance.height = instance.image:getHeight()

    instance.physics = {}
    instance.physics.body = love.physics.newBody(World, instance.x, instance.y, "static")
    instance.physics.shape = love.physics.newRectangleShape(instance.width, instance.height)
    instance.physics.fixture = love.physics.newFixture(instance.physics.body, instance.physics.shape)
    instance.physics.fixture:setSensor(true)
    table.insert(ActiveSpikes, instance)
end

function Spike:update(dt)

end

function Spike.updateAllSpikes(dt)
    for index, instance in ipairs(ActiveSpikes) do
        instance:update(dt)
    end
end

function Spike:draw()
    love.graphics.draw(self.image, self.x, self.y, 0, self.scaleX, 1, self.width / 2, self.height / 2)
end

function Spike.drawAllSpikes()
    for index, instance in ipairs(ActiveSpikes) do
        instance:draw()
    end
end

function Spike.beginContact(fixtureA, fixtureB, collision)
    for index, instance in ipairs(ActiveSpikes) do
        if fixtureA == instance.physics.fixture or fixtureB == instance.physics.fixture then
            if fixtureA == Player.physics.fixture or fixtureB == Player.physics.fixture then
                instance.isSpikeRemovable = true
                return true
            end
        end
    end
end
