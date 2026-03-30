# `spark` 与 `lark-cli` 对比报告

日期：2026-03-29

## 一、结论先行

如果把 `spark` 和 `lark-cli` 放在同一层比较，容易得出偏差结论。  
更准确的说法应该是：

- `spark` 当前更像一个面向具体机器人产品形态的业务系统，里面包含了飞书接入、Agent 运行时、工具体系和业务逻辑。
- `lark-cli` 更像一个飞书能力层，把开放平台能力整理成稳定命令、身份模型和技能说明，适合被人、脚本和 Agent 复用。

从长期建设的角度看，`lark-cli` 这套范式更适合做“共用底座”；`spark` 更适合继续承接业务场景和产品逻辑。

这不是说现在要直接替换 `spark`，而是说新增飞书能力、私有化适配和 Agent 接入，优先往 `lark-cli` 收敛，会更稳。

如果从老板更关心的几个维度看，这套范式的核心价值不在“命令行工具”本身，而在于它更接近一种 **AI-native 的能力组织方式**：

- 不需要我们预先把用户场景枚举得很细。
- 用户可以直接提出自己的目标。
- Agent 根据目标自行拆解动作、调用原子能力、在需要时提示用户补授权，再继续执行。

这和传统“先设计固定场景，再给机器人配固定流程”的思路不同。  
后者更适合稳定业务流程，前者更适合需求变化快、问题类型分散、用户预期不断抬高的场景。

## 二、当前对比对象

### 1. `spark`

主要参考：

- [AGENT_CONTEXT.md](/Users/wangqizhao/Developer/iflytek/spark/AGENT_CONTEXT.md)
- [CLAUDE.md](/Users/wangqizhao/Developer/iflytek/spark/CLAUDE.md)
- [send_msg.py](/Users/wangqizhao/Developer/iflytek/spark/SparkChannel/app/services/feishu_services/send_msg.py)
- [documents_services.py](/Users/wangqizhao/Developer/iflytek/spark/SparkChannel/app/services/feishu_services/documents_services.py)
- [get_history_messages.py](/Users/wangqizhao/Developer/iflytek/spark/SparkChannel/app/services/feishu_services/get_history_messages.py)

当前形态：

- `SparkChannel` 负责渠道接入、消息预处理、上下文整理
- `spark-pai` 负责 Agent 推理、工具调用、任务调度和结果发送

也就是说，`spark` 是一套运行中的机器人系统，不只是工具库。

### 2. `lark-cli`

主要参考：

- [README.md](/Users/wangqizhao/Developer/iflytek/cli/README.md)
- [PRIVATE_DEPLOY_REPORT.md](/Users/wangqizhao/Developer/iflytek/cli/PRIVATE_DEPLOY_REPORT.md)
- [PRIVATE_DEPLOY_TEST_REPORT.md](/Users/wangqizhao/Developer/iflytek/cli/PRIVATE_DEPLOY_TEST_REPORT.md)
- [演示测试总览](/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-29/README.md)

当前形态：

- 飞书能力按业务域组织
- 命令分三层：shortcut、API command、raw API
- 身份模型和 scope 校验是统一的
- 对 Agent 暴露的是稳定命令和技能，不是项目内部 service

## 三、关键差异

### 1. 能力组织方式不同

#### `spark`

在 `spark` 里，飞书能力更多是按场景封装在 service 中：

- 发消息是一个 service
- 读历史消息是一个 service
- 读文档是一个 service
- 新增机器人和工作流又是另一组 service

这类方式的优点是贴合业务，缺点是：

- 复用边界主要在当前项目内部
- 新机器人或新场景接入时，往往还要继续补 service
- Agent 使用这些能力时，要理解较多项目内部概念

#### `lark-cli`

在 `lark-cli` 里，飞书能力被整理成统一动作：

- `im +messages-send`
- `docs +create`
- `drive +upload`
- `calendar +create`
- `task +create`
- `base +record-upsert`

这种方式的好处是：

- 同一动作可以被人、脚本、Agent 共同使用
- 输入输出更稳定
- 权限和身份更容易统一治理
- 更适合作为平台能力层，而不是业务代码的一部分

