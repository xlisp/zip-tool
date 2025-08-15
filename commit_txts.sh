#!/bin/bash
set -e  # 有错误就停止

cd "$(dirname "$0")"   # 切换到脚本所在目录

# 遍历 output 目录里的所有 txt 文件
for file in output/*.txt; do
    echo "Committing $file ..."
    git add "$file"
    git commit -m "Update $(basename "$file")"
    git push origin main
done

echo "✅ All txt files committed."

