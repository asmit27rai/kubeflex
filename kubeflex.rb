# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kubeflex < Formula
  desc ""
  homepage "https://github.com/kubestellar/kubeflex"
  version "0.6.1"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.6.1/kubeflex_0.6.1_darwin_amd64.tar.gz"
      sha256 "ca3f08470d5f1f1b4472449dc6f35dfb256cae47eba81b5e437aad4247600a2a"

      def install
        bin.install "bin/kflex"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.6.1/kubeflex_0.6.1_darwin_arm64.tar.gz"
      sha256 "aa2ae194de459ece9618a84dc85af5160cd445cdefe9b1623f53480a6382c5a4"

      def install
        bin.install "bin/kflex"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.6.1/kubeflex_0.6.1_linux_amd64.tar.gz"
      sha256 "47e6cca55886fb0b7ba9d4d0da91d8d39b022b5f8e72322cfe5f673177ba6449"

      def install
        bin.install "bin/kflex"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.6.1/kubeflex_0.6.1_linux_arm64.tar.gz"
      sha256 "4b25b669e57ce663ea2c4ede54fbf99e1751e9c6348fc4a9d37b32699fe24690"

      def install
        bin.install "bin/kflex"
      end
    end
  end
end
