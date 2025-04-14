--- src/gui/Menu.lua
-- @class Menu
-- Handles the game menu overlay

local Menu = {}

function Menu:load()
    self.active = false

    -- Set up menu appearance
    self.width = 400
    self.height = 300
    self.x = (love.graphics.getWidth() - self.width) / 2
    self.y = (love.graphics.getHeight() - self.height) / 2

    -- Background color
    self.bgColor = {
        0.2,
        0.2,
        0.2,
        0.9,
    }

    -- Button properties
    self.buttonHeight = 60
    self.buttonWidth = 300
    self.buttonMargin = 20

    -- Create buttons
    self.buttons = {
        {
            text = "Resume",
            x = self.x + (self.width - self.buttonWidth) / 2,
            y = self.y + 50,
            width = self.buttonWidth,
            height = self.buttonHeight,
            action = function()
                self:toggle()
            end,
        },
        {
            text = "Settings",
            x = self.x + (self.width - self.buttonWidth) / 2,
            y = self.y + 50 + self.buttonHeight + self.buttonMargin,
            width = self.buttonWidth,
            height = self.buttonHeight,
            action = function()
                print("Settings menu (not implemented)")
            end,
        },
        {
            text = "Quit",
            x = self.x + (self.width - self.buttonWidth) / 2,
            y = self.y + 50 + (self.buttonHeight + self.buttonMargin) * 2,
            width = self.buttonWidth,
            height = self.buttonHeight,
            action = function()
                love.event.quit()
            end,
        },
    }

    -- Font for the menu
    self.font = love.graphics.newFont("assets/fonts/public-pixel-font.ttf", 24)
end

function Menu:update(dt)
    if not self.active then
        return
    end

    -- Get mouse position
    local mx, my = love.mouse.getPosition()

    -- Check button hover state
    for i, button in ipairs(self.buttons) do
        button.hover = mx >= button.x and mx <= button.x + button.width and my >= button.y and my <= button.y + button.height
    end
end

function Menu:draw()
    if not self.active then
        return
    end

    -- Store original font to restore later
    local originalFont = love.graphics.getFont()
    love.graphics.setFont(self.font)

    -- Draw semi-transparent background
    love.graphics.setColor(self.bgColor)
    love.graphics.rectangle("fill", self.x, self.y, self.width, self.height)

    -- Draw border
    love.graphics.setColor(1, 1, 1, 1)
    love.graphics.rectangle("line", self.x, self.y, self.width, self.height)

    -- Draw title
    love.graphics.setColor(1, 1, 1, 1)
    local title = "Menu"
    local titleWidth = self.font:getWidth(title)
    love.graphics.print(title, self.x + (self.width - titleWidth) / 2, self.y + 10)

    -- Draw buttons
    for i, button in ipairs(self.buttons) do
        -- Button background
        if button.hover then
            love.graphics.setColor(0.4, 0.4, 0.4, 1)
        else
            love.graphics.setColor(0.3, 0.3, 0.3, 1)
        end

        love.graphics.rectangle("fill", button.x, button.y, button.width, button.height)

        -- Button border
        love.graphics.setColor(1, 1, 1, 1)
        love.graphics.rectangle("line", button.x, button.y, button.width, button.height)

        -- Button text
        love.graphics.setColor(1, 1, 1, 1)
        local textWidth = self.font:getWidth(button.text)
        local textHeight = self.font:getHeight()
        love.graphics.print(button.text, button.x + (button.width - textWidth) / 2, button.y + (button.height - textHeight) / 2)
    end

    -- Restore original font
    love.graphics.setFont(originalFont)
end

function Menu:toggle()
    self.active = not self.active
end

function Menu:mousepressed(x, y, button)
    if not self.active or button ~= 1 then
        return
    end

    -- Check if any buttons were clicked
    for i, btn in ipairs(self.buttons) do
        if x >= btn.x and x <= btn.x + btn.width and y >= btn.y and y <= btn.y + btn.height then
            btn.action()
            return
        end
    end
end

return Menu
