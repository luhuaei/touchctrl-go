
# nix

```
sudo nix-channel --add https://github.com/luhuaei/touchctrl-go/archive/master.tar.gz touchctrl-go
nix-channel --update
```

Enable `touchctrl` service on your `configuration.nix`

```
{ config, pkgs, ... }:

{
  imports =
    [  <touchctrl-go/touchctrl.nix> ];
  services.touchctrl.enable = true;
}
```

# non-nix
Your need to download this repo to compile, and manually create `touchctrl` service.
