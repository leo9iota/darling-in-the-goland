Player = {}

function love.load()
    self.x = 100
    self.y = 0
    self.width = 20
    self.height = 60
    self.xVelocity = 0
    self.yVelocity = 0
    --[[
        The "maxSpeed" variable defines the player's maximum speed. The "acceleration" variable
        defines how fast the player is able to accelerate from a idling state. The "friction"
        variable defines how fast the player will slow down after no movement input is given.
    ]]
    self.maxSpeed = 200
    self.acceleration = 4000
    self.friction = 3500

    --[[
        Create physics table where all of the player's physic properties are stored. The
        "newBody()" function from LÖVE 2D creates a new physics body. We set the body to be
        non-rotating with the "setFixedRotation()" function and set it to true.
    ]]
    self.physics = {}
    self.physics.body = love.physics.newBody(World, self.x, self.y, "dynamic")
    self.physics.body:setFixedRotation(true)
    self.physics.shape = love.physics.newRectangleShape(self.width, self.height)
    self.physics.fixture = love.physics.newFixture(self.physics.body, self.physics.shape)
end

function love.update(dt)
    
end

function love.draw()
    love.graphics.rectangle()
end