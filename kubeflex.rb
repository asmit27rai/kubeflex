# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kubeflex < Formula
  desc ""
  homepage "https://github.com/kubestellar/kubeflex"
  version "0.6.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.6.0/kubeflex_0.6.0_darwin_amd64.tar.gz"
      sha256 "0aadf73e00bfbc5106500fb90a4bb1c82d12ea47ac39f78b13fa7d4cb65d415e"

      def install
        bin.install "bin/kflex"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.6.0/kubeflex_0.6.0_darwin_arm64.tar.gz"
      sha256 "65e256ad9798325d9a5b64d179f85eccc0526c0ebe1a5f6e4f2df9502b195814"

      def install
        bin.install "bin/kflex"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.6.0/kubeflex_0.6.0_linux_amd64.tar.gz"
      sha256 "08989ea14c2692f7534e52d916c260df276e04624d7c92ac63afc3388990f753"

      def install
        bin.install "bin/kflex"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.6.0/kubeflex_0.6.0_linux_arm64.tar.gz"
      sha256 "cfa4102106b260786dbb88ddaab12f23806dec97c687ebb4454b39cfb741800c"

      def install
        bin.install "bin/kflex"
      end
    end
  end
end
