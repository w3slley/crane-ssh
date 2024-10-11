# crane-ssh: a nicer way of configuring SSH keys
The common usecase for this tool, and in fact the reason I built it, is to make the process of setting SSH key pairs with platforms like Github or Gitlab way easier. You just run the tool and it takes care of setting up the public and private keys (using `ssh-keygen`). At the end, you'll have the public key copied to the clipboard and all that is left is to go to the platform and paste the public key there.

## Requirements
- Go (developed with v1.22.4)
- xclip or xsel (for Linux/Unix systems)

## Installation

To install, run:
```
curl -fsSL https://raw.githubusercontent.com/w3slley/crane-ssh/main/install.sh | bash
```

## Usage

Run:
```
crane-ssh generate --help
```

to see the necessary flags the command `generate` requires. If you want to create a new SSH key pair and add the public key to your Github account, you would do something like:
```
crane-ssh generate --host=github.com --alias=github.com --keyName=github
```
You could also simply run `crane-ssh generate` and add the values manually one by one. This would create a `github` and `github.pub` key pair in `.ssh/` and modify the `config` file to have:

```
Host github.com
  HostName github.com
  IdentityFile $HOME/.ssh/github
  Preferredauthentications publickey
  IdentitiesOnly yes
```

After the program finishes, you'll have the public key in your clipboard, so just go over Github and create a new SSH key. After that is done, check that it all works with:

```
ssh -T git@github.com
```

You should see `Hi <username>! You've successfully authenticated, but GitHub does not provide shell access.`. The same process applies to Bitbucket, Gitlab or any repository hosting manager tool or application that requires setting up SSH key pais as a method of authentication.

## Uninstall
To uninstall crane-ssh, run the `uninstall.sh` script. It will delete the executable from `.local/bin` and delete the `~/.crane-ssh` folder.
