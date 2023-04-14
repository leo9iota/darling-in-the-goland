Player = {}

function Player:load()
    self.x = 100
    self.y = 0
    self.width = 20
    self.height = 60
    self.xVelocity = 0
    self.yVelocity = 100
    --[[
        The "maxSpeed" variable defines the player's maximum speed. The "acceleration" variable
        defines how fast the player is able to accelerate from a idling state. The "friction"
        variable defines how fast the player will slow down after no movement input is given.
    ]]
    self.maxSpeed = 200
    self.acceleration = 2500
    self.friction = 3500
    self.gravity = 1500

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

function Player:update(dt)
    Player:syncPhysics()
    Player:movement(dt)
end

--[[
    This function is responsible for the movement of the player. Note that for the player to
    be moving in the left direction we have to subtract from the xVelocity. The "applyFriction()"
    function is only being called if no input is given by the user, this means the player will
    come to a standstill.
]]
function Player:movement(dt)
    if love.keyboard.isDown("d", "right") then
        if self.xVelocity < self.maxSpeed then
            if self.xVelocity + self.acceleration * dt < self.maxSpeed then
                self.xVelocity = self.xVelocity + self.acceleration * dt
                print(self.xVelocity)
            else
                self.xVelocity = self.maxSpeed
            end
        end
    elseif love.keyboard.isDown("a", "left") then
        if self.xVelocity > -self.maxSpeed then
            if self.xVelocity - self.acceleration * dt > -self.maxSpeed then
                self.xVelocity = self.xVelocity - self.acceleration * dt
                print(self.xVelocity)
            else
                self.xVelocity = -self.maxSpeed
            end
        end
    else
        Player:applyFriction(dt)
    end
end

--[[
    This function is responsible for bringing the player to a standstill. Without this function
    the player would be constantly move. The if-statements checks that if the velocity on the
    x-axis is greater than 0 that is subtracts specified, which is stored in the "friction"
    variable from the xVelocity, so that the player is able to come to a standstill.
]]
function Player:applyFriction(dt)
    if self.xVelocity > 0 then
        if self.xVelocity - self.friction * dt > 0 then
            self.xVelocity = self.xVelocity - self.friction * dt
        else
            self.xVelocity = 0
        end
    elseif self.xVelocity < 0 then
        if self.xVelocity + self.friction * dt < 0 then
            self.xVelocity = self.xVelocity + self.friction * dt
        else
            self.xVelocity = 0
        end
    end
end

--[[
    This function synchronizes the physical body with the player's x and y position. In Lua,
    you are able to set multiple values to be equal to the returned values from one function.
]]
function Player:syncPhysics()
    -- self.x = self.physics.body:getX()
    -- self.y = self.physics.body:getY()
    self.x, self.y = self.physics.body:getPosition()
    self.physics.body:setLinearVelocity(self.xVelocity, self.yVelocity)
end

--[[
    When creating rectangles we have to off-set the values, due to the origin points being
    calculated differently in Box2D and LÖVE 2D. In Box2D the origin point is in the center
    middle, while the origin point in LÖVE 2D is at the top left corner.
]]
function Player:draw()
    love.graphics.rectangle("fill", self.x - self.width / 2, self.y - self.height / 2, self.width, self.height)
end