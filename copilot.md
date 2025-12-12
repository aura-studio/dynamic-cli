# 任务列表

* 在RenderData的构建中，为OS、Arch和Compiler字段添加支持。
* 在builder.sh第8行上面，添加对OS，Arch和Compiler的检查，看和预期是否一致，如果不一致则报错退出。检查基于go env中的GOOS，GOARCH和GOVERSION环境变量。