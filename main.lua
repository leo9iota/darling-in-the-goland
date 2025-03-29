-- package.path = "./modules/share/lua/5.4/?.lua;" .. "./modules/share/lua/5.4/?/init.lua;" .. package.path
-- package.cpath = "./modules/lib/lua/5.4/?.so;" .. package.cpath

love.graphics.setDefaultFilter("nearest", "nearest") -- Set filter to have pixel esthetic

local Map = require("src.Map")
local Camera = require("src.Camera")
local Player = require("src.Player")
local Coin = require("src.Coin")
local GUI = require("src.GUI")
local Spike = require("src.Spike")
local Stone = require("src.Stone")
local Enemy = require("src.Enemy")

-- math.randomseed(os.time()) -- Generate truly random numbers

function love.load()
    Enemy.loadAssets()
    Map:load()

    background = love.graphics.newImage("assets/world/background.png")
    GUI:load()

    Player:load()
end

--[[
    Each of the functions have to be called inside of `main.lua`. This is because `main.lua`
    is the entry point of every LÖVE 2D game. 
]]
function love.update(dt)
    World:update(dt)
    Player:update(dt)
    Coin.updateAll(dt)
    Spike.updateAll(dt)
    Stone.updateAll(dt)
    Enemy.updateAll(dt)
    GUI:update(dt)
    Camera:setPosition(Player.x, 0)
    Map:update(dt)
end

function love.draw()
    love.graphics.draw(background)
    Map.level:draw(-Camera.x, -Camera.y, Camera.scale, Camera.scale)

    Camera:init()
    --[[
        The "push()" adds the scaling to the stack and the "pop()" function removes the scaling
        from the stack, which means everything outside of these functions will not be affected
        by the scaling
    ]]
    -- love.graphics.push()
    -- love.graphics.scale(2, 2)

    Player:draw()
    Coin.drawAll()
    Spike.drawAll()
    Stone.drawAll()
    Enemy.drawAll()
    Camera:remove()

    -- love.graphics.pop()

    --[[
        The GUI function is called outside of the scaling functions:
        - push()
        - scale()
        - pop()
        
        because the GUI should be static.
    ]]
    GUI:draw()
    -- love.graphics.printf("Hello World", 200, 300, 420, "justify")
end

--[[
    This function is responsible for listening for input that activates jumping, which is a
    callback function that only needs to know when a specific key is pressed.
]]
function love.keypressed(key)
    Player:jump(key)

    if key == "escape" then
        love.event.quit()
    end
end

--[[
    If the player collects a coin we skip the collision detection for the ground etc.
    
    --- IMPORTANT ---
    You're not allowed to make any changes to the physics World inside of any of these
    callbacks. This is due to the fact that Box2D "locks" the world.

    The workaround is to mark the object outside of the callback and then remove it.
]]
function beginContact(fixtureA, fixtureB, collision)
    if Coin.beginContact(fixtureA, fixtureB, collision) then
        return
    end
    if Spike.beginContact(fixtureA, fixtureB, collision) then
        return
    end
    Enemy.beginContact(fixtureA, fixtureB, collision) -- We want the collision callback fn to continue to the code that grounds the player
    Player:beginContact(fixtureA, fixtureB, collision)
end

function endContact(fixtureA, fixtureB, collision)
    Player:endContact(fixtureA, fixtureB, collision)
end

