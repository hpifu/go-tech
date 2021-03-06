Feature: PUT /article

    Scenario: case update success
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (456)
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"
        Given mysql 执行
            """
            INSERT INTO articles (id, title, author_id, content)
            VALUES (456, "标题1", 123, "hello world")
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 123,
                "email": "hatlonely@foxmail.com"
            }
            """
        When http 请求 PUT /article/456
            """
            {
                "header": {
                    "Authorization": "d571bda90c2d4e32a793b8a1ff4ff984"
                },
                "json": {
                    "tags": [
                        "c++",
                        "golang"
                    ]
                }
            }
            """
        Then http 检查 202
        Then mysql 检查 "SELECT * FROM articles WHERE id=456"
            """
            {
                "title": "标题1",
                "author_id": 123,
                "tags": "c++,golang",
                "content": "hello world"
            }
            """
        When http 请求 PUT /article/456
            """
            {
                "header": {
                    "Authorization": "d571bda90c2d4e32a793b8a1ff4ff984"
                },
                "json": {
                    "title": "标题abc",
                    "content": "hello golang"
                }
            }
            """
        Then http 检查 202
        Then mysql 检查 "SELECT * FROM articles WHERE id=456"
            """
            {
                "title": "标题abc",
                "author_id": 123,
                "tags": "c++,golang",
                "content": "hello golang"
            }
            """
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (456)
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: case token 不正确
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (456)
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"
        Given mysql 执行
            """
            INSERT INTO articles (id, title, author_id, content)
            VALUES (456, "标题1", 123, "hello world")
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 123,
                "email": "hatlonely@foxmail.com"
            }
            """
        When http 请求 PUT /article/456
            """
            {
                "header": {
                    "Authorization": "wrong token"
                },
                "json": {
                    "tags": [
                        "c++",
                        "golang"
                    ]
                }
            }
            """
        Then http 检查 403
            """
            {
                "text": "没有该资源权限"
            }
            """
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (456)
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: case 该资源不属于该用户
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (456)
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"
        Given mysql 执行
            """
            INSERT INTO articles (id, title, author_id, content)
            VALUES (456, "标题1", 123, "hello world")
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 124,
                "email": "hatlonely@foxmail.com"
            }
            """
        When http 请求 PUT /article/456
            """
            {
                "header": {
                    "Authorization": "d571bda90c2d4e32a793b8a1ff4ff984"
                },
                "json": {
                    "title": "标题abc",
                    "content": "hello golang"
                }
            }
            """
        Then http 检查 403
            """
            {
                "text": "没有该资源权限"
            }
            """
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (456)
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"