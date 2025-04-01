-- src/Player.lua
local Player = {} -- Create "Player" class table (OOP)

function Player:load()
    self.x = 100
    self.y = 0
    self.startX = self.x -- Start position x after player died
    self.startY = self.y -- Start position y after player died
    self.width = 20
    self.height = 60
    self.xVelocity = 0
    self.yVelocity = 100

    -- Player states
    self.isAlive = true
    self.isGrounded = false
    self.canDoubleJump = true
    self.animState = "idle"

    -- Player movement
    self.jumpForce = -650
    self.playerDirection = "right"
    --[[
        The "maxSpeed" variable defines the player's maximum speed. The "acceleration" variable
        defines how fast the player is able to accelerate from a idling state. The "friction"
        variable defines how fast the player will slow down after no movement input is given
        by the user.
    ]]
    self.maxSpeed = 200
    self.acceleration = 1500
    self.friction = 3500
    self.gravity = 1500

    self.coinCount = 0
    self.health = {
        current = 3,
        max = 3
    }

    self.color = {
        red = 1,
        green = 1,
        blue = 1,
        untintSpeed = 3
    }

    --[[
        These variables are responsible for making the player jump if he was grounded recently.
        If the player is unable to jump the next frame he isn't grounded, doesn't feel that
        smooth. The "jumpTimeFrame" is the time the player is able to jump after not being
        grounded.
    ]]
    self.jumpTimeFrame = 0
    self.timeFrameDuration = 0.25

    self:loadAssets()

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

    -- self:takeDamage(1)
    -- self:takeDamage(1)
end

--[[
    The player has 3 states: idle, run and air. This function is responsible for animating these
    3 states.
]]
function Player:loadAssets()
    self.animation = {
        timer = 0,
        rate = 0.1
    }

    -- Run player animation
    self.animation.run = {
        totalFrames = 6,
        currentFrame = 1,
        images = {}
    }
    for i = 1, self.animation.run.totalFrames do
        self.animation.run.images[i] = love.graphics.newImage("assets/player/run/zero-two-run-" .. i .. ".png")
    end

    -- Idle player animation
    self.animation.idle = {
        totalFrames = 4,
        currentFrame = 1,
        images = {}
    }
    for i = 1, self.animation.idle.totalFrames do
        self.animation.idle.images[i] = love.graphics.newImage("assets/player/idle/zero-two-idle-" .. i .. ".png")
    end

    -- Air player animation
    self.animation.air = {
        totalFrames = 4,
        currentFrame = 1,
        images = {}
    }
    for i = 1, self.animation.air.totalFrames do
        self.animation.air.images[i] = love.graphics.newImage("assets/player/idle/zero-two-idle-" .. i .. ".png")
    end

    --[[
        This variable will store the current animation that we want to draw. We also get the
        width and height of our player assets.
    ]]
    self.animation.draw = self.animation.idle.images[1]
    self.animation.width = self.animation.draw:getWidth()
    self.animation.height = self.animation.draw:getHeight()
end

function Player:takeDamage(amount)
    self:tintRed()

    -- Check if current health status subtracted by the amount taken is larger than 0
    if self.health.current - amount > 0 then
        -- if thats the case subtract the amount taken
        self.health.current = self.health.current - amount
    else
        self.health.current = 0
        self:kill()
    end
    print("Player health: ", self.health.current)
end

--[[
    This function is responsible for killing the player if the health drops below
    the specified value in the health table.
]]
function Player:kill()
    print("Player has been killed.")
    self.isAlive = false
end

--[[
    This function is responsible for respawning the player after the got killed.
]]
function Player:respawn()
    -- Check if player is not alive
    if not self.isAlive then
        self:resetPosition()
        -- Reset the health of the player to the specified max health value
        self.health.current = self.health.max
        self.isAlive = true
    end
end

function Player:resetPosition()
    -- If thats the case position the player back to the specified starting position
    self.physics.body:setPosition(self.startX, self.startY)
