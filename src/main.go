{{/*
        Counting command. Count and have fun!
        
	Regex: `\A(\d+|\()`
	
        Author: H1nr1 <https://github.com/H1nr1>
*/}}

{{/* configurable values */}}
{{ $countTwice := false }} {{/* allow users to count multiple times in a row; true/false */}}
{{ $correctRID := false }} {{/* correct Counting role ID; set to false to disable */}}
{{ $incorrectRID := false }} {{/* incorrect Counting role ID; set to false to disable */}}
{{ $errorCID := .Channel.ID }} {{/* ID of channel to send errors to */}}
{{ $secondChance := true }} {{/* second chance if wrong; true/false */}}
{{ $statsCC := true }} {{/* if you added the Stats CC; true/false */}}
{{ $reactions := true }} {{/* allow confirmative reactions on message; true false */}}
	{{ $reactionDelete := true }} {{/* toggle for reactions to delete from last message; true/false */}}
	{{ $correctEmoji := "✅" }} {{/* emoji to react with if number is correct; if custom, use format name:id */}}
	{{ $warningEmoji := "⚠️" }} {{/* emoji to react with if wrong number AND Second Chance is true/on; if custom, use format name:id */}}
	{{ $incorrectEmoji := "❌" }} {{/* emoji to react with if number is incorrect; if custom, use format name:id */}}
{{/* end of configurable values */}}

{{$db:=or 
	(dbGet 0 "counting").Value 
	(sdict 
		"last" (sdict 
			"user" .BotUser.ID 
			"msg" 0
		) 
		"next" 1 
		"highscore" (sdict 
			"user" .BotUser.ID 
			"num" 0 
			"time" currentTime
		) 
		"saves" 2
	)
}}

{{with .ExecData }}
	{{$msg:=getMessage nil .}}
	{{if not $msg}} {{/* check if message was deleted */}}
		{{sendMessage nil (cembed 
			"description" (printf "%s deleted their number which was correct!\nThe next number is %d" 
				(userArg $db.last.user).Mention $db.next
			) 
			"color" 30654
		)}}
	{{else if $msg.EditedTimestamp}} {{/* check if message was edited */}}
		{{sendMessage nil (cembed 
			"description" (printf "%s edited their message, be careful" (userArg $db.last.user).Mention) 
			"color" 30654
		)}}
	{{end}}
	{{return}}
{{end}}

{{$number:=index .Args 0}}
{{try}}{{$number =exec "calc" $number}}
{{catch}}Invalid Number. Please try again{{return}}{{end}}
{{$number =reFind `\d+` $number|toInt}}

{{if and (eq $db.last.user .User.ID) (not $countTwice)}} {{/* checks user */}}
	{{sendMessage nil (cembed 
		"description" (printf "You can't count twice in a row 🥲\nThe next number is %d" $db.next) 
		"color" 16744192
	)}}
	{{return}}
{{end}}

{{if eq $db.next $number}} {{/* checks if correct number */}}
	{{$db.Set "next" (add $db.next 1)}}
	{{if $reactions}}
		{{try}}
			{{addReactions $correctEmoji}}
			{{if and $reactionDelete $db.last.msg}}
				{{deleteMessageReaction nil $db.last.msg .BotUser.ID $correctEmoji}}
			{{end}}
			{{if mod $number 100|not}}{{addReactions "💯"}}{{end}}
		{{catch}}
			{{with $incorrectRID}}
				{{addRoleID .}}{{(toDuration "1d").Seconds|toInt|removeRoleID .}}
			{{end}}
			{{sendMessage $errorCID (printf "Counting: `%s`" .Error)}}
		{{end}}
	{{end}}
	{{with $correctRID}}{{takeRoleID $db.last.user .}}{{addRoleID .}}{{end}}
	{{$db.Set "last" (sdict "user" .User.ID "msg" .Message.ID)}}
	{{if $statsCC}} {{/* update leaderboard values */}}
		{{$s:=dbIncr .User.ID "countingCorrect" 1}}{{$s =dbIncr .User.ID "counting" 1}}
		{{if gt $number $db.highscore.num}}
			{{if $reactions}}{{try}}{{addReactions "🏆"}}{{catch}}{{end}}{{end}}
			{{$db.Set "highscore" (sdict "user" .User.ID "num" $number "time" currentTime)}}
		{{end}}
	{{end}}
	{{dbSet 0 "counting" $db}}
	{{execCC .CCID nil 15 .Message.ID}} {{/* call to check if message was deleted/edited */}}
		
{{else}} {{/* wrong number */}}
	{{$db.Set "saves" (sub $db.saves 1)}}
	{{with $correctRID}}{{takeRoleID $db.last.user .}}{{end}}
	{{with $incorrectRID}}{{addRoleID .}}{{(toDuration "3d").Seconds|toInt|removeRoleID .}}{{end}}
	{{if and $secondChance (gt $db.saves 0)}} {{/* second chance */}}
		{{if $reactions}}
			{{try}}{{addReactions $warningEmoji}}
			{{catch}}
				{{with $incorrectRID}}
					{{addRoleID .}}{{(toDuration "1d").Seconds|toInt|removeRoleID .}}
				{{end}}
				{{sendMessage $errorCID (printf "Counting: `%s`" .Error)}}
			{{end}}
		{{end}}
		{{$db.Set "last" (sdict "user" .User.ID "msg" .Message.ID)}}{{dbSet 0 "counting" $db}}
		{{sendMessage nil (cembed 
			"description" (printf "%s sent an incorrect number of %d\n**But**, second chance saved the count!\nThe next number is %d" .User.Username $number $db.next) 
			"color" 16744192
		)}}

	{{else}} {{/* reset count */}}
		{{sendMessage nil (cembed 
			"description" (printf "%s sent an incorrect number of %d\nCorrect number was %d\nStart over at 1 🙃" .User.Mention $number $db.next) 
			"color" 16711680
		)}}
		{{$db.Set "last" (sdict "user" .BotUser.ID "msg" 0)}}
		{{$db.Set "next" 1}}{{$db.Set "saves" 2}}{{dbSet 0 "counting" $db}}
		{{if $statsCC}}{{$s:=dbIncr .User.ID "counting" 1}}{{end}}
		{{if $reactions}}
			{{try}}{{addReactions $incorrectEmoji}}
			{{catch}}
				{{with $incorrectRID}}
					{{addRoleID .}}{{(toDuration "1d").Seconds|toInt|removeRoleID .}}
				{{end}}
				{{sendMessage $errorCID (printf "Counting: `%s`" .Error)}}
			{{end}}
		{{end}}
	{{end}}
{{end}}
