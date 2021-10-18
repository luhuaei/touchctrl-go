{ lib, pkgs, config, ... }:

with lib;

let
  cfg = config.services.touchctrl;
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
        User = "nobody";
        NoNewPrivileges = true;
        ExecStart = "${pkgs.touchctrl-go}/bin/touchctrl";
      };
    };

    environment.systemPackages = [ pkgs.touchctrl-go ];
  };
}
