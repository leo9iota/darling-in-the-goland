--- src/core/Camera.lua
-- @class Camera
-- Game world camera that moves with the player with smooth following
local Camera = {
    -- Position
    x = 0,
    y = 0,

    -- Target position (what the camera is following)
    targetX = 0,
    targetY = 0,

    -- Smoothing settings
    smoothing = {
        enabled = true,
        type = "damped", -- "none", "linear", or "damped"
        speed = 5, -- For linear smoothing
        stiffness = 10 -- For damped smoothing (higher = faster)
    },

    -- Deadzone (area where camera doesn't move if target is inside)
    deadzone = {
        enabled = false,
        x = 0,
        y = 0,
        width = 100,
        height = 100
    },

    -- Bounds (limits where camera can go)
    bounds = {
        enabled = true,
        minX = 0,
        minY = 0,
        maxX = 0, -- Will be set based on map width
        maxY = 0  -- Will be set based on map height
    },

    scale = 2 -- The value 2 translates to 200% zoom
}

--- Initialize the camera for rendering
function Camera:init()
    love.graphics.push()
    love.graphics.scale(self.scale, self.scale)
    -- Round camera position to whole pixels to prevent jiggling of sprites
    local roundedX = math.floor(self.x + 0.5)
    local roundedY = math.floor(self.y + 0.5)
    love.graphics.translate(-roundedX, roundedY) -- Used to move the camera within the game world
end

--- Remove camera transformations
function Camera:remove()
    love.graphics.pop()
end

--- Set camera bounds based on map dimensions
-- @param mapWidth Width of the map
-- @param mapHeight Height of the map (optional)
function Camera:setBounds(mapWidth, mapHeight)
    self.bounds.maxX = mapWidth - love.graphics.getWidth() / self.scale

    if mapHeight then
        self.bounds.maxY = mapHeight - love.graphics.getHeight() / self.scale
    end
end

--- Set the camera's deadzone
-- @param x X position of deadzone center
-- @param y Y position of deadzone center
-- @param width Width of deadzone
-- @param height Height of deadzone
function Camera:setDeadzone(x, y, width, height)
    self.deadzone.enabled = true
    self.deadzone.x = x or 0
    self.deadzone.y = y or 0
    self.deadzone.width = width or 100
    self.deadzone.height = height or 100
end

--- Disable the camera's deadzone
function Camera:disableDeadzone()
    self.deadzone.enabled = false
end

--- Set the camera's smoothing
-- @param type Type of smoothing ("none", "linear", or "damped")
-- @param amount Amount of smoothing (speed for linear, stiffness for damped)
function Camera:setSmoothing(type, amount)
    self.smoothing.enabled = (type ~= "none")
    self.smoothing.type = type or "damped"

    if type == "linear" then
        self.smoothing.speed = amount or 5
    elseif type == "damped" then
        self.smoothing.stiffness = amount or 10
    end
end

--- Immediately set camera position without smoothing
-- @param x X position
-- @param y Y position
function Camera:setPosition(x, y)
    -- Set target position
    self.targetX = x - love.graphics.getWidth() / 2 / self.scale
    self.targetY = y

    -- Apply bounds to target
    self:applyBounds()

    -- Immediately update camera position to target (no smoothing)
    self.x = self.targetX
    self.y = self.targetY
end

--- Set the target position for the camera to follow
-- @param x X position
-- @param y Y position
function Camera:follow(x, y)
    -- Set target position
    self.targetX = x - love.graphics.getWidth() / 2 / self.scale
    self.targetY = y

    -- Apply deadzone if enabled
    if self.deadzone.enabled then
        self:applyDeadzone()
    end

    -- Apply bounds to target
    self:applyBounds()

    -- If smoothing is disabled, immediately update camera position
    if not self.smoothing.enabled then
        self.x = self.targetX
        self.y = self.targetY
    end
end

--- Apply deadzone to target position
function Camera:applyDeadzone()
    local dx = self.targetX - self.x
    local dy = self.targetY - self.y

    -- Calculate deadzone boundaries
    local left = self.x - self.deadzone.width / 2
    local right = self.x + self.deadzone.width / 2
    local top = self.y - self.deadzone.height / 2
    local bottom = self.y + self.deadzone.height / 2

    -- Only move camera if target is outside deadzone
    if self.targetX < left then
        self.targetX = left
    elseif self.targetX > right then
        self.targetX = right
    else
        self.targetX = self.x -- Keep camera at current position
    end

    if self.targetY < top then
        self.targetY = top
    elseif self.targetY > bottom then
        self.targetY = bottom
    else
        self.targetY = self.y -- Keep camera at current position
    end
end

--- Apply bounds to target position
function Camera:applyBounds()
    if self.bounds.enabled then
        if self.targetX < self.bounds.minX then
            self.targetX = self.bounds.minX
        elseif self.targetX > self.bounds.maxX then
            self.targetX = self.bounds.maxX
        end

        if self.targetY < self.bounds.minY then
            self.targetY = self.bounds.minY
        elseif self.targetY > self.bounds.maxY then
            self.targetY = self.bounds.maxY
        end
    end
end

--- Update camera position with smoothing
-- @param dt Delta time
function Camera:update(dt)
    if not self.smoothing.enabled then return end

    if self.smoothing.type == "linear" then
        -- Linear interpolation
        local speed = self.smoothing.speed * dt
        self.x = self.x + (self.targetX - self.x) * speed
        self.y = self.y + (self.targetY - self.y) * speed
    elseif self.smoothing.type == "damped" then
        -- Damped spring physics
        local stiffness = self.smoothing.stiffness
        local dx = self.targetX - self.x
        local dy = self.targetY - self.y

        self.x = self.x + dx * stiffness * dt
        self.y = self.y + dy * stiffness * dt
    end
end

return Camera
