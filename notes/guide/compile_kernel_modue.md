独立编译内核模块
--------

### step1： 下载对应的内核源码

### 在内核源码下执行：
```
    make oldconfig
    make prepare
    make scripts
```
### 在对应内核模块的源码目录下，查看makefile文件 找到
    obj-$(xxx) 括号内的宏定义


### 在内核源码根目录下，执行编译命令
make xxx=m -C ${内核源码根目录} M=${内核模块源码的目录} modules


FAQ
---
+ 这一步可能会提示scripts/sign-file.c:23:30: fatal error: openssl/opensslv.h: No such file or directory
    
    #安装libssl-dev软件包即可。

+ 问题现象
    
    modinfo *.ko
    
    modinfo: can't open '/lib/modules/<linux version number>/modules.dep': No such file or directo
    
    解决办法：
    
    cd /lib/
    
    mkdir modules
    
    cd modules/
    
    mkdir <linux version number>
    
    cd <linux version number>
    
    cp *.ko ./
    
    depmod
    
    mv modules.dep.bb modules.dep
    
    modinfo *.ko
+ insmod: can't insert 'xx.ko': invalid module format
    问题背景和解决方法
    这个问题源于那时候我用于编译驱动的内核和运行在开发板上的内核镜像配置不同导致的。后来参考了网上的博客，导致这个问题的原因其实挺多的，这里总结如下：
    
    1.用于编译内核的交叉编译器和内核不匹配，导致有些参数不兼容；
    
    对于原因1，可以参考这篇博客的博主：
    
    http://blog.csdn.net/stephen_yu/article/details/24481489（文章引自他人博客，特此声明）
    
    2.编译驱动的内核和运行在开发板上的内核版本不匹配；
    
    对于原因2，解决方法是：保证编译驱动的内核版本和运行在开发板上的内核版本是一致的就可以；
    
    3.内核版本一致，但是内核配置文件.config不一样导致的错误；
    
    原因3也是我当初遇到的一个情况，这种情况下，只要把运行在开发板上的内核源码中的.config文件覆盖用于编译驱动的内核源码，然后重新编译驱动就可以，其实原因3完全可以避免，我当初遇到这种情况的原因是因为在公司每个人负责的模块不同，配置内核的需求不同，而我后来是在自己的环境下编译的驱动模块放到了别人的开发板上运行，导致的这个问题。
    
    4.内核的版本检测配置选项导致驱动加载不了：  
    对于原因4，解决方法是配置内核，去掉版本检测配置选项，具体操作如下步骤： 
    4.1 配置内核：
    
    linux-3.14.38$make menuconfig
    make menuconfig—>Enable loadable module support  --->[*]   Module versioning support，把“Module versioning support”前面的星号去掉，禁止版本检测选项既可。
    
+ 安装模块时出现：[root@FriendlyARM nfs]# insmod key2.ko

key2: version magic '2.6.32.2 mod_unload modversions ARMv4 ' should be '2.6.32.2-FriendlyARM mod_unload ARMv4 '

这一行的意思就是说，当前插入的模块xxx.ko的版本信息(version magic)与正运行的kernel的版本信息不一致！应该是'2.6.32.2-FriendlyARM mod_unload ARMv4 '，而实际上xxx.ko的版本信息却是：'2.6.32.2 mod_unload modversions ARMv4 '； 显然它们之间差别是很小的。实际上，根据上面安装的kernel源码来看，它们应该是没有什么差别的。 所以，下面采用了一种比较极端的方式，强制xxx.ko的版本信息与运行的kernel保持一致。

修改/home/haiyang/linux-2.6.32.2/include/linux/tsrelease.h文件中的宏定义

     ＃define UTS_RELEASE "2.6.32.2

     为

     ＃define UTS_RELEASE "2.6.32.2-FriendlyARM”

     然后重新编译xxx.ko模块，这时候，它与内核的版本信息应该就是一致的了！试验下来确实如此,xxx.ko已经可以正常工作了！

但是可能会再次出现：

key2: version magic '2.6.32.2 mod_unload modversions ARMv4 ' should be '2.6.32.2-FriendlyARM mod_unload ARMv4 '这是因为arm公司在linux版本中加了自己的标志

此时需要修改：/home/haiyang/linux-2.6.32.2/include/linux/vermagic.h

[root@localhost linux]# vi vermagic.h

 

#include <linux/utsrelease.h>

#include <linux/module.h>

 

 

#ifdef CONFIG_SMP

#define MODULE_VERMAGIC_SMP "SMP "

#else

#define MODULE_VERMAGIC_SMP ""

#endif

#ifdef CONFIG_PREEMPT

#define MODULE_VERMAGIC_PREEMPT "preempt "

#else

#define MODULE_VERMAGIC_PREEMPT ""

#endif

#ifdef CONFIG_MODULE_UNLOAD

#define MODULE_VERMAGIC_MODULE_UNLOAD "mod_unload "

#else

#define MODULE_VERMAGIC_MODULE_UNLOAD ""

#endif

#ifdef CONFIG_MODVERSIONS

#define MODULE_VERMAGIC_MODVERSIONS "modversions "

#else

#define MODULE_VERMAGIC_MODVERSIONS ""

#endif

#ifndef MODULE_ARCH_VERMAGIC

#define MODULE_ARCH_VERMAGIC ""

#endif

 

#define VERMAGIC_STRING                                                

        UTS_RELEASE " "                                                

        MODULE_VERMAGIC_SMP MODULE_VERMAGIC_PREEMPT                    

        MODULE_VERMAGIC_MODULE_UNLOAD MODULE_VERMAGIC_MODVERSIONS      

        MODULE_ARCH_VERMAGIC

 

"vermagic.h" 34L, 837C

 

其中，VERMAGIC_STRING就是内核的版本信息，每个kernel module的版本信息就是从源代码树中的该宏定义获取的。所以，编译模块的时候一定要和实际使用该模块的内核的源代码树保持一致！不要张冠李戴，否则就会在加载模块的时候出现上述问题！

 

我的修改是#define MODULE_VERMAGIC_MODVERSIONS "modversions "

#define MODULE_VERMAGIC_MODVERSIONS " "

 

NOTES: 本方法并不是正规的解决办法，我是由于不想重新编译安装linux kernel，并且能够确保当前源码树中的kernel版本与系统运行的kernel版本是相同的情况下采用的权宜之计。如果不能保证这一点，最好不要采用这种方法。以免产生一些莫名其妙的问题！！

 

实际上最好的办法是：

要在menuconfig中改

解决办法如下：

在配置单中添加如下信息
General setup   --->

· Prompt for development and/or incomplete code/drivers
(-EmbedSky) Local version - append to kernel release
内核版本的差异导致的。

例如

（做s3c2410_ts触摸屏的驱动，程序交叉编译后，insmod出现如下的错误：
s3c2410_ts: version magic '2.6.29.4 mod_unload modversions ARMv5 'should be '2.6.29.4-FriendlyARM mod_unload ARMv4 '
linux kernel 中文组的朋友说是编译器版本的问题。重新配置linux kernel，并编译，更新开发板的内核。使用同样的编译器编译module，用nfs mount到开发板，insmod就OK了。果然是内核与module有冲突。）
