#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
批量清除 Markdown 文件中类似 “2。 ”、“38.” 的序号残留
并统一替换为一个全角句号 “。”。

兼顾 Windows：
  • 处理 UTF-8 以及带 BOM 的 UTF-8（常见于 Windows 记事本）；
  • 结束时在控制台显示彩色统计并发出 Beep 声作为提醒。
"""

import re
import sys
import argparse
import shutil
from pathlib import Path
import winsound
try:
    # 仅用于彩色输出，不安装也不影响功能
    from colorama import init as colorama_init, Fore, Style
    colorama_init()
    GREEN, RED, CYAN, RESET = Fore.GREEN, Fore.RED, Fore.CYAN, Style.RESET_ALL
except ImportError:
    GREEN = RED = CYAN = RESET = ''

# ------------------------------------------------------------------
# 1. 正则：匹配 “数字 + （全角句号|半角句号）”，前后允许有空白
# ------------------------------------------------------------------
pattern = re.compile(r'\b\d+\s*[\.。]')
replacement = '。'

def read_text(path: Path) -> str:
    """
    自动识别 UTF-8 / UTF-8-BOM
    """
    data = path.read_bytes()
    # 检测 UTF-8 BOM
    if data.startswith(b'\xef\xbb\xbf'):
        return data[3:].decode('utf-8')
    return data.decode('utf-8')

def write_text(path: Path, text: str, keep_bom: bool) -> None:
    encoded = text.encode('utf-8')
    if keep_bom:
        encoded = b'\xef\xbb\xbf' + encoded
    path.write_bytes(encoded)

def clean_file(md_path: Path, backup: bool) -> bool:
    """
    清洗单个文件
    返回值：True = 文件被修改；False = 无需修改
    """
    raw_bytes = md_path.read_bytes()
    has_bom   = raw_bytes.startswith(b'\xef\xbb\xbf')
    original  = raw_bytes[3:].decode('utf-8') if has_bom else raw_bytes.decode('utf-8')

    cleaned = pattern.sub(replacement, original)

    if cleaned == original:
        return False

    if backup:
        shutil.copy2(md_path, md_path.with_suffix(md_path.suffix + '.bak'))

    write_text(md_path, cleaned, keep_bom=has_bom)
    return True

def main(folder: Path, backup: bool, recursive: bool):
    total, changed, unchanged = 0, 0, 0

    md_iter = folder.rglob('*.md') if recursive else folder.glob('*.md')

    for md in md_iter:
        total += 1
        try:
            modified = clean_file(md, backup)
            if modified:
                changed += 1
                print(f'{GREEN}✔ 已修改:{RESET} {md}')
            else:
                unchanged += 1
                print(f'{CYAN}· 无需修改:{RESET} {md}')
        except Exception as e:
            print(f'{RED}× 发生错误:{RESET} {md}\n    {e}')

    # 统计
    print('\n========== 处理完毕 ==========')
    print(f'总文件数: {total}')
    print(f'已修改  : {changed}')
    print(f'未改变  : {unchanged}')
    # 提示音：frequency=1000Hz, duration=300ms
    winsound.Beep(1000, 300)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description='批量删除 Markdown 中 “数字+句号” 的引用残留并统一替换为 “。”'
    )
    parser.add_argument('folder', help='要处理的文件夹')
    parser.add_argument('-b', '--backup', action='store_true',
                        help='为每个被修改的文件生成 .bak 备份')
    parser.add_argument('-n', '--no-recursive', action='store_true',
                        help='仅处理当前层级，不递归子目录')
    args = parser.parse_args()

    target = Path(args.folder).expanduser().resolve()
    if not target.is_dir():
        sys.exit('❌ 指定路径不是有效文件夹！')

    main(target, backup=args.backup, recursive=not args.no_recursive)
