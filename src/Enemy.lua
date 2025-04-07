--- src/Enemy.lua
-- @class Enemy
-- Enemy that deals damage if collides with player
local Enemy = {}
local ActiveEnemies = {}
Enemy.__index = Enemy

-- @import Player
local Player = require("src.Player")

-- Basically a getter method
function Enemy:getCount()
    return #ActiveEnemies
end

function Enemy.removeAll()
    for i, v in ipairs(ActiveEnemies) do v.physics.body:destroy() end

    ActiveEnemies = {}
end

function Enemy.new(x, y)
    local enemy = setmetatable({}, Enemy)
    enemy.x = x
    enemy.y = y
    enemy.offsetY = -8
    enemy.rotation = 0

    enemy.speed = 100
    enemy.speedMultiplier = 1
    enemy.xVel = enemy.speed

    enemy.rageCounter = 0
    enemy.rageTrigger = 3

    enemy.damage = 1

    enemy.state = "walk" -- Default state

    enemy.animation = {
        timer = 0,
        rate = 0.1
    }
    enemy.animation.run = {
        total = 4,
        current = 1,
        img = Enemy.runAnim
    }
    enemy.animation.walk = {
        total = 4,
        current = 1,
        img = Enemy.walkAnim
    }
    enemy.animation.draw = enemy.animation.walk.img[1]

    enemy.physics = {}
    enemy.physics.body = love.physics.newBody(World, enemy.x, enemy.y, "dynamic")
    enemy.physics.body:setFixedRotation(true)
    enemy.physics.shape = love.physics.newRectangleShape(enemy.width * 0.4, enemy.height * 0.75)
    enemy.physics.fixture = love.physics.newFixture(enemy.physics.body, enemy.physics.shape)
    enemy.physics.body:setMass(25)
    table.insert(ActiveEnemies, enemy)
end

function Enemy.loadAssets()
    Enemy.runAnim = {}
    for i = 1, 4 do Enemy.runAnim[i] = love.graphics.newImage("assets/enemy/run/" .. i .. ".png") end

    Enemy.walkAnim = {}
    for i = 1, 4 do Enemy.walkAnim[i] = love.graphics.newImage("assets/enemy/walk/" .. i .. ".png") end

    Enemy.width = Enemy.runAnim[1]:getWidth()
    Enemy.height = Enemy.runAnim[1]:getHeight()
end

function Enemy:update(dt)
    self:syncPhysics()
    self:animate(dt)
end

function Enemy:incrementRage()
    self.rageCounter = self.rageCounter + 1
    if self.rageCounter > self.rageTrigger then
        self.state = "run"
        self.speedMultiplier = 3
        self.rageCounter = 0
    else
        self.state = "walk"
        self.speedMultiplier = 1
    end
end

function Enemy:flipDirection()
    self.xVel = -self.xVel
end

function Enemy:animate(dt)
    self.animation.timer = self.animation.timer + dt
    if self.animation.timer > self.animation.rate then
        self.animation.timer = 0
        self:setNewFrame()
    end
end

function Enemy:setNewFrame()
    local anim = self.animation[self.state]
    if anim.current < anim.total then
        anim.current = anim.current + 1
    else
        anim.current = 1
    end
    self.animation.draw = anim.img[anim.current]
end

function Enemy:syncPhysics()
    self.x, self.y = self.physics.body:getPosition()
    self.physics.body:setLinearVelocity(self.xVel * self.speedMultiplier, 100)
end

function Enemy:draw()
    local scaleX = 1
    if self.xVel < 0 then scaleX = -1 end
    love.graphics.draw(self.animation.draw, self.x, self.y + self.offsetY, self.rotation, scaleX, 1, self.width / 2, self.height / 2)
end

function Enemy.updateAll(dt)
    for i, enemy in ipairs(ActiveEnemies) do enemy:update(dt) end
end

function Enemy.drawAll()
    for i, enemy in ipairs(ActiveEnemies) do enemy:draw() end
end

function Enemy.removeAll()
    for i, v in ipairs(ActiveEnemies) do v.physics.body:destroy() end

    ActiveEnemies = {}
end

function Enemy.beginContact(a, b, collision)
    for i, enemy in ipairs(ActiveEnemies) do
        if a == enemy.physics.fixture or b == enemy.physics.fixture then
            if a == Player.physics.fixture or b == Player.physics.fixture then Player:takeDamage(enemy.damage) end
            enemy:incrementRage()
            enemy:flipDirection()
        end
    end
end

return Enemy
