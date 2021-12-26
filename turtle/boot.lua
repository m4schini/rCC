local xturtle=require("m4/xturtle")

term.clear();
term.setCursorPos(1,1)


local m4_path = "/m4/"

local function installM4()
    if not fs.exists(m4_path) then
        print("Installing m4...")
        mkdir(m4_path)
    end
end

installM4()

local id = multishell.launch({
    turtle=xturtle

}, m4_path .. "main.lua")
multishell.setTitle(id, "m4")


--term.setCursorPos(1,1)


-- shell.switchTab(id)