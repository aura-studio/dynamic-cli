# 任务列表

* pusher模块中，根据Procedure的格式，生成NewTaskList，重新修改NewTaskList函数
* pusher模块中，修改PushForProcedure，根据NewTaskList生成的任务列表，执行任务
* build模块会按照builder.sh的语法生成四个文件，task_list需要能按照规则正确检索到这四个文件，并推送到对应的远程仓库
* runtime.Version() 不要出现在路径中