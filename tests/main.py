import requests

backend_url = "https://localhost:8080"
token = ""
headers = {}

def PrintResponse(resp):
  print(f"[{resp.status_code}]: {resp.text}\n")

def Register():
  body = {
    "username": "azusayn",
    "password": "azusayn"
  }
  resp = requests.post(f"{backend_url}/register", json=body)
  print(resp.json())
  
  
def RegisterAdmin():
  body = {
    "username": "timetom790",
    "password": "timetom790"
  }
  resp = requests.post(f"{backend_url}/register-admin", json=body, headers=headers)
  PrintResponse(resp)

def Login():
  global token, headers
  body = {
    "username": "azusayn",
    "password": "azusayn"
  }
  resp = requests.post(f"{backend_url}/login", json=body)
  token = resp.json()['token']
  headers = {
    "Authorization": f"Bearer {token}"
  }
  print(resp)

if __name__ == "__main__":
  # Register()
  Login()
  RegisterAdmin()