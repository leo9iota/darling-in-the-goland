local Map = {}

local STI = require "sti" -- Import STI library to import maps from Tiled
local Player = require "Player"
local Coin = require "Coin"
local Spike = require "Spike"
local Stone = require "Stone"
local Enemy = require "Enemy"

local TILE_SIZE = 16 -- Constant for the tile size in pixels

function Map:load()
    self.currentLevel = 1 -- Variable for storing current level

    --[[
        The STI library relies on the open-source Box2D physics engine. With the "newWorld()"
        function call we get gravity. The "initBox2D()" function loads all layers and
        objects that have the custom property "collidable" into the new world.
    ]]
    World = love.physics.newWorld(0, 2000)
    --[[
        This function takes 4 functions as arguments but the last two are optional. The
        "beginContact" callback function gets called when two fixtures collide and the
        "endContact" callback functions gets called when two fixtures that we are calling
        no longer collide.
    ]]
    World:setCallbacks(beginContact, endContact)

    self:init()
    self:next()
end

function Map:init()
    self.level = STI("maps/map-1.lua", {"box2d"}) -- FIX: Load level first
    self.level:initBox2D(World)

    self.solidLayer = self.level.layers.solid -- Var for Tiled solid layer
    self.groundLayer = self.level.layers.ground -- Var for Tiled ground layer
    self.entityLayer = self.level.layers.entity -- Var for Tiled entity layer
    self.solidLayer.visible = false -- Hide "solid" layer from Tiled
    self.entityLayer.visible = false -- Hide "entity" layer from Tiled 

    MapWidth = self.groundLayer.width * TILE_SIZE -- Prevent camera to go out of bounds on right and left side

    self:spawnEntities()
end

function Map:next()
    self:clean()
    self.currentLevel = self.currentLevel + 1
    self:init()
end

function Map:clean()
    self.level:removeLayerBox2D("solid")
end

--[[ 
    Loop through "object" table to get all the entities.
    NOTE: Different origin points for circles and rectangles in Tiled.
    - Rectangle: top left
    - Circle: center
]]
--- This function spawns all game entities 
function Map:spawnEntities()
    for i, v in ipairs(self.entityLayer.objects) do
        --[[ 
            NOTE: Since the latest update in Tiled, the field isn't called "type" anymore, its called "class", but since 
            I use the old version of the STI lib, its still "type", or until I update the lib.

            FIXME: Could cause potential bug in the future if ever update to a newer version of the STI lib.
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

return Map
