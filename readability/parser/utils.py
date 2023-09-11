import re


def replace(start: str, init: list[str], repl: list[str]) -> str:
    n = start
    for l, r in zip(init, repl):
        n = n.replace(l, r)
    return n

def preprocess_html_content(content: str) -> str:
    csub = replace(content, ["\"", r"&", r"'", r"<", r">"], ["&quot;", "&amp;", "&apos;", "&lt;", "&gt;"])
    res = re.sub(r'^\s+$', r' ', csub)
    res = re.sub(r'\s+', ' ', res)
    return res
