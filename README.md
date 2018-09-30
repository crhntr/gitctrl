# gitctrl
Some Scripts to Control Many Git Repositories

## Problems this tool trys to "Solve"

### Updating Remote URLs from HTTPs to SSH

Migrating from using HTTP(S) remote servers to using SSH is a good thing (I guess), 
but with many cloned repos going to every directory and running 

```bash
pushd $SOME_REPO
  git remote -v
  git remote rm origin
  git remote add origin git@github.com:crhntr/gitctrl
popd
```
can take a little time. So I decided to try out a cool looking git library, 
[go-git](https://github.com/src-d/go-git), to try to automate this problem.

### Viewing the Status of Many Repos

I'm not always the best at remembering to commit and push my work. 
I wanted to see what the status of all my repo's were so I could prioritise 
finishing WIP and push it. So this tool does this. _Also, I haven't written much 
concurent Go recently so after noticing how long some repos took to return a status
I decided this was an opertunity to practice using some go routines._

## Usage

### `gitctrl statuses`

Go to some directory and run this command. The tool with walk the filesystem and when 
it finds a directory with a git repo it will print out the status.

### `gitctrl remote-origin-must-ssh`

Go to some directory and run this command. The tool with walk the filesystem and when 
it finds a directory with a git repo that has a remote named ***origin*** and has the prefix
`https://github.com` it will convert the URL and log the change as follows:

```
https://github.com/crhntr/gitctrl --> git@github.com:crhntr/gitctrl
```


