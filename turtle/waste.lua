local function waste()
while true do
    turtle.up()
    turtle.down()
end
end

local function log()
    while true do
        local e = {os.pullEvent()}
        if e[1] == "turtle_response" then
            print(textutils.serialise(e))
        end
    end
end

parallel.waitForAll(
    waste,
    log
)
