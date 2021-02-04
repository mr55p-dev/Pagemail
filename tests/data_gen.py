from datetime import datetime
from uuid import uuid4

new_user = {"email": "new_user@example.com", "name": "Name", "password": "password"}

def get_mock_data():
    mock_users = [
        {
            "id": uuid4(),
            "name": "test 1",
            "email": "example1@examplemail.com",
            "password": "password",
            "date_added": datetime(2021, 1, 1),
            "is_active": True,
        },
        {
            "id": uuid4(),
            "name": "test 2",
            "email": "example2@examplemail.com",
            "password": "password",
            "date_added": datetime(2021, 1, 1),
            "is_active": False,
        },
        {
            "id": uuid4(),
            "name": "test 3",
            "email": "example3@examplemail.com",
            "password": "password",
            "date_added": datetime(2021, 1, 1),
            "is_active": True,
        }
    ]
    mock_pages = [
        {
            "id": uuid4(),
            "page_url": "https://www.example1.com",
            "date_added": datetime(2021, 1, 1),
            "user_id": mock_users[0]["id"]
        },
        {
            "id": uuid4(),
            "page_url": "https://www.example2.com",
            "date_added": datetime(2021, 1, 1),
            "user_id": mock_users[0]["id"]
        },
        {
            "id": uuid4(),
            "page_url": "https://www.example3.com",
            "date_added": datetime(2021, 1, 1),
            "user_id": mock_users[1]["id"]
        },
        {
            "id": uuid4(),
            "page_url": "https://www.example4.com",
            "date_added": datetime(2021, 1, 1),
            "user_id": mock_users[1]["id"]
        }
    ]
    return (mock_users, mock_pages)