end

function Player:tintRed()
    self.color.green = 0
    self.color.blue = 0
end

function Player:untintRed(dt)
    self.color.red = math.min(self.color.red + self.color.untintSpeed * dt, 1)
    self.color.green = math.min(self.color.green + self.color.untintSpeed * dt, 1)
    self.color.blue = math.min(self.color.blue + self.color.untintSpeed * dt, 1)
end

function Player:incrementCoinCount()
    self.coinCount = self.coinCount + 1
end

function Player:update(dt)
    self:untintRed(dt)
    self:respawn()
    self:setAnimState()
    self:setPlayerDirection()
    self:animate(dt)
    self:decreaseTimeFrame(dt)
    self:syncPhysics()
    self:movement(dt)
    self:applyGravity(dt)
end

--[[
    This function is responsible for determining the animation state of the player. The
    player can either be in air, idle or running.
]]
function Player:setAnimState()
    if not self.isGrounded then
        self.animState = "air"
    elseif self.xVelocity == 0 then
        self.animState = "idle"
    else
        self.animState = "run"
    end

    -- print(self.animState)
end

--[[
    This function is responsible for setting the player sprite in the correct direction. When
    "D" is pressed the player sprite faces right and when "A" is pressed the player sprite
    faces left.
]]
function Player:setPlayerDirection()
    if self.xVelocity < 0 then
        self.playerDirection = "left"
    elseif self.xVelocity > 0 then
        self.playerDirection = "right"
    end
end

--[[
    This function is responsible for animating the player images
]]
function Player:animate(dt)
    self.animation.timer = self.animation.timer + dt

    if self.animation.timer > self.animation.rate then
        self.animation.timer = 0
        self:setNewFrame()
    end
end

--[[
    This function is responsible for updating the different player images, to create
    animation
]]
function Player:setNewFrame()
    local anim = self.animation[self.animState]

    if anim.currentFrame < anim.totalFrames then
        anim.currentFrame = anim.currentFrame + 1
    else
        anim.currentFrame = 1
    end

    self.animation.draw = anim.images[anim.currentFrame]
end

--[[
    This function is responsible for decreasing the time frame the player has to activate a
    jump after not being grounded (touching ground) for a specific amount of time, which is
    stored in the "timeFrameDuration" variable.
]]
function Player:decreaseTimeFrame(dt)
    if not self.isGrounded then self.jumpTimeFrame = self.jumpTimeFrame - dt end
end

--[[
    This function is responsible for applying gravity to the world.
]]
function Player:applyGravity(dt)
    if not self.isGrounded then self.yVelocity = self.yVelocity + self.gravity * dt end
end

--[[
    This function is responsible for the movement of the player. Note that for the player to be
    able to move in the left direction we have to subtract friction from the xVelocity.
    The "applyFriction()" function is only being called if no input is given by the user,
    this means the player will come to a standstill.
]]
function Player:movement(dt)
    if love.keyboard.isDown("d", "right") then
        self.xVelocity = math.min(self.xVelocity + self.acceleration * dt, self.maxSpeed)
        -- print("xVelocity:", math.floor(self.xVelocity))
        -- if self.xVelocity < self.maxSpeed then
        --     if self.xVelocity + self.acceleration * dt < self.maxSpeed then
        --         self.xVelocity = self.xVelocity + self.acceleration * dt
        --         print("xVelocity:", math.floor(self.xVelocity))
        --     else
        --         self.xVelocity = self.maxSpeed
        --     end
        -- end
    elseif love.keyboard.isDown("a", "left") then
        self.xVelocity = math.max(self.xVelocity - self.acceleration * dt, -self.maxSpeed)
        -- print("xVelocity:", math.floor(self.xVelocity))
        -- if self.xVelocity > -self.maxSpeed then
        --     if self.xVelocity - self.acceleration * dt > -self.maxSpeed then
        --         self.xVelocity = self.xVelocity - self.acceleration * dt
        --         print("xVelocity:", math.floor(self.xVelocity))
        --     else
        --         self.xVelocity = -self.maxSpeed
        --     end
        -- end
    else
        Player:applyFriction(dt)
    end
