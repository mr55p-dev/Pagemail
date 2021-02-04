from API.db.connection import get_db
from fastapi.testclient import TestClient
from API.app import app
from tests.mock_database import override_get_db
from tests.data_gen import new_user

import pytest_depends
import pytest

client = TestClient(app)
app.dependency_overrides[get_db] = override_get_db
mock_user = {"username": "mock@username.com", "password": "mock_password"}

@pytest.mark.depends(on=['register'])
def test_save():
    response = client.post(
        "/page/save",

    )

def test_mypages():
    pass
