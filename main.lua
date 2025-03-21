love.graphics.setDefaultFilter("nearest", "nearest") -- Set filter to have pixel esthetic

-- Import STI library to import maps from Tiled
local STI = require("sti")

local Camera = require "Camera"
local Player = require "Player"
local Coin = require "Coin"
local GUI = require "GUI"
local Spike = require "Spike"
local Stone = require "Stone"
local Enemy = require "Enemy"

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

    Map.layers.solid.visible = false -- Hide "solid" layer from Tiled
    Map.layers.entity.visible = false -- Hide "entity" layer from Tiled

    MapWidth = Map.layers.ground.width * tileSizeInPixels -- Prevent camera to go out of bounds on right side

    background = love.graphics.newImage("assets/background.png")
    GUI:load()

    Enemy.loadAssets()

    Player:load()

    spawnEntities()
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
    Enemy.beginContact(fixtureA, fixtureB, collision) -- We want the collision callback fn to continue to the code that grounds the player
    Player:beginContact(fixtureA, fixtureB, collision)
end

function endContact(fixtureA, fixtureB, collision)
    Player:endContact(fixtureA, fixtureB, collision)
end

--[[ 
    Loop through "object" table to get all the entities.
    NOTE: Different origin points for circles and rectangles in Tiled.
        - Rectangle: top left
        - Circle: center
]]
function spawnEntities()
    for i, v in ipairs(Map.layers.entity.objects) do
        --[[ 
            NOTE: Since the latest update in Tiled, the field isn't called "type" anymore, its called "class", but since 
            I use the old version of the STI lib, its still "type", or until I update the lib.
        ]]
        if v.type == "spikes" then
            Spike.new(v.x + v.width / 2, v.y + v.height / 2) -- The origin point in Tiled is the top left corner, but origin point of the physics module is the center, which means we need an offset 
        elseif v.type == "stone" then
            Stone.new(v.x + v.width / 2, v.y + v.height / 2)
        elseif v.type == "enemy" then
            Enemy.new(v.x + v.width / 2, v.y + v.height / 2)
        elseif v.type == "coin" then
            Coin.new(v.x, v.y) -- In Tiled a circle's origin point is in the center, so we can just pass x and y without offsets
        end
    end
end
