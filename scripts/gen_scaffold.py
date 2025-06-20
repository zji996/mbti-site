#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
为 MBTI 站点批量生成 /content/ 下的人格与配对 Markdown 模板。
支持多语言；默认 zh（与旧脚本保持兼容）。

用法示例：
    # 只生成（或补齐）中文
    python gen_scaffold.py

    # 生成/覆盖 zh + en
    python gen_scaffold.py --langs zh,en --force
"""

from pathlib import Path
import itertools
import argparse
import sys

TYPE_CODES = [
    "INTJ", "INTP", "ENTJ", "ENTP",
    "INFJ", "INFP", "ENFJ", "ENFP",
    "ISTJ", "ISFJ", "ESTJ", "ESFJ",
    "ISTP", "ISFP", "ESTP", "ESFP",
]

# gender map per language
GENDER_STR = {
    "zh": {"m": "男性", "f": "女性"},
    "en": {"m": "male",  "f": "female"},
}

# ------------ 新增：人格的固定元数据 ------------ #
TYPE_DATA = {
    "zh": {
        "INTJ": {"name": "建筑师", "tagline": "深谋远虑，构筑未来的思想建筑师。"},
        "INTP": {"name": "逻辑学家", "tagline": "知识的探索者，沉醉于思想的无限可能。"},
        "ENTJ": {"name": "指挥官", "tagline": "天生的领导者，用决心和远见引领变革。"},
        "ENTP": {"name": "辩论家", "tagline": "思想的碰撞者，享受智慧交锋的乐趣。"},
        "INFJ": {"name": "提倡者", "tagline": "安静的理想主义者，用信念点亮世界。"},
        "INFP": {"name": "调停者", "tagline": "诗意的梦想家，在平淡中寻找美好。"},
        "ENFJ": {"name": "主人公", "tagline": "富有感染力的领袖，激励他人实现梦想。"},
        "ENFP": {"name": "竞选者", "tagline": "热情洋溢的社交家，总能找到快乐的理由。"},
        "ISTJ": {"name": "物流师", "tagline": "一丝不苟的实干家，可靠是我的代名词。"},
        "ISFJ": {"name": "守卫者", "tagline": "温暖的守护者，时刻准备着保护所爱之人。"},
        "ESTJ": {"name": "总经理", "tagline": "出色的管理者，擅长管理事务或人员。"},
        "ESFJ": {"name": "执政官", "tagline": "乐于助人的社交中心，极具同情心。"},
        "ISTP": {"name": "鉴赏家", "tagline": "大胆务实的实验家，热衷于探索和创造。"},
        "ISFP": {"name": "探险家", "tagline": "灵活迷人的艺术家，随时准备探索新事物。"},
        "ESTP": {"name": "企业家", "tagline": "敏锐的行动派，活在当下的风险爱好者。"},
        "ESFP": {"name": "表演者", "tagline": "充满活力的娱乐家，生活就是我的舞台。"},
    },
    "en": {
        "INTJ": {"name": "The Architect", "tagline": "Imaginative and strategic thinkers, with a plan for everything."},
        "INTP": {"name": "The Logician", "tagline": "Innovative inventors with an unquenchable thirst for knowledge."},
        "ENTJ": {"name": "The Commander", "tagline": "Bold, imaginative and strong-willed leaders, always finding a way."},
        "ENTP": {"name": "The Debater", "tagline": "Smart and curious thinkers who cannot resist an intellectual challenge."},
        "INFJ": {"name": "The Advocate", "tagline": "Quiet and mystical, yet very inspiring and tireless idealists."},
        "INFP": {"name": "The Mediator", "tagline": "Poetic, kind and altruistic people, always eager to help a good cause."},
        "ENFJ": {"name": "The Protagonist", "tagline": "Charismatic and inspiring leaders, able to mesmerize their listeners."},
        "ENFP": {"name": "The Campaigner", "tagline": "Enthusiastic, creative and sociable free spirits, who can always find a reason to smile."},
        "ISTJ": {"name": "The Logistician", "tagline": "Practical and fact-minded individuals, whose reliability cannot be doubted."},
        "ISFJ": {"name": "The Defender", "tagline": "Very dedicated and warm protectors, always ready to defend their loved ones."},
        "ESTJ": {"name": "The Executive", "tagline": "Excellent administrators, unsurpassed at managing things or people."},
        "ESFJ": {"name": "The Consul", "tagline": "Extraordinarily caring, social and popular people, always eager to help."},
        "ISTP": {"name": "The Virtuoso", "tagline": "Bold and practical experimenters, masters of all kinds of tools."},
        "ISFP": {"name": "The Adventurer", "tagline": "Flexible and charming artists, always ready to explore and experience something new."},
        "ESTP": {"name": "The Entrepreneur", "tagline": "Smart, energetic and very perceptive people, who truly enjoy living on the edge."},
        "ESFP": {"name": "The Entertainer", "tagline": "Spontaneous, energetic and enthusiastic people – life is never boring around them."},
    }
}


# ------------ 模板：每种语言 1 份 ------------ #
TYPE_TMPL = {
    "zh": """\
