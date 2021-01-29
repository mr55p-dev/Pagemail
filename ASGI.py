import uvicorn, os, dotenv
dotenv.load_dotenv()

if __name__ == '__main__':
    uvicorn.run("API.app:app", host="0.0.0.0", port=int(os.getenv('PORT')))