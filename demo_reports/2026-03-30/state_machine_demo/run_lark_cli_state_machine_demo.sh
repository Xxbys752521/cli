#!/usr/bin/env bash
set -euo pipefail

ROOT="/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-30/state_machine_demo"
TMP_DIR="$ROOT/tmp"
REPORT="$ROOT/demo_run_report.md"
mkdir -p "$ROOT" "$TMP_DIR"
: > "$REPORT"

log() {
  printf '%s\n' "$*" >> "$REPORT"
}

section() {
  log ""
  log "## $1"
  log ""
}

run_cmd() {
  local label="$1"
  local cmd="$2"
  local outfile="$TMP_DIR/${label}.json"
  local out

  out="$(bash -lc "$cmd")"
  printf '%s\n' "$out" > "$outfile"

  log "### $label"
  log ""
  log '```bash'
  log "$cmd"
  log '```'
  log ""
  log '```json'
  log "$out"
  log '```'
  log ""

  printf '%s' "$out"
}

json_field() {
  local json="$1"
  local expr="$2"
  printf '%s' "$json" | jq -r "$expr"
}

TS="$(date +%Y%m%d_%H%M%S)"
GROUP_NAME="状态机Demo群-$TS"
BASE_NAME="状态机Demo答疑表-$TS"
DOC_TITLE="状态机Demo知识库-$TS"

SELF_JSON="$(lark-cli contact +get-user --as user)"
SELF_OPEN_ID="$(printf '%s' "$SELF_JSON" | jq -r '.data.user.open_id')"
SELF_NAME="$(printf '%s' "$SELF_JSON" | jq -r '.data.user.name')"

log "# lark-cli 状态机 Demo 运行报告"
log ""
log "- 运行时间：$(date '+%Y-%m-%d %H:%M:%S %z')"
log "- 执行用户：$SELF_NAME"
log "- 目标：从建群、构造聊天记录、建答疑表和知识库，到按状态机跑完自动答疑和转人工两条链路"

section "0. 准备上下文"
log "- 当前用户 open_id：\`$SELF_OPEN_ID\`"
log "- Demo 群名：\`$GROUP_NAME\`"
log "- Demo Base：\`$BASE_NAME\`"
log "- Demo 知识库文档：\`$DOC_TITLE\`"

section "1. 创建知识库文档"
cat > "$TMP_DIR/kb_init.md" <<'EOF'
# 常见答疑知识库

## 案例 1：私有化环境知识库打不开

- 现象：登录飞书后点击知识库页面打不开，或提示无权限
- 适用范围：私有化部署环境、新员工或新租户初次接入
- 处理建议：
  1. 确认当前账号已完成飞书登录
  2. 确认知识库对当前用户或部门已授权
  3. 如果是私有化环境，检查是否使用了正确的环境域名和对应租户
  4. 如果页面空白，优先检查是否存在权限未同步或网关配置异常
- 输出口径：先引导补充环境信息，再给出权限检查路径
EOF
DOC_JSON="$(run_cmd doc_create "lark-cli docs +create --as user --title '$DOC_TITLE' --markdown \"$(cat "$TMP_DIR/kb_init.md")\"")"
DOC_ID="$(json_field "$DOC_JSON" '.data.doc_id')"
DOC_URL="$(json_field "$DOC_JSON" '.data.doc_url')"

section "2. 创建答疑 Base 和表"
BASE_JSON="$(run_cmd base_create "lark-cli base +base-create --as user --name '$BASE_NAME' --time-zone Asia/Shanghai")"
BASE_TOKEN="$(json_field "$BASE_JSON" '.data.base.base_token')"
BASE_URL="$(json_field "$BASE_JSON" '.data.base.url')"

