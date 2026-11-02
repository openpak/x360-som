-- Dash.Social.lua  -  OpenPak online Social channel (replaces BuiltIn.Social.Lua.Offline.xzp)
--
-- Loaded by the 17559 dash native Lua-app host (luadashapp.cpp, FUN_927cef98 /
-- FUN_927c0680 @0x927C0680). The host runs this chunk as a function, then calls
-- the GLOBAL entry points it leaves behind:
--     SetAppConfig(cfg)   -- called first if present   (dash str @0x92103a30)
--     Main(args)          -- REQUIRED; "Missing entry point 'Main'" otherwise
--                            (dash str @0x921021b4 / @0x921039a0)
--     BackgroundApp()     -- optional                  (dash str @0x921039bc)
-- There is NO content hash or RSA check on this chunk at the dash VM
-- (FUN_927cef98 only checks lua_type == LUA_TFUNCTION). The gate lives in XAM,
-- on the system.online.manifest container that lists this file (see
-- re-som-manifest.md Part A: that manifest IS XeKeysVerifyRSASignature-signed).
--
-- Framework + natives used below are the dash built-ins proven present in the
-- 17559 dump (strings/natives): module/import/require, App.* (BuiltIn.ContentApp.xzp),
-- Xbox.App, Xbox.UserImpl.GetFriends, EcShowGamerCard. API shapes of App.UI /
-- scene objects are [I] (inferred from the luac -l listings in store-re/lst/);
-- the entry-point contract (SetAppConfig/Main globals) is [V].

module("Dash.Social", package.seeall)

require("Xbox.App")          -- native table: user/friends/gamercard etc.
import("App.Extensions")     -- App.Extensions.Register{ AppInit = ... }
import("App.UI")             -- scene / list helpers

-- ------------------------------------------------------------------ friends
-- Xbox.UserImpl.GetFriends(userIndex) -> iterator/table of friend records.
-- Each record exposes GamerTag and Xuid (field names per the offline module).
local function GetFriendList(userIndex)
    local list = {}
    local ok, friends = pcall(Xbox.UserImpl.GetFriends, userIndex)
    if ok and friends then
        for _, f in ipairs(friends) do
            list[#list + 1] = { gamertag = f.GamerTag or f.gamertag, xuid = f.Xuid or f.xuid }
        end
    end
    return list
end

-- A on a row -> native gamer card (EcShowGamerCard, Ec cmd id 8; param = xuid).
local function OnActivate(row)
    if row and row.xuid then
        EcShowGamerCard(row.xuid)
    end
end

local function BuildFriendsScene(userIndex)
    local scene = App.UI.CreateScene("DashAppHostScene")   -- host scene @0x92102340
    local panel = App.UI.CreateList(scene, "FriendsList")
    for _, friend in ipairs(GetFriendList(userIndex)) do
        local item = App.UI.AddListItem(panel, friend.gamertag or "")
        item.data = friend
        item.OnPress = function() OnActivate(friend) end
    end
    return scene
end

-- ----------------------------------------------------------- entry points
local _cfg

-- [V] called by the host before Main, if defined.
function SetAppConfig(cfg)
    _cfg = cfg or {}
end

-- [V] required global; the host pcalls this to start the channel.
function Main(args)
    local userIndex = (args and args.userIndex) or (_cfg and _cfg.userIndex) or 0
    local scene = BuildFriendsScene(userIndex)
    App.UI.Show(scene)
    return scene
end

-- optional background entry; keep it a no-op so the host's optional call is safe.
function BackgroundApp()
end
