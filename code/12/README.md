# 通道与并发模式示例

对应第 12 章作业，包含：
- `PipelineDoubleThenAdd`：两阶段流水线（`x*2` 然后 `x+1`）。
- `FanOutSquare`：fan-out/fan-in worker 处理 jobs 并归并输出。
- `SendWithTimeout`：背压场景下，发送超时/取消的处理示例。

## 运行
```bash
cd code/12
go test ./...
```
