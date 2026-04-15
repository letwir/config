#!/bin/bash
# 私がいつも使用しているアプリをインストールする。
# 個人用なので、必要なものだけを入れてくださいな。
# -----
# 最低条件のインストール
sudo apt-get update
sudo apt-get -y install curl wget tar jq
# ----- 変数 -----
LOKI_URL="$1"	# LokiのURLを引数で指定してください。例: http://localhost
SMARTCTLPATH=/usr/sbin/smartctl	# smartctlのバイナリ位置。デフォルトはここ。
EXPORTERDIR=/usr/local/bin	# Exporterのインストール場所
NODEEX=`curl -s https://api.github.com/repos/prometheus/node_exporter/releases/latest | jq -r .tag_name| cut -c 2-`	# Node_exporterのバージョン。 よく更新されているため適宜変更
LMEX="0.1.1"	# lmsensors_exporterのバージョン。 更新する気はない。
SMARTEX=`curl -s https://api.github.com/repos/prometheus-community/smartctl_exporter/releases/latest | jq -r .tag_name | cut -c 2-`	# smartctl_exporterのバージョン。 半年程度で更新されている様子
# ----- 変数ここまで -----

# fedora:
#ARCH=$(rpm --print-architecture 2>/dev/null || uname -m)
ARCH=$(dpkg --print-architecture 2>/dev/null || uname -m)
case "$ARCH" in
  "amd64"|"x86_64")
    BINARY="amd64" ;;  # Intel/AMDの64bit
  " i386"|"i686")
	BINARY="i386" ;;    # Intel/AMDの32bit
  "arm64"|"aarch64")
    BINARY="arm64" ;;  # Pi 3/4の64bit OS, armbian SBCsなど
  "armhf")
    BINARY="armv7" ;; # Pi 3/4の32bit OSなど
  "armv7l")
    BINARY="armv7" ;; # Pi 3/4の32bit OSなど
  "armv6l")
    BINARY="armv6" ;;  # Pi Zero/1など
  "armel")
    BINARY="armv5" ;;  # 古いarmhf系OSなど
  *)
    echo "未知のアーキテクチャ: $ARCH ですわ。手動で選んでくださいな。" && exit 1 ;;
esac
ARCH=$BINARY
# -----
# 未知だった場合は下の＃消して、バイナリのアーキテクチャを入れてください。対応してれば入ります。
#ARCH=
# -----

# 取得用のアプリインストール
# Alloy,他必要アプリのインストール
sudo mkdir -p /etc/apt/keyrings
sudo wget -O /etc/apt/keyrings/grafana.asc https://apt.grafana.com/gpg-full.key
sudo chmod 644 /etc/apt/keyrings/grafana.asc
echo "deb [signed-by=/etc/apt/keyrings/grafana.asc] https://apt.grafana.com stable main" | sudo tee /etc/apt/sources.list.d/grafana.list
sudo apt-get update
sudo apt-get -y install alloy ethtool mmc-utils smartmontools lm-sensors i2c-tools moreutils
sudo modprobe drivetemp
sudo modprobe eeprom_93xx46
sudo modprobe eeprom_93cx6
# もしモジュールがあれば
#echo -e "# for sensors\ndrivetemp\neeprom_93xx46\neeprom_93cx6" | sudo tee -a /etc/modules
yes | sudo sensors-detect

# -----
sudo systemctl stop node_exporter.service
# Node_Exporterのサービスインストール
wget https://github.com/prometheus/node_exporter/releases/download/v$NODEEX/node_exporter-$NODEEX.linux-$ARCH.tar.gz
sudo tar zxf node_exporter-$NODEEX.linux-$ARCH.tar.gz -O node_exporter-$NODEEX.linux-$ARCH/node_exporter > $EXPORTERDIR/node_exporter
sudo chmod +x $EXPORTERDIR/node_exporter
rm node_exporter-$NODEEX.linux-$ARCH.tar.gz
# サービス化
sudo tee /etc/systemd/system/node_exporter.service >/dev/null <<EOF
[Unit]
Description=Node Exporter Service
After=network-online.target

[Service]
Type=simple
PIDFile=/run/node_exporter.pid
ExecStart=$EXPORTERDIR/node_exporter
Environment="SCRIPT_ARGS=--collector.ethtool.device-include=.* --collector.filesystem.mount-points-exclude=^/(dev|proc|sys|var/lib/docker/.+|var/lib/kubelet/.+)($|/) --collector.cpu --collector.diskstats --collector.filesystem --collector.loadavg --collector.meminfo --collector.netdev --collector.netstat --collector.stat --collector.uname --collector.vmstat "

User=root
Group=root
SyslogIdentifier=node_exporter

Restart=on-failure
RemainAfterExit=no
RestartSec=100ms

StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
sudo systemctl enable node_exporter.service
sudo systemctl start node_exporter.service

