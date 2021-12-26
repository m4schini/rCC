local args = {...}
 
local packages = {
   trc={"Turtle Remote Controller", "https://raw.githubusercontent.com/m4schini/rCC/main/turtle/m4/main.lua?token=ADDWYWUYU3GEDBCOXBNCKVDB2HBNY"},
   trc_bg={"Turtle Remote Controller [multishell]", ""},
   xturtle={"XTurtle Lib", ""},
   pocket={"PocketOS", ""}
}

function printC(color, ...)
   local bc = term.getTextColor()

   term.setTextColor(color)
   print(...)
   term.setTextColor(bc)

end

local error = nil
local instr = args[1]
if instr == "list" then
   for k, v in pairs(packages) do
       local div = "\t| "
       if #k < 5 then
           div = "\t\t\t\t| "
       end

       print(k, div, v[1])
   end
elseif instr == "install" then
   local data = packages[args[2]]
   if data ~= nil then
       local dest = args[3] or args[2]
       if fs.exists(dest) or fs.isDir(dest) or fs.isReadOnly(dest) then
           error = "destination is invalid"
       else
           print("installing", args[2], "to", dest)
           local file = nil
           local download = nil

           local function openFile()
               file = fs.open(dest, "w")
           end

           local function downloadURL()
               r = http.get(data[2])
               download = r.readAll()
               r.close()
           end

           parallel.waitForAll(openFile, downloadURL)
           file.write(download)
           file.close()
           print(args[2], "installation finished")
       end
   else
       error = args[2] .. " doesnt exist"
   end

   
end

if error ~= nil then
   printC(colors.red, error)
end