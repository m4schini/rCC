 local args = {...}
 
 local packages = {
    trc="",
    trc_bg="",
    pocket=""
 }


local instr = args[1]
if instr == "list" then
    for k, v in pairs(packages) do
        print(k)
    end
end