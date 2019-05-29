if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

-- hub exists, replace git command
local hubpath=nyagos.which("hub.exe")
if hubpath then
  nyagos.alias.git = "hub.exe"
end

share.git = {}

-- setup local branch listup
local branchlist = function()
  local gitbranches = {}
  local gitbranch_tmp = nyagos.eval('git for-each-ref  --format="%(refname:short)" refs/heads/ 2> nul')
  for line in gitbranch_tmp:gmatch('[^\n]+') do
    table.insert(gitbranches,line)
  end
  return gitbranches
end

--setup current branch string
local currentbranch = function()
  return nyagos.eval('git rev-parse --abbrev-ref HEAD 2> nul')
end

-- follow up Intermediate state
local currentstatus = function()
  local status = ''

  -- https://github.com/git/git/blob/master/contrib/completion/git-prompt.sh
  -- r="|REBASE"
    -- r="|REBASE-i"
    -- r="|REBASE-m"
  -- r="|AM"
  -- r="|AM/REBASE"
  -- r="|MERGING"
  -- r="|CHERRY-PICKING"
  -- r="|REVERTING"
  -- r="|BISECTING"


  return status
end

-- subcommands
local gitsubcommands={}

-- keyword
gitsubcommands["bisect"]={"start", "bad", "good", "skip", "reset", "visualize", "replay", "log", "run"}
gitsubcommands["notes"]={"add", "append", "copy", "edit", "list", "prune", "remove", "show"}
gitsubcommands["reflog"]={"show", "delete", "expire"}
gitsubcommands["rerere"]={"clear", "forget", "diff", "remaining", "status", "gc"}
gitsubcommands["stash"]={"save", "list", "show", "apply", "clear", "drop", "pop", "create", "branch"}
gitsubcommands["submodule"]={"add", "status", "init", "deinit", "update", "summary", "foreach", "sync"}
gitsubcommands["svn"]={"init", "fetch", "clone", "rebase", "dcommit", "log", "find-rev", "set-tree", "commit-diff", "info", "create-ignore", "propget", "proplist", "show-ignore", "show-externals", "branch", "tag", "blame", "migrate", "mkdirs", "reset", "gc"}
gitsubcommands["worktree"]={"add", "list", "lock", "prune", "unlock"}

-- branch
gitsubcommands["checkout"]=branchlist
gitsubcommands["reset"]=branchlist
gitsubcommands["merge"]=branchlist
gitsubcommands["rebase"]=branchlist

local gitvar=share.git
gitvar.subcommand=gitsubcommands
gitvar.branch=branchlist
gitvar.currentbranch=currentbranch
share.git=gitvar

if share.maincmds then
  if share.maincmds["git"] then
    -- git command complementation exists.
    local maincmds = share.maincmds

    -- build
    for key, cmds in pairs(gitsubcommands) do
      local gitcommand="git "..key
      maincmds[gitcommand]=cmds
    end

    -- replace
    share.maincmds = maincmds
  end
end

-- EOF
