# porpose

- this is diretory include vagrant settings as `Bootstrapping` ( not include settings of midleware )

# vagrant

- vagrant box list:https://atlas.hashicorp.com/boxes/search
- example: 
```
mkdir centos7 && cd $0
vagrant ini centos/7

[check&up]
vagrant box list
vagrant up
vagrant ssh
```

# ssh

- (this procedure should be automated)
```
sshkey-gen (name: vagrant_centos7)
ln -s /**/conf/ssh_config/ ~/.ssh/config.d/vagrant_conf
echo include config.d/vagrant_conf/*.conf >> ~/.ssh/config
ssh centos7
```

