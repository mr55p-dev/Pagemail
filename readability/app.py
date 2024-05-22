from flask import Flask, request
from bs4 import BeautifulSoup
from readability import Document

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

if __name__ == '__main__':
    app.run(debug=True)
