Feature: articles 测试

    Scenario: articles
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (1,2)
            """
        Given mysql 执行
            """
            INSERT INTO articles (id, title, author_id, author, content)
            VALUES (1, "标题1", 666, "hatlonely", "hello world")
            """
        Given mysql 执行
            """
            INSERT INTO articles (id, title, author_id, author, content)
            VALUES (2, "标题2", 666, "hatlonely", "hello world")
            """
        When http 请求 GET /article
            """
            {
                "params": {
                    "offset": 0,
                    "limit": 2
                }
            }
            """
        Then http 检查 200
            """
            {
                "json": [
                    {
                        "id": 1,
                        "title": "标题1",
                        "authorID": 666,
                        "author": "hatlonely"
                    },
                    {
                        "id": 2,
                        "title": "标题2",
                        "authorID": 666,
                        "author": "hatlonely"
                    }
                ]
            }
            """
        Given mysql 执行
            """
            DELETE FROM articles WHERE id IN (1,2)
            """