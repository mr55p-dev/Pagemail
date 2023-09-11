from typing import Any
import boto3 
import json
import datetime

class Dispatcher:
    def __init__(self, session: boto3.Session):
        self._session = session
        self._client = self._session.client("polly")

    def create_job(self, text: str) -> tuple[Any, Exception | None]:
        try:
            # https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/polly/client/start_speech_synthesis_task.html#
            response = self._client.start_speech_synthesis_task(
                Engine='standard',
                LanguageCode='en-GB',
                OutputFormat='mp3',
                OutputS3BucketName='pagemail-speechsynthesis',
                Text=text,
                TextType='ssml',
                VoiceId='Amy',
            )
            response["SynthesisTask"]["CreationTime"] = response["SynthesisTask"]["CreationTime"].astimezone(datetime.UTC).strftime("%Y-%m-%dT%H:%M:%SZ")
            return json.dumps(response), None
        except Exception as e:
            return None, e
