# ansible

```
[for mac]
brew install ansible

$ ansible --version
ansible 2.4.2.0
```

- inventory (target host): (ref)[https://qiita.com/t_nakayama0714/items/fe55ee56d6446f67113c]
```
ansible (hostname) -i hosts_all.yml -m ping
```

- playbook
```
ansible-playbook -i hosts_all.yml test.yml
ansible-playbook -i hosts_all.yml site.yml
```