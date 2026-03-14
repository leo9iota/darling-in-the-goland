-- conf.lua
function love.conf(t)
    t.title = "Darling in the Goland"
    t.version = "11.5"
    t.console = true
    t.window.width = 1280
    t.window.height = 720
    t.window.vsync = 0
    --[[ 
        Allow resizing of game window, if the window is resizable, you can set minimum dimensions:
    ]]
    -- t.window.resizable = false
    -- t.window.minwidth = 800
    -- t.window.minheight = 600
    --[[ 
        Set to fullscreen
    ]]
    t.window.fullscreen = false
    --[[ 
        Set a custom icon for the game window
    ]]
    -- t.window.icon = "assets/icons/game-icon.png"
    --[[ 
        Disable unused modules to optimize performance
    ]]
    -- t.modules.joystick = false
    -- t.modules.physics = true
    -- t.modules.audio = true
    -- t.modules.video = false
end
