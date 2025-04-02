import pprint
import requests

backend_url = "http://localhost:8080"
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

def CreateWorkOrders():
  body = {
    "room_id": 1,
    "problem": "pipe leakage"
  }
  resp = requests.post(f"{backend_url}/create-work-order", json=body, headers=headers)
  print(resp.json())
    
def GetRooms():
  resp = requests.get(f"{backend_url}/get-rooms")
  print(resp)
  if resp.status_code == 200:
    print(resp.json())
    
def ListOwnedRooms():
  resp = requests.get(f"{backend_url}/list-owned-rooms", headers=headers)
  print(resp)
  if resp.status_code == 200:
    print(resp.json())


if __name__ == "__main__":
  # Register()
  Login()
  # GetRooms()
  # RegisterAdmin()
  # CreateWorkOrders()
  ListOwnedRooms()
  
  