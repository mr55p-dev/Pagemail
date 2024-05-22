from flask import Flask, request
from bs4 import BeautifulSoup
from readability import Document
from boto3 import Session
from botocore.exceptions import BotoCoreError, ClientError, ParamValidationError
from contextlib import closing
import os


session = Session(region_name="eu-west-1")
polly = session.client("polly")

bucket_name = os.getenv("BUCKET_NAME")
if not bucket_name:
    raise ValueError("BUCKET_NAME environment variable is not set")

def polly_synthesize(text: str) -> dict: 
    try:
        response = polly.start_speech_synthesis_task(
            Engine="standard",
            LanguageCode="en-US",
            VoiceId="Joanna",
            OutputFormat="mp3",
            OutputS3BucketName="pagemail-readability",
            Text=text,
            TextType="text",
        )
    except (BotoCoreError, ClientError, ParamValidationError) as error:
        return {"error": str(error)}

    try :
        return {"task_id": response["SynthesisTask"]["TaskId"]}
    except KeyError:
        return {"error": response}

app = Flask(__name__)

def extract_text_from_html(html_content):
    doc = Document(html_content)
    soup = BeautifulSoup(doc.summary(), 'html.parser')
    return soup.get_text()

@app.route('/extract', methods=['POST'])
def extract():
    html_content = request.data.decode('utf-8')
    article_text = extract_text_from_html(html_content)
    print(article_text)
    
    return article_text

@app.route("/synthesize", methods=["POST"])
def synthesize():
    text = request.data.decode("utf-8")
    res = polly_synthesize(text)
    if "error" in res:
        return res, 500
    return res

if __name__ == '__main__':
    app.run(debug=True)
