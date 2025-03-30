local Slab = require("modules.slab")

local SlabGUI = {}
local windowOpen = true
local sliderValue = 50
local textInput = ""

function SlabGUI:load()
    -- Initialize Slab
    Slab.Initialize({
        DisplayFPS = true,
        MessageTime = 5,
        Font = "assets/fonts/public-pixel-font.ttf",
        FontSize = 16,
    })
end

function SlabGUI:update(dt)
    -- Call Slab's update function to handle input
    Slab.Update(dt)

    -- Create a window
    if windowOpen then
        Slab.BeginWindow("HelloWorld", {
            Title = "Hello Slab World",
            X = 100,
            Y = 100,
            W = 400,
            H = 300,
            AllowResize = true,
            AutoSizeWindow = false,
        })

        -- Add a text label
        Slab.Text("Welcome to Darling in the Goland!")
        Slab.Separator()

        -- Add a slider
        Slab.Text("Slider Value: " .. sliderValue)
        if Slab.SliderFloat("Adjust Value", sliderValue, 0, 100, {
            W = 200,
        }) then
            sliderValue = Slab.GetSliderValue()
        end
        Slab.Separator()

        -- Add a text input field
        Slab.Text("Enter your name:")
        if Slab.Input("TextInput", {
            Text = textInput,
            ReturnOnText = false,
            W = 200,
        }) then
            textInput = Slab.GetInputText()
        end
        Slab.Separator()

        -- Add a button
        if Slab.Button("Click Me!") then
            -- Show a message when clicked
            Slab.BeginTooltip()
            Slab.Text("Hello, " .. (textInput ~= "" and textInput or "Player") .. "!")
            Slab.EndTooltip()
        end

        -- Add a close button
        if Slab.Button("Close Window") then
            windowOpen = false
        end

        Slab.EndWindow()
    end

    -- Button to reopen the window if closed
    if not windowOpen then
        if Slab.Button("Open Window", {
            X = 20,
            Y = 20,
        }) then
            windowOpen = true
        end
    end
end

function SlabGUI:draw()
    -- Draw all the Slab UI elements
    Slab.Draw()
end

-- Handle key events (Slab needs this for text input)
function SlabGUI:keypressed(key, scancode, isrepeat)
    Slab.KeyPressed(key, scancode, isrepeat)
end

function SlabGUI:keyreleased(key)
    Slab.KeyReleased(key)
end

function SlabGUI:textinput(text)
    Slab.TextInput(text)
end

function SlabGUI:mousepressed(x, y, button)
    Slab.MousePressed(x, y, button)
end

function SlabGUI:mousereleased(x, y, button)
    Slab.MouseReleased(x, y, button)
end

function SlabGUI:mousemoved(x, y, dx, dy)
    Slab.MouseMoved(x, y, dx, dy)
end

return SlabGUI
