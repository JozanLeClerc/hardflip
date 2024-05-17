# hardflip

The best TUI tool to connect to your distant hosts. For Linux/BSD. Might work
on macOS.

## Dependencies

Install those if you need them:

+ `sshpass` if you are using passwords with SSH (honnestly use keys instead)
+ `xfreerdp` for RDP. Called `freerdp2-x11` on some distros
+ `openstack` for OpenStack CLI.
+ `gpg` to crypt passwords. You can store them in plain text if you prefer but
I wouldn't recommend that option.
+ `go`
+ GNU `make`

## Install

To install `hardflip`, run those commands in your shell:

```sh
git clone git://gitjoe.xyz/jozan/hardflip
cd hardflip
make
sudo make install
make clean
```

Change this line:

```make
DEST			:= /usr
```

if you want to install stuff some other place

## Config
