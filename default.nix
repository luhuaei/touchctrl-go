{ pkgs ? import <nixpkgs> { } }:

pkgs.buildGoModule {
  name = "touchctrl-go";
  version = "0.0.1";
  src = ./.;

  doCheck = false;

  buildPhase = ''
    go build -o touchctrl
  '';

  installPhase = ''
    install -Dm755 touchctrl -t $out/bin
  '';

  # nix-prefetch '{ sha256 }: (callPackage (import ./default.nix) { }).go-modules.overrideAttrs (_: { modSha256 = sha256; })'
  vendorSha256 = "155vd5cz7i1w3fmg4aadw20p5ynd8g00xw8rnpsrd6s923bdaya2";
}
