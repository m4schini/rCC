local term_width, term_height = term.getSize()

-- Main Window
local main = window.create(term.current(), 1, 1, term_width - 1, term_height)
main.setBackgroundColour(colours.black)
main.setTextColour(colours.white)

-- Status bar - Shows connection status to server (websocket)
local status = window.create(term.current(), term_width, 1, 2, term_height)
status.setBackgroundColour(colours.red)
status.setTextColour(colours.white)

function printC(color, ...)
    local bc = term.getTextColor()

    term.setTextColor(color)
    print(...)
    term.setTextColor(bc)

end

function setConnectionStatus(connected)
    if connected then
        status.setBackgroundColour(colors.green)
        status.clear()
    else
        status.setBackgroundColour(colors.red)
        status.clear()
    end
end

function openWebsocket()
    setConnectionStatus(false)
    printC(colors.orange, "=> establishing connection")
    local url = "ws://localhost:8080/socket"
    local ws, err = http.websocket(url)

    while not ws do
        local wait = 5
        printC(colors.red, "=> connection failed. Retry in", wait .. "s")

        sleep(wait)
        ws, err = http.websocket(url)
    end

    setConnectionStatus(true)
    printC(colors.green, "=> connection established")

    return ws, err
end

function getData()
    return {
        id = os.getComputerID(),
        label = os.getComputerLabel(),
        fuel = turtle.getFuelLevel(),
        fuelLimit = turtle.getFuelLimit(),
        position = {turtle.getPosition()},
        chunk = {turtle.getChunkPosition()}
    }
end

function controlTurtle(cmd)
    local d = {}
    local c = cmd.CMD
    if c == 0 then
        d = getData()
    elseif c == 1 then
        d = {turtle.forward()}
    elseif c == 2 then
        d = {turtle.back()}
    elseif c == 3 then
        d = {turtle.down()}
    elseif c == 4 then
        d = {turtle.up()}
    elseif c == 5 then
        d = {turtle.turnLeft()}
    elseif c == 6 then
        d = {turtle.turnRight()}
    elseif c == 7 then
        d = {turtle.dig()}
    elseif c == 8 then
        d = {turtle.digUp()}
    elseif c == 9 then
        d = {turtle.digDown()}
    elseif c == 100 then
        d = turtle.inspectAll()
    elseif c == 101 then
        d = turtle.scan()
    elseif c == 102 then
        d = turtle.detectAll()
    elseif c == 103 then
        d = {turtle.getPosition()}
    elseif c == 104 then
        d = {turtle.getHeading()}
    elseif c == 105 then
        d = turtle.getInventory()
    else
        printC(colors.red, "=> UNK INSTR")
    end

    return d
end

local function processPrimaryLoop()
    term.redirect(main)
    term.setCursorPos(1, 1)


    local ws, err = openWebsocket()

    local function sendResponse(data)
        -- data["meta"] = metadata
        local x, y, z, suc = turtle.getPosition()
        local head = {}
        head.fuel = turtle.getFuelLevel()
        if suc then
            local pos = {}
            pos.x = x
            pos.y = y
            pos.z = z
            head.position = pos
        end

        data = textutils.serialiseJSON({
            data = data,
            head = head,
        }) -- textutils.serialiseJSON(data)
        term.redirect(main)
        print("<=", data)

        ws.send(data)
    end

    while true do
        local e = {os.pullEvent()}

        if e[1] == "websocket_closed" then
            printC(colors.red, "=> connection closed")
            ws, err = openWebsocket()
        elseif e[1] == "websocket_message" then
            print("=>", e[3])

            local cmd = textutils.unserializeJSON(e[3])
            local result = controlTurtle(cmd)

            sendResponse(result)
        end
    end
end

local function processAux()
    local term = status

    local i = 0;
    while true do
        i = i + 1

        sleep(1)
    end
end

parallel.waitForAll(processAux, processPrimaryLoop)
