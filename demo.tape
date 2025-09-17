Output demo.gif

Require echo

Set Shell "zsh"
Set FontSize 18
Set Width 1200
Set Height 800

Hide

Type "setopt INTERACTIVE_COMMENTS" Enter
Type "clear" Enter

Show

Type "einsicht --help" Enter

Sleep 1s

Type "# Show mail content" Sleep 500ms Enter
Type "einsicht list < ../mail.long.eml" Sleep 500ms Enter

Sleep 2s

Type "# Print body or attachments" Sleep 500ms Enter
Type "einsicht print < ../mail.long.eml" Sleep 500ms Enter

Sleep 1s

Type "# Open the TUI" Sleep 500ms Enter
Type "einsicht -b html < ../mail.long.eml" Sleep 500ms Enter

Sleep 2s

Tab Sleep 2s

Type "h" Sleep 2s

Type "q" Sleep 1s
