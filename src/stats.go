{{/*
        Counting statistics

        Command: CStats

        Author: H1nr1 <https://github.com/H1nr1>
*/}}

{{/* configurable values */}}
{{ $LEADERBOARD_LENGTH := 10 }} {{/* How many members to show on leaderboard; MAX OF 100 */}}
{{/* end of configurable values */}}

{{$db:=(dbGet 0 "counting").Value}}
{{$args:=joinStr " " .CmdArgs}}

{{if not $args}} {{/* general stats */}}
	{{sendMessage nil (cembed 
		"author" (sdict 
			"icon_url" (.Guild.IconURL "512") 
			"name" "🔢 Counting Statistics"
		) 
		"description" (printf "⌚ __Current Score:__ %d\n🏅 __High Score:__ %d on %v by %s (%d)\n⏮️ __Last Counter:__ %s (%d)\n💾 __Saves Remaining:__ %d" 
			(sub $db.next 1) 
			$db.highscore.num (formatTime $db.highscore.time "01/02") (userArg $db.highscore.user).Username (userArg $db.highscore.user).ID 
			(userArg $db.last.user).Username (userArg $db.last.user).ID 
			$db.saves
		) 
		"footer" (sdict "text" "Use this command: -CStats") 
		"color" 30654
	)}}

{{else if $u:=reFind `\d{17,19}` $args|toInt|userArg}} {{/* user mention/ID */}}
	{{if dbGet $u.ID "counting"|not}} {{/* no stats */}}
		{{sendMessage nil (cembed 
			"title" (print "No available stats") 
			"description" (printf "%s has yet to count ☹️\nMaybe give them a heads-up to come join?" 
				$u.Username
			) 
			"footer" (sdict "text" "Use this command: -CStats [User: @/ID]") 
			"color" 16711680
		)}}
	{{else}} {{/* mentioned user's stats */}}
		{{$uCount =(dbGet $u.ID "counting").Value}}
		{{$uCorrect =(dbGet $u.ID "countingCorrect").Value}}
		{{sendMessage nil (cembed 
			"title" (printf "🔢 %s's Counting Statistics" $u.Username) 
			"description" (printf "%s has counted **%d times**\n%d of those were correct\nThis makes %s's average **%d%**" 
				$u.Mention $uCount $uCorrect $u.Username 
				(div $uCorrect $uCount|mult 10000.0|round|mult 0.01)
			) 
			"footer" (sdict "text" "Use this command: -CStats [User: @/ID]") 
			"color" 30654
		)}}
	{{end}}

{{else if reFind `(?i)l(eader)?b` $args}} {{/* leaderboard */}}
	{{$desc:=""}}{{$pos:=1}}
	{{range dbTopEntries "countingCorrect" $LEADERBOARD_LENGTH 0}}
		{{- $desc =printf "%s\n#%-3d %4d - %-4s" $desc $pos (toInt .Value) (or (userArg .UserID) (str .UserID))}}
		{{- $pos =add $pos 1 -}}
	{{end}}
	{{sendMessage nil (cembed 
		"author" (sdict 
			"icon_url" (.Guild.IconURL "512") 
			"name" "Counting Leaderboard"
		) 
		"description" (printf "```Pos    ✅  User\n%s```" $desc) 
		"footer" (sdict "text" "Use this command: -CStats [Leaderboard/LB]") 
		"color" 30654
	)}}

{{else}} {{/* invalid syntax */}}
	{{sendMessage nil (cembed 
		"title" "Invalid Syntax" 
		"description" "For server counting statistics: `-CStats`\nFor a member's statistics: `-CStats <User: @/ID>`\nFor server leaderboard: `-CStats Leaderboard`" 
		"color" 16744192
	)}}
{{end}}
