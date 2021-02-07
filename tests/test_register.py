import json
import pytest
import pytest_depends

from API.app import app
from API.db.connection import get_db
from fastapi.testclient import TestClient

from tests.data_gen import get_mock_data, new_user
from tests.mock_database import override_get_db

from uuid import UUID
from datetime import datetime


client = TestClient(app)
app.dependency_overrides[get_db] = override_get_db
(mock_users, mock_pages) = get_mock_data()
headers = {"Content-Type": "application/x-www-form-urlencoded"}

@pytest.mark.depends(name="register")
def test_register_good():
    response = client.post(
        '/user/register',
        headers=headers,
        data=new_user
    )
    assert response.status_code == 200, response.text
    data = response.json()
    try:
        UUID(data["user"]["id"])
    except ValueError as e:
        raise e
    assert data["user"]["name"] == new_user["name"]
    assert data["user"]["email"] == new_user["email"]
    assert data["user"]["date_added"]

@pytest.mark.depends(on=['register'])
def test_register_duplicate_email():
    response = client.post(
        '/user/register',
        headers=headers,
        data=new_user
    )
    data = response.json()
    assert response.status_code == 400
    assert data["detail"] == "Username already exists."

def test_register_bad_password():
    test_user = new_user.copy()
    test_user["password"] = None
    response = client.post(
        '/user/register',
        headers=headers,
        data=test_user
    )
    assert response.status_code == 422

def test_register_bad_email():
    test_user = new_user.copy()
    test_user["email"] = None
    response = client.post(
        '/user/register',
        headers=headers,
        data=test_user
    )
    assert response.status_code == 422
def test_register_bad_name():
    test_user = new_user.copy()
    test_user["name"] = None
    response = client.post(
        '/user/register',
        headers=headers,
        data=test_user
    )
    assert response.status_code == 422

