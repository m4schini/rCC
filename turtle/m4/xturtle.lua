local expect = require("cc.expect").expect

-- native turtle function.
-- if you overwrite function but still want to use them,
-- save them here.
nativeTurtle = {
    dig = turtle.dig,
    digDown = turtle.digDown,
    digUp = turtle.digUp
}

XTurtle = turtle

function XTurtle.getPosition()
    local x, y, z = gps.locate()

    if x == nil then
        return 0, 0, 0, false
    else
        return x, y, z, true
    end
end

function XTurtle.getChunkPosition()
    local chunkSize = 16
    local x, y, z, ok = XTurtle.getPosition()


    if ok then
        x = (x % chunkSize)
        y = (y % chunkSize)
        z = (z % chunkSize)

        return x, y, z
    else
        return 0, 0, 0
    end
end

function XTurtle.getChunk()
    local chunkSize = 16
    local x, y, z, ok = XTurtle.getPosition()


    if ok then
        x = (x - (x % chunkSize)) / chunkSize
        y = (y - (y % chunkSize)) / chunkSize
        z = (z - (z % chunkSize)) / chunkSize

        return x, y, z
    else
        return 0, 0, 0
    end
end

--[[orientation will be:
-x = 1
-z = 2
+x = 3
+z = 4
This matches exactly with orientation in game, except that Minecraft uses 0 for +z instead of 4.
--]]
-- src: https://www.computercraft.info/forums2/index.php?/topic/1704-get-the-direction-the-turtle-face/
function XTurtle.getHeading()
    local restorePos = false
    local loc1 = vector.new(gps.locate(2, false))
    if not turtle.forward() then
        if turtle.back() then
            loc1 = vector.new(gps.locate(2, false))

            XTurtle.dig() -- in case of gravel and sand
            turtle.forward()
        else
            return -1, "nope"
        end
    else
        restorePos = true
    end
    local loc2 = vector.new(gps.locate(2, false))
    if restorePos then
        turtle.back()
    end
    local heading = loc2 - loc1

    local face = ((heading.x + math.abs(heading.x) * 2) + (heading.z + math.abs(heading.z) * 3))
    local direction = nil
    
    if face == 1 then
        direction = "west"
    elseif face == 2 then
        direction = "north"
    elseif face == 3 then
        direction = "east"
    elseif face == 4 then
        direction = "south"
    end

    return face, direction
end

function XTurtle.detectAll()
    local d = {
        up = turtle.detectUp(),
        front = turtle.detect(),
        down = turtle.detectDown()
    }

    return d
end

function XTurtle.inspectAll()
    local d = {}
    local block, up = turtle.inspectUp()
    if block then
        d["up"] = up
    else
        d["up"] = {
            name = "minecraft:air"
        }
    end

    local block, front = turtle.inspect()
    if block then
        d["front"] = front
    else
        d["front"] = {
            name = "minecraft:air"
        }
    end

    local block, down = turtle.inspectDown()
    if block then
        d["down"] = down
    else
        d["down"] = {
            name = "minecraft:air"
        }
    end

    return d
end

function XTurtle.scan()
    local d = XTurtle.inspectAll()

    turtle.turnRight()
    local block, front = turtle.inspect()
    if block then
        d["right"] = front
    else
        d["right"] = {
            name = "minecraft:air"
        }
    end

    turtle.turnRight()
    local block, front = turtle.inspect()
    if block then
        d["back"] = front
    else
        d["back"] = {
            name = "minecraft:air"
        }
    end

    turtle.turnRight()
    local block, front = turtle.inspect()
    if block then
        d["left"] = front
    else
        d["left"] = {
            name = "minecraft:air"
        }
    end

    turtle.turnRight()

    return d
end

-- Digs the block in front of the turtle. 
-- **Keeps Digging** until there is no block in front of turtle
function XTurtle.dig()
    while turtle.detect() do
        nativeTurtle.dig()
    end

    return true
end

-- Digs the block below the turtle. 
-- **Keeps Digging** until there is no block below the turtle
function XTurtle.digDown()
    while turtle.detectDown() do
        nativeTurtle.digDown()
    end

    return true
end

-- Digs the block in above the turtle. 
-- **Keeps Digging** until there is no block above the turtle
function XTurtle.digUp()
    while turtle.detectUp() do
        nativeTurtle.digUp()
    end

    return true
end

-- Full Inventory overview (or from specified startSlot to endSlot)
-- @param startslot (default: 1) first slot
-- @param endslot (default: 16) last slot
function XTurtle.getInventory(startslot, endslot)
    startslot = startslot or 1
    endslot = endslot or 16

    local inv = {}

    for i = startslot, endslot, 1 do
        inv[i] = turtle.getItemDetail(i) or {}
    end

    return inv
end

function XTurtle.isInventoryFull(startSlot, endSlot)
    for i = startSlot or 1, endSlot or 16, 1 do
        if turtle.getItemCount(i) == 0 then
            return false
        end
    end
    return true
end

return XTurtle
