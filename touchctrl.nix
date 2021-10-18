{ lib, pkgs, config, ... }:

with lib;

let
  cfg = config.services.touchctrl;
  touchctrl-go = pkgs.callPackage ./. { };
in
{
  options = {
    services.touchctrl = {
      enable = mkOption {
        default = false;
        type = types.bool;
        description = "Whether to run touchctrl server.";
      };
    };
  };

  config = mkIf cfg.enable {
    systemd.services.touchctrl = {
      description = "Use trouchpad press left ctrl";
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];

      serviceConfig = {
        User = "root";
        NoNewPrivileges = true;
        ExecStart = "${touchctrl-go}/bin/touchctrl";
      };
    };

    environment.systemPackages = [ touchctrl-go ];
  };
}
