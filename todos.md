1. BUG: Solve the issue with the writer(It should stop writing the loaded file input to the command.txt file)
2. Work on the tabspace feature -- autocomplete
3. Implement pipes features
4. Implement signal handling
5. Implement I/O redirection

# BUGS:
1. When I press spacebar after tab autocomplete the cursor goes back 2 characters back
2. Signal handling seems not to be working in cases where I want to exit from something like typing even though no text is being executed
3. When large text is being written in the terminal of the shell it reprints the prompt path and appends the text/commands and arguments being typed.
4. There also seems to be a bug such that my shell does not recognize some git commands expecially for large arguments such as commits



# FIRST TIME SETUP LOGIC TASKS
1. Implement session/remember the user everytime they start the shell to use it instead of prompting the user 
    a username and password evertime they exit and restart the shell
2. Append the username + host then the path on the full path: Example: username@host:~$
3. Implement max retries for user wrong password inputs
4. For special commands such as deleting a directory and its contents prompt the user for a password

# FEATURES TO BUILD
1. Environment variables
2. I/O redirection
3. Pipes
4. SuperUser commands

# FEATURE: ENVIRONMENT VARIABLES
1. Add logic for $$ execution as a command in my shell
2. Positional parameters
3. ${VAR} -- syntax with curly braces -- There is a bug in the system such that the prefix and suffix checking in the env variables is not finding the env variable even though I am parsing it correct

# FEATURE: INPUT/OUPUT REDIRECTION
1. Handle > for input redirection. Replaces the existing content in a file with a new input
2. Handle >> for output redirection with appending. Appends data to the end of a file.

# Strategy on building the I/O redirection
- Read the token and check for the >, >>, <, < symbols in the token
- For the part before the symbol execute the command
- For the part after the symbol open/read/append into the file
- Handle the above three logic for all the symbols.