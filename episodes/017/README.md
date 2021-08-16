# eCHO episode 17: CVE 2021 3490
[Live](https://youtu.be/VZ1V2nMvQH4)

### Headlines
[eBPF Summit!](https://ebpf.io/summit-2021/)
[eBPF Foundation!!](https://www.isovalent.com/blog/post/2021-08-ebpf-foundation-announcement)
GitHub picks Friday 13th to kill off password-based Git authentication

### Outline
* ebpf ctf coming with Tabitha.

* talk about research in general.
* review the CVE and the fix. 
* Discuss the blog post and why it's interesting.
* [Github link](https://github.com/chompie1337/Linux_LPE_eBPF_CVE-2021-3490)
* [Valentina's Blog](https://www.graplsecurity.com/post/kernel-pwning-with-ebpf-a-love-story)
* I've got a couple of vms prepared to play with things.
* Things needed to make this work.
    * very specific kernel versions
    * sysctl -a | grep bpf 
        * kernel.unprivileged_bpf_disabled = 1 <- user space bpf disabled. 
    * a userspace bpf program that enables you to escalate to root.

* Fun with capsh and kind and docker.



### References
* [why does android give this exploit such a low sev rating when nist gave it a severe?](https://twitter.com/jeffvanderstoep/status/1422771606309335043)

* [The original CVE](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-3490)
* [Manfred's blog on the original work](https://www.zerodayinitiative.com/blog/2020/4/8/cve-2020-8835-linux-kernel-privilege-escalation-via-improper-ebpf-program-verification)
* [The alu32 big fixed by Daniel](https://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf.git/commit/?id=049c4e13714ecbca567b4d5f6d563f05d431c80e)

* [changing the defaults in the kernel](https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=08389d888287c3823f80b0216766b71e17f0aba5)
* Two exciting DefCon talks on eBPF [eBPF - I thought we were friends](https://www.youtube.com/watch?v=5zixNDolLrg) and [Warping Reality: Creating and Countering the Next Generation of Linux Rootkits](https://youtu.be/g6SKWT7sROQ)


More context on the safety of eBPF: 
[Safe Programs The Foundation of BPF - Alexei Starovoitov, Facebook](https://www.youtube.com/watch?v=AV8xY318rtc)

Also by Alexei [CAP_BPF](https://lwn.net/Articles/820560/)