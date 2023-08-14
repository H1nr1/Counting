{{/*
        Counting statistics
    
        Author: H1nr1 <https://github.com/H1nr1>
*/}}

{{/* Configurable Values */}}
{{ $LBLength := 10 }} {{/* How many members to show on leaderboard; MAX OF 100 */}}
{{/* End of configurable values */}}

{{/* No Touchy */}}
{{/* Initializing variables */}}
{{ $db := (dbGet 0 "Counting").Value }}
{{ $CCount := (dbGet .User.ID "CCount").Value }}{{ $CCorrect := (dbGet .User.ID "CCorrect").Value }}

{{/* Actual Code */}}
{{ $Args := parseArgs 0 "" (carg "string" "") }}
{{ if not ($Args.IsSet 0) }} {{/* Server Stats */}}
	{{ sendMessage nil (cembed "author" (sdict "icon_url" (print "https://cdn.discordapp.com/icons/" .Guild.ID "/" .Guild.Icon) "name" "🔢 Counting Statistics") "description" (print "⌚ __Current Score:__ " (sub $db.Next 1) "\n🏅 __High Score:__ " $db.HighScore.Num " on " (formatTime $db.HighScore.Time "01/02") " by " (userArg $db.HighScore.User) "\n⏮️ __Last Counter:__ " (userArg $db.Last.User) "\n💾 __Saves Remaining:__ " $db.SecondChance) "footer" (sdict "text" "Use this command: -CStats") "color" 30654 ) }}

{{ else if (inFold (cslice "Me" "My" "0") ($Args.Get 0)) }} {{/* Triggering user's stats */}}
	{{ sendMessage nil (cembed "title" (print "**🔢 " .User.Username "'s Counting Statistics**") "description" (print .User.Mention " has counted a __total__ of **" $CCount " times**\n" .User.Mention " has counted __correctly__ **" $CCorrect " times**\nThis makes " .User.Mention "'s __average__ **" (div (round (mult (div $CCorrect $CCount) 10000)) 100) "%**") "footer" (sdict "text" "Use this command: -CStats [Me/My/0]") "color" 30654 ) }}

{{ else if ($User := userArg (toInt (reFind `\d{17,19}` ($Args.Get 0)))) }} {{/* User Mention */}}
	{{ if not (dbGet $User.ID "CCount") }} {{/* No Stats */}}
		{{ sendMessage nil (cembed "title" (print "No available stats") "description" (print $User " has yet to count ☹️\nMaybe give them a heads-up to come join?") "footer" (sdict "text" "Use this command: -CStats [User: @/ID]")) }}
	{{ else }} {{/* Mentioned user's stats */}}
		{{ $CCount = (dbGet $User.ID "CCount").Value }}{{ $CCorrect = (dbGet $User.ID "CCorrect").Value }}
		{{ sendMessage nil (cembed "title" (print "**🔢 " $User "'s Counting Statistics**") "description" (print $User.Mention " has counted a __total__ of **" $CCount " times**\n" $User.Mention " has counted __correctly__ **" $CCorrect " times**\nThis makes" $User.Mention "'s __average__ **" (div (round (mult (div $CCorrect $CCount) 10000)) 100) "%**") "footer" (sdict "text" "Use this command: -CStats [User: @/ID]") "color" 30654 ) }}
	{{ end }}

{{ else if (inFold (cslice "Leaderboard" "LB") ($Args.Get 0)) }} {{/* Server leaderboard */}}
	{{ $Desc := "" }}{{ $Place := 1 }}
	{{ range (dbTopEntries "CCorrect" $LBLength 0) }}
		{{- $Desc = (joinStr "\n" $Desc (printf "#%-3d %4d - %-4s" $Place (toInt .Value) (or (userArg .UserID) (str .UserID)))) }}
		{{- $Place = add $Place 1 -}}
	{{ end }}
	{{ sendMessage nil (cembed "author" (sdict "icon_url" (print "https://cdn.discordapp.com/icons/" .Guild.ID "/" .Guild.Icon) "name" "Counting Leaderboard") "description" (print "```Pos    ✅  User\n" $Desc "```") "footer" (sdict "text" "Use this command: -CStats [Leaderboard/LB]") "color" 30654 ) }}

{{ else }} {{/* Invalid syntax */}}
	{{ sendMessage nil (cembed "title" "Invalid Syntax" "description" "For server counting statistics: `-CStats`\nFor your statistics: `-CStats <Me/My/0>`\nFor another member's statistics: `-CStats <User: @/ID>`\nFor server leaderboard: `-CStats <Leaderboard/LB>`" "color" 16744192 ) }}
{{ end }}

{{ deleteTrigger 5 }}