这里最关键的一点是：  
当飞书能力被收敛成稳定原子动作后，Agent 的工作重点会从“猜系统里有没有某个预设场景”变成“围绕用户目标组合动作”。

这会直接改变用户体验：

- 用户不需要先学会产品预设的入口。
- 用户不需要先知道某个需求属于“会议场景”“知识库场景”还是“任务场景”。
- 用户只要表达目标，Agent 就可以自己判断需要调哪些能力。

这就是为什么说它更 AI-native。  
AI-native 不等于界面里加了大模型，而是能力层本身就是围绕“目标驱动 + 动态组合”设计的。

### 2. AI 接入成本不同

`spark` 现在已经在做 Agent 架构升级，相关材料在：

- [findings.md](/Users/wangqizhao/Developer/iflytek/spark/docs/agent-migration/findings.md)
- [progress.md](/Users/wangqizhao/Developer/iflytek/spark/docs/agent-migration/progress.md)
- [langgraph-migration-plan.md](/Users/wangqizhao/Developer/iflytek/spark/docs/agent-migration/langgraph-migration-plan.md)

从这些文档可以看出，当前 `spark` 仍在处理几类问题：

- tool calling 的稳定性
- 多步任务衔接
- 缓存和冷启动
- function calling 主路径收敛

这说明 `spark` 的主矛盾之一，是“如何让一个业务机器人系统更适合 Agent”。  
而 `lark-cli` 的出发点相反：先把飞书能力做成 Agent 易于调用的形态，再让业务系统去消费。

这两条路都能走通，但从平台建设角度，后者的复用价值更高。

换成更直白的说法：

- `spark` 当前更像“先有业务机器人，再逐步补齐 Agent 能力”。
- `lark-cli` 更像“先有 Agent 可以稳定调用的能力层，再让业务机器人去组合它”。

从迭代速度和灵活度看，后者通常更占优，原因有三点：

1. **不需要先定义完整场景**
   - 传统做法往往要先想清楚“用户会问什么、流程怎么走、卡在哪一步”。
   - `lark-cli` 这种原子能力层更适合让用户直接提要求，再由 Agent 临场组织动作。

2. **新需求不必先做成新场景**
   - 只要底层原子能力已经具备，很多新需求只是组合方式变化，不一定需要新开一个项目迭代。

3. **授权链路天然可中断可继续**
   - Agent 可以先执行，遇到权限不足时提示用户补授权，然后继续往下做。
   - 这比“场景全写死、遇到权限问题整条流程报错退出”更符合真实使用过程。

### 3. 私有化适配状态不同

你本地的 `cli` 已经做了私有化适配，且不是停留在概念阶段。

当前已完成的工作包括：

- 端点覆盖
- Authorization Code Flow 适配
- scope 兼容处理
- 一轮真实环境测试

关键文件：

- [login.go](/Users/wangqizhao/Developer/iflytek/cli/cmd/auth/login.go)
- [authcode_flow.go](/Users/wangqizhao/Developer/iflytek/cli/internal/auth/authcode_flow.go)
- [scope.go](/Users/wangqizhao/Developer/iflytek/cli/internal/auth/scope.go)
- [config.go](/Users/wangqizhao/Developer/iflytek/cli/internal/core/config.go)
- [types.go](/Users/wangqizhao/Developer/iflytek/cli/internal/core/types.go)
- [remote.go](/Users/wangqizhao/Developer/iflytek/cli/internal/registry/remote.go)

这说明 `cli` 已经处在“可以继续投入收敛”的阶段，不是从零起步。

### 4. 治理和留痕能力不同

`spark` 目前对飞书能力的治理更偏项目内治理。  
`lark-cli` 则天然带有：

- `auth status`
- `doctor`
- `schema`
- 显式 `bot/user`
- 结构化错误输出

这对生产排障和演示都很重要。

对用户满意度也有直接影响。  
用户通常并不关心底层是哪个 service 或哪个 API，他们更在意三件事：

1. 我说一个目标，系统能不能接住。
2. 中间缺权限时，系统能不能说清楚下一步该干什么。
3. 授权之后，系统能不能继续，而不是让我重新来一遍。

在这三点上，统一能力层通常比场景化硬编码流程更有优势。

