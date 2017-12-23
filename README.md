# env

- vagrant and Virtualbox
```
$ vagrant -v
Vagrant 2.0.1

$ VBoxManage -v
5.2.4r119785
```
- vagrant box list:https://atlas.hashicorp.com/boxes/search
- ref: `vagrant ini centos/7`

# ansible

```
brew install ansible

$ ansible --version
ansible 2.4.2.0
```
- ssh for vagrant: (ref)[https://qiita.com/deco/items/d33395e30875b2923fac]
```
vagrant ssh-config >> ~/.ssh/config
ssh default

> ssh-keygen -t rsa
>
> # check port
> # vagrant ssh-config 
> scp -P (port num) -i ~/.ssh/(your private) ~/.ssh/(your pub) (your_vagrant:your_dir)
>
> vim VagrantFile
> config.ssh.private_key_path = "~/.ssh/(your private)"
>
> ssh root@192.168.(your vagrant)
```

- inventory (target host): (ref)[https://qiita.com/t_nakayama0714/items/fe55ee56d6446f67113c]
```
mkdir inventory
cd inventory

vim inventory/hosts

> [targets]
> 192.168.100.20
>
> ansible all -i inventory/hosts -m ping


mkdir /etc/ansible

vim /etc/ansible/ansible.cnf
[ssh_connection]
ssh_args = -F /**/**/(yourdir)/ssh_conf

ansible (hostname) -i inventory/hosts -m ping
```

- playbook 
```
vim group_vars/targets.yml

message: "Hello Ansible !"
fruits:
  apples:
    amount: 10
  bananas:
    amount: 20
  oranges:
    amount: 30


vim test.yml

- hosts: targets
  user: root
  tasks:
    - name: output message.
      debug: msg="{{ message }}"
    - name: output fruits
      debug: msg="We want {{ item.value.amount }} {{ item.key }} !"
      with_dict: "{{ fruits }}"


ansible-playbook -i inventory/hosts test.yml
```