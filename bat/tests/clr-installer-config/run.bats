#!/usr/bin/env bats

# shared test functions
load ../../lib/mixerlib

setup() {
  global_setup
}

@test "Build image using default clr-installer config" {
  mixer $MIXARGS init --clear-version $CLRVER --mix-version 10
  mixer $MIXARGS bundle add editors

  # `build image` should fail because the prerequisite commands are not executed.
  # This is intentional in order to reduce test execution time.
  # The goal is to unit test the creation of clr-installer config file with relevant bundles.
  run sudo mixer $MIXARGS build image
  [[ "$status" -ne 0 ]]
  [[ "$output" =~ "release-image-config.yaml not found" ]]
  [[ "$output" =~ "Updating image bundle list based on mixbundles" ]]
}

@test "Update clr-installer config file with mix bundle list" {
  run cat release-image-config.yaml
  [[ "$status" -eq 0 ]]
  [[ "$output" =~ "os-core" ]]
  [[ "$output" =~ "os-core-update" ]]
  [[ "$output" =~ "kernel-native" ]]
  [[ "$output" =~ "bootloader" ]]
  [[ "$output" =~ "editors" ]]
}

# vi: ft=sh ts=8 sw=2 sts=2 et tw=80
