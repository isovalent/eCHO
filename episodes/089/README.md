# Episode 89: Stack walking

## Headlines

* [KubeCon videos are up](https://youtube.com/playlist?list=PLj6h78yzYM2PyrvCoOii4rAopBswfz1p7)
* [CiliumCon too](https://youtube.com/playlist?list=PLj6h78yzYM2Meb36FX-bKd-3fpNvtlzpE)
* [KubeCon wrap-up blog](https://isovalent.com/blog/post/kubecon-europe-2023-wrap-up)
* [BPF dev stats for 6.4](https://lore.kernel.org/bpf/ZFAOojsT93ZxwNu3@google.com/t/#u)

## Stack walking with Parca

* [Parca, the continuous profiling project for cloud-native applications](https://www.parca.dev/).
* [In-depth blogpost on unwinding native code without frame pointers](https://www.polarsignals.com/blog/posts/2022/11/29/profiling-without-frame-pointers/).
* [BPF code for our native unwinder](https://github.com/parca-dev/parca-agent/tree/1dc4360bd006da325653a821ef6b84c5b33da3a3/bpf).
* [rbperf, a BPF profiler and tracer for Ruby](https://github.com/javierhonduco/rbperf).
* [Design Decisions of a Continuous Profiler](https://www.polarsignals.com/blog/posts/2022/12/14/design-of-continuous-profilers/).
* [CppCon 2017: Dave Watson “C++ Exceptions and Stack Unwinding”](https://www.youtube.com/watch?v=_Ivd3qzgT7U&t=2114s).
* [DWARF Debugging Information Format Version 5](https://dwarfstd.org/doc/DWARF5.pdf).
* [System V Application Binary Interface, AMD64 Architecture Processor Supplement](https://refspecs.linuxbase.org/elf/x86_64-abi-0.99.pdf).
* [Fast and Reliable DWARF Unwinding, and Beyond](https://fzn.fr/projects/frdwarf/frdwarf-oopsla19.pdf).
* [BPF Features by Linux Kernel Version](https://github.com/iovisor/bcc/blob/275aa3f3ea5236c6b556cbb94c583db105cb92fd/docs/kernel-versions.md).
* [Fedora Decides After All To Allow Default Compiler Flag To Help Debugging/Profiling](https://www.phoronix.com/news/F38-fno-omit-frame-pointer);
* [Writing ARM64 code for Apple platforms "The frame pointer register (x29) must always address a valid frame record."](https://developer.apple.com/documentation/xcode/writing-arm64-code-for-apple-platforms)
* [Golang enabling frame pointers by default](https://github.com/golang/go/issues/15840).
* Kernel unwinding:
    * [ORC](https://docs.kernel.org/x86/orc-unwinder.html)
    * [SFrame, WIP](https://lore.kernel.org/linux-toolchains/20230501200410.3973453-1-indu.bhagat@oracle.com/T/#t)

## Early BPF overhead metrics

Walking stacks of a host running Postgres, CPython, Ruby (MRI) applications (some with >90 frames)
````
P50: 285ns
P90: 370ns
Max: 428ns
````

(kernel 6.0.1`8 with Intel i7-8700K (late ‘17) )

We'll get more metrics this year.