FIELDS_JSON='[
  {"name":"问题ID","type":"text"},
  {"name":"当前状态","type":"select","multiple":false,"options":[{"name":"待补充"},{"name":"待检索"},{"name":"已自动回答"},{"name":"已转人工"},{"name":"人工已处理"},{"name":"已回写知识库"},{"name":"已结束"}]},
  {"name":"提问人","type":"text"},
  {"name":"问题内容","type":"text"}
]'
TABLE_JSON="$(run_cmd table_create "lark-cli base +table-create --as user --base-token $BASE_TOKEN --name '答疑表' --fields '$FIELDS_JSON'")"
TABLE_ID="$(json_field "$TABLE_JSON" '.data.table.table_id // .data.table.id // .data.table_id')"
run_cmd table_add_field_case "sleep 2; lark-cli base +field-create --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --json '{\"type\":\"select\",\"name\":\"是否命中案例\",\"multiple\":false,\"options\":[{\"name\":\"待定\"},{\"name\":\"是\"},{\"name\":\"否\"}]}'"
run_cmd table_add_field_result "sleep 4; lark-cli base +field-create --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --json '{\"type\":\"text\",\"name\":\"处理结论\"}'"
run_cmd table_add_field_kb "sleep 6; lark-cli base +field-create --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --json '{\"type\":\"text\",\"name\":\"知识库链接\"}'"

section "3. 创建 Demo 群"
GROUP_JSON="$(run_cmd group_create "lark-cli im +chat-create --as bot --type private --name '$GROUP_NAME' --users '$SELF_OPEN_ID' --set-bot-manager")"
CHAT_ID="$(json_field "$GROUP_JSON" '.data.chat_id // .data.chat.chat_id // .data.chat.id')"

section "4. Case A：自动答疑闭环"
run_cmd case_a_msg_01 "lark-cli im +messages-send --as bot --chat-id $CHAT_ID --text '【模拟提问人A】登录飞书后打不开知识库，帮我看看'"
CASE_A_CREATE_JSON="$(run_cmd case_a_record_create "lark-cli base +record-upsert --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --json '{\"问题ID\":\"CASE-A\",\"当前状态\":\"待补充\",\"提问人\":\"模拟提问人A\",\"问题内容\":\"登录飞书后打不开知识库\",\"是否命中案例\":\"待定\",\"处理结论\":\"\",\"知识库链接\":\"$DOC_URL\"}'")"
CASE_A_RECORD_ID="$(json_field "$CASE_A_CREATE_JSON" '.data.record.record_id_list[0]')"

run_cmd case_a_msg_02 "lark-cli im +messages-send --as bot --chat-id $CHAT_ID --text '【机器人】需要补充两个信息：1）是否是私有化环境；2）具体报错是无权限、空白页还是打不开链接。'"
run_cmd case_a_msg_03 "lark-cli im +messages-send --as bot --chat-id $CHAT_ID --text '【模拟提问人A补充】是私有化环境，点击后提示无权限。'"
run_cmd case_a_record_update_01 "lark-cli base +record-upsert --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --record-id $CASE_A_RECORD_ID --json '{\"当前状态\":\"待检索\",\"问题内容\":\"私有化环境中点击知识库提示无权限\",\"是否命中案例\":\"待定\"}'"
run_cmd case_a_kb_fetch "lark-cli docs +fetch --as user --doc $DOC_ID"
run_cmd case_a_msg_04 "lark-cli im +messages-send --as bot --chat-id $CHAT_ID --text '【机器人】已命中相似案例：先确认当前账号是否已被知识库授权；若是私有化环境，再检查是否用了正确的环境域名和租户。若页面仍异常，再排查权限同步和网关配置。'"
run_cmd case_a_record_update_02 "lark-cli base +record-upsert --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --record-id $CASE_A_RECORD_ID --json '{\"当前状态\":\"已自动回答\",\"是否命中案例\":\"是\",\"处理结论\":\"命中私有化环境知识库权限案例，已给出授权与环境校验建议。\"}'"
run_cmd case_a_msg_05 "lark-cli im +messages-send --as bot --chat-id $CHAT_ID --text '【模拟提问人A反馈】已解决，满意。'"
run_cmd case_a_record_update_03 "lark-cli base +record-upsert --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --record-id $CASE_A_RECORD_ID --json '{\"当前状态\":\"已结束\"}'"
run_cmd case_a_kb_update "lark-cli docs +update --as user --doc $DOC_ID --mode append --markdown '## 回写记录：CASE-A\n- 结果：自动答疑命中历史案例并解决\n- 备注：用户补充私有化环境和无权限信息后可直接命中知识库。'"

