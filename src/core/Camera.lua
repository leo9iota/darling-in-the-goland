--- src/core/Camera.lua
-- @class Camera
-- HUMP camera implementation
local humpCamera = require("hump.camera") -- Import hump camera module

local Camera = {}

-- Camera instance
local cam = nil

-- Camera bounds
local bounds = {
    left = 0,
    right = 0
}

-- Initialize the camera
function Camera:load()
    -- Create a new camera instance with default position and zoom
    cam = humpCamera(0, 0, 1)

    -- Default scale
    self.scale = 1

    -- Default position
    self.x = 0
    self.y = 0
end

-- Set camera bounds to prevent it from showing areas outside the map
-- @param mapWidth The width of the map in pixels
function Camera:setBounds(mapWidth)
    bounds.right = mapWidth
    -- We use the screen width to calculate the left bound to prevent showing empty space
    bounds.left = love.graphics.getWidth() / 2
end

-- Set the camera smoothing type
-- @param smoothType The type of smoothing ("none", "linear", or "damped")
-- @param value The smoothing value (speed for linear, stiffness for damped)
function Camera:setSmoothing(smoothType, value)
    if not cam then
        self:load()
    end

    if smoothType == "none" then
        cam.smoother = humpCamera.smooth.none()
    elseif smoothType == "linear" then
        cam.smoother = humpCamera.smooth.linear(value or 5)
    elseif smoothType == "damped" then
        cam.smoother = humpCamera.smooth.damped(value or 5)
    end
end

-- Make the camera follow a target
-- @param x The x position to follow
-- @param y The y position to follow (optional)
function Camera:follow(x, y)
    if not cam then
        self:load()
    end

    -- Apply bounds to the target position
    local targetX = math.max(bounds.left, math.min(x, bounds.right - love.graphics.getWidth() / 2))

    -- Store the current position for external access
    self.x = cam.x
    self.y = cam.y

    -- Use lockX for horizontal following only
    if y then
        cam:lockPosition(targetX, y)
    else
        cam:lockX(targetX)
    end
end

-- Update the camera
-- @param dt Delta time
function Camera:update(dt)
    if not cam then
        self:load()
    end

    -- Update the stored position for external access
    self.x = cam.x
    self.y = cam.y
end

-- Initialize the camera for drawing (attach)
function Camera:init()
    if not cam then
        self:load()
    end

    cam:attach()
end

-- Remove the camera transformations (detach)
function Camera:remove()
    if not cam then
        self:load()
    end

    cam:detach()
end

-- Get the mouse position in world coordinates
function Camera:mousePosition()
    if not cam then
        self:load()
    end

    return cam:mousePosition()
end

-- Convert screen coordinates to world coordinates
function Camera:screenToWorld(x, y)
    if not cam then
        self:load()
    end

    return cam:worldCoords(x, y)
end

-- Convert world coordinates to screen coordinates
function Camera:worldToScreen(x, y)
    if not cam then
        self:load()
    end

    return cam:cameraCoords(x, y)
end

-- Set the camera zoom level
function Camera:setZoom(zoom)
    if not cam then
        self:load()
    end

    cam:zoomTo(zoom)
    self.scale = zoom
end

return Camera