# -----
sudo systemctl stop lmsensors_exporter.service
# lmsensors_exporterをインストール
sudo wget https://github.com/letwir/lmsensors_exporter/releases/download/$LMEX/lmsensors_exporter-$ARCH -O $EXPORTERDIR/lmsensors_exporter
sudo chmod +x $EXPORTERDIR/lmsensors_exporter
# サービス化
sudo tee /etc/systemd/system/lmsensors_exporter.service >/dev/null << EOF
[Unit]
Description=Lm_sensors Exporter Service
After=network-online.target

[Service]
Type=simple
PIDFile=/run/lmsensors_exporter.pid
ExecStart=$EXPORTERDIR/lmsensors_exporter

User=root
Group=root
SyslogIdentifier=lmsensors_exporter

Restart=on-failure
RemainAfterExit=no
RestartSec=100ms

StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
sudo systemctl enable lmsensors_exporter.service
sudo systemctl start lmsensors_exporter.service
# -----
sudo systemctl stop smartctl_exporter.service
# smartctl_exporterをインストール
wget https://github.com/prometheus-community/smartctl_exporter/releases/download/v$SMARTEX/smartctl_exporter-$SMARTEX.linux-$ARCH.tar.gz
sudo tar zxf smartctl_exporter-$SMARTEX.linux-$ARCH.tar.gz -O smartctl_exporter-$SMARTEX.linux-$ARCH/smartctl_exporter > $EXPORTERDIR/smartctl_exporter
sudo chmod +x $EXPORTERDIR/smartctl_exporter
rm smartctl_exporter-$SMARTEX.linux-$ARCH.tar.gz
# サービス化
sudo tee /etc/systemd/system/smartctl_exporter.service >/dev/null << EOF
[Unit]
Description=Smartctl_sensors Exporter Service
After=network-online.target

[Service]
Type=simple
PIDFile=/run/smartctl_exporter.pid
ExecStart=$EXPORTERDIR/smartctl_exporter
Environment="SCRIPT_ARGS=--smartctl.path=$SMARTCTLPATH --smartctl.device-include=/dev/disk/* --smartctl.scan-device-type=by-id "

User=root
Group=root
SyslogIdentifier=smartctl_exporter

Restart=on-failure
RemainAfterExit=no
RestartSec=100ms

StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
sudo systemctl enable smartctl_exporter.service
sudo systemctl start smartctl_exporter.service

# =====
# Alloyのインストール
# ------------------------------
# アンインストール時にはリポジトリファイルが残るので削除
#sudo rm -i /etc/apt/sources.list.d/grafana.list
# ------------------------------
# alloyユーザに必要グループを付与
sudo usermod -aG docker alloy
sudo usermod -aG adm alloy
sudo usermod -aG systemd-journal alloy

# ------------------------------
curl https://raw.githubusercontent.com/letwir/config/refs/heads/main/alloy/linux.alloy | sed "s@LOKI_URL@$LOKI_URL@g" | sudo tee /etc/alloy/config.alloy

# サービス有効化
sudo systemctl enable alloy
sudo systemctl start alloy


# 2. 個人向けアプリのインストール
sudo apt-get -y install \
curl nano wget git gh \
zip unzip 7zip xz-utils \
cifs-utils net-tools samba smbclient ethtool openssh-server nmap iperf3 \
python3-pip python3-venv \
build-essential cmake make gcc bison automake autoconf pkgconf libncurses-dev \
pipx translate-shell sudo dkms

cd
python3 -m venv .venv
curl https://sh.rustup.rs/ | bash -s -- --default-toolchain nightly -y --profile minimal
curl https://raw.githubusercontent.com/letwir/config/refs/heads/main/rustup/cargo-linux.toml | tee ~/.cargo/config
source $HOME/.cargo/env
curl https://raw.githubusercontent.com/cargo-bins/cargo-binstall/main/install-from-binstall-release.sh | bash -s -- -y
curl -sS https://starship.rs/install.sh | sh
cargo binstall git-delta bat sd vivid cargo-cache
cargo install fd-find lsd frs ripgrep du-dust hexyl choose lms starship --all-features
cargo cache -a


curl https://github.com/letwir.keys | tee .ssh/authorized_keys


# ロケールの設定
sed -i "s/# ja_JP.UTF-8 UTF-8/ja_JP.UTF-8 UTF-8/g" /etc/locale.gen
locale-gen
sudo localectl set-locale LANG=ja_JP.UTF-8
sudo timedatectl set-timezone Asia/Tokyo
sudo localectl set-keymap jp106
sudo localectl set-x11-keymap jp

# bashrcの設定
# bashrcのバックアップ
cp ~/.bashrc ~/.bashrc.bak
curl https://raw.githubusercontent.com/letwir/config/refs/heads/main/.bashrc | tee ~/.bashrc

# nanoのインストール
curl https://raw.githubusercontent.com/scopatz/nanorc/master/install.sh | sh
curl https://raw.githubusercontent.com/letwir/config/refs/heads/main/.nanorc | tee ~/.nanorc

# tmuxのインストール
sudo apt-get -y install tmux
curl https://raw.githubusercontent.com/letwir/config/refs/heads/main/.tmux.conf | tee ~/.tmux.conf
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm

# ネットワーク最適化
IF=`ip a | grep UP | head -n2 | tail -n1 | choose 1:1`
IF=${IF/%?/}
