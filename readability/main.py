from parser import Parser, Dispatcher

import sys
import traceback
import json
import boto3


def run(session, input: str):
    parser = Parser()
    dispatch = Dispatcher(session)

    data = json.loads(input)
    text = data["content"]
    
    parser.feed(text)
    parsed = parser.output_stream

    out, err = dispatch.create_job(parsed)
    if err:
        traceback.print_exception(err, file=sys.stderr)
        sys.exit(1)

    sys.stdout.write(out)
    sys.exit(0)

def main():
    input = sys.stdin.read()
    session = boto3.Session()

    run(session, input)

    
if __name__ == "__main__":
    main()
