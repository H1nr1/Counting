{{/*
        Counting command. Count and have fun!
        
	See <https://yagpdb-cc.github.io/fun/counting-v2/overview>
	
        Author: H1nr1 <https://github.com/H1nr1>
*/}}

{{/* Configurable Values */}}
{{ $CountTwice := false }} {{/* Allow users to count multiple times in a row; true/false */}}
{{ $CorrectRID := false }} {{/* Correct Counting role ID; set to false to disable */}}
{{ $IncorrectRID := flase }} {{/* Incorrect Counting role ID; set to false to disable */}}
{{ $ErrorCID := 123 }} {{/* Channel ID to send errors to */}}
{{ $SecondChance := true }} {{/* Second chance if wrong; true/false */}}
{{ $StatsCC := true }} {{/* If you added the Stats CC; true/false */}}
{{ $Reactions := true }} {{/* Allow confirmative reactions on message; true false */}}
	{{ $ReactionDelete := true }} {{/* Toggle for reactions to delete from last message; true/false */}}
	{{ $CorrectEmoji := "✅" }} {{/* Emoji to react with if number is correct; if custom, use format name:id */}}
	{{ $WarningEmoji := "⚠️" }} {{/* Emoji to react with if wrong number AND Second Chance is true/on; if custom, use format name:id */}}
	{{ $IncorrectEmoji := "❌" }} {{/* Emoji to react with if number is incorrect; if custom, use format name:id */}}
{{/* End of configurable values */}}

{{/* No Touchy */}}
{{ $db := or (dbGet 0 "Counting").Value (sdict "Last" (sdict "User" 204255221017214977 "Msg" 0) "Next" 1 "HighScore" (sdict "User" 204255221017214977 "Num" 1 "Time" currentTime) "SecondChance" 2) }}

{{ with .ExecData }}
	{{ if not (getMessage nil .) }} {{/* Check if number was deleted */}}
		{{ sendMessage nil (cembed "description" (print (userArg $db.Last.User).Mention " deleted their number which was correct!\nThe next number is " $db.Next) "color" 30654) }}
	{{ return }}{{ end }}
{{ return }}{{ end }}

{{ $Number = toInt (round (slice ($Number := (exec "calc" (index .Args 0))) 9 (sub (len $Number) 1))) }}

{{ if or (ne $db.Last.User .User.ID) $CountTwice }} {{/* Checks user */}}
	{{ if eq $db.Next $Number }} {{/* Checks if correct number */}}
		{{ $db.Set "Next" (add $db.Next 1) }}
		{{ try }}
			{{ if $Reactions }}
				{{ addReactions $CorrectEmoji }}
				{{ if and $ReactionDelete (ne $db.Last.Msg 0) }}
					{{ deleteMessageReaction nil $db.Last.Msg 204255221017214977 $CorrectEmoji }}
				{{ end }}
				{{ if not (mod $Number 100) }}
					{{ addReactions "💯" }}
				{{ end }}
			{{ end }}
		{{ catch }}
			{{ sendMessage $ErrorCID (print "Counting: `" .Error "`") }}
		{{ end }}
		{{ with $CorrectRID }} {{/* Correct Counting role */}}
			{{ takeRoleID $db.Last.User . }}{{ giveRoleID $.User.ID . }}
		{{ end }}
		{{ $db.Set "Last" (sdict "User" .User.ID "Msg" .Message.ID) }}
		{{ if $StatsCC }} {{/* Update leaderboard values */}}
			{{ $Silence := dbIncr .User.ID "CCorrect" 1 }}{{ $Silence = dbIncr .User.ID "CCount" 1 }}
			{{ if gt $Number $db.HighScore.Num }}{{ $db.Set "HighScore" (sdict "User" .User.ID "Num" $Number "Time" currentTime) }}{{ end }}
		{{ end }}
		{{ dbSet 0 "Counting" $db }}
		{{ execCC .CCID nil 10 .Message.ID }} {{/* Call to check if message was deleted */}}
		
	{{ else }} {{/* Wrong number */}}
		{{ $db.Set "SecondChance" (sub $db.SecondChance 1) }}
		{{ with $CorrectRID }}{{ takeRoleID $db.Last.User . }}{{ end }}
		{{ with $IncorrectRID }}{{ addRoleID . }}{{ removeRoleID . 259200 }}{{ end }}
		{{ if and $SecondChance (gt $db.SecondChance 0) }} {{/* Second Chance */}}
			{{ try }}{{ if $Reactions }}{{ addReactions $WarningEmoji }}{{ end }}
			{{ catch }}{{ sendMessage $ErrorCID (print "Counting: `" .Error "`") }}{{ end }}
			{{ $db.Set "Last" (sdict "User" .User.ID "Msg" .Message.ID) }}{{ dbSet 0 "Counting" $db }}
			{{ sendMessage nil (cembed "description" (print .User.Username " sent an incorrect number of " $Number "\n**But**, second chance saved the count!\nThe next number is " $db.Next) "color" 16744192) }}

		{{ else }} {{/* Reset count */}}
			{{ sendMessage nil (cembed "description" (print .User.Mention " sent an incorrect number of " $Number "\nCorrect number was " $db.Next "\nStart over at 1 🙃") "color" 16711680) }}
			{{ $db.Set "Last" (sdict "User" 204255221017214977 "Msg" 0) }}{{ $db.Set "Next" 1 }}{{ $db.Set "SecondChance" 2 }}
			{{ dbSet 0 "Counting" $db }}
			{{ if $StatsCC }}{{ $Silence := dbIncr .User.ID "CCount" 1 }}{{ end }}
			{{ try }}{{ if $Reactions }}{{ addReactions $IncorrectEmoji }}{{ end }}
			{{ catch }}{{ sendMessage $ErrorCID (print "Counting: `" .Error "`") }}{{ end }}
		{{ end }}
	{{ end }}

{{ else }} {{/* Same user */}}
	{{ sendMessage nil (cembed "description" (print "You can't count twice in a row 🥲\nThe next number is " $db.Next) "color" 16744192) }}
{{ end }}
