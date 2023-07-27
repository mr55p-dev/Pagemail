import sys
from pprint import pprint
import json
from typing import Self, TextIO
import boto3
# import bs4
import logging
from html.parser import HTMLParser
from contextlib import closing

class Node:
    def __init__(self, tag: str):
        self.tag = tag
        self.children: list[Self | str] = []

    def isTextNode(self) -> bool:
        return len(self.children) == 1 and isinstance(self.children[0], str)

class Parser(HTMLParser):
    IGNORE = ["div", "main", "article", "a", "header", ""]
    # IGNORE = []

    def __init__(self, *, convert_charrefs: bool = True) -> None:
        self.root = Node("root")
        self.current_tag = [self.root]
        super().__init__(convert_charrefs=convert_charrefs)

    def use_tag(self, tag):
        return tag not in Parser.IGNORE

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        if not self.use_tag(tag):
            return

        new = Node(tag)
        self.current_tag[-1].children.append(new)
        self.current_tag.append(new)
        return

    def handle_endtag(self, tag: str) -> None:
        if not self.use_tag(tag):
            return

        self.current_tag.pop()
        return

    def handle_data(self, data: str) -> None:
        ctag = self.current_tag[-1]
        if len(ctag.children) and isinstance(ctag.children[-1], str):
            ctag.children[-1] += data.replace("\n", "")
        else:
            ctag.children.append(data.replace("\n", ""))

    def recurse_tree(self):
        self._recurse(self.root, 0)

    def _recurse(self, node: Node, indent_level):
        space = "  " * indent_level
        space_txt = "  " * (indent_level + 1)
        if node.isTextNode():
            print(space_txt + "".join(node.children))
            return

        print(space + "<" + node.tag + ">")
        for child in node.children:
            if isinstance(child, str):
                print(space_txt + child)
            else:
                self._recurse(child, indent_level + 1)
        print(space + "</" + node.tag + ">")

def replace(start: str, init: list[str], repl: list[str]) -> str:
    n = start
    for l, r in zip(init, repl):
        n = start.replace(l, r)
    return n

def preprocess_html_content(content: str) -> str:
    csub = replace(content, ['"', "&", "'", "<", ">"], ["&quot;", "&amp;", "&apos;", "&lt;", "&gt;"])
    return csub

def main(text = None):
    session = boto3.Session()
    client = session.client("polly")

    # Load the data 
    if not text:
        text = sys.stdin.read()
        data = json.loads(text)
        with open("bootstrap_py", "w") as f:
            f.write(json.dumps(data))

    # Parse the json 
    data = json.loads(text)
    with open("test.html", "w") as f:
        f.write(data["content"])
    
    parser = Parser()
    parser.feed(data["content"])
    rval = parser.recurse_tree()
    pprint(rval)

    # res = client.synthesize_speech(Text="Hello, world!", OutputFormat="mp3", VoiceId="Joanna")
    # with closing(res["AudioStream"]) as stream, open("out.mp3", "wb") as f:
    #     f.write(stream.read())

    
if __name__ == "__main__":
    data = None
    if len(sys.argv) > 1:
        with open("bootstrap_py", "r") as f:
            data = f.read()

    main(data)
