--- src/gui/DebugGUI.lua
-- @class DebugGUI
-- Display metrics, such as fps, etc., for debugging performance bottlenecks
local DebugGUI = {}

-- Default state is hidden
DebugGUI.active = true -- Set to true by default for testing

-- Metrics storage
DebugGUI.fps = 0
DebugGUI.memoryUsage = 0
DebugGUI.drawCalls = 0
DebugGUI.entityCounts = {
    coins = 0,
    enemies = 0,
    spikes = 0,
    stones = 0,
    total = 0,
}

function DebugGUI:load()
    -- Set up font for debug info
    self.font = love.graphics.newFont(14)

    -- Color for debug text (yellow with slight transparency)
    self.textColor = {
        1,
        1,
        0,
        0.75,
    }

    -- Background color for debug panel
    self.bgColor = {
        0,
        0,
        0,
        0.6,
    }
end

function DebugGUI:update(dt)
    if not self.active then
        return
    end

    -- Update metrics
    self.fps = love.timer.getFPS()
    self.memoryUsage = collectgarbage("count") / 1024 -- Convert KB to MB

    -- Get the current stats (including drawcalls)
    local stats = love.graphics.getStats()
    self.drawCalls = stats.drawcalls
end

function DebugGUI:updateEntityCount(coins, enemies, spikes, stones)
    if not self.active then
        return
    end

    self.entityCounts.coins = coins or 0
    self.entityCounts.enemies = enemies or 0
    self.entityCounts.spikes = spikes or 0
    self.entityCounts.stones = stones or 0
    self.entityCounts.total = self.entityCounts.coins + self.entityCounts.enemies + self.entityCounts.spikes + self.entityCounts.stones
end

function DebugGUI:draw()
    if not self.active then
        return
    end

    -- Save current graphics state
    local prevFont = love.graphics.getFont()
    local r, g, b, a = love.graphics.getColor()

    -- Set up for debug drawing
    love.graphics.setFont(self.font)

    -- Draw background panel
    love.graphics.setColor(self.bgColor)
    love.graphics.rectangle("fill", 10, 10, 220, 185)

    -- Draw metrics text
    love.graphics.setColor(self.textColor)

    local y = 15
    love.graphics.print(string.format("FPS: %d", self.fps), 20, y)
    y = y + 20

    love.graphics.print(string.format("Memory: %.2f MB", self.memoryUsage), 20, y)
    y = y + 20

    love.graphics.print(string.format("Draw calls: %d", self.drawCalls), 20, y)
    y = y + 20

    love.graphics.print("Entities:", 20, y)
    y = y + 20

    love.graphics.print(string.format("    Coins: %d", self.entityCounts.coins), 20, y)
    y = y + 20

    love.graphics.print(string.format("    Enemies: %d", self.entityCounts.enemies), 20, y)
    y = y + 20

    love.graphics.print(string.format("    Spikes: %d", self.entityCounts.spikes), 20, y)
    y = y + 20

    love.graphics.print(string.format("    Stones: %d", self.entityCounts.stones), 20, y)
    y = y + 20

    love.graphics.print(string.format("Total: %d", self.entityCounts.total), 20, y)

    -- Restore graphics state
    love.graphics.setFont(prevFont)
    love.graphics.setColor(r, g, b, a)
end

function DebugGUI:toggle()
    self.active = not self.active
end

return DebugGUI
