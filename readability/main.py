import re
import sys
import json
import boto3
import base64
import time
from html.parser import HTMLParser
import logging

class Tag():
    OPEN = 1
    CLOSE = 2
    OPENCLOSE = 3
    def __init__(self, tagname: str, isopen: int) -> None:
        self.name = tagname
        self.isopen = isopen

    def __repr__(self) -> str:
        match self.isopen:
            case Tag.CLOSE:
                return f"</{self.name}>"
            case Tag.OPENCLOSE:
                return f"<{self.name}/>"
            case Tag.OPEN:
                return f"<{self.name}>"
            case _:
                return f"<{self.name}>"


class Parser(HTMLParser):
    KEEP = ["p"]

    REMAP_OPEN = {
        "section": Tag("break", Tag.OPENCLOSE),
        "article": Tag("speak", Tag.OPEN),
        "strong": Tag("emphasis", Tag.OPEN),
        "em": Tag("emphasis", Tag.OPEN),
        "h1": Tag("p", Tag.OPEN),
        "h2": Tag("p", Tag.OPEN),
        "h3": Tag("p", Tag.OPEN),
        "h4": Tag("p", Tag.OPEN),
        "h5": Tag("p", Tag.OPEN),
        "h6": Tag("p", Tag.OPEN),
    }
    
    REMAP_CLOSE = {
        "article": Tag("speak", Tag.CLOSE),
        "strong": Tag("emphasis", Tag.CLOSE),
        "em": Tag("emphasis", Tag.CLOSE),
        "h1": Tag("p", Tag.CLOSE),
        "h2": Tag("p", Tag.CLOSE),
        "h3": Tag("p", Tag.CLOSE),
        "h4": Tag("p", Tag.CLOSE),
        "h5": Tag("p", Tag.CLOSE),
        "h6": Tag("p", Tag.CLOSE),
    }

    def __init__(self, *, convert_charrefs: bool = True) -> None:
        super().__init__(convert_charrefs=convert_charrefs)
        self._output_stream = []
        self.tagstack = []

    def handle_starttag(self, tag: str, _: list[tuple[str, str | None]]) -> None:
        if tag in Parser.KEEP:
            self._output_stream.append(Tag(tag, Tag.OPEN))
        elif tag in Parser.REMAP_OPEN:
            replacement = Parser.REMAP_OPEN[tag]
            self._output_stream.append(replacement)
            if replacement.isopen == Tag.OPENCLOSE:
                self.tagstack.append(False)
                return 
        else:
            self.tagstack.append(False)
            return
        self.tagstack.append(True)

    def handle_endtag(self, tag: str) -> None:
        if not self.tagstack.pop():
            return

        if tag in Parser.REMAP_CLOSE:
            self._output_stream.append(Parser.REMAP_CLOSE[tag])
        else:
            self._output_stream.append(Tag(tag, Tag.CLOSE))

    def handle_data(self, data: str) -> None:
        processed_data = preprocess_html_content(data)
        if len(self._output_stream) and isinstance(self._output_stream[-1], str):
            self._output_stream[-1] += processed_data
        elif processed_data:
            self._output_stream.append(processed_data)

    @property
    def output_stream(self):
        return "".join(str(i) for i in self._output_stream)

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

def main():
    session = boto3.Session()
    client = session.client("polly")

    # Load the data 
    text = sys.stdin.read()
    data = json.loads(text)

    # Parse the json 
    data = json.loads(text)
    
    parser = Parser()
    parser.feed(data["content"])
    inp = parser.output_stream

    time.sleep(15)

    try:
        # https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/polly/client/start_speech_synthesis_task.html#
        response = client.start_speech_synthesis_task(
            Engine='standard',
            LanguageCode='en-GB',
            OutputFormat='mp3',
            OutputS3BucketName='pagemail-speechsynthesis',
            Text=inp,
            TextType='ssml',
            VoiceId='Amy',
        )
        sys.stdout.write(json.dumps(response["SynthesisTask"], default=str))
        sys.stdout.flush()
        sys.stdout.write("""{
"ResponseMetadata": {
    "RequestId": "6b17d1a6-def2-4f22-9d64-0e71101f8c13",
    "HTTPStatusCode": 200,
    "HTTPHeaders": { "x-amzn-requestid": "6b17d1a6-def2-4f22-9d64-0e71101f8c13", "content-type": "application/json", "content-length": "472", "date": "Thu, 27 Jul 2023 23:48:54 GMT" },
    "RetryAttempts": 0
},
"SynthesisTask": {
    "Engine": "standard",
    "TaskId": "d75be692-5f58-4534-b7fe-4d6e51c53a51",
    "TaskStatus": "scheduled",
    "OutputUri": "https://s3.eu-west-2.amazonaws.com/pagemail-speechsynthesis/d75be692-5f58-4534-b7fe-4d6e51c53a51.mp3",
    "RequestCharacters": 1398,
    "OutputFormat": "mp3",
    "TextType": "text",
    "VoiceId": "Amy",
    "LanguageCode": "en-GB"
}}""")
        sys.stdout.flush()
    except Exception as e:
        logging.exception(e)
        logging.fatal("Failed to create new speech synthesis task")
        sys.exit(1)

    logging.info("Finished job")


    
if __name__ == "__main__":
    logging.basicConfig(filename="parser.log")
    main()
