--- main.lua
-- Entry point of the Love2D game
-- package.path = "./modules/share/lua/5.4/?.lua;" .. "./modules/share/lua/5.4/?/init.lua;" .. package.path
-- package.cpath = "./modules/lib/lua/5.4/?.so;" .. package.cpath
love.graphics.setDefaultFilter("nearest", "nearest") -- Set filter to have pixel esthetic

-- @import Map, Camera, Player, Coin, HUD, Menu, DebugGUI, Spike, Stone, Enemy, Background
local Map = require("src.map.Map")
local Camera = require("src.core.Camera")

-- GUI import statements
local HUD = require("src.gui.HUD")
local Menu = require("src.gui.Menu")
local DebugGUI = require("src.gui.DebugGUI")

-- Entity import statements
local Player = require("src.entities.Player")
local Spike = require("src.entities.Spike")
local Stone = require("src.entities.Stone")
local Enemy = require("src.entities.Enemy")
local Coin = require("src.entities.Coin")

-- Visual import statements
local Background = require("src.visuals.Background")

-- math.randomseed(os.time()) -- Generate truly random numbers

--- love.load() function
-- Use to load objects and assets
function love.load()
    Enemy.loadAssets()
    Map:load()

    -- Load parallax background with 4 layers
    -- Define paths to background images (from closest to farthest, following the project convention)
    local backgroundLayers = {
        "assets/background/jungle/jungle-background-1.png", -- Closest layer (foreground elements)
        "assets/background/jungle/jungle-background-2.png", -- Middle-close layer
        "assets/background/jungle/jungle-background-3.png", -- Middle-far layer
        "assets/background/jungle/jungle-background-4.png", -- Farthest layer (sky/background)
    }

    -- Define parallax factors (smaller = slower movement)
    -- Closest layer moves fastest, farthest layer moves slowest
    local parallaxFactors = {
        0.7,
        0.5,
        0.3,
        0.1,
    }

    -- Initialize background with layers
    Background:load(backgroundLayers, parallaxFactors)

    -- Initialize camera with smooth following
    Camera:setBounds(MapWidth) -- Set camera bounds based on map width
    Camera:setSmoothing("damped", 5) -- Use damped smoothing with medium stiffness

    HUD:load()
    Menu:load()
    DebugGUI:load()
    Player:load()
end

--[[
    Each of the functions have to be called inside of `main.lua`. This is because `main.lua`
    is the entry point of every LÖVE 2D game.
]]
--- love.update(dt)
-- @param dt Delta time is used to calculate the time between the previous and current frame
function love.update(dt)
    if not Menu.active then -- Only update game if menu is not active
        World:update(dt)
        Player:update(dt)
        Coin.updateAll(dt)
        Spike.updateAll(dt)
        Stone.updateAll(dt)
        Enemy.updateAll(dt)
        HUD:update(dt)

        -- Update camera to smoothly follow the player
        Camera:follow(Player.x, 0)
        Camera:update(dt)

        Map:update(dt)

        -- Update background with camera position
        Background:update(dt, Camera.x)

        -- Get entity counts directly using the # operator on the internal tables
        DebugGUI:updateEntityCount(Coin.getCount(), Enemy.getCount(), Spike.getCount(), Stone.getCount())
    end

    Menu:update(dt) -- Always update menu
    DebugGUI:update(dt)
end

function love.draw()
    -- Draw parallax background
    Background:draw()

    -- Draw map
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
    HUD:draw()
    Menu:draw()
    DebugGUI:draw()
    -- love.graphics.printf("Hello World", 200, 300, 420, "justify")
end

--[[
    This function is responsible for listening for input that activates jumping, which is a
    callback function that only needs to know when a specific key is pressed.
]]
function love.keypressed(key)
    Player:jump(key)

    if key == "escape" then
        Menu:toggle() -- Toggle menu instead of quitting
    end

    if key == "f3" then -- Common debug toggle key
        DebugGUI:toggle()
    end
end

--[[
    Add mouse press handler for menu button interaction
]]
function love.mousepressed(x, y, button)
    Menu:mousepressed(x, y, button)
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

