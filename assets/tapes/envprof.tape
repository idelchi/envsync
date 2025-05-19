Output gifs/envprof.gif

Set Theme "CLRS"
Set FontSize 16
Set Width 1200
Set Height 800
Set TypingSpeed 100ms

Env SHELL "bash"

Type 'envprof ls'

Enter
Enter
Sleep 1s

Type 'envprof ls dev -v'

Enter
Enter
Sleep 1s

Type 'echo HOST=$HOST'

Enter
Enter
Sleep 1s

Type 'envprof shell dev'

Enter
Sleep 1s


Type 'echo HOST=$HOST'

Enter
Sleep 1s

Type 'exit'

Enter
Sleep 1s

Type 'envprof export dev dev.env'

Enter
Enter
Sleep 1s

Type 'cat dev.env'

Enter
Enter
Sleep 1s

Sleep 3s
