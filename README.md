# touchctrl

When you finger tap touchpad `left top area` or `right top area`, touchctrl will press `ctrl`, so that you can tap touchpad press some combine action,  press `<tap>+c` equal `<ctrl>+c`;

`left top area` and `right top area` define is
```
var (
	// touch area absolute position
	LeftRect = Rect{
		TopLeft:     Point{X: 1, Y: 1},
		RightBottom: Point{X: 250, Y: 250},
	}
	RightRect = Rect{
		TopLeft:     Point{X: 1000, Y: 1},
		RightBottom: Point{X: 1500, Y: 250},
	}
)
```
You can change it on `manager.go` to adjust hot area.

# Install
## nix

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

## non-nix
Your need to download this repo to compile, and manually copy `touchctrl` service.

```
cd touchctrl-go
go build
sudo install -D -m 0755 touchctrl /usr/bin/
sudo install -D touchctrl.service /usr/lib/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start touchctrl.service
```
