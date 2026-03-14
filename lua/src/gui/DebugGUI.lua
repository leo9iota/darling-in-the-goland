--- src/gui/DebugGUI.lua
-- @class DebugGUI
-- Display metrics, such as fps, memory, draw calls, etc., for debugging performance bottlenecks
local DebugGUI = {}

DebugGUI.active = true -- Visible by default for testing

-- Metrics
DebugGUI.fps = 0
DebugGUI.frameTime = 0
DebugGUI.memoryUsage = 0
DebugGUI.drawCalls = 0

DebugGUI.entityCounts = {
    coins = 0,
    enemies = 0,
    spikes = 0,
    stones = 0,
    total = 0,
}

-- Constants
DebugGUI.lineHeight = 20
DebugGUI.margin = 10
DebugGUI.warningDrawCallThreshold = 1000

--- load()
-- Conventional function naming, e.g. t:load(), t:update(), t:draw().
-- REMEMBER: They are user defined and don't come from Love2D.
function DebugGUI:load()
    self.font = love.graphics.newFont(14)

    self.textColor = {
        1,
        1,
        0,
        0.75,
    }
    self.warningColor = {
        1,
        0.2,
        0.2,
        1,
    } -- Red for warnings
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

    self.fps = love.timer.getFPS()
    self.frameTime = dt * 1000 -- milliseconds
    self.memoryUsage = collectgarbage("count") / 1024
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

    -- Get the current draw calls before we start drawing our debug info
    local stats = love.graphics.getStats()
    self.drawCalls = stats.drawcalls

    -- Save current graphics state
    local prevFont = love.graphics.getFont()
    local r, g, b, a = love.graphics.getColor()

    love.graphics.setFont(self.font)

    -- Calculate panel height based on lines
    local numLines = 10 -- total lines printed (change if you add more)
    local panelHeight = numLines * self.lineHeight + self.margin

    -- Draw background
    love.graphics.setColor(self.bgColor)
    love.graphics.rectangle("fill", 10, 10, 250, panelHeight)

    -- Draw text
    love.graphics.setColor(self.textColor)
    local y = 15

    love.graphics.print(string.format("FPS: %d", self.fps), 20, y)
    y = y + self.lineHeight

    love.graphics.print(string.format("Frame Time: %.2f ms", self.frameTime), 20, y)
    y = y + self.lineHeight

    love.graphics.print(string.format("Memory: %.2f MB", self.memoryUsage), 20, y)
    y = y + self.lineHeight

    -- Warning color for high draw calls
    if self.drawCalls > self.warningDrawCallThreshold then
        love.graphics.setColor(self.warningColor)
    else
        love.graphics.setColor(self.textColor)
    end
    love.graphics.print(string.format("Draw Calls: %d", self.drawCalls), 20, y)
    y = y + self.lineHeight

    -- Reset to normal color for entity counts
    love.graphics.setColor(self.textColor)
    love.graphics.print("Entities:", 20, y)
    y = y + self.lineHeight

    love.graphics.print(string.format("  Coins: %d", self.entityCounts.coins), 20, y)
    y = y + self.lineHeight

    love.graphics.print(string.format("  Enemies: %d", self.entityCounts.enemies), 20, y)
    y = y + self.lineHeight

    love.graphics.print(string.format("  Spikes: %d", self.entityCounts.spikes), 20, y)
    y = y + self.lineHeight

    love.graphics.print(string.format("  Stones: %d", self.entityCounts.stones), 20, y)
    y = y + self.lineHeight

    love.graphics.print(string.format("Total: %d", self.entityCounts.total), 20, y)

    -- Restore previous graphics state
    love.graphics.setFont(prevFont)
    love.graphics.setColor(r, g, b, a)
end

function DebugGUI:toggle()
    self.active = not self.active
end

return DebugGUI
