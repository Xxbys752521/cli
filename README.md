# xfchat_cli

`xfchat_cli` is the private-deployment `lark-cli` fork used in this repository.
It keeps the `xfchat_cli` binary and `~/.xfchat_cli` config namespace, adds
private endpoint and OAuth handling, and ships bundled agent skills for external
AI runtimes.

Current runtime guidance is aligned to the private deployment:

- enabled: IM, Docs, Drive, Base, Sheets, Calendar, Task, Contact, Wiki, VC
- bundled skills: 19 skills are packaged with the CLI
- disabled by default: `mail` remains as historical skill documentation only

See [README.zh.md](./README.zh.md) for installation, authentication, and the
private-deployment capability matrix.
