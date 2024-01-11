from html.parser import HTMLParser
from parser.utils import preprocess_html_content

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
        self._depth = 0

    def handle_starttag(self, tag: str, _: list[tuple[str, str | None]]) -> None:
        self._depth += 1
        if tag in Parser.KEEP:
            self._output_stream.append(Tag(tag, Tag.OPEN))
        elif tag in Parser.REMAP_OPEN:
            replacement = Parser.REMAP_OPEN[tag]
            self._output_stream.append(replacement)

    def handle_endtag(self, tag: str) -> None:
        self._depth -= 1
        if tag in Parser.KEEP:
            self._output_stream.append(Tag(tag, Tag.CLOSE))
        if tag in Parser.REMAP_CLOSE:
            replacement = Parser.REMAP_CLOSE[tag]
            if replacement.isopen == Tag.OPENCLOSE:
                return
            self._output_stream.append(replacement)

    def handle_data(self, data: str) -> None:
        processed_data = preprocess_html_content(data)
        if not data or not data.strip():
            return
        elif len(self._output_stream) and isinstance(self._output_stream[-1], str):
            self._output_stream[-1] += processed_data
        elif processed_data:
            self._output_stream.append(processed_data)

    @property
    def output_stream(self):
        return "".join(["<speak>", *(str(i) for i in self._output_stream), "</speak>"])
