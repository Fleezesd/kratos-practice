- kratos log
```go
给出Kratos集成logrus和zap两个日志库的示例代码:

集成logrus:

import (
  "github.com/sirupsen/logrus"
  "github.com/go-kratos/kratos/contrib/log/logrus"

  "github.com/go-kratos/kratos/v2/log"
)

// 创建logrus日志对象
logger := logrus.New() 

// 创建logrus适配器
logrusLogger := logrus.NewLogger(logger)

// 生成Kratos日志
kratosLog := log.NewHelper(logrusLogger)

kratosLog.Info("info log")
集成zap:

Copy
import (
  "go.uber.org/zap"
  "github.com/go-kratos/kratos/contrib/log/zap"
  
  "github.com/go-kratos/kratos/v2/log"  
)

// 创建zap日志
zapLogger, _ := zap.NewProduction()

// 创建zap适配器
zapLogger := zap.NewLogger(zapLogger)

// 生成Kratos日志 
kratosLog := log.NewHelper(zapLogger)

kratosLog.Info("info log")
通过创建适配器作为桥接,可以方便的在Kratos中集成不同的日志库,而业务代码只需要使用Kratos日志接口即可,不需要变更。
```