---
lang: zh
code: {code}
gender: {gkey}
name: {name}
tagline: {tagline}
summary:        
---

[DR] 请概述 {code} {gcn} 人格特点, 深度画像, 常见盲点, 成长建议

""",
    "en": """\
---
lang: en
code: {code}
gender: {gkey}
name: {name}
tagline: {tagline}
summary:       # [GPT] ≤50-word overview
---

[DR] Summarize the characteristics of {code} {gcn}, deep profile, common blind spots and growth advice

"""
}

PAIR_TMPL = {
    "zh": """\
---
lang: zh
male_code: {mcode}
female_code: {fcode}
compatibility_score:       # [GPT] 请给出 1–10 的契合度评分
---

[DR] 请分析 {mcode} 男与 {fcode} 女的契合度, 互补优势, 潜在冲突, 相处建议等

""",
    "en": """\
---
lang: en
male_code: {mcode}
female_code: {fcode}
compatibility_score:       # [GPT] Give a 1–10 compatibility score
---

[DR] Analyze the compatibility, complementary strengths, possible conflicts and tips for {mcode} male and {fcode} female

"""
}

# ------------ 目录 ------------ #
ROOT = Path(__file__).resolve().parent.parent
CONTENT_DIR = ROOT / "content"      # 兼容旧路径

# ------------------------------------------------ #
def write(path: Path, text: str, force: bool):
    path.parent.mkdir(parents=True, exist_ok=True)
    if path.exists() and not force:
        return
    path.write_text(text, encoding="utf-8")
    print(("overwrite" if path.exists() else "create").ljust(10), path.relative_to(ROOT))

def gen_types(lang: str, force: bool):
    gmap = GENDER_STR[lang]
    base_dir = CONTENT_DIR / lang / "types"
    for code, gkey in itertools.product(TYPE_CODES, gmap):
        fname = f"{code.lower()}_{gkey}.md"
        path  = base_dir / fname
        
        # 从 TYPE_DATA 中获取 name 和 tagline
        type_info = TYPE_DATA[lang][code]

        # 格式化模板，并传入新增的数据
        txt   = TYPE_TMPL[lang].format(
            code=code,
            gkey=gkey,
            gcn=gmap[gkey],
            name=type_info["name"],
            tagline=type_info["tagline"]
        )
        write(path, txt, force)

def gen_pairs(lang: str, force: bool):
    base_dir = CONTENT_DIR / lang / "pairings"
    for male, female in itertools.product(TYPE_CODES, TYPE_CODES):
        fname = f"{male.lower()}_m__{female.lower()}_f.md"
        path  = base_dir / fname
        txt   = PAIR_TMPL[lang].format(mcode=male, fcode=female)
        write(path, txt, force)

# ------------------------------------------------ #
def parse_langs(val: str):
    langs = [x.strip() for x in val.split(",") if x.strip()]
    for l in langs:
        if l not in TYPE_TMPL:
            sys.exit(f"[ERR] unsupported lang: {l}")
    return langs

def main():
    ap = argparse.ArgumentParser(description="Generate content scaffold (multi-language)")
    ap.add_argument("--force", action="store_true", help="overwrite existing files")
    ap.add_argument(
        "--langs", default="zh",
        help="comma-separated languages to generate, e.g. zh,en (default: zh)"
    )
    args = ap.parse_args()
    langs = parse_langs(args.langs)

    for lang in langs:
        print(f"\n>>> Generating language = {lang}")
        gen_types(lang=lang, force=args.force)
        gen_pairs(lang=lang, force=args.force)

    print("\nDone.")

if __name__ == "__main__":
    main()