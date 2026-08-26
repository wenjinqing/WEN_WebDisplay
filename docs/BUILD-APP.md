# 安卓 APP 打包指南

本项目通过 Capacitor 套壳生成安卓 APP,网页代码不变。

## 环境(本服务器已装好)

- JDK 21:`/usr/lib/jvm/java-21-openjdk`
- Android SDK:`/opt/android-sdk`(platform-36 / build-tools 36.0.0)
- 签名密钥:`/root/catcafe.keystore`,密码在 `/root/.catcafe-keystore-pass`
  （同一份密钥也加密存放在 GitHub Actions secrets 里,换机器打包务必用同一把钥匙,否则老用户无法覆盖安装升级）

## 打包步骤

服务器上直接跑脚本即可(包含下列全部步骤):

```bash
/root/build-catcafe-apk.sh
```

手动步骤:

```bash
cd /root/workspace/WEN_WebDisplay
npm run build
npx cap sync android
cp /root/catcafe.keystore android/app/catcafe.keystore   # .gitignore 已排除,不会入库
cd android
export JAVA_HOME=/usr/lib/jvm/java-21-openjdk
export ANDROID_HOME=/opt/android-sdk
export KEYSTORE_PASSWORD=$(cat /root/.catcafe-keystore-pass)
export KEY_PASSWORD=$KEYSTORE_PASSWORD
export KEY_ALIAS=catcafe
./gradlew assembleRelease --no-daemon
# 产物:android/app/build/outputs/apk/release/app-release.apk
cp app/build/outputs/apk/release/app-release.apk /var/www/downloads/app/catcafe.apk
```

发布后用户在官网导航栏点「下载应用」(安卓)即可下载安装。

## 服务器内存注意

本机只有 2G 内存,Gradle 容易 OOM。已采取:`gradle.properties` 限堆(-Xmx896m /
kotlin 640m / workers=1)+ 4G swap 文件 `/swapfile-app`(重启后需 `swapon /swapfile-app`)。
若还是偶发被杀,直接重跑脚本,增量编译会很快走完。

## 发版 checklist(APP 内「检查更新」依赖)

1. `android/app/build.gradle`:versionCode 整数 +1,versionName 顺手改
2. 打包并把 APK 覆盖到 `/var/www/downloads/app/catcafe.apk`
3. 更新 `/var/www/downloads/app/version.json` 里的 versionCode / versionName / notes
   (老版本 APP 的「我的 → 检查更新」就是读这个文件,versionCode 变大才会提示更新)

## 版本号

升级 APP 时改 `android/app/build.gradle` 里的 `versionCode`(整数,必须递增)和 `versionName`。

## Windows 本地打包

装好 JDK 21 + Android Studio(或独立 SDK)后步骤同上;Gradle 走腾讯云镜像
(`android/gradle/wrapper/gradle-wrapper.properties` 已配置,国内机器都快)。

## 云端打包(备用)

`deploy/build-apk.yml` 是现成的 GitHub Actions 工作流(push 自动出 APK 并上传到
Release)。启用方法:把它移到 `.github/workflows/` 并推送——注意推送 token 需要
勾选 `workflow` 权限。签名密钥的 4 个 secrets 已经配置好了。
