# BuildForProcedure中，从Procedure结构构造RenderData结构

* Name由Target的Namespace，Package，Version通过下划线连接而成
* Module由Source的Module直接赋值
* Version由Source的Version直接赋值
* Dir由三部分使用/组成，第一部分是House，由 Warehouse的Local直接赋值，第二部分是Enviroment，Enviroment由Toolchain的OS，Arch，Compiler, Varian通过下划线连接而成, 第三部分是Name,由Target的Namespace，Package，Version通过下划线连接而成
* Variant由Toolchain的Variant直接赋值
* 其他没有提及的字段删除