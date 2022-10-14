To build the static libraries for each platform, run cmake and make:

$ cmake -DBUILD_SHARED_LIBS=OFF -DHAVE_QSORT_R_GNU=0
$ make

And then copy the libgit2.a over the existing file, corresponding to
the platform built.

Note, if compiling for Linux amd64, the environment should have
OpenSSL 1.1+ installed. If not, the compiled library might not work
inside Alpine container images. Similarly, disabling libc provided
qsort is important to guarantee compability with Alpine's musl libc.