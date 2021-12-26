local API_URL = "http://localhost:8080/turtle/"
local term_width, term_height = term.getSize()

-- Main Window
local main = window.create(term.current(), 1, 1, term_width - 1, term_height)
main.setBackgroundColour(colours.black)
main.setTextColour(colours.white)


-- Status bar - Shows connection status to server (websocket)
local status = window.create(term.current(), term_width, 1, 2, term_height)
status.setBackgroundColour(colours.red)
status.setTextColour(colours.white)

function updateStatus(fuelLevel)
    p = fuelLevel / 100000

    if p > 0.8 then
        status.setBackgroundColour(colors.blue)
    elseif p > 0.5 then
        status.setBackgroundColour(colors.green)
    elseif p > 0.1 then
        status.setBackgroundColour(colors.orange)
    else
        status.setBackgroundColour(colors.red)
    end

    status.clear()
end

function sendCmd(turtle, cmd)
    http.request(API_URL .. turtle .. "/" .. cmd)
end

function printInfo(data)
    local x, y, z = gps.locate()
    term.redirect(main)

    term.clear()
    term.setCursorPos(1,1)

    print("PLAYER:\t\t\t\t TURTLE:")
    print("X", x, "\t\t\t\t", data["position"]["X"])
    print("Y", y, "\t\t\t\t", data["position"]["Y"])
    print("Z", z, "\t\t\t\t", data["position"]["Z"])
    print()
    print(textutils.serialise(data))

    updateStatus(data["fuel"])
end

while true do
    local event = {os.pullEvent()}

    if event[1] == "char" then
        local c = event[2]

        local turtleId = 1

        if c == "w" then
            sendCmd(turtleId, "forward")
        elseif c == "s" then
            sendCmd(turtleId, "back")
        elseif c == "a" then
            sendCmd(turtleId, "turnLeft")
        elseif c == "d" then
            sendCmd(turtleId, "turnRight")
        elseif c == "q" then
            sendCmd(turtleId, "up")
        elseif c == "e" then
            sendCmd(turtleId, "down")
        elseif c == "r" then
            sendCmd(turtleId, "digUp")
        elseif c == "f" then
            sendCmd(turtleId, "dig")
        elseif c == "v" then
            sendCmd(turtleId, "digDown")
        end
    elseif event[1] == "http_success" then

        local json = textutils.unserialiseJSON(event[3].readAll())

        printInfo(json)
    end
end
