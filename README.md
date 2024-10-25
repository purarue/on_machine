## on_machine

A helper CLI tool to determine which computer you're currently on. Often in ones dotfiles or personal scripts, you do something like:

```bash
case "$(uname -s)" in
  Linux) command_on_linux ;;
  Darwin) command_on_mac ;;
  *) echo 'unknown...' ;;
esac
```

...to run particular commands on different machines/installs.

This is fine if you're always going to have two installs, but as you add more it starts to get more and more complicated. For example, I have:

- multiple linux installs (for example one on `Ubuntu` and another on `Arch`)
- [termux](https://termux.com/) on my phone
- [`wsl`](https://docs.microsoft.com/en-us/windows/wsl/install) will also return `Linux`, when you likely want to do something custom on windows

So, `on_machine` generates a unique-enough fingerprint of your system (which you can tune to be as simple or complicated as you want), so you can do:

```bash
case "$(on_machine)" in
  linux_arch_*) command_on_arch ;;
  linux_ubuntu_*) command_on_ubuntu ;;
  android_termux_*) command_on_termux ;;
  windows_*) command_on_wsl ;;
  mac_*) command_on_mac ;;
esac
```

This borrows a lot of ideas from tools like [`neofetch`](https://github.com/dylanaraps/neofetch) to figure out what operating system/distribution one is using

### Install

Using `go install` to put it on your `$GOBIN`:

`go install github.com/purarue/on_machine/cmd/on_machine@latest`

I recommend you have both `uname` and `lsb_release` installed if possible on `linux`, that makes distribution detection much nicer. Otherwise, this defaults to the `golang` `runtime` module defaults

To manually build:

```bash
git clone https://github.com/purarue/on_machine
cd ./on_machine
go build ./cmd/on_machine
# copy binary somewhere on your $PATH
sudo cp ./on_machine /usr/local/bin
```

## Usage

```
usage: on_machine [-h] [-cmd <print|match>] [OPTIONS] [PATTERN]

Tool to determine which operating system/machine you're on.

Commands:
print [default]: prints the computed fingerprint after interpolating the pattern '%o_%d_%h'
match: does directory/path matching based on the pattern, changes the default pattern to '%o/%d/%h'

print
---
PATTERN is a printf-styled format string, supporting the following sequences:

%o - Operating System (using uname)
%d - Distro (using lsb_release)
%h - Hostname (name of the computer)
%a - Arch (detected by golang)
%O - Golang OS (unmodified golang detected operating system)

By default, this uses '%o_%d_%h'

match
---
Directory/path matching, Uses the pattern to match directory structures.
Can provide the base path to use with -base, that replaces '/' with
OS-specific path separator in the pattern. For more information, see the docs:
https://github.com/purarue/on_machine

Options:

  -base string
    	base directory to use to match paths
  -cmd string
    	on_machine command to run (default "print")
  -delimiter string
    	delimiter to print between matches (default "\n")
  -filter string
    	filter matches to either 'dir' or 'file'
  -json
    	print matches as a JSON array
  -print0
    	use the null character as the delimiter
  -skip-last-delim
    	dont print the delimiter after the last match
```

For basic usage, can just run `on_machine` with no arguments, using a `case` statements

```bash
case "$(on_machine)" in
  linux_arch_*) command_on_arch ;;
  linux_ubuntu_*) command_on_ubuntu ;;
  android_termux_*) command_on_termux ;;
  windows_*) command_on_wsl ;;
  mac_*) command_on_mac ;;
esac

# or, for example, to skip running a command on android and run everywhere else
case "$(on_machine)" in
  android_*) ;;
  *) command_to_run_everywhere_else ;;
esac
```

Sometimes the `hostname` cannot be determined (e.g. on android), so if you'd like you can set something custom, you can use the `ON_MACHINE_HOSTNAME`, e.g., in your shell startup:

```bash
ON_OS="$(on_machine)"
[[ "$ON_OS" =~ ^android.* ]] && export ON_MACHINE_HOSTNAME="phone"
```

```bash
$ on_machine
android_termux_phone
```

### match

This can be thought of as an alternative to the above, its a different way to figure out what code to run on different machines, by placing scripts in a particular directory structure

You can manually do case/regex statements in bash (and often that is enough), but in some cases that can become complicated. I use this to match [background processes](https://github.com/purarue/bgproc) scripts on different machines -- to figure out which bash scripts to run. Those are organized like:

```
matching_examples/dir_based
├── all
│   └── this_dir_is_matched_everytime
├── android
│   └── on_android
├── linux
│   ├── arch
│   │   ├── home
│   │   │   └── on_arch_when_hostname_is_home
│   │   ├── only_on_arch
│   │   └── work
│   │       └── on_arch_when_hostname_is_work
│   ├── matched_on_any_linux_install
│   └── ubuntu
│       └── only_on_ubuntu
└── mac
    └── on_mac
```

Then, say my hostname is `home`, and I'm on my `arch` machine. The `match` command computes the following:

```bash
$ on_machine -cmd match -base ./matching_examples/dir_based '%o/%d/%h'
/home/username/Repos/on_machine/matching_examples/dir_based/all
/home/username/Repos/on_machine/matching_examples/dir_based/linux
/home/username/Repos/on_machine/matching_examples/dir_based/linux/arch
/home/username/Repos/on_machine/matching_examples/dir_based/linux/arch/home
```

Note: `all` is like a `*`, its always matched -- so that would be where I store shared jobs.

[`bgproc_on_machine`](https://github.com/purarue/bgproc/blob/master/bgproc_on_machine) uses multiple directories to organize different jobs. That uses on_machine internally like:

```bash
MATCHES=()
while read -r -d $'\0' match; do
	MATCHES+=("$match")
done < <(on_machine -cmd match -print0 -filter dir -base "$base" '%o/%d')
```

... which figures out which directories/scripts to include when running jobs:

```
$ bgproc_on_machine -o
1655993990:Searching for jobs in:
1655993990:/home/username/data/jobs/all
1655993990:/home/username/data/jobs/linux
1655993990:/home/username/.local/scripts/supervisor_jobs/all
1655993990:/home/username/.local/scripts/supervisor_jobs/linux
1655993990:/home/username/Repos/HPI-personal/jobs/all
1655993990:/home/username/Repos/HPI-personal/jobs/linux
```

For a real example of one of these directory structures, see [my HPI jobs](https://github.com/purarue/HPI-personal/tree/master/jobs)

If the pattern includes an extension, this extracts that and tries to match at each level going down. For example, for setting up your `~/.zshrc` setup: I want some code that runs everywhere (`all.zsh`), some that runs on `android`, some that runs on `linux`, and then additional code that runs on `linux` and `arch`. So, given:

```
./matching_examples/with_extensions
├── all.zsh
├── android.zsh
├── linux
│   ├── arch.zsh
│   └── ubuntu.zsh
├── linux.zsh
└── mac.zsh
```

On my arch machine, using the pattern `%o/%d.zsh`, this matches:

```bash
$ on_machine -cmd match -base ./matching_examples/with_extensions -filter file '%o/%d.zsh
/home/username/Repos/on_machine/matching_examples/with_extensions/all.zsh
/home/username/Repos/on_machine/matching_examples/with_extensions/linux/arch.zsh
/home/username/Repos/on_machine/matching_examples/with_extensions/linux.zsh
```

Essentially, this lets me pick what scripts to run organized as a directory, instead of ever-growing `case` statement in a `bash` script somewhere.

## Contributing

I test this on all my machines, but it gets increasingly difficult to test on systems I don't have access to.

If you're able to test this on some operating system not listed above, or use one this doesn't support, am happy to accept PRs for new operating systems, or possible strategies to detect new systems based on some command output/metadata file that exists
