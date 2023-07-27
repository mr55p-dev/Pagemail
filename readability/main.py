from functools import reduce
import re
import sys
from pprint import pprint
import json
from typing import Self, TextIO
import boto3
# import bs4
import logging
from html.parser import HTMLParser
from contextlib import closing

class Tag():
    def __init__(self, tagname: str, isopen: bool) -> None:
        self.name = tagname
        self.isopen = isopen

    def __repr__(self) -> str:
        if self.isopen:
            return f"<{self.name}>"
        else:
            return f"</{self.name}>"


class Parser(HTMLParser):
    IGNORE = ["div", "main", "article", "a", "header", ""]
    KEEP = ["section", "p", "article", "h2"]
    REMAP = {
        "article": "speak",
        "strong": "emphasis",
        "em": "emphasis",
    }

    def __init__(self, *, convert_charrefs: bool = True) -> None:
        super().__init__(convert_charrefs=convert_charrefs)
        self.output_stream = []
        self.tagstack = []

    def use_tag(self, tag):
        return tag not in Parser.IGNORE

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        if tag in Parser.KEEP:
            self.output_stream.append(Tag(tag, True))
            self.tagstack.append(True)
        else:
            self.tagstack.append(False)

    def handle_endtag(self, tag: str) -> None:
        if self.tagstack.pop():
            self.output_stream.append(Tag(tag, False))

    def handle_data(self, data: str) -> None:
        processed_data = preprocess_html_content(data)
        if len(self.output_stream) and isinstance(self.output_stream[-1], str):
            self.output_stream[-1] += processed_data
        elif processed_data:
            self.output_stream.append(processed_data)

def replace(start: str, init: list[str], repl: list[str]) -> str:
    n = start
    for l, r in zip(init, repl):
        n = n.replace(l, r)
    return n


def preprocess_html_content(content: str) -> str:
    csub = replace(content, ["\"", r"&", r"'", r"<", r">"], ["&quot;", "&amp;", "&apos;", "&lt;", "&gt;"])
    res = re.sub(r'^\s+$', r'', csub)
    return res

def main(text = None):
    # session = boto3.Session()
    # client = session.client("polly")

    # Load the data 
    if not text:
        text = sys.stdin.read()
        data = json.loads(text)
        with open("bootstrap_py", "w") as f:
            f.write(json.dumps(data))

    # Parse the json 
    data = json.loads(text)
    # with open("test.html", "w") as f:
    #     f.write(data["content"])
    
    parser = Parser()
    parser.feed(data["content"])
    print("\n".join(str(i) for i in parser.output_stream))

    # res = client.synthesize_speech(Text="Hello, world!", OutputFormat="mp3", VoiceId="Joanna")
    # with closing(res["AudioStream"]) as stream, open("out.mp3", "wb") as f:
    #     f.write(stream.read())

    
if __name__ == "__main__":
    data = None
    # if len(sys.argv) > 1:
    #     with open("bootstrap_py", "r") as f:
    #         data = f.read()

    main(data)
