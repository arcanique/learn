手动编译升级
在安装之前请确保自己的linux系统有足够的虚拟内存（建议1G）
增加swap虚拟内存

获取安装包并解压

地址是：
http://ftp.gnu.org/gnu/gcc/
https://mirrors.ustc.edu.cn/gnu/gcc/(国内)
解压缩
tar -jxvf gcc-9.1.0.tar.bz2
或则
tar -zxvf gcc-9.1.0.tar.gz
这两个版本选择下载时的版本进行不同的解压操作

2.下载供编译需求的依赖项
编译gcc 需要
[GMP](https://gmplib.org/)，
[MPFR](http://www.mpfr.org/),
[MPC](http://www.multiprecision.org/)
我们可以执行安装包里面的脚本去安装
cd gcc-9.1.0
./contrib/download_prerequisites　
如果不成功可以手动下载依赖
yum install gmp-devel                              // 編譯依賴此庫  
yum install mpfr-devel                             // 編譯依賴此庫  
yum install libmpc-devel                         // 編譯依賴此庫  
当然以上库可以在ftp://gcc.gnu.org/pub/gcc/infrastructure/下载

isl 选择安装(我执行./contrib/download_prerequisites成功，没有安装这个文件)
下载isl-0.18.tar.bz2                                  // yum 没有这个库源码安装
wget ftp://gcc.gnu.org/pub/gcc/infrastructure/isl-0.18.tar.bz2
tar -jxvf isl-0.18.tar.bz2
cd isl-0.18 
./configure   
make  
make install   


建立一个目录供编译出的文件存放

mkdir gcc-build-9
cd gcc-build-9

4.编译安装
yum groupinstall "Development Tools"
yum install glibc-static libstdc++-static（这两个是必要的开发环境）
../configure --prefix=/usr/local/gcc9 --enable-languages=c,c++,go  --disable-multilib // 指定gcc9安裝地址，指定所需安裝語言，不支持32位  
make -j4 (-j4选项是make对多核处理器的优化，如果不成功请使用 make，相关优化选项可以移步至参考文献[2]。建议不要使用make -j来编译，虽然可以缩短编译时间，但极大可能会编译失败)
make install((安装需要root权限!))

5.删除旧的gcc和g++
  rm -f /usr/bin/gcc
  rm -f /usr/bin/g++


链接新的gcc和g++

ln -s /usr/local/gcc9/bin/gcc /usr/bin/gcc
ln -s /usr/local/gcc9/bin/g++ /usr/bin/g++

7.更新libstdc++
通过 ls -lrt /usr/lib64/libstdc++.so.6 可以看到不是链接我们最新的libstdc++.so.6.0.26
删除
rm -f /usr/lib64/libstdc++.so.6
更新
ln -s /usr/local/lib64/libstdc++.