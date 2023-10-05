// Code generated by running "go generate". DO NOT EDIT.

// Copyright 2019 The Wuffs Authors.
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// https://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or https://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package main

const usageStr = `ractool manipulates Random Access Compression (RAC) files.

Random access means that it is possible to reconstruct part of the decompressed
file, starting at a given offset into the decompressed file, without always
having to first decompress all of the preceding data.

In comparison to some other popular compression formats, all four of the Zlib,
Brotli, LZ4 and Zstandard specifications explicitly contain the identical
phrase: "the data format defined by this specification does not attempt to
allow random access to compressed data".

See the RAC specification for more details:
https://github.com/google/wuffs/blob/main/doc/spec/rac-spec.md

Usage:

    ractool [flags] [input_filename]

If no input_filename is given, stdin is used. Either way, output is written to
stdout.

The flags should include exactly one of -decode or -encode.

By default, a RAC file's chunks are decoded in parallel, using more total CPU
time to substantially reduce the real (wall clock) time taken. Batch (instead
of interactive) processing of many RAC files may want to pass -singlethreaded
to prefer minimizing total CPU time.

When encoding, the input is partitioned into chunks and each chunk is
compressed independently. You can specify the target chunk size in terms of
either its compressed size or decompressed size. By default (if both
-cchunksize and -dchunksize are zero), a 64KiB -dchunksize is used.

You can also specify a -cpagesize, which is similar to but not exactly the same
concept as alignment. If non-zero, padding is inserted into the output to
minimize the number of pages that each chunk occupies. Look for "CPageSize" in
the "package rac" documentation for more details:
https://godoc.org/github.com/google/wuffs/lib/rac

A RAC file consists of an index and the chunks. The index may be either at the
start or at the end of the file. At the start results in slightly smaller and
slightly more efficient RAC files, but the encoding process needs more memory
or temporary disk space.

Examples:

    ractool -decode foo.rac | sha256sum
    ractool -decode -drange=400..500 foo.rac
    ractool -encode foo.dat > foo.rac
    ractool -encode -codec=zlib -dchunksize=256k foo.dat > foo.rac

The "400..500" flag value means the 100 bytes ranging from a DSpace offset
(offset in terms of decompressed bytes, not compressed bytes) of 400
(inclusive) to 500 (exclusive). Either or both bounds may be omitted, similar
to Rust slice syntax. A "400.." flag value would mean ranging from 400
(inclusive) to the end of the decompressed file.

The "256k" flag value means 256 kibibytes (262144 bytes), as does "256K".
Similarly, "1m" and "1M" both mean 1 mebibyte (1048576 bytes).

General Flags:

    -decode
        whether to decode the input
    -encode
        whether to encode the input
    -quiet
        whether to suppress messages

Decode-Related Flags:

    -drange
        the "i..j" range to decompress, "..8" means the first 8 bytes
    -singlethreaded
        whether to decode on a single execution thread

Encode-Related Flags:

    -cchunksize
        the chunk size (in CSpace)
    -codec
        the compression codec (default "zstd")
    -cpagesize
        the page size (in CSpace)
    -dchunksize
        the chunk size (in DSpace)
    -indexlocation
        the index location, "start" or "end" (default "start")
    -resources
        comma-separated list of resource files, such as shared dictionaries
    -tmpdir
        directory (e.g. $TMPDIR) for intermediate work; empty means in-memory

Codecs:

    lz4
    zlib
    zstd

Only zlib is fully supported. The others will work for the flags' default
values, but they (1) don't support -cchunksize, only -dchunksize, and (2) don't
support -resources. See https://github.com/google/wuffs/issues/23 for more
details.

Installation:

Like any other implemented-in-Go program, to install the ractool program:

    go install github.com/google/wuffs/cmd/ractool

Extended Example:

    --------
    $ # Fetch and unzip the enwik8 test file, a sample of Wikipedia.
    $ wget http://mattmahoney.net/dc/enwik8.zip
    $ unzip enwik8.zip

    $ # Also zstd-encode it, as a reference point. Using compression level 15,
    $ # instead of the default of 3, matches what ractool uses.
    $ zstd -15 enwik8

    $ # Create a shared dictionary. Using zstd-the-program produces a
    $ # dictionary that is especially useful for zstd-the-format, but it can
    $ # also be used by other formats as a 'raw' prefix dictionary.
    $ zstd -15 --train -B64K --maxdict=32K -o dict.dat enwik8

    $ # RAC-encode it with various codecs, with and without that dictionary.
    $ ractool -encode -codec=zlib -resources=dict.dat enwik8 > zlib.withdict.rac
    $ ractool -encode -codec=zlib                     enwik8 > zlib.sansdict.rac
    $ ractool -encode -codec=zstd -resources=dict.dat enwik8 > zstd.withdict.rac
    $ ractool -encode -codec=zstd                     enwik8 > zstd.sansdict.rac
    $ ractool -encode -codec=lz4                      enwik8 > lz4.sansdict.rac

    $ # The size overhead (comparing RAC+Xxx to Xxx) is about 0.2% (with) or
    $ # 4.8% (sans) for zlib/zip and about 13% (with) or 28% (sans) for zstd,
    $ # depending on whether we used a shared dictionary (with or sans).
    $ ls -l
    total 362080
    -rw-r----- 1 tao tao     32768 Oct 25 10:10 dict.dat
    -rw-r----- 1 tao tao 100000000 Jun  2  2011 enwik8
    -rw-r----- 1 tao tao  36445475 Sep  2  2011 enwik8.zip
    -rw-r----- 1 tao tao  29563109 Jun  2  2011 enwik8.zst
    -rw-r----- 1 tao tao  58813316 Oct 25 10:17 lz4.sansdict.rac
    -rw-r----- 1 tao tao  38185178 Oct 25 10:16 zlib.sansdict.rac
    -rw-r----- 1 tao tao  36505786 Oct 25 10:16 zlib.withdict.rac
    -rw-r----- 1 tao tao  37820491 Oct 25 10:17 zstd.sansdict.rac
    -rw-r----- 1 tao tao  33386395 Oct 25 10:17 zstd.withdict.rac

    $ # Check that the decompressed forms all match.
    $ cat enwik8                            | sha256sum
    2b49720ec4d78c3c9fabaee6e4179a5e997302b3a70029f30f2d582218c024a8  -
    $ unzip -p enwik8.zip                   | sha256sum
    2b49720ec4d78c3c9fabaee6e4179a5e997302b3a70029f30f2d582218c024a8  -
    $ unzstd --stdout enwik8.zst            | sha256sum
    2b49720ec4d78c3c9fabaee6e4179a5e997302b3a70029f30f2d582218c024a8  -
    $ for f in *.rac; do ractool -decode $f | sha256sum; done
    2b49720ec4d78c3c9fabaee6e4179a5e997302b3a70029f30f2d582218c024a8  -
    2b49720ec4d78c3c9fabaee6e4179a5e997302b3a70029f30f2d582218c024a8  -
    2b49720ec4d78c3c9fabaee6e4179a5e997302b3a70029f30f2d582218c024a8  -
    2b49720ec4d78c3c9fabaee6e4179a5e997302b3a70029f30f2d582218c024a8  -
    2b49720ec4d78c3c9fabaee6e4179a5e997302b3a70029f30f2d582218c024a8  -

    $ # Compare how long it takes to produce 8 bytes from the middle of
    $ # the decompressed file, which happens to be the word "Business".
    $ time unzip -p enwik8.zip | dd if=/dev/stdin status=none \
    >     iflag=skip_bytes,count_bytes skip=50000000 count=8
    Business
    real    0m0.379s
    user    0m0.410s
    sys     0m0.080s
    $ time unzstd --stdout enwik8.zst | dd if=/dev/stdin status=none \
    >     iflag=skip_bytes,count_bytes skip=50000000 count=8
    Business
    real    0m0.172s
    user    0m0.141s
    sys     0m0.103s
    $ time ractool -decode -drange=50000000..50000008 zstd.withdict.rac
    Business
    real    0m0.004s
    user    0m0.005s
    sys     0m0.001s

    $ # A RAC file's chunks can be decoded in parallel, unlike ZIP,
    $ # substantially reducing the real (wall clock) time taken even
    $ # though both of these files use DEFLATE (RFC 1951) compression.
    $ #
    $ # Comparing the -singlethreaded time suggests that zlib-the-library's
    $ # DEFLATE implementation is faster than unzip's.
    $ time unzip -p                        enwik8.zip        > /dev/null
    real    0m0.711s
    user    0m0.690s
    sys     0m0.021s
    $ time ractool -decode -singlethreaded zlib.withdict.rac > /dev/null
    real    0m0.519s
    user    0m0.513s
    sys     0m0.017s
    $ time ractool -decode                 zlib.withdict.rac > /dev/null
    real    0m0.052s
    user    0m0.678s
    sys     0m0.036s

    $ # A similar comparison can be made for Zstandard.
    $ time unzstd --stdout                 enwik8.zst        > /dev/null
    real    0m0.203s
    user    0m0.187s
    sys     0m0.016s
    $ time ractool -decode -singlethreaded zstd.withdict.rac > /dev/null
    real    0m0.235s
    user    0m0.206s
    sys     0m0.033s
    $ time ractool -decode                 zstd.withdict.rac > /dev/null
    real    0m0.037s
    user    0m0.374s
    sys     0m0.080s

    $ # For reference, LZ4 numbers.
    $ time ractool -decode -singlethreaded lz4.sansdict.rac  > /dev/null
    real    0m0.072s
    user    0m0.053s
    sys     0m0.021s
    $ time ractool -decode                 lz4.sansdict.rac  > /dev/null
    real    0m0.024s
    user    0m0.097s
    sys     0m0.034s
    --------
`
