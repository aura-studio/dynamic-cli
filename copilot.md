1.将builder cleaner pusher 三个模块的模块明改为build clean push，以保持和cmd目录下命令一致。

2. 三个模块分别调用config.CreateProcedure函数来创建procedure对象，然后使用build.BuildForProcedure,clean.CleanForProcedure, push.PushForProcedure函数来执行相应的操作。