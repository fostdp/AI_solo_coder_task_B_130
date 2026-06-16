import json
import re

with open('scripts/init_new_features.sql', 'r', encoding='utf-8') as f:
    content = f.read()

matches = re.findall(r"'(\{[^']*source_archaeology[^']*\})'", content)

print(f'找到 {len(matches)} 个包含考古字段的JSON:')
for i, match in enumerate(matches, 1):
    try:
        data = json.loads(match)
        print(f'  [{i}] ✓ JSON格式正确')
        print(f'       source_archaeology: {data.get("source_archaeology", "缺失")[:30]}...')
        print(f'       uncertainty_range: {data.get("uncertainty_range", "缺失")}')
        print(f'       experimental_method: {data.get("experimental_method", "缺失")[:30]}...')
    except json.JSONDecodeError as e:
        print(f'  [{i}] ✗ JSON格式错误: {e}')
        print(f'       内容: {match[:50]}...')