end

--[[
    This function is responsible for bringing the player to a standstill. Without this function
    the player would be constantly moving. The if-statements check, that if the velocity on the
    x-axis is greater than 0 that it subtracts the specified value, which is stored in the
    "friction" variable, from the xVelocity so that the player comes to a standstill.
]]
function Player:applyFriction(dt)
    -- If the player is moving in the right direction
    if self.xVelocity > 0 then
        self.xVelocity = math.max(self.xVelocity - self.friction * dt, 0)
        -- and the xVel is still higher even after adding friction
        -- if self.xVelocity - self.friction * dt > 0 then
        --     -- then apply the specified friction amount
        --     self.xVelocity = self.xVelocity - self.friction * dt
        -- else
        --     -- else set xVel to 0
        --     self.xVelocity = 0
        -- end
        -- If the player is moving the in the left direction
    elseif self.xVelocity < 0 then
        self.xVelocity = math.min(self.xVelocity + self.friction * dt, 0)
        -- Do the same thing as in the first if-statement
        -- if self.xVelocity + self.friction * dt < 0 then
        --     self.xVelocity = self.xVelocity + self.friction * dt
        -- else
        --     self.xVelocity = 0
        -- end
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
    These function are responsible for the player-side of collision
]]
function Player:beginContact(fixtureA, fixtureB, collision)
    -- If the player is already touching ground then skip the whole function
    if self.isGrounded == true then return end

    local nx, ny = collision:getNormal()

    if fixtureA == self.physics.fixture then
        if ny > 0 then
            self:land(collision)
            -- Check if ny < 0, then prevent player from getting stuck under collidable object
        elseif ny < 0 then
            self.yVelocity = 0
        end
    elseif fixtureB == self.physics.fixture then
        if ny < 0 then
            self:land(collision)
        elseif ny > 0 then
            -- Check if ny < 0, then prevent player from getting stuck under collidable object
            self.yVelocity = 0
        end
    end
end

function Player:land(collision)
    self.currentGroundCollision = collision
    self.yVelocity = 0
    self.isGrounded = true
    self.canDoubleJump = true
    self.jumpTimeFrame = self.timeFrameDuration
end

--[[
    This function is responsible for making the player jump, which is also defined in the
    callback function in main.lua (love.keypressed). "jumpTimeFrame" will allow the player to
    jump even if the entity isn't grounded.
]]
function Player:jump(key)
    if key == "w" or key == "up" then
        if self.isGrounded or self.jumpTimeFrame > 0 then
            self.yVelocity = self.jumpForce
            self.isGrounded = false
            self.jumpTimeFrame = 0
        elseif self.canDoubleJump then
            self.canDoubleJump = false
            self.yVelocity = self.jumpForce * 0.75
        end
    end
end

function Player:endContact(fixtureA, fixtureB, collision)
    if fixtureA == self.physics.fixture or fixtureB == self.physics.fixture then
        if self.currentGroundCollision == collision then self.isGrounded = false end
    end
end

--[[
    When creating rectangles we have to off-set the values, due to the origin points being
    calculated differently in Box2D and LÖVE 2D. In Box2D the origin point is in the center
    middle, while the origin point in LÖVE 2D is at the top left corner.
]]
function Player:draw()
    local scaleX = 1
    if self.playerDirection == "left" then scaleX = -1 end
    -- love.graphics.rectangle("fill", self.x - self.width / 2, self.y - self.height / 2, self.width, self.height)
    love.graphics.setColor(self.color.red, self.color.green, self.color.blue) -- Draw player color
    love.graphics.draw(self.animation.draw, self.x, self.y, 0, scaleX, 1, self.animation.width / 2, self.animation.height / 2)
    love.graphics.setColor(1, 1, 1, 1)
end

return Player
