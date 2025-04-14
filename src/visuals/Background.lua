--- src/visuals/Background.lua
-- @class Background
-- Background system with parallax effect
-- Implements a 4-layer parallax background where layer 1 is closest and layer 4 is farthest
local Background = {}

-- Table to store all background layers
Background.layers = {}

--- Load background layers
-- @param layerPaths Table of image paths for each layer (from closest to farthest)
-- @param parallaxFactors Table of parallax factors for each layer (smaller values = slower movement)
function Background:load(layerPaths, parallaxFactors)
    self.layers = {}

    -- Load each background layer
    for i, path in ipairs(layerPaths) do
        table.insert(self.layers, {
            image = love.graphics.newImage(path),
            parallaxFactor = parallaxFactors[i] or (1 - (i - 1) * 0.2), -- Default decreasing factors if not provided
            x = 0,
            y = 0, -- Allow for vertical positioning
        })
    end

    -- Get screen dimensions
    self.screenWidth = love.graphics.getWidth()
    self.screenHeight = love.graphics.getHeight()

    -- Position layers vertically to align with the bottom of the screen
    -- This ensures backgrounds with different heights align properly
    for _, layer in ipairs(self.layers) do
        local imgHeight = layer.image:getHeight()
        -- Position the image so it aligns with the bottom of the screen
        layer.y = self.screenHeight - imgHeight
    end
end

--- Update background positions based on camera movement
-- @param dt Delta time
-- @param cameraX Camera X position
function Background:update(dt, cameraX)
    for _, layer in ipairs(self.layers) do
        -- Update layer position based on camera and parallax factor
        -- The lower the parallax factor, the slower the layer moves
        layer.x = -cameraX * layer.parallaxFactor
    end
end

--- Draw all background layers
function Background:draw()
    -- Draw from farthest to closest (reverse order from our layers table)
    for i = #self.layers, 1, -1 do
        local layer = self.layers[i]
        local img = layer.image
        local imgWidth = img:getWidth()

        -- Calculate how many copies we need to cover the screen
        local copies = math.ceil(self.screenWidth / imgWidth) + 1

        -- Calculate the starting position for tiling
        local startX = math.floor(layer.x % imgWidth) - imgWidth

        -- Draw enough copies to cover the screen
        for j = 0, copies do
            love.graphics.draw(img, startX + (j * imgWidth), layer.y)
        end
    end
end

return Background
