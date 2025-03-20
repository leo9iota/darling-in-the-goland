love.graphics.setDefaultFilter("nearest", "nearest") -- Set filter to have pixel esthetic

-- Import STI library to import maps from Tiled
local STI = require("sti")

local Player = require("Player")
local Coin = require("Coin")
local GUI = require("GUI")
local Spike = require("Spike")
local Stone = require("Stone")
local Camera = require("Camera")

-- math.randomseed(os.time()) -- Generate truly random numbers

function love.load()
    local tileSizeInPixels = 16

    --[[
        The STI library relies on the open-source Box2D physics engine. With the "newWorld()"
        function call we get gravity. The "initBox2D()" function loads all layers and
        objects that have the custom property "collidable" into the new world.
    ]]
    Map = STI("maps/map-1.lua", {"box2d"})
    World = love.physics.newWorld(0, 2000)

    --[[
        This function takes 4 functions as arguments but the last two are optional. The
        "beginContact" callback function gets called when two fixtures collide and the
        "endContact" callback functions gets called when two fixtures that we are calling
        no longer collide.
    ]]
    World:setCallbacks(beginContact, endContact)

    Map:initBox2D(World)
    Map.layers.solid.visible = false

    MapWidth = Map.layers.ground.width * tileSizeInPixels -- Prevent camera to go out of bounds on right side

    background = love.graphics.newImage("assets/background.png")
    GUI:load()
    Player:load()

    Coin.new(160, 180)
    Coin.new(320, 150)
    Coin.new(370, 150)

    Spike.new(495, 305)
    Spike.new(460, 305)
    Spike.new(425, 305)
    Spike.new(390, 305)
    Spike.new(355, 305)
end

--[[
    Each of the functions have to be called inside of `main.lua`. This is because `main.lua`
    is the entry point of every LÖVE 2D game. 
]]
function love.update(dt)
    World:update(dt)
    Player:update(dt)
    Coin.updateAllCoins(dt)
    Spike.updateAllSpikes(dt)
    GUI:update(dt)
    Camera:setPosition(Player.x, 0)
end

function love.draw()
    love.graphics.draw(background)
    Map:draw(-Camera.x, -Camera.y, Camera.scale, Camera.scale)

    Camera:init()
    --[[
        The "push()" adds the scaling to the stack and the "pop()" function removes the scaling
        from the stack, which means everything outside of these functions will not be affected
        by the scaling
    ]]
    -- love.graphics.push()
    -- love.graphics.scale(2, 2)

    Player:draw()
    Coin.drawAllCoins()
    Spike.drawAllSpikes()

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

    if key == "escape" then love.event.quit() end
end

--[[
    If the player collects a coin we skip the collision detection for the ground etc.
    
    --- IMPORTANT ---
    You're not allowed to make any changes to the physics World inside of any of these
    callbacks. This is due to the fact that Box2D "locks" the world.

    The workaround is to mark the object outside of the callback and then remove it.
]]
function beginContact(fixtureA, fixtureB, collision)
    if Coin.beginContact(fixtureA, fixtureB, collision) then return end
    if Spike.beginContact(fixtureA, fixtureB, collision) then return end
    Player:beginContact(fixtureA, fixtureB, collision)
end

function endContact(fixtureA, fixtureB, collision)
    Player:endContact(fixtureA, fixtureB, collision)
end
