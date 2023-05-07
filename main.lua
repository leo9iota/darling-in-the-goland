-- Import STI library to import maps from Tiled
local STI = require("sti")

require("Player")
require("Coin")

love.graphics.setDefaultFilter("nearest", "nearest") -- Set filter to have pixel esthetic
-- math.randomseed(os.time()) -- Generate truly random numbers

function love.load()
    --[[
        The STI library relies on the open-source Box2D physics engine. With the "newWorld()"
        function call we get gravity. The "initBox2D()" function loads all layers and
        objects that have the custom property "collidable" into the new world.
    ]]
    Map = STI("maps/map-1.lua", { "box2d" })
    World = love.physics.newWorld(0, 0)
    
    --[[
        This function takes 4 functions as arguments but the last two are optional. The
        "beginContact" callback function gets called when two fixtures collide and the
        "endContact" callback functions gets called when two fixtures that were colling
        no longer collide.
    ]]
    World:setCallbacks(beginContact, endContact)

    Map:initBox2D(World)
    Map.layers.Solid.visible = false
    background = love.graphics.newImage("assets/background.png")
    Player:load()
    Coin.new(160, 180)
    Coin.new(320, 150)
    Coin.new(370, 150)
end

function love.update(dt)
    World:update(dt)
    Player:update(dt)
    Coin:updateAllCoins(dt)
end

function love.draw()
    love.graphics.draw(background)
    Map:draw(0, 0, 2, 2)

    --[[
        The "push()" adds the scaling to the stack and the "pop()" function removes the scaling
        from the stack, which means everything outside of these functions will not be affected
        by the scaling
    ]]
    love.graphics.push()
    love.graphics.scale(2, 2)
    Player:draw()
    Coin:drawAllCoins()
    love.graphics.pop()
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
    if Coin.beginContact(fixtureA, fixtureB, collision) then return end
    Player:beginContact(fixtureA, fixtureB, collision)
end

function endContact(fixtureA, fixtureB, collision)
   Player:endContact(fixtureA, fixtureB, collision) 
end
