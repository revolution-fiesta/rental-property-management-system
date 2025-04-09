from locust import HttpUser, task, between
import json

class MyLocust(HttpUser):

    @task
    def login(self):
        # 请求体内容
        headers = {'Content-Type': 'application/json'}
        data = {
            "username": "ysj123",
            "password": "123456"
        }
        
        # 使用 POST 请求发送 JSON 数据
        self.client.post("/login", json=data, headers=headers)
