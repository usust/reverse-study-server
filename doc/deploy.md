




## 创建 service

创建服务文件 `/etc/systemd/system/reverse-study-server.service`

```ini
[Unit]
Description=Reverse Study Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/deploy-server/reverse
ExecStart=/opt/deploy-server/reverse-study-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

配置服务
```sh
sudo systemctl daemon-reload
sudo systemctl start reverse-study-server.service
sudo systemctl status reverse-study-server.service
sudo systemctl enable reverse-study-server.service
```

查看日志：
```sh
journalctl -u actions-runner -f
```