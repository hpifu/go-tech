Feature: POST /article

    Scenario: article case1
        Given mysql 执行
            """
            DELETE FROM articles WHERE title IN ('标题1')
            """
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 123,
                "email": "hatlonely@foxmail.com"
            }
            """
        When http 请求 POST /article
            """
            {
                "header": {
                    "Authorization": "d571bda90c2d4e32a793b8a1ff4ff984"
                },
                "json": {
                    "title": "标题1",
                    "tags": [
                        "c++",
                        "java"
                    ],
                    "content": "hello world"
                }
            }
            """
        Then http 检查 201
        Then mysql 检查 "SELECT * FROM articles WHERE title='标题1'"
            """
            {
                "title": "标题1",
                "author_id": 123,
                "tags": "c++,java",
                "author": "hatlonely",
                "content": "hello world"
            }
            """
        Given mysql 执行
            """
            DELETE FROM articles WHERE title IN ('标题1')
            """
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"