section "5. Case B：未命中案例，转人工后回写知识库"
run_cmd case_b_msg_01 "lark-cli im +messages-send --as bot --chat-id $CHAT_ID --text '【模拟提问人B】新员工入职后，外部知识库链接打开是空白页，重新登录也不行。'"
CASE_B_CREATE_JSON="$(run_cmd case_b_record_create "lark-cli base +record-upsert --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --json '{\"问题ID\":\"CASE-B\",\"当前状态\":\"待检索\",\"提问人\":\"模拟提问人B\",\"问题内容\":\"新员工入职后，外部知识库链接打开空白页，重新登录也不行\",\"是否命中案例\":\"待定\",\"处理结论\":\"\",\"知识库链接\":\"$DOC_URL\"}'")"
CASE_B_RECORD_ID="$(json_field "$CASE_B_CREATE_JSON" '.data.record.record_id_list[0]')"
run_cmd case_b_kb_fetch "lark-cli docs +fetch --as user --doc $DOC_ID"
run_cmd case_b_record_update_01 "lark-cli base +record-upsert --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --record-id $CASE_B_RECORD_ID --json '{\"当前状态\":\"已转人工\",\"是否命中案例\":\"否\",\"处理结论\":\"未命中现有案例，转人工排查。\"}'"
run_cmd case_b_msg_02 "lark-cli im +messages-send --as bot --chat-id $CHAT_ID --text '【机器人】当前知识库未命中相似案例，已转人工排查，请等待。'"
run_cmd case_b_msg_03 "lark-cli im +messages-send --as bot --chat-id $CHAT_ID --text '【模拟人工】排查结果：该员工未加入外部知识库可见范围，且浏览器缓存保留了旧租户跳转地址。已补权限并清理缓存，现已恢复。'"
run_cmd case_b_record_update_02 "lark-cli base +record-upsert --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --record-id $CASE_B_RECORD_ID --json '{\"当前状态\":\"人工已处理\",\"处理结论\":\"人工确认是知识库可见范围缺失叠加旧租户缓存导致，补权限并清缓存后恢复。\"}'"
run_cmd case_b_kb_update "lark-cli docs +update --as user --doc $DOC_ID --mode append --markdown '## 案例 2：新员工打开外部知识库为空白页\n- 现象：外部知识库链接打开空白，重新登录无效\n- 原因：知识库可见范围未包含该员工，同时浏览器保留旧租户缓存\n- 处理建议：补充知识库可见范围，清理浏览器缓存后重新登录\n- 来源：CASE-B 人工处理回写'"
run_cmd case_b_record_update_03 "lark-cli base +record-upsert --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --record-id $CASE_B_RECORD_ID --json '{\"当前状态\":\"已回写知识库\",\"知识库链接\":\"$DOC_URL\"}'"
run_cmd case_b_msg_04 "lark-cli im +messages-send --as bot --chat-id $CHAT_ID --text '【机器人】人工处理结果已同步，并已将该问题回写到知识库。'"
run_cmd case_b_record_update_04 "lark-cli base +record-upsert --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --record-id $CASE_B_RECORD_ID --json '{\"当前状态\":\"已结束\"}'"

section "6. 拉取最终证据"
run_cmd final_chat_messages "lark-cli im +chat-messages-list --as user --chat-id $CHAT_ID --sort asc --page-size 50"
run_cmd final_record_a "lark-cli base +record-get --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --record-id $CASE_A_RECORD_ID"
run_cmd final_record_b "lark-cli base +record-get --as user --base-token $BASE_TOKEN --table-id $TABLE_ID --record-id $CASE_B_RECORD_ID"
run_cmd final_kb_doc "lark-cli docs +fetch --as user --doc $DOC_ID"

section "7. 结论"
log "- 自动答疑链路已覆盖：提问 -> 补充信息 -> 检索知识库 -> 命中案例 -> 自动回答 -> 用户满意 -> 结束"
log "- 转人工链路已覆盖：提问 -> 检索未命中 -> 转人工 -> 人工处理 -> 回写知识库 -> 结束"
log "- Demo 群链接（需在飞书内打开）：chat_id \`$CHAT_ID\`"
log "- Demo Base：[$BASE_NAME]($BASE_URL)"
log "- Demo 知识库：[$DOC_TITLE]($DOC_URL)"

printf 'Demo finished. Report written to %s\n' "$REPORT"