## 四、你本地 `cli` 适配项目的进度判断

当前可以把状态概括为三句话：

### 1. 已经证明方向可行

私有化端点、私有化 OAuth、scope 兼容和真实命令执行都已经跑通。

### 2. 还没有收敛成标准版本

现在仍是本地工作区改动，还不是内部统一发行版，也没有形成标准接入方式。

### 3. 最关键的现实风险还在环境侧

不是所有私有化飞书环境都能达到标准飞书的演示效果。  
当前最明显的两个问题是：

- **MCP 支持缺失或不完整**
  - 这会直接影响文档类 shortcut 的可用程度。
  - `docs +create` 这类能力在标准飞书里表现更完整，但在私有化环境可能受阻。

- **User 身份能力受限**
  - 一部分接口只对 `bot` 开放或对 `user` 支持不完整。
  - 即便命令层设计得完整，实际 demo 也可能做不到标准飞书那种完整链路。

这一点必须和老板提前说清楚，否则容易把环境限制误判成方案问题。

## 五、为什么说 `lark-cli` 更适合做底座

不是因为它“更新”或者“看起来更 AI”，而是因为它在工程上更接近底座的几个特征：

1. 能力域边界清楚
2. 身份和权限模型统一
3. 输入输出协议稳定
4. 人、脚本、Agent 共享同一套接口
5. 私有化适配只需要集中处理一层

这意味着以后如果再做第二个、第三个机器人，很多飞书能力不需要重新封装。

如果进一步放到老板更关心的指标上，这套方式更值得关注的是下面几个方面：

### 1. 迭代速度

在场景式方案里，新需求经常意味着：

- 新写一个场景
- 新接一个 service
- 新补 prompt
- 新测一遍链路

在原子能力层方案里，新需求更可能只是：

- 复用现有能力
- 调整能力组合
- 只补缺的那个原子动作

这两种方式的研发节奏完全不一样。  
前者更像项目制，后者更像能力积累制。

### 2. 灵活度

场景式方案更擅长处理高频、明确、流程固定的问题。  
AI-native 的原子能力层更擅长处理“需求表达不固定、组合关系多变、用户自己也说不清要走哪个入口”的问题。

这类需求在实际使用里往往越来越多。

### 3. 用户满意度

用户满意度不只是回答质量，还包括：

- 需求是否能被接住
- 失败是否说得明白
- 是否能少走弯路
- 是否能在授权后接着做

从这个角度看，`lark-cli` 更像是在给 Agent 一套可以持续完成任务的操作面，而不是一组离散功能按钮。

### 4. 产品增长空间

如果后续机器人能力继续扩展，场景式方案的压力会越来越大，因为每多一个场景，就多一套维护成本。  
原子能力层的增长方式更平缓：底层能力越全，Agent 能解决的问题越多，不一定需要同步增加大量新场景。

## 六、建议的建设方向

不建议把结论表达成“应该尽快替换 spark”。  
更稳妥的表达是：

1. `spark` 继续承接现有业务逻辑和产品形态。
2. 新增飞书能力优先往 `lark-cli` 收敛。
3. 私有化适配先在 `cli` 这一层稳定下来。
4. 后续再逐步让 `spark` 通过 `cli` 消费飞书能力，而不是继续内部重复封装。

这种说法更符合当前代码和环境现状，也更容易争取资源。

如果要对外表达得更聚焦一些，可以把重点放在下面这句话上：

**这不是在换一套实现，而是在把飞书能力从“预设场景驱动”调整为“用户目标驱动”。**

用户提需求，Agent 负责拆解和执行；  
权限不够时，Agent 提示授权；  
授权完成后，流程继续推进。  

从产品体验和迭代效率看，这比先替用户设想一大堆场景，再把需求塞进固定入口，更符合 AI 产品下一阶段的使用方式。

## 七、附：这份结论的边界

这份报告只回答“哪种范式更适合做长期能力层”，不回答“哪个项目今天更完整”。  
如果只看今天的业务闭环，`spark` 显然更完整。  
如果看后续复用、私有化适配、Agent 接入成本和跨项目共享，`lark-cli` 更值得继续投入。
