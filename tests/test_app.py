from fastapi.testclient import TestClient
from API.app import app

client = TestClient(app)

def test_main_root():
    response = client.get("/")
    assert response.status_code == 200
    assert response.text == '"Hello World."'