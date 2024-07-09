class Inblog < Formula
  desc "Write a blog from your email"
  homepage "https://github.com/weebney/inblog"
  url "https://github.com/weebney/inblog/archive/v1.0.tar.gz"
  sha256 "the_actual_sha256_of_your_tarball"
  license "BSD-2-Clause"

  depends_on "go" => :build
  depends_on "make" => :build

  def install
    system "make"
    bin.install "inblog"
  end
end
